package task

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDB(t *testing.T) (*TaskDB, string) {
	t.Helper()
	tmpDir := t.TempDir()
	tasksRoot := filepath.Join(tmpDir, "tasks")
	return NewTaskDB(tasksRoot), tasksRoot
}

func createTaskFile(t *testing.T, tasksRoot, id, title string) {
	t.Helper()
	taskDir := filepath.Join(tasksRoot, id)
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := `---
type: ""
role: dev
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-01-30T00:00:00Z
date_edited: 2026-01-30T00:00:00Z
owner_approval: false
completed: false
---

# ` + title + `

Task body content.
`

	taskFile := filepath.Join(taskDir, id+".md")
	if err := os.WriteFile(taskFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write task: %v", err)
	}
}

func TestTaskDB_SetParent(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "P1aaa-parent", "Parent Task")
	createTaskFile(t, tasksRoot, "C1bbb-child", "Child Task")

	// Set parent relationship
	if err := db.SetParent("C1bbb-child", "P1aaa-parent"); err != nil {
		t.Fatalf("SetParent failed: %v", err)
	}

	// Verify child has parent set
	child, err := db.Get("C1bbb-child")
	if err != nil {
		t.Fatalf("Get child: %v", err)
	}
	if child.Meta.Parent != "P1aaa-parent" {
		t.Errorf("expected parent P1aaa-parent, got %q", child.Meta.Parent)
	}
	if !child.Dirty {
		t.Error("child should be marked dirty")
	}

	// Verify children lookup works
	children := db.GetChildren("P1aaa-parent")
	if len(children) != 1 {
		t.Fatalf("expected 1 child, got %d", len(children))
	}
	if children[0].ID != "C1bbb-child" {
		t.Errorf("expected child C1bbb-child, got %s", children[0].ID)
	}
}

func TestTaskDB_SetParent_CycleDetection(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "T1aaa-a", "Task A")
	createTaskFile(t, tasksRoot, "T2bbb-b", "Task B")
	createTaskFile(t, tasksRoot, "T3ccc-c", "Task C")

	// Create chain: A -> B -> C
	if err := db.SetParent("T2bbb-b", "T1aaa-a"); err != nil {
		t.Fatalf("SetParent B->A: %v", err)
	}
	if err := db.SetParent("T3ccc-c", "T2bbb-b"); err != nil {
		t.Fatalf("SetParent C->B: %v", err)
	}

	// Try to create cycle: A -> C (would create A -> B -> C -> A)
	err := db.SetParent("T1aaa-a", "T3ccc-c")
	if err == nil {
		t.Fatal("expected cycle detection error")
	}
}

func TestTaskDB_AddBlocker(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "T1aaa-blocked", "Blocked Task")
	createTaskFile(t, tasksRoot, "T2bbb-blocker", "Blocker Task")

	// Add blocker relationship
	if err := db.AddBlocker("T1aaa-blocked", "T2bbb-blocker"); err != nil {
		t.Fatalf("AddBlocker failed: %v", err)
	}

	// Verify blocked task has blocker
	blocked, _ := db.Get("T1aaa-blocked")
	if len(blocked.Meta.Blockers) != 1 || blocked.Meta.Blockers[0] != "T2bbb-blocker" {
		t.Errorf("expected blocker T2bbb-blocker, got %v", blocked.Meta.Blockers)
	}

	// Verify blocker task has blocks
	blocker, _ := db.Get("T2bbb-blocker")
	if len(blocker.Meta.Blocks) != 1 || blocker.Meta.Blocks[0] != "T1aaa-blocked" {
		t.Errorf("expected blocks T1aaa-blocked, got %v", blocker.Meta.Blocks)
	}

	// Both should be dirty
	if !blocked.Dirty || !blocker.Dirty {
		t.Error("both tasks should be marked dirty")
	}
}

func TestTaskDB_RemoveBlocker(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "T1aaa-blocked", "Blocked Task")
	createTaskFile(t, tasksRoot, "T2bbb-blocker", "Blocker Task")

	// Add then remove blocker
	if err := db.AddBlocker("T1aaa-blocked", "T2bbb-blocker"); err != nil {
		t.Fatalf("AddBlocker failed: %v", err)
	}

	if err := db.RemoveBlocker("T1aaa-blocked", "T2bbb-blocker"); err != nil {
		t.Fatalf("RemoveBlocker failed: %v", err)
	}

	// Verify relationship removed from both sides
	blocked, _ := db.Get("T1aaa-blocked")
	if len(blocked.Meta.Blockers) != 0 {
		t.Errorf("expected no blockers, got %v", blocked.Meta.Blockers)
	}

	blocker, _ := db.Get("T2bbb-blocker")
	if len(blocker.Meta.Blocks) != 0 {
		t.Errorf("expected no blocks, got %v", blocker.Meta.Blocks)
	}
}

