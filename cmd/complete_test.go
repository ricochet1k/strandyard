package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/ricochet1k/strandyard/pkg/activity"
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
