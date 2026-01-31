package task_test

import (
	"fmt"
	"log"

	"github.com/ricochet1k/strandyard/pkg/task"
)

// Example demonstrating basic TaskDB usage
func ExampleTaskDB_basic() {
	// Create a new TaskDB
	db := task.NewTaskDB("tasks")

	// Load all tasks from disk
	if err := db.LoadAll(); err != nil {
		log.Fatal(err)
	}

	// Get a specific task (lazy loads if not already loaded)
	t, err := db.Get("T1234-example")
	if err != nil {
		log.Fatal(err)
	}

	// Modify the task
	t.SetTitle("Updated Title")

	// Save dirty tasks
	count, err := db.SaveDirty()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Saved %d tasks\n", count)
}

// Example demonstrating parent-child relationships
func ExampleTaskDB_parentChild() {
	db := task.NewTaskDB("tasks")
	db.LoadAll()

	// Set parent-child relationship
	// This automatically validates that both tasks exist
	if err := db.SetParent("C1234-child", "P5678-parent"); err != nil {
		log.Fatal(err)
	}

	// Get all children of a parent
	children := db.GetChildren("P5678-parent")
	for _, child := range children {
		fmt.Println(child.ID)
	}

	// Clear parent relationship
	if err := db.ClearParent("C1234-child"); err != nil {
		log.Fatal(err)
	}

	db.SaveDirty()
}

// Example demonstrating blocker relationships
func ExampleTaskDB_blockers() {
	db := task.NewTaskDB("tasks")
	db.LoadAll()

	// Add a blocker relationship
	// After this, T1 will have T2 in its blockers, and T2 will have T1 in its blocks
	if err := db.AddBlocker("T1111-blocked", "T2222-blocker"); err != nil {
		log.Fatal(err)
	}

	// Remove a blocker relationship (maintains bidirectional consistency)
	if err := db.RemoveBlocker("T1111-blocked", "T2222-blocker"); err != nil {
		log.Fatal(err)
	}

	// Sync blockers from children
	// This ensures all parents are blocked by their incomplete children
	count, err := db.SyncBlockersFromChildren()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %d tasks\n", count)

	db.SaveDirty()
}

// Example demonstrating task completion
func ExampleTaskDB_completion() {
	db := task.NewTaskDB("tasks")
	db.LoadAll()

	// Mark a task as completed
	if err := db.SetCompleted("T1234-task", true); err != nil {
		log.Fatal(err)
	}

	// Update blocker relationships after completion
	// This removes the completed task from the blockers of tasks it was blocking
	if err := db.UpdateBlockersAfterCompletion("T1234-task"); err != nil {
		log.Fatal(err)
	}

	// You can also sync all blockers from children, which will:
	// - Remove completed children from parent blockers
	// - Add incomplete children to parent blockers
	db.SyncBlockersFromChildren()

	db.SaveDirty()
}

// Example demonstrating validation and repair
func ExampleTaskDB_validation() {
	db := task.NewTaskDB("tasks")
	db.LoadAll()

	// Validate all tasks
	errors := db.Validate()
	for _, err := range errors {
		fmt.Printf("Validation error: %v\n", err)
	}

	// Fix missing references (removes references to non-existent tasks)
	notices := db.FixMissingReferences()
	for _, notice := range notices {
		fmt.Printf("Fixed: %v\n", notice)
	}

	// Fix blocker relationships (ensures bidirectional consistency)
	count := db.FixBlockerRelationships()
	fmt.Printf("Fixed %d tasks\n", count)

	db.SaveDirty()
}
