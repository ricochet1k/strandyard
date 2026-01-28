package task

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

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

// Validator verifys tasks and their relationships
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

// Validate runs all validations, auto-fixes relationships, and returns errors
func (v *Validator) Validate() []ValidationError {
	v.errors = []ValidationError{}

	// First pass: verify basic task properties
	for id, task := range v.tasks {
		v.verifyID(id, task)
		v.verifyRole(id, task)
		v.verifyPriority(id, task)
		v.verifyParent(id, task)
		v.verifyBlockers(id, task)
		v.verifyTaskLinks(id, task)
	}

	// Auto-fix bidirectional blocker relationships
	v.fixBlockerRelationships()

	// Second pass: verify bidirectional relationships are now correct
	for id, task := range v.tasks {
		v.verifyBidirectionalBlockers(id, task)
	}

	return v.errors
}

// FixMissingReferences removes references to tasks that no longer exist.
// Returns notices describing each fix that was applied.
func (v *Validator) FixMissingReferences() []ValidationError {
	notices := []ValidationError{}
	now := time.Now()

	for id, task := range v.tasks {
		changed := false

		if task.Meta.Parent != "" {
			if _, exists := v.tasks[task.Meta.Parent]; !exists {
				notices = append(notices, ValidationError{
					TaskID:  id,
					File:    task.FilePath,
					Message: fmt.Sprintf("parent task %s does not exist", task.Meta.Parent),
				})
				task.Meta.Parent = ""
				changed = true
			}
		}

		blockers, missingBlockers := filterExistingTaskIDs(task.Meta.Blockers, v.tasks)
		for _, blocker := range missingBlockers {
			notices = append(notices, ValidationError{
				TaskID:  id,
				File:    task.FilePath,
				Message: fmt.Sprintf("blocker task %s does not exist", blocker),
			})
		}
		if !equalStringSlices(task.Meta.Blockers, blockers) {
			task.Meta.Blockers = blockers
			changed = true
		}

		blocks, missingBlocks := filterExistingTaskIDs(task.Meta.Blocks, v.tasks)
		for _, blocked := range missingBlocks {
			notices = append(notices, ValidationError{
				TaskID:  id,
				File:    task.FilePath,
				Message: fmt.Sprintf("blocks non-existent task %s", blocked),
			})
		}
		if !equalStringSlices(task.Meta.Blocks, blocks) {
			task.Meta.Blocks = blocks
			changed = true
		}

		if changed {
			task.Meta.DateEdited = now
			task.MarkDirty()
		}
	}

	return notices
}

// fixBlockerRelationships automatically fixes missing bidirectional blocker relationships
func (v *Validator) fixBlockerRelationships() {
	now := time.Now()
	for taskID, task := range v.tasks {
		for _, blockerID := range task.Meta.Blockers {
			if blockerID == "" {
				continue
			}
			blocker, exists := v.tasks[blockerID]
			if !exists {
				continue // Will be caught by validation
			}

			updated, changed := addUniqueSorted(blocker.Meta.Blocks, taskID)
			if changed {
				blocker.Meta.Blocks = updated
				blocker.Meta.DateEdited = now
				blocker.MarkDirty()
			}
		}

		for _, blockedID := range task.Meta.Blocks {
			if blockedID == "" {
				continue
			}
			blocked, exists := v.tasks[blockedID]
			if !exists {
				continue // Will be caught by validation
			}

			updated, changed := addUniqueSorted(blocked.Meta.Blockers, taskID)
			if changed {
				blocked.Meta.Blockers = updated
				blocked.Meta.DateEdited = now
				blocked.MarkDirty()
			}
		}
	}
}

// verifyID checks if the task ID follows the correct format
func (v *Validator) verifyID(id string, task *Task) {
	if !v.idPattern.MatchString(id) {
		v.errors = append(v.errors, ValidationError{
			TaskID:  id,
			File:    task.FilePath,
			Message: fmt.Sprintf("malformed ID: must be <PREFIX><4-lowercase-alphanumeric>-<slug> (e.g., T3k7x-example)"),
		})
	}
}

