package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ricochet1k/strandyard/pkg/activity"
	"github.com/ricochet1k/strandyard/pkg/task"
)

func TestCompleteWritesToActivityLog(t *testing.T) {
	repo, _ := setupTestEnv(t)
	if err := runInit(io.Discard, initOptions{ProjectName: "", StorageMode: storageLocal}); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Create a developer role file
	roleFile := filepath.Join(repo, ".strand", "roles", "developer.md")
	roleContent := `# Developer

## Role
Developer (human or AI) — implements tasks, writes code, and produces working software.

## Responsibilities
- Implement tasks assigned by the Architect
- Write clean, maintainable code following project conventions
- Add tests for new functionality
- Document code and update relevant documentation
- Fix bugs and address issues
- Ensure code passes validation and tests before marking tasks complete
`
	if err := os.WriteFile(roleFile, []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	// Create a buffer to capture output
	var output bytes.Buffer

	// Create a test task
	taskID := "T" + testToken(t.Name()) + "-test-task"
	taskDir := filepath.Join(repo, ".strand", "tasks", taskID)

	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatalf("failed to create task dir: %v", err)
	}

	taskFile := filepath.Join(taskDir, taskID+".md")
	taskContent := `---
type: implement
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T13:43:58Z
owner_approval: false
completed: false
---

# Test Task

## Summary
A simple test task for activity log verification.

## Acceptance Criteria
- Task completes successfully
- Activity log records the completion
`

	if err := os.WriteFile(taskFile, []byte(taskContent), 0o644); err != nil {
		t.Fatalf("failed to write task file: %v", err)
	}

	// Run strand complete
	if err := runComplete(&output, "", taskID, 0, "developer", "Test completion report"); err != nil {
		t.Logf("Output:\n%s", output.String())
		t.Fatalf("runComplete failed: %v", err)
	}

	// Verify activity log has the entry
	activityLog, err := activity.Open(filepath.Join(repo, ".strand"))
	if err != nil {
		t.Fatalf("failed to open activity log: %v", err)
	}
	defer activityLog.Close()

	entries, err := activityLog.ReadEntries()
	if err != nil {
		t.Fatalf("failed to read activity log: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 activity log entry, got %d", len(entries))
	}

	entry := entries[0]
	if entry.TaskID != taskID {
		t.Errorf("expected task ID %s, got %s", taskID, entry.TaskID)
	}

	if entry.Type != activity.EventTaskCompleted {
		t.Errorf("expected event type %s, got %s", activity.EventTaskCompleted, entry.Type)
	}

	if entry.Report != "Test completion report" {
		t.Errorf("expected report 'Test completion report', got '%s'", entry.Report)
	}
}

func TestCompleteViaLastTodoWritesToActivityLog(t *testing.T) {
	repo, _ := setupTestEnv(t)
	if err := runInit(io.Discard, initOptions{ProjectName: "", StorageMode: storageLocal}); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Create a developer role file
	roleFile := filepath.Join(repo, ".strand", "roles", "developer.md")
	roleContent := `# Developer

## Role
Developer (human or AI) — implements tasks, writes code, and produces working software.

## Responsibilities
- Implement tasks assigned by the Architect
- Write clean, maintainable code following project conventions
- Add tests for new functionality
- Document code and update relevant documentation
- Fix bugs and address issues
- Ensure code passes validation and tests before marking tasks complete
`
	if err := os.WriteFile(roleFile, []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	// Create a test task with one todo
	taskID := "T" + testToken(t.Name()) + "-test-todo-task"
	taskDir := filepath.Join(repo, ".strand", "tasks", taskID)

	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatalf("failed to create task dir: %v", err)
	}

	taskFile := filepath.Join(taskDir, taskID+".md")
	taskContent := `---
type: implement
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T13:43:58Z
owner_approval: false
completed: false
---

# Test Todo Task

## Summary
A test task with a todo item for activity log verification.

## Acceptance Criteria
- Todo completes successfully
- Activity log records the task completion

## TODOs
- [ ] (role: developer) Complete this todo
`

	if err := os.WriteFile(taskFile, []byte(taskContent), 0o644); err != nil {
		t.Fatalf("failed to write task file: %v", err)
	}

	// Run strand complete with --todo flag
	if err := runComplete(io.Discard, "", taskID, 1, "developer", "Todo completion report"); err != nil {
		t.Fatalf("runComplete failed: %v", err)
	}

	// Verify activity log has the entry
	activityLog, err := activity.Open(filepath.Join(repo, ".strand"))
	if err != nil {
		t.Fatalf("failed to open activity log: %v", err)
	}
	defer activityLog.Close()

	entries, err := activityLog.ReadEntries()
	if err != nil {
		t.Fatalf("failed to read activity log: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 activity log entry, got %d", len(entries))
	}

	entry := entries[0]
	if entry.TaskID != taskID {
		t.Errorf("expected task ID %s, got %s", taskID, entry.TaskID)
	}

	if entry.Type != activity.EventTaskCompleted {
		t.Errorf("expected event type %s, got %s", activity.EventTaskCompleted, entry.Type)
	}

	if entry.Report != "Todo completion report" {
		t.Errorf("expected report 'Todo completion report', got '%s'", entry.Report)
	}
}

