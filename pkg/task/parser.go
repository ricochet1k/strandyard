package task

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v3"
)

// Parser handles parsing task files.
type Parser struct {
}

// NewParser creates a new task parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile parses a task file and returns a Task
func (p *Parser) ParseFile(filePath string) (*Task, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	fileName := filepath.Base(filePath)
	id := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	t, err := p.ParseString(string(data), id)
	if err != nil {
		var fmErr *InvalidFrontmatterError
		if errors.As(err, &fmErr) && fmErr.Path == "" {
			fmErr.Path = filePath
			return nil, fmErr
		}
		// Add file path to FrontmatterParseError if needed
		var parseErr *FrontmatterParseError
		if errors.As(err, &parseErr) && parseErr.Path == "" {
			parseErr.Path = filePath
		}
		return nil, err
	}
	t.FilePath = filePath
	t.Dir = filepath.Dir(filePath)
	return t, nil
}

// ParseStandaloneFile parses a markdown file that is not in a task directory
func (p *Parser) ParseStandaloneFile(filePath string) (*Task, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	fileName := filepath.Base(filePath)
	id := strings.TrimSuffix(fileName, filepath.Ext(fileName))

	t, err := p.ParseString(string(data), id)
	if err != nil {
		var fmErr *InvalidFrontmatterError
		if errors.As(err, &fmErr) && fmErr.Path == "" {
			fmErr.Path = filePath
			return nil, fmErr
		}
		// Add file path to FrontmatterParseError if needed
		var parseErr *FrontmatterParseError
		if errors.As(err, &parseErr) && parseErr.Path == "" {
			parseErr.Path = filePath
		}
		return nil, err
	}
	t.FilePath = filePath
	t.Dir = filepath.Dir(filePath)
	return t, nil
}

// ParseString parses a string into a Task
func (p *Parser) ParseString(content string, id string) (*Task, error) {
	var meta Metadata
	frontmatterText, body, hasFrontmatter, err := splitFrontmatter(content)
	if err != nil {
		return nil, err
	}
	if hasFrontmatter {
		if err := yaml.Unmarshal([]byte(frontmatterText), &meta); err != nil {
			// Extract line number from YAML error message if available
			lineNum := extractLineNumberFromYAMLError(err)
			return nil, &FrontmatterParseError{
				LineNumber: lineNum,
				YAMLError:  err.Error(),
			}
		}
	}

	t := &Task{
		ID:   id,
		Meta: meta,
	}

	// err := ast.Walk(parsed, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
	// 	switch n := n.(type) {
	// 	case *ast.Document:
	// 		return ast.WalkContinue, nil

	// 	case *ast.Text:
	// 		fmt.Printf("node: %v %v %q\n", n.Kind(), n.Attributes(), n.)
	// 		return ast.WalkContinue, nil

	// 	default:
	// 		fmt.Printf("node: %v %v %v\n", n.Kind(), n.Attributes(), n.ChildCount())
	// 		return ast.WalkContinue, nil
	// 	}
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// panic("todo")

	// Split content into sections
	sections := SplitByHeadings(body)

	for _, section := range sections {
		h := strings.ToLower(section.Heading)
		if section.Level == 1 {
			t.TitleContent = section.Heading
			if section.Content != "" {
				t.BodyContent += section.Content + "\n\n"
			}
		} else if h == "todos" || h == "tasks" || h == "subtasks" {
			items := ParseTaskItems(section.Content)
			for _, item := range items {
				if item.SubtaskID != "" {
					t.SubsItems = append(t.SubsItems, item)
				} else {
					t.TodoItems = append(t.TodoItems, item)
				}
			}
		} else if h == "progress" {
			t.ProgressContent = section.Content
		} else if h == "" {
			if t.BodyContent != "" {
				t.BodyContent += "\n\n"
			}
			t.BodyContent += section.Content
		} else {
			if t.BodyContent != "" {
				t.BodyContent += "\n\n"
			}
			t.BodyContent += "## " + section.Heading + "\n" + section.Content
		}
	}
	t.BodyContent = strings.TrimSpace(t.BodyContent)

	return t, nil
}