func TestTaskDB_SetCompleted(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "T1aaa-task", "Test Task")

	// Mark as completed
	if err := db.SetCompleted("T1aaa-task", true); err != nil {
		t.Fatalf("SetCompleted failed: %v", err)
	}

	task, _ := db.Get("T1aaa-task")
	if !task.Meta.Completed {
		t.Error("task should be completed")
	}
	if !task.Dirty {
		t.Error("task should be marked dirty")
	}

	// Mark as incomplete
	if err := db.SetCompleted("T1aaa-task", false); err != nil {
		t.Fatalf("SetCompleted false failed: %v", err)
	}

	task, _ = db.Get("T1aaa-task")
	if task.Meta.Completed {
		t.Error("task should not be completed")
	}
}

func TestTaskDB_UpdateBlockersAfterCompletion(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "T1aaa-blocked", "Blocked Task")
	createTaskFile(t, tasksRoot, "T2bbb-blocker", "Blocker Task")

	// Set up blocker relationship
	if err := db.AddBlocker("T1aaa-blocked", "T2bbb-blocker"); err != nil {
		t.Fatalf("AddBlocker failed: %v", err)
	}

	// Complete the blocker
	if err := db.SetCompleted("T2bbb-blocker", true); err != nil {
		t.Fatalf("SetCompleted failed: %v", err)
	}

	// Update blockers after completion
	if err := db.UpdateBlockersAfterCompletion("T2bbb-blocker"); err != nil {
		t.Fatalf("UpdateBlockersAfterCompletion failed: %v", err)
	}

	// Verify blocker removed from blocked task
	blocked, _ := db.Get("T1aaa-blocked")
	if len(blocked.Meta.Blockers) != 0 {
		t.Errorf("expected no blockers after completion, got %v", blocked.Meta.Blockers)
	}
}

func TestTaskDB_FixBlockerRelationships(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "T1aaa-blocked", "Blocked Task")
	createTaskFile(t, tasksRoot, "T2bbb-blocker", "Blocker Task")

	// Manually create inconsistent state (only set blockers, not blocks)
	blocked, _ := db.Get("T1aaa-blocked")
	blocked.Meta.Blockers = []string{"T2bbb-blocker"}
	blocked.MarkDirty()

	// Fix relationships
	modified := db.FixBlockerRelationships()
	if modified == 0 {
		t.Error("expected at least one task to be modified")
	}

	// Verify both sides are now consistent
	blocker, _ := db.Get("T2bbb-blocker")
	if len(blocker.Meta.Blocks) != 1 || blocker.Meta.Blocks[0] != "T1aaa-blocked" {
		t.Errorf("expected blocks to be fixed, got %v", blocker.Meta.Blocks)
	}
}

func TestTaskDB_SyncBlockersFromChildren(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "P1aaa-parent", "Parent Task")
	createTaskFile(t, tasksRoot, "C1bbb-child1", "Child 1")
	createTaskFile(t, tasksRoot, "C2ccc-child2", "Child 2")

	// Set up parent-child relationships
	if err := db.SetParent("C1bbb-child1", "P1aaa-parent"); err != nil {
		t.Fatalf("SetParent child1: %v", err)
	}
	if err := db.SetParent("C2ccc-child2", "P1aaa-parent"); err != nil {
		t.Fatalf("SetParent child2: %v", err)
	}

	// Sync blockers from children
	modified, err := db.SyncBlockersFromChildren()
	if err != nil {
		t.Fatalf("SyncBlockersFromChildren failed: %v", err)
	}
	if modified == 0 {
		t.Error("expected at least one task to be modified")
	}

	// Verify parent is now blocked by both children
	parent, _ := db.Get("P1aaa-parent")
	if len(parent.Meta.Blockers) != 2 {
		t.Fatalf("expected 2 blockers, got %d: %v", len(parent.Meta.Blockers), parent.Meta.Blockers)
	}

	// Complete one child
	if err := db.SetCompleted("C1bbb-child1", true); err != nil {
		t.Fatalf("SetCompleted: %v", err)
	}

	// Sync again
	if _, err := db.SyncBlockersFromChildren(); err != nil {
		t.Fatalf("SyncBlockersFromChildren after completion: %v", err)
	}

	// Verify parent now only blocked by incomplete child
	parent, _ = db.Get("P1aaa-parent")
	if len(parent.Meta.Blockers) != 1 {
		t.Fatalf("expected 1 blocker after child completion, got %d: %v", len(parent.Meta.Blockers), parent.Meta.Blockers)
	}
	if parent.Meta.Blockers[0] != "C2ccc-child2" {
		t.Errorf("expected blocker C2ccc-child2, got %s", parent.Meta.Blockers[0])
	}
}

