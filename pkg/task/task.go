package task

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

// Metadata represents the YAML frontmatter of a task
type Metadata struct {
	Kind          string    `yaml:"kind"`
	Role          string    `yaml:"role"`
	Priority      string    `yaml:"priority"`
	Parent        string    `yaml:"parent"`
	Blockers      []string  `yaml:"blockers"`
	Blocks        []string  `yaml:"blocks"`
	DateCreated   time.Time `yaml:"date_created"`
	DateEdited    time.Time `yaml:"date_edited"`
	OwnerApproval bool      `yaml:"owner_approval"`
	Completed     bool      `yaml:"completed"`
}

// Task represents a complete task with metadata and content
type Task struct {
	ID       string
	Dir      string
	FilePath string
	Meta     Metadata
	Content  string
	Document ast.Node
}

// Title returns the first level-1 heading text, or empty string if not found.
func (t *Task) Title() string {
	lines := strings.Split(t.Content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

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

	// Parse the markdown with frontmatter
	var meta Metadata
	ctx := parser.NewContext()
	doc := p.md.Parser().Parse(text.NewReader(data), parser.WithContext(ctx))

	// Extract frontmatter
	fm := frontmatter.Get(ctx)
	if fm != nil {
		if err := fm.Decode(&meta); err != nil {
			return nil, fmt.Errorf("failed to decode frontmatter in %s: %w", filePath, err)
		}
	}

	// Extract task ID from directory name
	dir := filepath.Dir(filePath)
	id := filepath.Base(dir)

	task := &Task{
		ID:       id,
		Dir:      dir,
		FilePath: filePath,
		Meta:     meta,
		Content:  string(data),
		Document: doc,
	}

	return task, nil
}

// LoadTasks walks the tasks directory and loads all tasks
func (p *Parser) LoadTasks(tasksRoot string) (map[string]*Task, error) {
	tasks := make(map[string]*Task)

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
			filepath.Join(path, "README.md"),
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

		// Parse the task
		task, err := p.ParseFile(taskFile)
		if err != nil {
			return fmt.Errorf("failed to parse task %s: %w", taskFile, err)
		}

		// Store by ID
		tasks[task.ID] = task

		return nil
	})

	return tasks, err
}

// ExtractFirstTodoRole extracts the role from the first TODO item in task content
// Format: - [ ] (role: developer) Do something
func ExtractFirstTodoRole(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Look for TODO pattern: - [ ] (role: xxx)
		if strings.HasPrefix(line, "- [ ]") || strings.HasPrefix(line, "- [x]") {
			// Find (role: xxx) pattern
			start := strings.Index(line, "(role:")
			if start == -1 {
				continue
			}
			start += 6 // len("(role:")
			end := strings.Index(line[start:], ")")
			if end == -1 {
				continue
			}
			role := strings.TrimSpace(line[start : start+end])
			if role != "" {
				return role
			}
		}
	}
	return ""
}

// GetEffectiveRole returns the task's role, checking metadata first, then first TODO
func (t *Task) GetEffectiveRole() string {
	if t.Meta.Role != "" {
		return t.Meta.Role
	}
	return ExtractFirstTodoRole(t.Content)
}