func splitFrontmatter(content string) (string, string, bool, error) {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return "", "", false, nil
	}
	if !isFrontmatterDelimiter(lines[0]) {
		return "", strings.TrimSpace(content), false, nil
	}

	end := -1
	for i := 1; i < len(lines); i++ {
		if isFrontmatterDelimiter(lines[i]) {
			end = i
			break
		}
	}
	if end == -1 {
		return "", "", true, &InvalidFrontmatterError{Path: ""}
	}

	frontmatterLines := make([]string, 0, end-1)
	for i := 1; i < end; i++ {
		frontmatterLines = append(frontmatterLines, strings.TrimSuffix(lines[i], "\r"))
	}
	bodyLines := make([]string, 0, len(lines)-end-1)
	for i := end + 1; i < len(lines); i++ {
		bodyLines = append(bodyLines, strings.TrimSuffix(lines[i], "\r"))
	}

	frontmatterText := strings.Join(frontmatterLines, "\n")
	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))
	return frontmatterText, body, true, nil
}

func isFrontmatterDelimiter(line string) bool {
	trimmed := strings.TrimSuffix(line, "\r")
	if len(trimmed) < 3 {
		return false
	}
	for _, ch := range trimmed {
		if ch != '-' {
			return false
		}
	}
	return true
}

// extractLineNumberFromYAMLError extracts the line number from a YAML error message.
// YAML error format: "yaml: line N: <error message>"
func extractLineNumberFromYAMLError(err error) int {
	errMsg := err.Error()
	// Look for "line N:" pattern in the error message
	// Example: "yaml: line 2: could not find expected ':'"
	parts := strings.Split(errMsg, "line ")
	if len(parts) > 1 {
		// Extract the number after "line "
		rest := parts[1]
		colonIdx := strings.Index(rest, ":")
		if colonIdx > 0 {
			numStr := rest[:colonIdx]
			// Try to parse as integer
			if n, err := strconv.Atoi(strings.TrimSpace(numStr)); err == nil {
				return n
			}
		}
	}
	return 0 // Return 0 if we can't extract the line number
}

// Section represents a markdown section.
type Section struct {
	Level   int
	Heading string
	Content string
}

// SplitByHeadings splits markdown content into sections.
func SplitByHeadings(body string) []Section {
	var sections []Section
	lines := strings.Split(body, "\n")
	var currentSection *Section

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			if currentSection != nil {
				sections = append(sections, *currentSection)
			}
			currentSection = &Section{
				Level:   1,
				Heading: strings.TrimSpace(trimmed[2:]),
			}
		} else if strings.HasPrefix(trimmed, "## ") {
			if currentSection != nil {
				sections = append(sections, *currentSection)
			}
			currentSection = &Section{
				Level:   2,
				Heading: strings.TrimSpace(trimmed[3:]),
			}
		} else {
			if currentSection == nil {
				currentSection = &Section{Level: 0}
			}
			if currentSection.Content != "" {
				currentSection.Content += "\n"
			}
			currentSection.Content += line
		}
	}

	if currentSection != nil {
		sections = append(sections, *currentSection)
	}

	for i := range sections {
		sections[i].Content = strings.TrimSpace(sections[i].Content)
	}

	return sections
}

// ExtractTitle finds the first H1 in the content.
func ExtractTitle(content string) string {
	sections := SplitByHeadings(content)
	for _, s := range sections {
		if s.Level == 1 {
			return s.Heading
		}
	}
	return ""
}

// LoadTasks walks the tasks directory and loads all tasks, in parallel.
func (p *Parser) LoadTasks(tasksRoot string) (map[string]*Task, error) {
	tasksChan := make(chan *Task)
	errChan := make(chan error)

	go func() {
		var eg errgroup.Group
		eg.SetLimit(10)

		err := filepath.WalkDir(tasksRoot, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// Skip the root directory
			if path == tasksRoot {
				return nil
			}

			// Skip directories
			if d.IsDir() {
				return nil
			}

			// Only process .md files
			if filepath.Ext(path) != ".md" {
				return nil
			}

			// Skip master lists
			fileName := d.Name()
			if fileName == "root-tasks.md" || fileName == "free-tasks.md" {
				return nil
			}

			eg.Go(func() error {
				// Parse the task
				task, err := p.ParseFile(path)
				if err != nil {
					return fmt.Errorf("failed to parse task %s: %w", path, err)
				}

				tasksChan <- task
				return nil
			})

			return nil
		})

		wgErr := eg.Wait()
		close(tasksChan)

		errChan <- errors.Join(err, wgErr)
	}()

	tasks := make(map[string]*Task)
	for task := range tasksChan {
		tasks[task.ID] = task
	}

	err := <-errChan

	return tasks, err
}
