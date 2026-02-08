package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestIsActiveStatus(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{"open", true},
		{"in_progress", true},
		{"done", false},
		{"cancelled", false},
		{"duplicate", false},
		{"", true}, // Empty defaults to active for backward compatibility
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status=%q", tt.status), func(t *testing.T) {
			result := IsActiveStatus(tt.status)
			if result != tt.expected {
				t.Errorf("IsActiveStatus(%q) = %v, expected %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGenerateMasterListsExcludesNonActiveStatus(t *testing.T) {
	tmpDir := t.TempDir()
	tasksRoot := filepath.Join(tmpDir, "tasks")
	if err := os.MkdirAll(tasksRoot, 0o755); err != nil {
		t.Fatalf("failed to create tasks root: %v", err)
	}

	parser := NewParser()
	roleName := testRoleName(t, "exclude")
	roleDir := filepath.Join(tmpDir, "roles")
	if err := os.MkdirAll(roleDir, 0o755); err != nil {
		t.Fatalf("failed to create role dir: %v", err)
	}

	roleContent := fmt.Sprintf("# %s\n\nRole for testing.", strings.Title(roleName))
	if err := os.WriteFile(filepath.Join(roleDir, roleName+".md"), []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	tasks := make(map[string]*Task)

	// Create task with done status - should NOT be in free-list
	t1Content := fmt.Sprintf(`---
role: %s
priority: high
status: done
---

# Completed Task

This task is done.
`, roleName)
	t1, _ := parser.ParseString(t1Content, "T1aaa-done")
	t1.FilePath = "tasks/T1aaa-done/T1aaa-done.md"
	t1.SetTitle("Completed Task")
	tasks[t1.ID] = t1

	// Create task with cancelled status - should NOT be in free-list
	t2Content := fmt.Sprintf(`---
role: %s
priority: high
status: cancelled
---

# Cancelled Task

This task is cancelled.
`, roleName)
	t2, _ := parser.ParseString(t2Content, "T2bbb-cancelled")
	t2.FilePath = "tasks/T2bbb-cancelled/T2bbb-cancelled.md"
	t2.SetTitle("Cancelled Task")
	tasks[t2.ID] = t2

	// Create task with duplicate status - should NOT be in free-list
	t3Content := fmt.Sprintf(`---
role: %s
priority: high
status: duplicate
---

# Duplicate Task

This task is a duplicate.
`, roleName)
	t3, _ := parser.ParseString(t3Content, "T3ccc-duplicate")
	t3.FilePath = "tasks/T3ccc-duplicate/T3ccc-duplicate.md"
	t3.SetTitle("Duplicate Task")
	tasks[t3.ID] = t3

	// Create task with open status - SHOULD be in free-list
	t4Content := fmt.Sprintf(`---
role: %s
priority: high
status: open
---

# Open Task

This task is open.
`, roleName)
	t4, _ := parser.ParseString(t4Content, "T4ddd-open")
	t4.FilePath = "tasks/T4ddd-open/T4ddd-open.md"
	t4.SetTitle("Open Task")
	tasks[t4.ID] = t4

	// Create task with in_progress status - SHOULD be in free-list
	t5Content := fmt.Sprintf(`---
role: %s
priority: medium
status: in_progress
---

# In Progress Task

This task is in progress.
`, roleName)
	t5, _ := parser.ParseString(t5Content, "T5eee-in-progress")
	t5.FilePath = "tasks/T5eee-in-progress/T5eee-in-progress.md"
	t5.SetTitle("In Progress Task")
	tasks[t5.ID] = t5

	// Create task with empty status (backward compatibility) - SHOULD be in free-list
	t6Content := fmt.Sprintf(`---
role: %s
priority: low
---

# Legacy Task

This task has no status field.
`, roleName)
	t6, _ := parser.ParseString(t6Content, "T6fff-legacy")
	t6.FilePath = "tasks/T6fff-legacy/T6fff-legacy.md"
	t6.SetTitle("Legacy Task")
	tasks[t6.ID] = t6

	rootsFile := filepath.Join(tasksRoot, "root-tasks.md")
	freeFile := filepath.Join(tasksRoot, "free-tasks.md")

	if err := GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		t.Fatalf("GenerateMasterLists failed: %v", err)
	}

	content, err := os.ReadFile(freeFile)
	if err != nil {
		t.Fatalf("failed to read free-tasks.md: %v", err)
	}
	freeContent := string(content)

	// Verify that only open, in_progress, and empty status tasks are in free-list
	tests := []struct {
		taskID      string
		shouldExist bool
	}{
		{"T1aaa-done", false},
		{"T2bbb-cancelled", false},
		{"T3ccc-duplicate", false},
		{"T4ddd-open", true},
		{"T5eee-in-progress", true},
		{"T6fff-legacy", true},
	}

	for _, tt := range tests {
		inList := strings.Contains(freeContent, tt.taskID)
		if inList != tt.shouldExist {
			t.Errorf("Task %s: expected in list = %v, got %v", tt.taskID, tt.shouldExist, inList)
		}
	}
}

func TestUpdateFreeListIncrementallyWithStatus(t *testing.T) {
	tmpDir := t.TempDir()
	tasksRoot := filepath.Join(tmpDir, "tasks")
	if err := os.MkdirAll(tasksRoot, 0o755); err != nil {
		t.Fatalf("failed to create tasks root: %v", err)
	}

	parser := NewParser()
	roleName := testRoleName(t, "incremental")
	roleDir := filepath.Join(tmpDir, "roles")
	if err := os.MkdirAll(roleDir, 0o755); err != nil {
		t.Fatalf("failed to create role dir: %v", err)
	}

	roleContent := fmt.Sprintf("# %s\n\nRole for testing.", strings.Title(roleName))
	if err := os.WriteFile(filepath.Join(roleDir, roleName+".md"), []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	tasks := make(map[string]*Task)

	// Create initial tasks
	t1Content := fmt.Sprintf(`---
role: %s
priority: high
status: open
---

# Task 1

First task.
`, roleName)
	t1, _ := parser.ParseString(t1Content, "T1aaa-task1")
	t1.FilePath = "tasks/T1aaa-task1/T1aaa-task1.md"
	t1.SetTitle("Task 1")
	tasks[t1.ID] = t1

	t2Content := fmt.Sprintf(`---
role: %s
priority: high
status: in_progress
---

# Task 2

Second task.
`, roleName)
	t2, _ := parser.ParseString(t2Content, "T2bbb-task2")
	t2.FilePath = "tasks/T2bbb-task2/T2bbb-task2.md"
	t2.SetTitle("Task 2")
	tasks[t2.ID] = t2

	freeFile := filepath.Join(tasksRoot, "free-tasks.md")
	rootsFile := filepath.Join(tasksRoot, "root-tasks.md")

	// Generate initial free-list
	if err := GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		t.Fatalf("GenerateMasterLists failed: %v", err)
	}

	content, _ := os.ReadFile(freeFile)
	initialContent := string(content)

	// Verify both tasks are in free-list initially
	if !strings.Contains(initialContent, "T1aaa-task1") {
		t.Error("T1aaa-task1 should be in initial free-list")
	}
	if !strings.Contains(initialContent, "T2bbb-task2") {
		t.Error("T2bbb-task2 should be in initial free-list")
	}

	// Update task 1 status to done
	t1.Meta.Status = "done"
	t1.MarkDirty()
	tasks[t1.ID] = t1

	// Calculate incremental update (remove T1aaa-task1)
	update := IncrementalFreeListUpdate{
		RemoveTaskIDs: []string{"T1aaa-task1"},
		AddTasks:      []*Task{},
	}

	if err := UpdateFreeListIncrementally(tasks, freeFile, update); err != nil {
		t.Fatalf("UpdateFreeListIncrementally failed: %v", err)
	}

	updatedContent, _ := os.ReadFile(freeFile)
	updatedStr := string(updatedContent)

	// Verify T1aaa-task1 is removed from free-list
	if strings.Contains(updatedStr, "T1aaa-task1") {
		t.Error("T1aaa-task1 should not be in updated free-list")
	}

	// Verify T2bbb-task2 is still in free-list
	if !strings.Contains(updatedStr, "T2bbb-task2") {
		t.Error("T2bbb-task2 should still be in updated free-list")
	}
}

func TestCalculateIncrementalFreeListUpdateWithStatus(t *testing.T) {
	parser := NewParser()
	roleName := testRoleName(t, "calculate")

	tasks := make(map[string]*Task)

	// Create blocker task
	blockerContent := fmt.Sprintf(`---
role: %s
priority: high
status: open
---

# Blocker Task

This task blocks others.
`, roleName)
	blocker, _ := parser.ParseString(blockerContent, "Tblkr-blocker")
	blocker.FilePath = "tasks/Tblkr-blocker/Tblkr-blocker.md"
	blocker.SetTitle("Blocker Task")
	tasks[blocker.ID] = blocker

	// Create task blocked by the blocker
	blockedContent := fmt.Sprintf(`---
role: %s
priority: high
blockers:
  - Tblkr-blocker
status: open
---

# Blocked Task

This task is blocked.
`, roleName)
	blocked, _ := parser.ParseString(blockedContent, "Tblkd-blocked")
	blocked.FilePath = "tasks/Tblkd-blocked/Tblkd-blocked.md"
	blocked.SetTitle("Blocked Task")
	tasks[blocked.ID] = blocked

	// Calculate update when blocker is completed
	blocker.Meta.Status = "done"
	blocker.Meta.Completed = true
	tasks[blocker.ID] = blocker

	update, err := CalculateIncrementalFreeListUpdate(tasks, blocker.ID)
	if err != nil {
		t.Fatalf("CalculateIncrementalFreeListUpdate failed: %v", err)
	}

	// Verify blocker is in RemoveTaskIDs
	if len(update.RemoveTaskIDs) != 1 || update.RemoveTaskIDs[0] != blocker.ID {
		t.Errorf("expected RemoveTaskIDs to contain %q, got %v", blocker.ID, update.RemoveTaskIDs)
	}

	// Verify blocked task is added to AddTasks (since it's no longer blocked)
	if len(update.AddTasks) != 1 || update.AddTasks[0].ID != blocked.ID {
		t.Errorf("expected AddTasks to contain %q, got %v", blocked.ID, update.AddTasks)
	}
}

func TestFreeListDoesNotIncludeInactiveStatusWhenBlocked(t *testing.T) {
	parser := NewParser()
	roleName := testRoleName(t, "blocked-inactive")

	tasks := make(map[string]*Task)

	// Create blocker task
	blockerContent := fmt.Sprintf(`---
role: %s
priority: high
status: open
---

# Blocker Task

This task blocks others.
`, roleName)
	blocker, _ := parser.ParseString(blockerContent, "Tblkr-inactive")
	blocker.FilePath = "tasks/Tblkr-inactive/Tblkr-inactive.md"
	blocker.SetTitle("Blocker Task")
	tasks[blocker.ID] = blocker

	// Create task with non-active status that's also blocked
	inactiveContent := fmt.Sprintf(`---
role: %s
priority: high
blockers:
  - Tblkr-inactive
status: cancelled
---

# Cancelled Blocked Task

This task is cancelled and blocked.
`, roleName)
	inactive, _ := parser.ParseString(inactiveContent, "Tcnld-inactive")
	inactive.FilePath = "tasks/Tcnld-inactive/Tcnld-inactive.md"
	inactive.SetTitle("Cancelled Blocked Task")
	tasks[inactive.ID] = inactive

	// Calculate update when blocker is completed
	blocker.Meta.Status = "done"
	blocker.Meta.Completed = true
	tasks[blocker.ID] = blocker

	update, err := CalculateIncrementalFreeListUpdate(tasks, blocker.ID)
	if err != nil {
		t.Fatalf("CalculateIncrementalFreeListUpdate failed: %v", err)
	}

	// Verify that cancelled task is NOT added to free-list even though blocker is removed
	for _, task := range update.AddTasks {
		if task.ID == inactive.ID {
			t.Errorf("inactive status task should not be added to free-list")
		}
	}
}

func TestBackwardCompatibilityWithCompletedField(t *testing.T) {
	tmpDir := t.TempDir()
	tasksRoot := filepath.Join(tmpDir, "tasks")
	if err := os.MkdirAll(tasksRoot, 0o755); err != nil {
		t.Fatalf("failed to create tasks root: %v", err)
	}

	parser := NewParser()
	roleName := testRoleName(t, "backward")
	roleDir := filepath.Join(tmpDir, "roles")
	if err := os.MkdirAll(roleDir, 0o755); err != nil {
		t.Fatalf("failed to create role dir: %v", err)
	}

	roleContent := fmt.Sprintf("# %s\n\nRole for testing.", strings.Title(roleName))
	if err := os.WriteFile(filepath.Join(roleDir, roleName+".md"), []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	tasks := make(map[string]*Task)

	// Create task with old completed: true field
	oldContent := fmt.Sprintf(`---
role: %s
priority: high
completed: true
---

# Old Completed Task

This task uses the old completed field.
`, roleName)
	oldTask, _ := parser.ParseString(oldContent, "Told0-legacy")
	oldTask.FilePath = "tasks/Told0-legacy/Told0-legacy.md"
	oldTask.SetTitle("Old Completed Task")
	tasks[oldTask.ID] = oldTask

	// Create task with old completed: false field and no status
	activeContent := fmt.Sprintf(`---
role: %s
priority: high
completed: false
---

# Old Active Task

This task uses the old completed field but is active.
`, roleName)
	activeTask, _ := parser.ParseString(activeContent, "Tact0-legacy-active")
	activeTask.FilePath = "tasks/Tact0-legacy-active/Tact0-legacy-active.md"
	activeTask.SetTitle("Old Active Task")
	tasks[activeTask.ID] = activeTask

	rootsFile := filepath.Join(tasksRoot, "root-tasks.md")
	freeFile := filepath.Join(tasksRoot, "free-tasks.md")

	if err := GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		t.Fatalf("GenerateMasterLists failed: %v", err)
	}

	content, _ := os.ReadFile(freeFile)
	freeContent := string(content)

	// Verify completed: true task is excluded
	if strings.Contains(freeContent, "Told0-legacy") {
		t.Error("Task with completed: true should not be in free-list")
	}

	// Verify completed: false task is included (since status is empty and defaults to active)
	if !strings.Contains(freeContent, "Tact0-legacy-active") {
		t.Error("Task with completed: false and no status should be in free-list")
	}
}

func TestFreeListStatusValidation(t *testing.T) {
	tmpDir := t.TempDir()
	tasksRoot := filepath.Join(tmpDir, "tasks")
	if err := os.MkdirAll(tasksRoot, 0o755); err != nil {
		t.Fatalf("failed to create tasks root: %v", err)
	}

	parser := NewParser()
	roleName := testRoleName(t, "validation")
	roleDir := filepath.Join(tmpDir, "roles")
	if err := os.MkdirAll(roleDir, 0o755); err != nil {
		t.Fatalf("failed to create role dir: %v", err)
	}

	roleContent := fmt.Sprintf("# %s\n\nRole for testing.", strings.Title(roleName))
	if err := os.WriteFile(filepath.Join(roleDir, roleName+".md"), []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	tasks := make(map[string]*Task)

	// Create various status tasks
	statuses := []string{"open", "in_progress", "done", "cancelled", "duplicate", ""}
	for i, status := range statuses {
		taskID := fmt.Sprintf("Tval%d-status-%d", i, i)
		var content string
		if status == "" {
			content = fmt.Sprintf(`---
role: %s
priority: high
---

# Task %d

Status: (empty)
`, roleName, i)
		} else {
			content = fmt.Sprintf(`---
role: %s
priority: high
status: %s
---

# Task %d

Status: %s
`, roleName, status, i, status)
		}
		task, _ := parser.ParseString(content, taskID)
		task.FilePath = fmt.Sprintf("tasks/%s/%s.md", taskID, taskID)
		task.SetTitle(fmt.Sprintf("Task %d", i))
		tasks[task.ID] = task
	}

	rootsFile := filepath.Join(tasksRoot, "root-tasks.md")
	freeFile := filepath.Join(tasksRoot, "free-tasks.md")

	if err := GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		t.Fatalf("GenerateMasterLists failed: %v", err)
	}

	content, _ := os.ReadFile(freeFile)
	freeContent := string(content)

	// Verify all active-status tasks are in free-list
	expectedInList := map[string]bool{
		"Tval0-status-0": true,  // open
		"Tval1-status-1": true,  // in_progress
		"Tval2-status-2": false, // done
		"Tval3-status-3": false, // cancelled
		"Tval4-status-4": false, // duplicate
		"Tval5-status-5": true,  // empty (defaults to open)
	}

	for taskID, shouldExist := range expectedInList {
		inList := strings.Contains(freeContent, taskID)
		if inList != shouldExist {
			t.Errorf("Task %s: expected in list = %v, got %v", taskID, shouldExist, inList)
		}
	}
}

func TestParseListWithStatusField(t *testing.T) {
	parser := NewParser()
	roleName := testRoleName(t, "parse-list")

	tasks := make(map[string]*Task)

	// Create a task with status field
	content := fmt.Sprintf(`---
role: %s
priority: high
status: open
date_created: 2026-01-01T00:00:00Z
date_edited: 2026-01-01T00:00:00Z
---

# Test Task

This task has a status field.
`, roleName)

	task, _ := parser.ParseString(content, "Tlist-status")
	task.FilePath = "tasks/Tlist-status/Tlist-status.md"
	task.SetTitle("Test Task")
	tasks[task.ID] = task

	// Create a simple free-list content string
	listContent := `# Free tasks

## High

- [Test Task](tasks/Tlist-status/Tlist-status.md)
`

	// Parse the free-list
	parsed := ParseFreeList(listContent, tasks)

	if len(parsed.TaskIDs) != 1 || parsed.TaskIDs[0] != "Tlist-status" {
		t.Errorf("expected TaskIDs to contain Tlist-status, got %v", parsed.TaskIDs)
	}

	if parsed.Title != "Free tasks" {
		t.Errorf("expected title 'Free tasks', got %q", parsed.Title)
	}
}

func TestPriorityGroupsWithStatusField(t *testing.T) {
	tmpDir := t.TempDir()
	tasksRoot := filepath.Join(tmpDir, "tasks")
	if err := os.MkdirAll(tasksRoot, 0o755); err != nil {
		t.Fatalf("failed to create tasks root: %v", err)
	}

	parser := NewParser()
	roleName := testRoleName(t, "priority-status")
	roleDir := filepath.Join(tmpDir, "roles")
	if err := os.MkdirAll(roleDir, 0o755); err != nil {
		t.Fatalf("failed to create role dir: %v", err)
	}

	roleContent := fmt.Sprintf("# %s\n\nRole for testing.", strings.Title(roleName))
	if err := os.WriteFile(filepath.Join(roleDir, roleName+".md"), []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	tasks := make(map[string]*Task)

	// Create tasks with different priorities and statuses
	priorities := []string{"high", "medium", "low"}
	for _, priority := range priorities {
		// Active task with this priority
		activeContent := fmt.Sprintf(`---
role: %s
priority: %s
status: open
---

# %s Priority Open Task

This task is %s priority and open.
`, roleName, priority, strings.Title(priority), priority)
		activeTask, _ := parser.ParseString(activeContent, fmt.Sprintf("Tpri%s-open", priority[:1]))
		activeTask.FilePath = fmt.Sprintf("tasks/Tpri%s-open/Tpri%s-open.md", priority[:1], priority[:1])
		activeTask.SetTitle(fmt.Sprintf("%s Priority Open Task", strings.Title(priority)))
		tasks[activeTask.ID] = activeTask

		// Inactive task with this priority
		inactiveContent := fmt.Sprintf(`---
role: %s
priority: %s
status: done
---

# %s Priority Done Task

This task is %s priority and done.
`, roleName, priority, strings.Title(priority), priority)
		inactiveTask, _ := parser.ParseString(inactiveContent, fmt.Sprintf("Tpri%s-done", priority[:1]))
		inactiveTask.FilePath = fmt.Sprintf("tasks/Tpri%s-done/Tpri%s-done.md", priority[:1], priority[:1])
		inactiveTask.SetTitle(fmt.Sprintf("%s Priority Done Task", strings.Title(priority)))
		tasks[inactiveTask.ID] = inactiveTask
	}

	rootsFile := filepath.Join(tasksRoot, "root-tasks.md")
	freeFile := filepath.Join(tasksRoot, "free-tasks.md")

	if err := GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		t.Fatalf("GenerateMasterLists failed: %v", err)
	}

	content, _ := os.ReadFile(freeFile)
	freeContent := string(content)

	// Verify priority sections are maintained and only active tasks are listed
	highSection := false
	mediumSection := false
	lowSection := false

	lines := strings.Split(freeContent, "\n")
	for _, line := range lines {
		if strings.Contains(line, "## High") {
			highSection = true
		}
		if strings.Contains(line, "## Medium") {
			mediumSection = true
		}
		if strings.Contains(line, "## Low") {
			lowSection = true
		}
	}

	if !highSection || !mediumSection || !lowSection {
		t.Error("free-list should have High, Medium, and Low priority sections")
	}

	// Verify only open tasks are in each priority section
	for _, priority := range priorities {
		openID := fmt.Sprintf("Tpri%s-open", priority[:1])
		doneID := fmt.Sprintf("Tpri%s-done", priority[:1])

		if !strings.Contains(freeContent, openID) {
			t.Errorf("open task %s should be in free-list", openID)
		}
		if strings.Contains(freeContent, doneID) {
			t.Errorf("done task %s should not be in free-list", doneID)
		}
	}
}

func TestFreeListPreservesParentSubtaskOrder(t *testing.T) {
	parser := NewParser()

	parent, _ := parser.ParseString(`# Parent

## Subtasks
- [ ] (subtask: T2bbb) Second child
- [ ] (subtask: T1aaa) First child
`, "P1parent")
	parent.FilePath = "tasks/P1parent.md"

	first, _ := parser.ParseString("# First child\n", "T1aaa-first")
	first.Meta.Parent = "P1parent"
	first.Meta.Priority = PriorityHigh
	first.FilePath = "tasks/T1aaa-first.md"

	second, _ := parser.ParseString("# Second child\n", "T2bbb-second")
	second.Meta.Parent = "P1parent"
	second.Meta.Priority = PriorityHigh
	second.FilePath = "tasks/T2bbb-second.md"

	tasks := map[string]*Task{
		"P1parent":     parent,
		"T1aaa-first":  first,
		"T2bbb-second": second,
	}

	tmp := t.TempDir()
	rootsFile := filepath.Join(tmp, "root-tasks.md")
	freeFile := filepath.Join(tmp, "free-tasks.md")

	if err := GenerateMasterLists(tasks, "tasks", rootsFile, freeFile); err != nil {
		t.Fatalf("GenerateMasterLists failed: %v", err)
	}

	contentBytes, err := os.ReadFile(freeFile)
	if err != nil {
		t.Fatalf("read free list: %v", err)
	}
	content := string(contentBytes)

	secondIdx := strings.Index(content, "T2bbb-second")
	firstIdx := strings.Index(content, "T1aaa-first")
	if secondIdx == -1 || firstIdx == -1 {
		t.Fatalf("expected both subtasks in free list:\n%s", content)
	}
	if secondIdx > firstIdx {
		t.Fatalf("expected second child before first child based on parent order:\n%s", content)
	}
}

func TestFreeListGenerationTiming(t *testing.T) {
	tmpDir := t.TempDir()
	tasksRoot := filepath.Join(tmpDir, "tasks")
	if err := os.MkdirAll(tasksRoot, 0o755); err != nil {
		t.Fatalf("failed to create tasks root: %v", err)
	}

	parser := NewParser()
	roleName := testRoleName(t, "timing")
	roleDir := filepath.Join(tmpDir, "roles")
	if err := os.MkdirAll(roleDir, 0o755); err != nil {
		t.Fatalf("failed to create role dir: %v", err)
	}

	roleContent := fmt.Sprintf("# %s\n\nRole for testing.", strings.Title(roleName))
	if err := os.WriteFile(filepath.Join(roleDir, roleName+".md"), []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	origWd, _ := os.Getwd()
	defer os.Chdir(origWd)
	os.Chdir(tmpDir)

	tasks := make(map[string]*Task)

	// Create a task
	taskContent := fmt.Sprintf(`---
role: %s
priority: high
status: open
date_created: 2026-01-01T00:00:00Z
date_edited: 2026-01-01T00:00:00Z
---

# Status Transition Task

This task transitions between statuses.
`, roleName)
	task, _ := parser.ParseString(taskContent, "Ttmng-status-transition")
	task.FilePath = "tasks/Ttmng-status-transition/Ttmng-status-transition.md"
	task.SetTitle("Status Transition Task")
	tasks[task.ID] = task

	rootsFile := filepath.Join(tasksRoot, "root-tasks.md")
	freeFile := filepath.Join(tasksRoot, "free-tasks.md")

	// Initial generation - task should be in free-list
	if err := GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		t.Fatalf("initial GenerateMasterLists failed: %v", err)
	}

	content, _ := os.ReadFile(freeFile)
	initialContent := string(content)
	if !strings.Contains(initialContent, "Ttmng-status-transition") {
		t.Error("task should be in initial free-list with status: open")
	}

	// Update task status to done
	task.Meta.Status = "done"
	task.Meta.DateEdited = time.Now().UTC()
	tasks[task.ID] = task

	// Regenerate - task should NOT be in free-list
	if err := GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		t.Fatalf("second GenerateMasterLists failed: %v", err)
	}

	content, _ = os.ReadFile(freeFile)
	updatedContent := string(content)
	if strings.Contains(updatedContent, "Ttmng-status-transition") {
		t.Error("task should not be in free-list with status: done")
	}

	// Update task status to in_progress
	task.Meta.Status = "in_progress"
	task.Meta.DateEdited = time.Now().UTC()
	tasks[task.ID] = task

	// Regenerate - task should be in free-list
	if err := GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		t.Fatalf("third GenerateMasterLists failed: %v", err)
	}

	content, _ = os.ReadFile(freeFile)
	finalContent := string(content)
	if !strings.Contains(finalContent, "Ttmng-status-transition") {
		t.Error("task should be in free-list with status: in_progress")
	}
}
