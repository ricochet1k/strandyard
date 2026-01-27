package task

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	TaskID  string
	File    string
	Message string
}

func (e ValidationError) Error() string {
	if e.TaskID != "" {
		return fmt.Sprintf("Task %s: %s", e.TaskID, e.Message)
	}
	return e.Message
}

// Validator validates tasks and their relationships
type Validator struct {
	tasks     map[string]*Task
	errors    []ValidationError
	idPattern *regexp.Regexp
}

// NewValidator creates a new validator
func NewValidator(tasks map[string]*Task) *Validator {
	// ID pattern: <PREFIX><4-lowercase-alphanumeric>-<slug>
	// PREFIX is single uppercase letter
	// Token is 4 lowercase base36 characters (0-9, a-z)
	// Slug is 1+ alphanumeric/hyphen characters
	idPattern := regexp.MustCompile(`^[A-Z][0-9a-z]{4}-[a-zA-Z0-9-]+$`)

	return &Validator{
		tasks:     tasks,
		errors:    []ValidationError{},
		idPattern: idPattern,
	}
}

// Validate runs all validations and returns errors
func (v *Validator) Validate() []ValidationError {
	v.errors = []ValidationError{}

	for id, task := range v.tasks {
		v.validateID(id, task)
		v.validateRole(id, task)
		v.validateParent(id, task)
		v.validateBlockers(id, task)
	}

	return v.errors
}

// validateID checks if the task ID follows the correct format
func (v *Validator) validateID(id string, task *Task) {
	if !v.idPattern.MatchString(id) {
		v.errors = append(v.errors, ValidationError{
			TaskID:  id,
			File:    task.FilePath,
			Message: fmt.Sprintf("malformed ID: must be <PREFIX><4-lowercase-alphanumeric>-<slug> (e.g., T3k7x-example)"),
		})
	}
}

// validateRole checks if the role exists
func (v *Validator) validateRole(id string, task *Task) {
	role := task.GetEffectiveRole()
	if role == "" {
		v.errors = append(v.errors, ValidationError{
			TaskID:  id,
			File:    task.FilePath,
			Message: "missing role in frontmatter and no role found in first TODO",
		})
		return
	}

	// Check if role file exists
	rolePath := filepath.Join("roles", role+".md")
	if _, err := os.Stat(rolePath); os.IsNotExist(err) {
		v.errors = append(v.errors, ValidationError{
			TaskID:  id,
			File:    task.FilePath,
			Message: fmt.Sprintf("role file %s does not exist", rolePath),
		})
	}
}

// validateParent checks if parent task exists
func (v *Validator) validateParent(id string, task *Task) {
	if task.Meta.Parent == "" {
		return // Root task, no parent to validate
	}

	if _, exists := v.tasks[task.Meta.Parent]; !exists {
		v.errors = append(v.errors, ValidationError{
			TaskID:  id,
			File:    task.FilePath,
			Message: fmt.Sprintf("parent task %s does not exist", task.Meta.Parent),
		})
	}
}

// validateBlockers checks if blocker tasks exist
func (v *Validator) validateBlockers(id string, task *Task) {
	for _, blocker := range task.Meta.Blockers {
		if blocker == "" {
			continue
		}

		if _, exists := v.tasks[blocker]; !exists {
			v.errors = append(v.errors, ValidationError{
				TaskID:  id,
				File:    task.FilePath,
				Message: fmt.Sprintf("blocker task %s does not exist", blocker),
			})
		}
	}
}

// GenerateMasterLists creates root-tasks.md and free-tasks.md
func GenerateMasterLists(tasks map[string]*Task, tasksRoot, rootsFile, freeFile string) error {
	roots := []string{}
	free := []string{}

	for _, task := range tasks {
		// Task file path is already repo-relative; keep it stable in lists.
		rel := filepath.ToSlash(task.FilePath)

		// Root tasks have no parent and are not completed
		if task.Meta.Parent == "" && !task.Meta.Completed {
			roots = append(roots, rel)
		}

		// Free tasks have no blockers and are not completed
		if len(task.Meta.Blockers) == 0 && !task.Meta.Completed {
			free = append(free, rel)
		}
	}

	// Sort for deterministic output
	sort.Strings(roots)
	sort.Strings(free)

	// Write files
	if err := writeListFile(rootsFile, "Root tasks", roots); err != nil {
		return err
	}
	if err := writeListFile(freeFile, "Free tasks", free); err != nil {
		return err
	}

	return nil
}

func writeListFile(path, title string, entries []string) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	var sb strings.Builder
	sb.WriteString("# ")
	sb.WriteString(title)
	sb.WriteString("\n\n")
	for _, e := range entries {
		sb.WriteString("- ")
		sb.WriteString(e)
		sb.WriteString("\n")
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