func TestActivityLogIntegration(t *testing.T) {
	repo, _ := setupTestEnv(t)
	if err := runInit(io.Discard, initOptions{ProjectName: "", StorageMode: storageLocal}); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Create a developer role file
	roleFile := filepath.Join(repo, ".strand", "roles", "developer.md")
	roleContent := `# Developer

## Role
Developer (human or AI) — implements tasks, writes code, and produces working software.

## Responsibilities
- Implement tasks assigned by the Architect
- Write clean, maintainable code following project conventions
- Add tests for new functionality
- Document code and update relevant documentation
- Fix bugs and address issues
- Ensure code passes validation and tests before marking tasks complete
`
	if err := os.WriteFile(roleFile, []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	// Create multiple tasks and complete them
	tasks := []struct {
		id     string
		title  string
		report string
	}{
		{testToken(t.Name()) + "-task-1", "Task 1", "Completed task 1"},
		{testToken(t.Name()) + "-task-2", "Task 2", "Completed task 2"},
		{testToken(t.Name()) + "-task-3", "Task 3", "Completed task 3"},
	}

	for _, tt := range tasks {
		taskID := "T" + tt.id
		taskDir := filepath.Join(repo, ".strand", "tasks", taskID)

		if err := os.MkdirAll(taskDir, 0o755); err != nil {
			t.Fatalf("failed to create task dir: %v", err)
		}

		taskFile := filepath.Join(taskDir, taskID+".md")
		taskContent := fmt.Sprintf(`---
type: implement
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T13:43:58Z
owner_approval: false
completed: false
---

# %s

## Summary
A test task for integration testing.

## Acceptance Criteria
- Task completes successfully
`, tt.title)

		if err := os.WriteFile(taskFile, []byte(taskContent), 0o644); err != nil {
			t.Fatalf("failed to write task file: %v", err)
		}

		// Complete the task
		if err := runComplete(io.Discard, "", taskID, 0, "developer", tt.report); err != nil {
			t.Fatalf("runComplete failed for %s: %v", taskID, err)
		}
	}

	// Verify activity log has all completions
	activityLog, err := activity.Open(filepath.Join(repo, ".strand"))
	if err != nil {
		t.Fatalf("failed to open activity log: %v", err)
	}
	defer activityLog.Close()

	entries, err := activityLog.ReadEntries()
	if err != nil {
		t.Fatalf("failed to read activity log: %v", err)
	}

	if len(entries) != len(tasks) {
		t.Fatalf("expected %d activity log entries, got %d", len(tasks), len(entries))
	}

	// Verify each entry
	taskIDs := make(map[string]bool)
	for _, tt := range tasks {
		taskIDs["T"+tt.id] = false
	}

	for _, entry := range entries {
		if entry.Type != activity.EventTaskCompleted {
			t.Errorf("expected event type %s, got %s", activity.EventTaskCompleted, entry.Type)
		}
		if _, ok := taskIDs[entry.TaskID]; !ok {
			t.Errorf("unexpected task ID %s in activity log", entry.TaskID)
		} else {
			taskIDs[entry.TaskID] = true
		}
	}

	// Verify all tasks were logged
	for taskID, found := range taskIDs {
		if !found {
			t.Errorf("task ID %s not found in activity log", taskID)
		}
	}

	// Test counting completions since a time
	now := time.Now().UTC()
	hoursAgo := now.Add(-1 * time.Hour)
	count, err := activityLog.CountCompletionsSince(hoursAgo)
	if err != nil {
		t.Fatalf("failed to count completions: %v", err)
	}

	// All completions should be within the last hour
	if count != len(tasks) {
		t.Errorf("expected %d completions since %v, got %d", len(tasks), hoursAgo, count)
	}

	// Test counting completions for specific task
	taskID := "T" + tasks[0].id
	taskCount, err := activityLog.CountCompletionsForTaskSince(taskID, hoursAgo)
	if err != nil {
		t.Fatalf("failed to count completions for task: %v", err)
	}

	if taskCount != 1 {
		t.Errorf("expected 1 completion for task %s, got %d", taskID, taskCount)
	}
}

func TestCompleteTodoUpdatesFreeList(t *testing.T) {
	repo, _ := setupTestEnv(t)
	if err := runInit(io.Discard, initOptions{ProjectName: "", StorageMode: storageLocal}); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Create a developer role file
	roleFile := filepath.Join(repo, ".strand", "roles", "developer.md")
	roleContent := `# Developer

## Role
Developer (human or AI) — implements tasks, writes code, and produces working software.

## Responsibilities
- Implement tasks assigned by the Architect
- Write clean, maintainable code following project conventions
- Add tests for new functionality
- Document code and update relevant documentation
- Fix bugs and address issues
- Ensure code passes validation and tests before marking tasks complete
`
	if err := os.WriteFile(roleFile, []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	// Create a test task with one todo
	taskID := "T" + testToken(t.Name()) + "-todo-freelist"
	taskDir := filepath.Join(repo, ".strand", "tasks", taskID)

	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatalf("failed to create task dir: %v", err)
	}

	taskFile := filepath.Join(taskDir, taskID+".md")
	taskContent := `---
type: implement
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T13:43:58Z
owner_approval: false
completed: false
---

# Test Todo Task for Free List

## Summary
A test task with a todo item to verify free-list updates.

## Acceptance Criteria
- Todo completes successfully
- Free-list is updated to remove the completed task

## TODOs
- [ ] (role: developer) Complete this todo
`

	if err := os.WriteFile(taskFile, []byte(taskContent), 0o644); err != nil {
		t.Fatalf("failed to write task file: %v", err)
	}

	// Get project paths
	paths, err := resolveProjectPaths("")
	if err != nil {
		t.Fatalf("failed to resolve project paths: %v", err)
	}

	// Run strand repair first to create the initial free-list
	if err := runRepair(io.Discard, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
		t.Fatalf("runRepair failed: %v", err)
	}

	// Verify the task is in the free-list before completion
	contentBefore, err := os.ReadFile(paths.FreeTasksFile)
	if err != nil {
		t.Fatalf("failed to read free-list before: %v", err)
	}

	if !bytes.Contains(contentBefore, []byte(taskID)) {
		t.Fatalf("task %s should be in free-list before completion", taskID)
	}

	// Run strand complete with --todo flag for the last todo
	if err := runComplete(io.Discard, "", taskID, 1, "developer", "Completed the todo"); err != nil {
		t.Fatalf("runComplete failed: %v", err)
	}

	// Verify the task is removed from the free-list after completion
	contentAfter, err := os.ReadFile(paths.FreeTasksFile)
	if err != nil {
		t.Fatalf("failed to read free-list after: %v", err)
	}

	if bytes.Contains(contentAfter, []byte(taskID)) {
		t.Fatalf("task %s should not be in free-list after completion", taskID)
	}
}


// TestAtomicityOfFreeListCalculation tests that CalculateIncrementalFreeListUpdate
// correctly identifies newly-unblocked tasks. This is a unit test that doesn't
// require complex setup.
func TestAtomicityOfFreeListCalculation(t *testing.T) {
	// Create test tasks: blocker and blocked
	blocker := &task.Task{
		ID: "T1blocker",
		Meta: task.Metadata{
			Completed: true, // This task is now completed
			Blockers:  []string{},
		},
	}
	
	blocked := &task.Task{
		ID: "T2blocked",
		Meta: task.Metadata{
			Completed: false,
			Blockers:  []string{"T1blocker"}, // This task was blocked
		},
	}
	
	alsoBlocked := &task.Task{
		ID: "T3alsoBlocked",
		Meta: task.Metadata{
			Completed: false,
			Blockers:  []string{"T1blocker", "T4other"}, // This task has multiple blockers
		},
	}
	
	otherBlocker := &task.Task{
		ID: "T4other",
		Meta: task.Metadata{
			Completed: false, // This blocker is NOT yet completed
			Blockers:  []string{},
		},
	}
	
	tasks := map[string]*task.Task{
		"T1blocker":      blocker,
		"T2blocked":      blocked,
		"T3alsoBlocked":  alsoBlocked,
		"T4other":        otherBlocker,
	}
	
	// Calculate what should be added to free-list when T1blocker is completed
	update, err := task.CalculateIncrementalFreeListUpdate(tasks, "T1blocker")
	if err != nil {
		t.Fatalf("CalculateIncrementalFreeListUpdate failed: %v", err)
	}
	
	// Verify T1blocker is removed
	if len(update.RemoveTaskIDs) != 1 || update.RemoveTaskIDs[0] != "T1blocker" {
		t.Errorf("expected T1blocker in RemoveTaskIDs, got %v", update.RemoveTaskIDs)
	}
	
	// Verify T2blocked is added (now unblocked)
	addedIDs := make(map[string]bool)
	for _, t := range update.AddTasks {
		addedIDs[t.ID] = true
	}
	
	if !addedIDs["T2blocked"] {
		t.Errorf("T2blocked should be added (all blockers completed)")
	}
	
	// Verify T3alsoBlocked is NOT added (T4other still blocks it)
	if addedIDs["T3alsoBlocked"] {
		t.Errorf("T3alsoBlocked should NOT be added (still has blocker T4other)")
	}
}
