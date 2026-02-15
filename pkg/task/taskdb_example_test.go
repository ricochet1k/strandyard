package task_test

import (
	"fmt"

	"github.com/ricochet1k/strandyard/pkg/task"
)

// ExampleTaskDB_workflow shows the recommended load -> mutate -> reconcile -> save
// flow for command handlers.
func ExampleTaskDB_workflow() {
	db := task.NewTaskDB("tasks")

	if err := db.LoadAllIfEmpty(); err != nil {
		return
	}

	// Resolve a user-facing short or full ID, then fetch the task.
	t, resolvedID, err := db.GetResolved("T48or")
	if err != nil {
		return
	}

	if err := db.MarkInProgress(resolvedID); err != nil {
		return
	}

	// Prefer TaskDB mutators over direct metadata writes.
	if err := db.AddTodo(resolvedID, "draft godoc for TaskDB mental model"); err != nil {
		return
	}

	// Reconcile after relationship/status edits when global consistency matters.
	if _, err := db.ReconcileBlockerRelationships(); err != nil {
		return
	}

	count, err := db.SaveDirty()
	if err != nil {
		return
	}
	fmt.Printf("updated %s, saved %d task(s)\n", t.ID, count)
}

// ExampleTaskDB_parentChild explains parent updates and the reconciliation step
// required for canonical blocker state.
func ExampleTaskDB_parentChild() {
	db := task.NewTaskDB("tasks")
	if err := db.LoadAll(); err != nil {
		return
	}

	// SetParent validates IDs and prevents cycles, but does not rebuild all
	// blocker edges by itself.
	if err := db.SetParent("C1234-child", "P5678-parent"); err != nil {
		return
	}

	if _, err := db.ReconcileBlockerRelationships(); err != nil {
		return
	}

	if err := db.ClearParent("C1234-child"); err != nil {
		return
	}

	if _, err := db.ReconcileBlockerRelationships(); err != nil {
		return
	}

	_, _ = db.SaveDirty()
}

// ExampleTaskDB_blockers shows explicit blocker mutations and automatic
// bidirectional updates.
func ExampleTaskDB_blockers() {
	db := task.NewTaskDB("tasks")
	if err := db.LoadAll(); err != nil {
		return
	}

	if err := db.AddBlocker("T1111-blocked", "T2222-blocker"); err != nil {
		return
	}

	if err := db.RemoveBlocker("T1111-blocked", "T2222-blocker"); err != nil {
		return
	}

	count, err := db.ReconcileBlockerRelationships()
	if err != nil {
		return
	}
	fmt.Printf("Updated %d tasks\n", count)

	_, _ = db.SaveDirty()
}

// ExampleTaskDB_statusLifecycle demonstrates status transitions with a report.
func ExampleTaskDB_statusLifecycle() {
	db := task.NewTaskDB("tasks")
	if err := db.LoadAll(); err != nil {
		return
	}

	if err := db.SetStatusWithReport("T1234-task", task.StatusDone, "implemented and tested"); err != nil {
		return
	}

	if err := db.UpdateBlockersAfterCompletion("T1234-task"); err != nil {
		return
	}

	_, _ = db.SaveDirty()
}

// ExampleTaskDB_notAllowed documents common API misuse and safer alternatives.
func ExampleTaskDB_notAllowed() {
	db := task.NewTaskDB("tasks")
	if err := db.LoadAll(); err != nil {
		return
	}

	t, err := db.Get("T1234-task")
	if err != nil {
		return
	}

	// Avoid direct relationship edits like: t.Meta.Parent = "other".
	// They bypass invariants and may leave blockers/blocks inconsistent.

	// Prefer mutators that enforce invariants.
	_ = db.SetParent(t.ID, "P5678-parent")
	_, _ = db.ReconcileBlockerRelationships()
	_, _ = db.SaveDirty()
}

// ExampleTaskDB_validation demonstrates explicit validation and missing-reference
// cleanup before save.
func ExampleTaskDB_validation() {
	db := task.NewTaskDB("tasks")
	if err := db.LoadAll(); err != nil {
		return
	}

	errors := db.Validate()
	for _, err := range errors {
		fmt.Printf("Validation error: %v\n", err)
	}

	notices := db.FixMissingReferences()
	for _, notice := range notices {
		fmt.Printf("Fixed: %v\n", notice)
	}

	count, _ := db.ReconcileBlockerRelationships()
	fmt.Printf("Fixed %d tasks\n", count)

	_, _ = db.SaveDirty()
}
