package task

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
	"golang.org/x/sync/errgroup"
)

// Parser handles parsing task files using goldmark
type Parser struct {
	md goldmark.Markdown
}

// NewParser creates a new task parser
func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
	)
	return &Parser{md: md}
}

// ParseFile parses a task file and returns a Task
func (p *Parser) ParseFile(filePath string) (*Task, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	dir := filepath.Dir(filePath)
	id := filepath.Base(dir)

	t, err := p.ParseString(string(data), id)
	if err != nil {
		return nil, err
	}
	t.FilePath = filePath
	t.Dir = dir
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
		return nil, err
	}
	t.FilePath = filePath
	t.Dir = filepath.Dir(filePath)
	return t, nil
}

// ParseString parses a string into a Task
func (p *Parser) ParseString(content string, id string) (*Task, error) {
	// Parse the markdown with frontmatter
	var meta Metadata
	ctx := parser.NewContext()
	_ = p.md.Parser().Parse(text.NewReader([]byte(content)), parser.WithContext(ctx))

	// Extract frontmatter
	fm := frontmatter.Get(ctx)
	if fm != nil {
		if err := fm.Decode(&meta); err != nil {
			return nil, fmt.Errorf("failed to decode frontmatter: %w", err)
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
	body := ""
	parts := strings.SplitN(content, "---", 3)
	if len(parts) >= 3 {
		body = strings.TrimSpace(parts[2])
	} else {
		body = strings.TrimSpace(content)
	}

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

			// Only process directories
			if !d.IsDir() {
				return nil
			}

			// Look for task file in this directory
			// Priority: <task-id>.md, then task.md, then README.md
			dirName := filepath.Base(path)
			taskFile := ""

			candidates := []string{
				filepath.Join(path, dirName+".md"),
				filepath.Join(path, "task.md"),
			}

			for _, candidate := range candidates {
				if _, err := os.Stat(candidate); err == nil {
					taskFile = candidate
					break
				}
			}

			if taskFile == "" {
				// No task file found, skip this directory
				return nil
			}

			eg.Go(func() error {
				// Parse the task
				task, err := p.ParseFile(taskFile)
				if err != nil {
					return fmt.Errorf("failed to parse task %s: %w", taskFile, err)
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