// verifyRole checks if the role exists
func (v *Validator) verifyRole(id string, task *Task) {
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

// verifyPriority checks if the priority is a known value (or empty).
func (v *Validator) verifyPriority(id string, task *Task) {
	if IsValidPriority(task.Meta.Priority) {
		return
	}

	v.errors = append(v.errors, ValidationError{
		TaskID:  id,
		File:    task.FilePath,
		Message: fmt.Sprintf("invalid priority %q: must be high, medium, or low", task.Meta.Priority),
	})
}

// verifyParent checks if parent task exists
func (v *Validator) verifyParent(id string, task *Task) {
	if task.Meta.Parent == "" {
		return // Root task, no parent to verify
	}

	if _, exists := v.tasks[task.Meta.Parent]; !exists {
		v.errors = append(v.errors, ValidationError{
			TaskID:  id,
			File:    task.FilePath,
			Message: fmt.Sprintf("parent task %s does not exist", task.Meta.Parent),
		})
	}
}

// verifyBlockers checks if blocker tasks exist
func (v *Validator) verifyBlockers(id string, task *Task) {
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

// verifyTaskLinks scans task content for references to other tasks and verifies they exist
func (v *Validator) verifyTaskLinks(id string, task *Task) {
	ast.Walk(task.Document, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// Look for link nodes
		if link, ok := n.(*ast.Link); ok {
			destination := string(link.Destination)
			taskID := extractTaskIDFromPath(destination)
			if taskID != "" && taskID != id { // Don't verify self-references
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

// verifyBidirectionalBlockers ensures blocker relationships are bidirectional
func (v *Validator) verifyBidirectionalBlockers(id string, task *Task) {
	// For each blocker that this task lists, verify it lists this task in its blocks field
	for _, blockerID := range task.Meta.Blockers {
		if blockerID == "" {
			continue
		}

		blocker, exists := v.tasks[blockerID]
		if !exists {
			// Already caught by verifyBlockers
			continue
		}

		// Check if blocker lists this task in its blocks field
		if !contains(blocker.Meta.Blocks, id) {
			v.errors = append(v.errors, ValidationError{
				TaskID:  id,
				File:    task.FilePath,
				Message: fmt.Sprintf("task has blocker %s, but %s doesn't list this task in blocks field", blockerID, blockerID),
			})
		}
	}

	// For each task that this one blocks, verify it lists this task in its blockers field
	for _, blockedID := range task.Meta.Blocks {
		if blockedID == "" {
			continue
		}

		blocked, exists := v.tasks[blockedID]
		if !exists {
			v.errors = append(v.errors, ValidationError{
				TaskID:  id,
				File:    task.FilePath,
				Message: fmt.Sprintf("blocks non-existent task %s", blockedID),
			})
			continue
		}

		// Check if blocked task lists this in its blockers field
		if !contains(blocked.Meta.Blockers, id) {
			v.errors = append(v.errors, ValidationError{
				TaskID:  id,
				File:    task.FilePath,
				Message: fmt.Sprintf("task blocks %s, but %s doesn't list this task in blockers field", blockedID, blockedID),
			})
		}
	}
}

// contains checks if a string slice contains a value
func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func addUniqueSorted(slice []string, val string) ([]string, bool) {
	if val == "" {
		return slice, false
	}
	if contains(slice, val) {
		return slice, false
	}
	slice = append(slice, val)
	sort.Strings(slice)
	return slice, true
}

func filterExistingTaskIDs(ids []string, tasks map[string]*Task) ([]string, []string) {
	kept := []string{}
	missing := []string{}
	seen := map[string]struct{}{}
	missingSeen := map[string]struct{}{}

	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, exists := tasks[id]; !exists {
			if _, ok := missingSeen[id]; !ok {
				missing = append(missing, id)
				missingSeen[id] = struct{}{}
			}
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		kept = append(kept, id)
	}

	sort.Strings(kept)
	sort.Strings(missing)
	return kept, missing
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

// IncrementalFreeListUpdate represents changes to make to free-tasks.md
type IncrementalFreeListUpdate struct {
	RemoveTaskIDs []string // Tasks to remove (completed tasks)
	AddTasks      []*Task  // Tasks to add (newly unblocked tasks)
}

// UpdateFreeListIncrementally updates free-tasks.md incrementally
func UpdateFreeListIncrementally(tasks map[string]*Task, freeFile string, update IncrementalFreeListUpdate) error {
	// Read current free-tasks.md content
	content, err := os.ReadFile(freeFile)
	if err != nil {
		return fmt.Errorf("failed to read free tasks file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Parse existing entries
	entriesByPriority := map[string][]listEntry{
		PriorityHigh:   {},
		PriorityMedium: {},
		PriorityLow:    {},
	}
	other := []listEntry{}

	var currentSection string
	var title string

	// Parse the file structure
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			title = strings.TrimPrefix(line, "# ")
			continue
		}
		if strings.HasPrefix(line, "## ") {
			currentSection = strings.TrimPrefix(line, "## ")
			continue
		}
		if strings.HasPrefix(line, "- [") && strings.Contains(line, "](") {
			// Extract task info from markdown link
			endLabel := strings.Index(line, "](")
			if endLabel > 2 {
				label := line[2:endLabel]
				startPath := endLabel + 2
				endPath := strings.LastIndex(line, ")")
				if endPath > startPath {
					path := line[startPath:endPath]

					entry := listEntry{Label: label, Path: path}

					// Determine which priority section this belongs to
					switch currentSection {
					case "High":
						entriesByPriority[PriorityHigh] = append(entriesByPriority[PriorityHigh], entry)
					case "Medium":
						entriesByPriority[PriorityMedium] = append(entriesByPriority[PriorityMedium], entry)
					case "Low":
						entriesByPriority[PriorityLow] = append(entriesByPriority[PriorityLow], entry)
					case "Other":
						other = append(other, entry)
					default:
						// If no section yet, it's in the old format or we need to detect
						// Try to find the task to determine its priority
						for _, task := range tasks {
							taskRelPath := filepath.ToSlash(task.FilePath)
							if strings.Contains(path, taskRelPath) || strings.Contains(taskRelPath, path) {
								priority := NormalizePriority(task.Meta.Priority)
								switch priority {
								case PriorityHigh:
									entriesByPriority[PriorityHigh] = append(entriesByPriority[PriorityHigh], entry)
								case PriorityMedium:
									entriesByPriority[PriorityMedium] = append(entriesByPriority[PriorityMedium], entry)
								case PriorityLow:
									entriesByPriority[PriorityLow] = append(entriesByPriority[PriorityLow], entry)
								default:
									other = append(other, entry)
								}
								break
							}
						}
					}
				}
			}
		}
	}

	// Remove completed tasks
	removeSet := make(map[string]bool)
	for _, taskID := range update.RemoveTaskIDs {
		removeSet[taskID] = true
	}

	// Filter out removed tasks
	filterEntries := func(entries []listEntry) []listEntry {
		filtered := make([]listEntry, 0, len(entries))
		for _, entry := range entries {
			// Check if this entry corresponds to a removed task
			shouldRemove := false
			for taskID, task := range tasks {
				taskRelPath := filepath.ToSlash(task.FilePath)
				if removeSet[taskID] && (strings.Contains(entry.Path, taskRelPath) || strings.Contains(taskRelPath, entry.Path)) {
					shouldRemove = true
					break
				}
			}
			if !shouldRemove {
				filtered = append(filtered, entry)
			}
		}
		return filtered
	}

	entriesByPriority[PriorityHigh] = filterEntries(entriesByPriority[PriorityHigh])
	entriesByPriority[PriorityMedium] = filterEntries(entriesByPriority[PriorityMedium])
	entriesByPriority[PriorityLow] = filterEntries(entriesByPriority[PriorityLow])
	other = filterEntries(other)

	// Add newly unblocked tasks
	for _, task := range update.AddTasks {
		title := task.Title()
		if title == "" {
			title = task.ID
		}

		rel := filepath.ToSlash(task.FilePath)
		entry := listEntry{Path: rel, Label: title}

		switch NormalizePriority(task.Meta.Priority) {
		case PriorityHigh:
			entriesByPriority[PriorityHigh] = append(entriesByPriority[PriorityHigh], entry)
		case PriorityMedium:
			entriesByPriority[PriorityMedium] = append(entriesByPriority[PriorityMedium], entry)
		case PriorityLow:
			entriesByPriority[PriorityLow] = append(entriesByPriority[PriorityLow], entry)
		default:
			other = append(other, entry)
		}
	}

	// Sort for deterministic output
	sort.Slice(entriesByPriority[PriorityHigh], func(i, j int) bool {
		return entriesByPriority[PriorityHigh][i].Path < entriesByPriority[PriorityHigh][j].Path
	})
	sort.Slice(entriesByPriority[PriorityMedium], func(i, j int) bool {
		return entriesByPriority[PriorityMedium][i].Path < entriesByPriority[PriorityMedium][j].Path
	})
	sort.Slice(entriesByPriority[PriorityLow], func(i, j int) bool {
		return entriesByPriority[PriorityLow][i].Path < entriesByPriority[PriorityLow][j].Path
	})
	sort.Slice(other, func(i, j int) bool { return other[i].Path < other[j].Path })

	// Write the updated file
	if err := writePriorityListFile(freeFile, title, entriesByPriority, other); err != nil {
		return fmt.Errorf("failed to write updated free tasks file: %w", err)
	}

	return nil
}

// CalculateIncrementalFreeListUpdate determines what changes need to be made to free-tasks.md
// when a task is completed
func CalculateIncrementalFreeListUpdate(tasks map[string]*Task, completedTaskID string) (IncrementalFreeListUpdate, error) {
	_, exists := tasks[completedTaskID]
	if !exists {
		return IncrementalFreeListUpdate{}, fmt.Errorf("completed task not found: %s", completedTaskID)
	}

	update := IncrementalFreeListUpdate{
		RemoveTaskIDs: []string{completedTaskID},
		AddTasks:      []*Task{},
	}

	// Find tasks that were blocked by the completed task
	for _, task := range tasks {
		if task.Meta.Completed {
			continue // Skip completed tasks
		}

		// Check if this task was blocked by the completed task
		wasBlocked := false
		for _, blocker := range task.Meta.Blockers {
			if blocker == completedTaskID {
				wasBlocked = true
				break
			}
		}

		if wasBlocked {
			// Check if this task is now free (no remaining blockers)
			allBlockersCompleted := true
			for _, blocker := range task.Meta.Blockers {
				if blocker == completedTaskID {
					continue // This is the task we're completing
				}
				blockerTask, blockerExists := tasks[blocker]
				if !blockerExists || !blockerTask.Meta.Completed {
					allBlockersCompleted = false
					break
				}
			}

			if allBlockersCompleted {
				update.AddTasks = append(update.AddTasks, task)
			}
		}
	}

	return update, nil
}
