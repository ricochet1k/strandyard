package task

import (
	"fmt"
	"os"
	"slices"
	"sort"
)

// TaskDB lazy-loads and manages tasks with strict relationship integrity.
// All operations through TaskDB automatically maintain bidirectional relationships
// between parent/children, blockers/blocks, and completion status.
type TaskDB struct {
	tasksRoot string
	parser    *Parser
	tasks     map[string]*Task
}

// NewTaskDB creates a new TaskDB instance.
func NewTaskDB(tasksRoot string) *TaskDB {
	return &TaskDB{
		tasksRoot: tasksRoot,
		parser:    NewParser(),
		tasks:     make(map[string]*Task),
	}
}

// Get retrieves a task by ID, lazy-loading from disk if needed.
func (db *TaskDB) Get(id string) (*Task, error) {
	if task, ok := db.tasks[id]; ok {
		return task, nil
	}
	return db.Load(id)
}

// Load forces a reload of a task from disk.
func (db *TaskDB) Load(id string) (*Task, error) {
	// Try to find the task file for this ID
	// Priority: <id>/<id>.md, then <id>/task.md
	candidates := []string{
		fmt.Sprintf("%s/%s/%s.md", db.tasksRoot, id, id),
		fmt.Sprintf("%s/%s/task.md", db.tasksRoot, id),
	}

	var task *Task
	var err error
	for _, path := range candidates {
		task, err = db.parser.ParseFile(path)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load task %s: %w", id, err)
	}

	db.tasks[id] = task
	return task, nil
}

// LoadAll loads all tasks from the tasks root directory.
func (db *TaskDB) LoadAll() error {
	tasks, err := db.parser.LoadTasks(db.tasksRoot)
	if err != nil {
		return err
	}
	db.tasks = tasks
	return nil
}

// Save writes a specific task to disk if dirty.
func (db *TaskDB) Save(id string) error {
	task, ok := db.tasks[id]
	if !ok {
		return fmt.Errorf("task %s not loaded", id)
	}
	if !task.Dirty {
		return nil
	}
	return task.Write()
}

// SaveDirty writes all dirty tasks to disk.
func (db *TaskDB) SaveDirty() (int, error) {
	return WriteDirtyTasks(db.tasks)
}

// SaveAll writes all tasks to disk regardless of dirty status.
func (db *TaskDB) SaveAll() (int, error) {
	return WriteAllTasks(db.tasks)
}

// GetAll returns all loaded tasks.
func (db *TaskDB) GetAll() map[string]*Task {
	return db.tasks
}

// SetParent sets the parent of a child task, ensuring the parent exists.
// This does NOT automatically add the child to the parent's blockers - call
// SyncBlockersFromChildren for that.
func (db *TaskDB) SetParent(childID, parentID string) error {
	if childID == parentID {
		return fmt.Errorf("task cannot be its own parent")
	}

	child, err := db.Get(childID)
	if err != nil {
		return fmt.Errorf("child task not found: %w", err)
	}

	if parentID == "" {
		return db.ClearParent(childID)
	}

	parent, err := db.Get(parentID)
	if err != nil {
		return fmt.Errorf("parent task not found: %w", err)
	}

	// Check for cycles
	if db.wouldCreateCycle(childID, parentID) {
		return fmt.Errorf("would create parent cycle")
	}

	if child.Meta.Parent != parentID {
		child.Meta.Parent = parentID
		child.MarkDirty()
	}

	// Note: We don't update parent's blockers here - that's done via SyncBlockersFromChildren
	_ = parent // Ensure parent exists

	return nil
}

// wouldCreateCycle checks if setting parentID as the parent of childID would create a cycle.
func (db *TaskDB) wouldCreateCycle(childID, parentID string) bool {
	visited := make(map[string]bool)
	current := parentID

	for current != "" {
		if current == childID {
			return true
		}
		if visited[current] {
			return false // Already checked this path
		}
		visited[current] = true

		task, ok := db.tasks[current]
		if !ok {
			return false
		}
		current = task.Meta.Parent
	}

	return false
}

// ClearParent removes the parent relationship from a child task.
func (db *TaskDB) ClearParent(childID string) error {
	child, err := db.Get(childID)
	if err != nil {
		return fmt.Errorf("child task not found: %w", err)
	}

	if child.Meta.Parent != "" {
		child.Meta.Parent = ""
		child.MarkDirty()
	}

	return nil
}

// GetChildren returns all tasks that have the given task as their parent.
func (db *TaskDB) GetChildren(parentID string) []*Task {
	var children []*Task
	for _, task := range db.tasks {
		if task.Meta.Parent == parentID {
			children = append(children, task)
		}
	}
	return children
}

