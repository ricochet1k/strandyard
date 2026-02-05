package task

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Metadata represents the YAML frontmatter of a task
type Metadata struct {
	Type          string    `yaml:"type"`
	Role          string    `yaml:"role"`
	Priority      string    `yaml:"priority"`
	Parent        string    `yaml:"parent"`
	Blockers      []string  `yaml:"blockers"`
	Blocks        []string  `yaml:"blocks"`
	DateCreated   time.Time `yaml:"date_created"`
	DateEdited    time.Time `yaml:"date_edited"`
	OwnerApproval bool      `yaml:"owner_approval"`
	Completed     bool      `yaml:"completed"`
	Every         []string  `yaml:"every,omitempty"`
	Description   string    `yaml:"description"`
}

// Task represents a complete task with metadata and content
type Task struct {
	ID              string
	Dir             string
	FilePath        string
	Meta            Metadata
	TitleContent    string
	BodyContent     string
	TodoItems       []TaskItem
	SubsItems       []TaskItem
	ProgressContent string
	OtherContent    string
	Dirty           bool
}

// SetTitle updates the task title.
func (t *Task) SetTitle(newTitle string) {
	if t.TitleContent == newTitle {
		return
	}
	t.TitleContent = newTitle
	t.MarkDirty()
}

// SetBody replaces the body content, preserving title and special sections.
func (t *Task) SetBody(newBody string) {
	sections := SplitByHeadings(newBody)
	var cleanBody strings.Builder
	for _, s := range sections {
		h := strings.ToLower(s.Heading)
		if h == "todos" || h == "tasks" || h == "subtasks" || h == "progress" {
			continue
		}
		if s.Level == 1 {
			// Skip title in SetBody as requested
			continue
		}
		if s.Heading != "" {
			cleanBody.WriteString(fmt.Sprintf("## %s\n", s.Heading))
		}
		cleanBody.WriteString(s.Content)
		cleanBody.WriteString("\n\n")
	}
	newBodyContent := strings.TrimSpace(cleanBody.String())
	if t.BodyContent == newBodyContent {
		return
	}
	t.BodyContent = newBodyContent
	t.MarkDirty()
}

// MarkDirty marks the task as modified.
func (t *Task) MarkDirty() {
	if !t.Dirty {
		// if _, file, line, ok := runtime.Caller(1); ok {
		// 	fmt.Printf("MarkDirty called from %v:%v\n", file, line)
		// }
		t.Meta.DateEdited = time.Now().UTC()
	}
	t.Dirty = true
}

// Write persists updated metadata to the task file.
func (t *Task) Write() error {
	newContent := t.Content()
	if err := os.WriteFile(t.FilePath, []byte(newContent), 0o644); err != nil {
		return err
	}

	t.Dirty = false
	return nil
}

// Title returns the task title.
func (t *Task) Title() string {
	return t.TitleContent
}

// Content returns the full task content as it would be written to file.
func (t *Task) Content() string {
	var sb strings.Builder

	frontmatterBytes, _ := yaml.Marshal(&t.Meta)
	sb.WriteString("---\n")
	sb.Write(frontmatterBytes)
	sb.WriteString("---\n\n")

	if t.TitleContent != "" {
		sb.WriteString("# ")
		sb.WriteString(t.TitleContent)
		sb.WriteString("\n\n")
	}

	if t.BodyContent != "" {
		sb.WriteString(t.BodyContent)
		sb.WriteString("\n\n")
	}

	if len(t.TodoItems) > 0 {
		sb.WriteString("## TODOs\n")
		sb.WriteString(FormatTodoItems(t.TodoItems))
		sb.WriteString("\n\n")
	}

	if len(t.SubsItems) > 0 {
		sb.WriteString("## Subtasks\n")
		sb.WriteString(FormatSubtaskItems(t.SubsItems))
		sb.WriteString("\n\n")
	}

	if t.ProgressContent != "" {
		sb.WriteString("## Progress\n")
		sb.WriteString(t.ProgressContent)
		sb.WriteString("\n\n")
	}

	if t.OtherContent != "" {
		sb.WriteString(t.OtherContent)
		sb.WriteString("\n")
	}

	return strings.TrimRight(sb.String(), "\n") + "\n"
}

// GetEffectiveRole returns the task's role, checking metadata first, then first TODO
func (t *Task) GetEffectiveRole() string {
	for _, item := range t.TodoItems {
		if !item.Checked && item.Role != "" {
			return item.Role
		}
	}
	if t.Meta.Role != "" {
		return t.Meta.Role
	}
	return ""
}

// WriteAllTasks writes any tasks dirty or not.
func WriteAllTasks(tasks map[string]*Task) (int, error) {
	updated := 0
	for _, task := range tasks {
		if err := task.Write(); err != nil {
			return updated, err
		}
		updated++
	}

	return updated, nil
}

// WriteDirtyTasks writes any tasks marked as dirty.
func WriteDirtyTasks(tasks map[string]*Task) (int, error) {
	updated := 0
	for _, task := range tasks {
		if !task.Dirty {
			continue
		}
		if err := task.Write(); err != nil {
			return updated, err
		}
		updated++
	}

	return updated, nil
}

// InvalidFrontmatterError indicates a task file is missing frontmatter delimiters.
type InvalidFrontmatterError struct {
	Path string
}

func (e *InvalidFrontmatterError) Error() string {
	return "invalid task file format: missing frontmatter delimiters in " + e.Path
}

// FrontmatterParseError indicates an error parsing YAML frontmatter with line number info.
type FrontmatterParseError struct {
	Path       string
	LineNumber int
	YAMLError  string
}

func (e *FrontmatterParseError) Error() string {
	if e.Path != "" && e.LineNumber > 0 {
		return fmt.Sprintf("invalid YAML in %s at line %d: %s", e.Path, e.LineNumber, e.YAMLError)
	}
	if e.Path != "" {
		return fmt.Sprintf("invalid YAML in %s: %s", e.Path, e.YAMLError)
	}
	return fmt.Sprintf("invalid YAML: %s", e.YAMLError)
}
