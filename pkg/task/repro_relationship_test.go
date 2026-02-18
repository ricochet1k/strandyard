package task

import (
	"slices"
	"testing"
)

func TestReproSetParentDoesNotUpdateBlockers(t *testing.T) {
	tmp := t.TempDir()
	db := NewTaskDB(tmp)

	// Create parent
	parent, err := db.GetOrCreate("P1aaa-parent")
	if err != nil {
		t.Fatalf("create parent: %v", err)
	}
	parent.Meta.Role = "dev"
	db.Save(parent.ID)

	// Create child
	child, err := db.GetOrCreate("T1bbb-child")
	if err != nil {
		t.Fatalf("create child: %v", err)
	}
	child.Meta.Role = "dev"
	db.Save(child.ID)

	// Set parent
	if err := db.SetParent(child.ID, parent.ID); err != nil {
		t.Fatalf("SetParent: %v", err)
	}

	// Verify child has parent set
	if child.Meta.Parent != parent.ID {
		t.Errorf("expected parent %s, got %s", parent.ID, child.Meta.Parent)
	}

	// VERIFY ISSUE: Parent should be blocked by child, but SetParent currently doesn't do it
	if slices.Contains(parent.Meta.Blockers, child.ID) {
		t.Logf("Unexpectedly passed: parent is blocked by child")
	} else {
		t.Logf("Confirmed: parent is NOT blocked by child after SetParent")
	}

	// Also check child.Meta.Blocks
	if slices.Contains(child.Meta.Blocks, parent.ID) {
		t.Logf("Unexpectedly passed: child blocks parent")
	} else {
		t.Logf("Confirmed: child does NOT block parent after SetParent")
	}
}

func TestReproSetCompletedDoesNotUpdateBlockers(t *testing.T) {
	tmp := t.TempDir()
	db := NewTaskDB(tmp)

	// Create parent and child
	parent, _ := db.GetOrCreate("P1aaa-parent")
	child, _ := db.GetOrCreate("T1bbb-child")

	parent.Meta.Role = "dev"
	child.Meta.Role = "dev"

	// Manually set up blocking relationship (or use Reconcile)
	// We need to set Parent field first because Reconcile depends on it
	child.Meta.Parent = parent.ID

	// Explicitly setup blockers/blocks to simulate a reconciled state
	child.Meta.Blocks = []string{parent.ID}
	parent.Meta.Blockers = []string{child.ID}

	db.Save(parent.ID)
	db.Save(child.ID)

	// Verify initial state
	if !slices.Contains(parent.Meta.Blockers, child.ID) {
		t.Fatalf("setup failed: parent should be blocked")
	}

	// Set child as completed
	if err := db.SetCompleted(child.ID, true); err != nil {
		t.Fatalf("SetCompleted: %v", err)
	}

	// VERIFY ISSUE: Parent should NO LONGER be blocked by child
	if slices.Contains(parent.Meta.Blockers, child.ID) {
		t.Logf("Confirmed: parent is STILL blocked by child after SetCompleted")
	} else {
		t.Logf("Unexpectedly passed: parent is unblocked")
	}
}
