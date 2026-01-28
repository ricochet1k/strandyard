package task

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/yuin/goldmark/ast"
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

type listEntry struct {
	Path  string
	Label string
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
		v.validatePriority(id, task)
		v.validateParent(id, task)
		v.validateBlockers(id, task)
		v.validateTaskLinks(id, task)
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

// validatePriority checks if the priority is a known value (or empty).
func (v *Validator) validatePriority(id string, task *Task) {
	if IsValidPriority(task.Meta.Priority) {
		return
	}

	v.errors = append(v.errors, ValidationError{
		TaskID:  id,
		File:    task.FilePath,
		Message: fmt.Sprintf("invalid priority %q: must be high, medium, or low", task.Meta.Priority),
	})
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

// validateTaskLinks scans task content for references to other tasks and verifies they exist
func (v *Validator) validateTaskLinks(id string, task *Task) {
	ast.Walk(task.Document, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// Look for link nodes
		if link, ok := n.(*ast.Link); ok {
			destination := string(link.Destination)
			taskID := extractTaskIDFromPath(destination)
			if taskID != "" && taskID != id { // Don't validate self-references
				if _, exists := v.tasks[taskID]; !exists {
					v.errors = append(v.errors, ValidationError{
						TaskID:  id,
						File:    task.FilePath,
						Message: fmt.Sprintf("broken link: task %s does not exist", taskID),
					})
				}
			}
		}

		return ast.WalkContinue, nil
	})
}

// extractTaskIDFromPath extracts a task ID from a file or directory path
// Examples:
//   - "tasks/T3k7x-example/T3k7x-example.md" -> "T3k7x-example"
//   - "../T5h7w-task/task.md" -> "T5h7w-task"
//   - "T2k9p-other" -> "T2k9p-other"
func extractTaskIDFromPath(path string) string {
	// Clean path and split by directory separator
	path = filepath.Clean(path)
	parts := strings.Split(filepath.ToSlash(path), "/")

	// Task ID pattern: <PREFIX><4-lowercase-alphanumeric>-<slug>
	idPattern := regexp.MustCompile(`^[A-Z][0-9a-z]{4}-[a-zA-Z0-9-]+$`)

	// Scan path components for ID pattern
	for _, part := range parts {
		if idPattern.MatchString(part) {
			return part
		}
	}

	return ""
}

// GenerateMasterLists creates root-tasks.md and free-tasks.md
func GenerateMasterLists(tasks map[string]*Task, tasksRoot, rootsFile, freeFile string) error {
	roots := []listEntry{}
	freeByPriority := map[string][]listEntry{
		PriorityHigh:   {},
		PriorityMedium: {},
		PriorityLow:    {},
	}
	freeOther := []listEntry{}

	for _, task := range tasks {
		// Task file path is repo-relative; convert to list-relative when writing.
		rel := filepath.ToSlash(task.FilePath)
		title := task.Title()
		if title == "" {
			title = task.ID
		}

		// Root tasks have no parent and are not completed
		if task.Meta.Parent == "" && !task.Meta.Completed {
			roots = append(roots, listEntry{Path: rel, Label: title})
		}

		// Free tasks have no blockers and are not completed
		if len(task.Meta.Blockers) == 0 && !task.Meta.Completed {
			switch NormalizePriority(task.Meta.Priority) {
			case PriorityHigh:
				freeByPriority[PriorityHigh] = append(freeByPriority[PriorityHigh], listEntry{Path: rel, Label: title})
			case PriorityMedium:
				freeByPriority[PriorityMedium] = append(freeByPriority[PriorityMedium], listEntry{Path: rel, Label: title})
			case PriorityLow:
				freeByPriority[PriorityLow] = append(freeByPriority[PriorityLow], listEntry{Path: rel, Label: title})
			default:
				freeOther = append(freeOther, listEntry{Path: rel, Label: title})
			}
		}
	}

	// Sort for deterministic output
	sort.Slice(roots, func(i, j int) bool { return roots[i].Path < roots[j].Path })
	for key := range freeByPriority {
		items := freeByPriority[key]
		sort.Slice(items, func(i, j int) bool { return items[i].Path < items[j].Path })
		freeByPriority[key] = items
	}
	sort.Slice(freeOther, func(i, j int) bool { return freeOther[i].Path < freeOther[j].Path })

	// Write files
	if err := writeListFile(rootsFile, "Root tasks", roots); err != nil {
		return err
	}
	if err := writePriorityListFile(freeFile, "Free tasks", freeByPriority, freeOther); err != nil {
		return err
	}

	return nil
}

func writeListFile(path, title string, entries []listEntry) error {
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
		sb.WriteString("[")
		sb.WriteString(e.Label)
		sb.WriteString("](")
		sb.WriteString(relativeListPath(path, e.Path))
		sb.WriteString(")")
		sb.WriteString("\n")
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func writePriorityListFile(path, title string, entries map[string][]listEntry, other []listEntry) error {
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

	writeSection := func(name string, items []listEntry) {
		sb.WriteString("## ")
		sb.WriteString(name)
		sb.WriteString("\n\n")
		for _, e := range items {
			sb.WriteString("- ")
			sb.WriteString("[")
			sb.WriteString(e.Label)
			sb.WriteString("](")
			sb.WriteString(relativeListPath(path, e.Path))
			sb.WriteString(")")
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	writeSection("High", entries[PriorityHigh])
	writeSection("Medium", entries[PriorityMedium])
	writeSection("Low", entries[PriorityLow])
	if len(other) > 0 {
		writeSection("Other", other)
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func relativeListPath(listFile, target string) string {
	base := filepath.Dir(listFile)
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return filepath.ToSlash(target)
	}
	return filepath.ToSlash(rel)
}
