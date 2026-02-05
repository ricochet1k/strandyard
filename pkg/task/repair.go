package task

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
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

// Validator verifys tasks and their relationships
type Validator struct {
	tasks     map[string]*Task
	errors    []ValidationError
	idPattern *regexp.Regexp
	rolesDir  string
}

type listEntry struct {
	Path  string
	Label string
}

// NewValidator creates a new validator
func NewValidator(tasks map[string]*Task) *Validator {
	return NewValidatorWithRoles(tasks, "roles")
}

// NewValidatorWithRoles creates a validator with a custom roles directory.
func NewValidatorWithRoles(tasks map[string]*Task, rolesDir string) *Validator {
	// ID pattern: <PREFIX><4-lowercase-alphanumeric>-<slug>
	// PREFIX is single uppercase letter
	// Token is 4 lowercase base36 characters (0-9, a-z)
	// Slug is 1+ alphanumeric/hyphen characters
	idPattern := regexp.MustCompile(`^[A-Z][0-9a-z]{4,6}-[a-zA-Z0-9-]+$`)

	return &Validator{
		tasks:     tasks,
		errors:    []ValidationError{},
		idPattern: idPattern,
		rolesDir:  rolesDir,
	}
}

// Validate runs all validations, auto-fixes relationships, and returns errors
func (v *Validator) ValidateAndRepair() []ValidationError {
	v.errors = []ValidationError{}

	// First pass: verify basic task properties
	for id, task := range v.tasks {
		v.verifyID(id, task)
		v.verifyRole(id, task)
		v.verifyPriority(id, task)
		v.verifyParent(id, task)
		v.verifyTaskLinks(id, task)
	}

	v.fixSubtaskTextTitles()

	if _, err := UpdateAllParentTodoEntries(v.tasks); err != nil {
		v.errors = append(v.errors, ValidationError{
			Message: fmt.Sprintf("failed to update parent TODO entries: %v", err),
		})
	}

	return v.errors
}

// FixMissingReferences removes references to tasks that no longer exist.
// Returns notices describing each fix that was applied.
func (v *Validator) FixMissingReferences() []ValidationError {
	notices := []ValidationError{}

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
		if !slices.Equal(task.Meta.Blockers, blockers) {
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
		if !slices.Equal(task.Meta.Blocks, blocks) {
			task.Meta.Blocks = blocks
			changed = true
		}

		if changed {
			task.MarkDirty()
		}
	}

	return notices
}

func (v *Validator) fixSubtaskTextTitles() {
	for _, task := range v.tasks {
		for _, subtask := range task.SubsItems {
			if subtask.SubtaskID != "" {
				if subtaskTask, exists := v.tasks[subtask.SubtaskID]; exists {
					// Update subtask text to match the subtask's title if different
					if subtask.Text != subtaskTask.TitleContent {
						subtask.Text = subtaskTask.TitleContent
						task.MarkDirty()
					}
				}
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
			Message: "malformed ID: must be <PREFIX><4-lowercase-alphanumeric>-<slug> (e.g., T3k7x-example)",
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
	roleDir := v.rolesDir
	if strings.TrimSpace(roleDir) == "" {
		roleDir = "roles"
	}
	rolePath := filepath.Join(roleDir, role+".md")
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

// verifyTaskLinks scans task content for references to other tasks and verifies they exist
func (v *Validator) verifyTaskLinks(id string, task *Task) {
	// Task ID pattern: <PREFIX><4-lowercase-alphanumeric>-<slug>
	// Simplified regex for identifying task links in markdown [text](path)
	linkPattern := regexp.MustCompile(`\[[^\]]*\]\(([^)]+)\)`)

	allContent := task.TitleContent + "\n" + task.BodyContent + "\n" + FormatTodoItems(task.TodoItems) + "\n" + FormatSubtaskItems(task.SubsItems) + "\n" + task.OtherContent
	matches := linkPattern.FindAllStringSubmatch(allContent, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		destination := match[1]
		targetID := extractTaskIDFromPath(destination)
		if targetID != "" && targetID != id {
			if _, exists := v.tasks[targetID]; !exists {
				v.errors = append(v.errors, ValidationError{
					TaskID:  id,
					File:    task.FilePath,
					Message: fmt.Sprintf("broken link: task %s does not exist", targetID),
				})
			}
		}
	}
}

func addUniqueSorted(slice []string, val string) ([]string, bool) {
	if val == "" {
		return slice, false
	}
	if slices.Contains(slice, val) {
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
	idPattern := regexp.MustCompile(`^[A-Z][0-9a-z]{4,6}-[a-zA-Z0-9-]+$`)

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
		if task.Meta.Parent == "" && !task.Meta.Completed && IsActiveStatus(task.Meta.Status) {
			roots = append(roots, listEntry{Path: rel, Label: title})
		}

		// Free tasks have no blockers and are not completed
		if len(task.Meta.Blockers) == 0 && !task.Meta.Completed && IsActiveStatus(task.Meta.Status) {
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

	parsed := ParseFreeList(string(content), tasks)

	entriesByPriority := map[string][]listEntry{
		PriorityHigh:   {},
		PriorityMedium: {},
		PriorityLow:    {},
	}
	other := []listEntry{}

	title := parsed.Title

	// Remove completed tasks
	removeSet := make(map[string]bool)
	for _, taskID := range update.RemoveTaskIDs {
		removeSet[taskID] = true
	}

	kept := make(map[string]bool)
	for _, taskID := range parsed.TaskIDs {
		if removeSet[taskID] {
			continue
		}
		task, exists := tasks[taskID]
		if !exists || task.Meta.Completed || !IsActiveStatus(task.Meta.Status) {
			continue
		}
		kept[taskID] = true
	}

	// Add newly unblocked tasks
	for _, task := range update.AddTasks {
		if task.Meta.Completed || !IsActiveStatus(task.Meta.Status) {
			continue
		}
		kept[task.ID] = true
	}

	// Build entries from tasks map for deterministic labels and paths
	for taskID := range kept {
		task, exists := tasks[taskID]
		if !exists {
			continue
		}
		label := task.Title()
		if label == "" {
			label = task.ID
		}
		entry := listEntry{Path: filepath.ToSlash(task.FilePath), Label: label}

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
		if task.Meta.Completed || !IsActiveStatus(task.Meta.Status) {
			continue // Skip completed tasks or inactive tasks
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
				if !blockerExists || !blockerTask.Meta.Completed || !IsActiveStatus(blockerTask.Meta.Status) {
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