func TestTaskDB_LazyLoading(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "T1aaa-task", "Test Task")

	// Task should not be loaded yet
	if len(db.tasks) != 0 {
		t.Error("no tasks should be loaded initially")
	}

	// Get should trigger lazy load
	task, err := db.Get("T1aaa-task")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if task.ID != "T1aaa-task" {
		t.Errorf("expected task T1aaa-task, got %s", task.ID)
	}

	// Task should now be in cache
	if len(db.tasks) != 1 {
		t.Error("task should be cached after Get")
	}

	// Second Get should use cache (not reload)
	task2, err := db.Get("T1aaa-task")
	if err != nil {
		t.Fatalf("second Get failed: %v", err)
	}
	if task != task2 {
		t.Error("should return same cached instance")
	}
}

func TestTaskDB_SaveAndReload(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "T1aaa-task", "Original Title")

	// Load and modify task
	task, _ := db.Get("T1aaa-task")
	task.SetTitle("Modified Title")

	// Save
	if err := db.Save("T1aaa-task"); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Create new DB instance and reload
	db2 := NewTaskDB(tasksRoot)
	reloaded, err := db2.Get("T1aaa-task")
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}

	if reloaded.TitleContent != "Modified Title" {
		t.Errorf("expected title 'Modified Title', got %q", reloaded.TitleContent)
	}
}

func TestTaskDB_GetOrCreate(t *testing.T) {
	db, _ := setupTestDB(t)

	// Create new task that doesn't exist on disk
	task, err := db.GetOrCreate("T1aaa-new")
	if err != nil {
		t.Fatalf("GetOrCreate failed: %v", err)
	}

	if task.ID != "T1aaa-new" {
		t.Errorf("expected ID T1aaa-new, got %s", task.ID)
	}
	if !task.Dirty {
		t.Error("new task should be marked dirty")
	}

	// Second call should return same instance
	task2, err := db.GetOrCreate("T1aaa-new")
	if err != nil {
		t.Fatalf("second GetOrCreate failed: %v", err)
	}
	if task != task2 {
		t.Error("should return same instance")
	}
}

func TestTaskDB_ClearParent(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	createTaskFile(t, tasksRoot, "P1aaa-parent", "Parent Task")
	createTaskFile(t, tasksRoot, "C1bbb-child", "Child Task")

	// Set then clear parent
	if err := db.SetParent("C1bbb-child", "P1aaa-parent"); err != nil {
		t.Fatalf("SetParent failed: %v", err)
	}

	if err := db.ClearParent("C1bbb-child"); err != nil {
		t.Fatalf("ClearParent failed: %v", err)
	}

	child, _ := db.Get("C1bbb-child")
	if child.Meta.Parent != "" {
		t.Errorf("expected no parent, got %q", child.Meta.Parent)
	}

	children := db.GetChildren("P1aaa-parent")
	if len(children) != 0 {
		t.Errorf("expected no children, got %d", len(children))
	}
}

func TestTaskDB_GetAncestors(t *testing.T) {
	db, tasksRoot := setupTestDB(t)

	// Create task hierarchy: Root -> Child -> GrandChild
	createTaskFile(t, tasksRoot, "R1aaa-root", "Root Task")
	createTaskFile(t, tasksRoot, "C1bbb-child", "Child Task")
	createTaskFile(t, tasksRoot, "G1ccc-grandchild", "GrandChild Task")

	// Set up parent relationships
	if err := db.SetParent("C1bbb-child", "R1aaa-root"); err != nil {
		t.Fatalf("SetParent failed: %v", err)
	}
	if err := db.SetParent("G1ccc-grandchild", "C1bbb-child"); err != nil {
		t.Fatalf("SetParent failed: %v", err)
	}

	// Test GetAncestors for root task (no parent)
	ancestors := db.GetAncestors("R1aaa-root")
	if len(ancestors) != 0 {
		t.Fatalf("expected no ancestors for root, got %d", len(ancestors))
	}

	// Test GetAncestors for child task
	ancestors = db.GetAncestors("C1bbb-child")
	if len(ancestors) != 1 {
		t.Fatalf("expected 1 ancestor for child, got %d", len(ancestors))
	}
	if ancestors[0][0] != "R1aaa" || ancestors[0][1] != "Root Task" {
		t.Errorf("expected [R1aaa Root Task], got [%s %s]", ancestors[0][0], ancestors[0][1])
	}

	// Test GetAncestors for grandchild task
	ancestors = db.GetAncestors("G1ccc-grandchild")
	if len(ancestors) != 2 {
		t.Fatalf("expected 2 ancestors for grandchild, got %d", len(ancestors))
	}
	if ancestors[0][0] != "C1bbb" || ancestors[0][1] != "Child Task" {
		t.Errorf("expected first ancestor [C1bbb Child Task], got [%s %s]", ancestors[0][0], ancestors[0][1])
	}
	if ancestors[1][0] != "R1aaa" || ancestors[1][1] != "Root Task" {
		t.Errorf("expected second ancestor [R1aaa Root Task], got [%s %s]", ancestors[1][0], ancestors[1][1])
	}

	// Test GetAncestors for non-existent task
	ancestors = db.GetAncestors("NONEXISTENT")
	if len(ancestors) != 0 {
		t.Fatalf("expected no ancestors for non-existent task, got %d", len(ancestors))
	}
}