// AddBlocker adds a blocker to a task and maintains bidirectional relationship.
// After this operation:
// - taskID will have blockerID in its Blockers list
// - blockerID will have taskID in its Blocks list
func (db *TaskDB) AddBlocker(taskID, blockerID string) error {
	if taskID == blockerID {
		return fmt.Errorf("task cannot block itself")
	}

	task, err := db.Get(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	blocker, err := db.Get(blockerID)
	if err != nil {
		return fmt.Errorf("blocker task not found: %w", err)
	}

	// Add blockerID to task's blockers
	if !slices.Contains(task.Meta.Blockers, blockerID) {
		task.Meta.Blockers = append(task.Meta.Blockers, blockerID)
		sort.Strings(task.Meta.Blockers)
		task.MarkDirty()
	}

	// Add taskID to blocker's blocks
	if !slices.Contains(blocker.Meta.Blocks, taskID) {
		blocker.Meta.Blocks = append(blocker.Meta.Blocks, taskID)
		sort.Strings(blocker.Meta.Blocks)
		blocker.MarkDirty()
	}

	return nil
}

// RemoveBlocker removes a blocker from a task and maintains bidirectional relationship.
func (db *TaskDB) RemoveBlocker(taskID, blockerID string) error {
	task, err := db.Get(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	blocker, err := db.Get(blockerID)
	if err != nil {
		// If blocker doesn't exist, just remove it from task's blockers
		if idx := slices.Index(task.Meta.Blockers, blockerID); idx != -1 {
			task.Meta.Blockers = slices.Delete(task.Meta.Blockers, idx, idx+1)
			task.MarkDirty()
		}
		return nil
	}

	// Remove blockerID from task's blockers
	if idx := slices.Index(task.Meta.Blockers, blockerID); idx != -1 {
		task.Meta.Blockers = slices.Delete(task.Meta.Blockers, idx, idx+1)
		task.MarkDirty()
	}

	// Remove taskID from blocker's blocks
	if idx := slices.Index(blocker.Meta.Blocks, taskID); idx != -1 {
		blocker.Meta.Blocks = slices.Delete(blocker.Meta.Blocks, idx, idx+1)
		blocker.MarkDirty()
	}

	return nil
}

// AddBlocked marks that this task blocks another task (inverse of AddBlocker).
func (db *TaskDB) AddBlocked(taskID, blockedID string) error {
	return db.AddBlocker(blockedID, taskID)
}

// RemoveBlocked removes a task from this task's blocks list (inverse of RemoveBlocker).
func (db *TaskDB) RemoveBlocked(taskID, blockedID string) error {
	return db.RemoveBlocker(blockedID, taskID)
}

// SetCompleted marks a task as completed or incomplete.
// This does NOT automatically update blocker relationships - call
// SyncBlockersFromChildren or UpdateBlockersAfterCompletion for that.
func (db *TaskDB) SetCompleted(taskID string, completed bool) error {
	task, err := db.Get(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	if task.Meta.Completed != completed {
		task.Meta.Completed = completed
		task.MarkDirty()
	}

	return nil
}

// SyncBlockersFromChildren updates all parent task blockers based on their
// incomplete children. This ensures:
// - Each parent is blocked by all incomplete children
// - Each parent is blocked by all tasks that explicitly block it
// Returns the number of tasks modified.
func (db *TaskDB) SyncBlockersFromChildren() (int, error) {
	return UpdateBlockersFromChildren(db.tasks)
}

// UpdateBlockersAfterCompletion should be called after marking a task complete.
// It removes the completed task from the blockers lists of tasks it was blocking.
func (db *TaskDB) UpdateBlockersAfterCompletion(taskID string) error {
	task, err := db.Get(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	if !task.Meta.Completed {
		return fmt.Errorf("task %s is not completed", taskID)
	}

	// Remove this task from all tasks it blocks
	for _, blockedID := range task.Meta.Blocks {
		blocked, err := db.Get(blockedID)
		if err != nil {
			continue // Skip if blocked task doesn't exist
		}

		if idx := slices.Index(blocked.Meta.Blockers, taskID); idx != -1 {
			blocked.Meta.Blockers = slices.Delete(blocked.Meta.Blockers, idx, idx+1)
			blocked.MarkDirty()
		}
	}

	// If this task has a parent, it should be removed from parent's blockers
	if task.Meta.Parent != "" {
		parent, err := db.Get(task.Meta.Parent)
		if err == nil {
			if idx := slices.Index(parent.Meta.Blockers, taskID); idx != -1 {
				parent.Meta.Blockers = slices.Delete(parent.Meta.Blockers, idx, idx+1)
				parent.MarkDirty()
			}
		}
	}

	return nil
}

// Validate runs validation checks on all loaded tasks.
func (db *TaskDB) Validate() []ValidationError {
	validator := NewValidator(db.tasks)
	return validator.ValidateAndRepair()
}

// FixMissingReferences removes references to tasks that no longer exist.
func (db *TaskDB) FixMissingReferences() []ValidationError {
	validator := NewValidator(db.tasks)
	return validator.FixMissingReferences()
}

// FixBlockerRelationships ensures all blocker/blocks relationships are bidirectional.
// Returns the number of tasks modified.
func (db *TaskDB) FixBlockerRelationships() int {
	modified := 0

	// Build a map of what should be in each task's blocks list
	shouldBlock := make(map[string]map[string]bool)
	for taskID, task := range db.tasks {
		for _, blockerID := range task.Meta.Blockers {
			if shouldBlock[blockerID] == nil {
				shouldBlock[blockerID] = make(map[string]bool)
			}
			shouldBlock[blockerID][taskID] = true
		}
	}

	// Fix blocks lists (lazy-load blocker tasks as needed)
	for blockerID, blockedTasks := range shouldBlock {
		blocker, err := db.Get(blockerID)
		if err != nil {
			continue // Blocker task doesn't exist
		}

		// Build desired blocks list
		desired := make([]string, 0, len(blockedTasks))
		for blockedID := range blockedTasks {
			desired = append(desired, blockedID)
		}
		sort.Strings(desired)

		if !slices.Equal(blocker.Meta.Blocks, desired) {
			blocker.Meta.Blocks = desired
			blocker.MarkDirty()
			modified++
		}
	}

	// Also check for blocks that shouldn't be there (where the blockers list doesn't match)
	for taskID, task := range db.tasks {
		changed := false
		newBlocks := make([]string, 0, len(task.Meta.Blocks))

		for _, blockedID := range task.Meta.Blocks {
			blocked, err := db.Get(blockedID)
			if err != nil {
				// Blocked task doesn't exist, remove it
				changed = true
				continue
			}

			// Check if blocked task has this task in its blockers
			if slices.Contains(blocked.Meta.Blockers, taskID) {
				newBlocks = append(newBlocks, blockedID)
			} else {
				// Add this task to blocked's blockers
				blocked.Meta.Blockers = append(blocked.Meta.Blockers, taskID)
				sort.Strings(blocked.Meta.Blockers)
				blocked.MarkDirty()
				newBlocks = append(newBlocks, blockedID)
				modified++
			}
		}

		if changed || !slices.Equal(task.Meta.Blocks, newBlocks) {
			task.Meta.Blocks = newBlocks
			task.MarkDirty()
			modified++
		}
	}

	return modified
}

// GetOrCreate gets a task by ID, or creates a new empty task with that ID if it doesn't exist.
// The created task is not persisted until Save is called.
func (db *TaskDB) GetOrCreate(id string) (*Task, error) {
	if task, ok := db.tasks[id]; ok {
		return task, nil
	}

	// Try to load from disk first
	task, err := db.Load(id)
	if err == nil {
		return task, nil
	}

	// Create new task
	task = &Task{
		ID:       id,
		FilePath: fmt.Sprintf("%s/%s/%s.md", db.tasksRoot, id, id),
		Dir:      fmt.Sprintf("%s/%s", db.tasksRoot, id),
		Meta:     Metadata{},
	}
	db.tasks[id] = task
	task.MarkDirty()

	return task, nil
}

// LoadAllIfEmpty loads all tasks from disk if the task map is empty.
// This is a convenience method for commands that need all tasks loaded.
func (db *TaskDB) LoadAllIfEmpty() error {
	if len(db.tasks) > 0 {
		return nil
	}
	return db.LoadAll()
}

// ResolveID resolves a short ID, prefix, or full ID to a canonical task ID.
// Requires tasks to be loaded (calls LoadAllIfEmpty internally).
func (db *TaskDB) ResolveID(input string) (string, error) {
	if err := db.LoadAllIfEmpty(); err != nil {
		return "", err
	}
	return ResolveTaskID(db.tasks, input)
}

// ResolveIDs resolves a list of task ID inputs, de-duplicates, and sorts them.
func (db *TaskDB) ResolveIDs(inputs []string) ([]string, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	if err := db.LoadAllIfEmpty(); err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	resolved := make([]string, 0, len(inputs))
	for _, input := range inputs {
		id, err := ResolveTaskID(db.tasks, input)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		resolved = append(resolved, id)
	}
	sort.Strings(resolved)
	return resolved, nil
}

// GetResolved combines ResolveID and Get into a single call.
// Returns the task, its resolved ID, and any error.
func (db *TaskDB) GetResolved(input string) (*Task, string, error) {
	id, err := db.ResolveID(input)
	if err != nil {
		return nil, "", err
	}
	task, err := db.Get(id)
	if err != nil {
		return nil, id, err
	}
	return task, id, nil
}

// Has returns true if a task with the given ID is loaded.
func (db *TaskDB) Has(id string) bool {
	_, ok := db.tasks[id]
	return ok
}

// ReadRaw reads the raw file contents for a task.
// Useful for the "show" command which displays the original markdown.
func (db *TaskDB) ReadRaw(id string) ([]byte, error) {
	task, err := db.Get(id)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(task.FilePath)
}

// CompleteTodoResult contains the result of completing a todo item.
type CompleteTodoResult struct {
	TaskCompleted       bool
	RemainingIncomplete int
}

// CompleteTodo marks a todo item as complete on a task.
// todoNum is 1-based. If this was the last incomplete todo, the task is marked complete.
func (db *TaskDB) CompleteTodo(taskID string, todoNum int, report string) (*CompleteTodoResult, error) {
	task, err := db.Get(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	if todoNum <= 0 || todoNum > len(task.TodoItems) {
		return nil, fmt.Errorf("invalid todo number %d, task has %d todo items", todoNum, len(task.TodoItems))
	}

	todoIndex := todoNum - 1
	if task.TodoItems[todoIndex].Checked {
		return &CompleteTodoResult{
			TaskCompleted:       false,
			RemainingIncomplete: countIncompleteTodos(task),
		}, nil
	}

	task.TodoItems[todoIndex].Checked = true
	if report != "" {
		task.TodoItems[todoIndex].Report = report
	}
	task.MarkDirty()

	remaining := countIncompleteTodos(task)
	result := &CompleteTodoResult{
		TaskCompleted:       remaining == 0,
		RemainingIncomplete: remaining,
	}

	if result.TaskCompleted {
		task.Meta.Completed = true
		task.MarkDirty()
	}

	return result, nil
}

func countIncompleteTodos(task *Task) int {
	count := 0
	for _, todo := range task.TodoItems {
		if !todo.Checked {
			count++
		}
	}
	return count
}

// AppendCompletionReport appends a completion report to the task's content.
func (db *TaskDB) AppendCompletionReport(taskID, report string) error {
	task, err := db.Get(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	if report != "" {
		if task.OtherContent != "" {
			task.OtherContent += "\n\n"
		}
		task.OtherContent += "## Completion Report\n" + report
		task.MarkDirty()
	}

	return nil
}

// CompleteTask marks a task as completed and optionally appends a report.
func (db *TaskDB) CompleteTask(taskID, report string) error {
	if err := db.SetCompleted(taskID, true); err != nil {
		return err
	}
	return db.AppendCompletionReport(taskID, report)
}

// GetIncompleteTodos returns the incomplete todo items for a task.
func (db *TaskDB) GetIncompleteTodos(taskID string) ([]TaskItem, error) {
	task, err := db.Get(taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	var incomplete []TaskItem
	for _, todo := range task.TodoItems {
		if !todo.Checked {
			incomplete = append(incomplete, todo)
		}
	}
	return incomplete, nil
}

// UpdateParentTodos updates the TODO entries for a parent task based on its children.
// Returns true if the parent was modified.
func (db *TaskDB) UpdateParentTodos(parentID string) (bool, error) {
	return UpdateParentTodoEntries(db.tasks, parentID)
}

// UpdateParentTodosForChild updates the parent's TODO entries after a child change.
// If the child has no parent, this is a no-op.
// Returns true if the parent was modified.
func (db *TaskDB) UpdateParentTodosForChild(childID string) (bool, error) {
	task, err := db.Get(childID)
	if err != nil {
		return false, fmt.Errorf("task not found: %w", err)
	}
	if task.Meta.Parent == "" {
		return false, nil
	}
	return db.UpdateParentTodos(task.Meta.Parent)
}

// TasksRoot returns the root directory for tasks.
func (db *TaskDB) TasksRoot() string {
	return db.tasksRoot
}

// SetRole sets the role for a task.
func (db *TaskDB) SetRole(taskID, role string) error {
	task, err := db.Get(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	if task.Meta.Role != role {
		task.Meta.Role = role
		task.MarkDirty()
	}
	return nil
}

// SetPriority sets the priority for a task.
func (db *TaskDB) SetPriority(taskID, priority string) error {
	task, err := db.Get(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	normalized := NormalizePriority(priority)
	if !IsValidPriority(normalized) {
		return fmt.Errorf("invalid priority: %s", priority)
	}
	if task.Meta.Priority != normalized {
		task.Meta.Priority = normalized
		task.MarkDirty()
	}
	return nil
}

// SetTitle sets the title for a task.
func (db *TaskDB) SetTitle(taskID, title string) error {
	task, err := db.Get(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	task.SetTitle(title)
	return nil
}

// SetBody sets the body content for a task.
func (db *TaskDB) SetBody(taskID, body string) error {
	task, err := db.Get(taskID)
	if err != nil {
		return fmt.Errorf("task not found: %w", err)
	}
	task.SetBody(body)
	return nil
}
