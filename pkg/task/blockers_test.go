package task

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateBlockersFromChildren(t *testing.T) {
	tmp := t.TempDir()
	tasksRoot := filepath.Join(tmp, "tasks")

	parentID := "E1aaa-parent"
	childID := "T1bbb-child"
	parentDir := filepath.Join(tasksRoot, parentID)
	childDir := filepath.Join(parentDir, childID)
	if err := os.MkdirAll(childDir, 0o755); err != nil {
		t.Fatalf("mkdir child dir: %v", err)
	}

	parentFile := filepath.Join(parentDir, parentID+".md")
	childFile := filepath.Join(childDir, childID+".md")

	writeTask(t, parentFile, "architect", "high", "", []string{"Z9zzz-other"}, false)
	writeTask(t, childFile, "developer", "medium", parentID, nil, false)

	parser := NewParser()
	tasks, err := parser.LoadTasks(tasksRoot)
	if err != nil {
		t.Fatalf("load tasks: %v", err)
	}

	if _, err := UpdateBlockersFromChildren(tasks); err != nil {
		t.Fatalf("update blockers: %v", err)
	}

	updated, err := parser.ParseFile(parentFile)
	if err != nil {
		t.Fatalf("parse parent: %v", err)
	}

	if len(updated.Meta.Blockers) != 2 {
		t.Fatalf("expected 2 blockers, got %v", updated.Meta.Blockers)
	}
	if updated.Meta.Blockers[0] != "T1bbb-child" && updated.Meta.Blockers[1] != "T1bbb-child" {
		t.Fatalf("expected child blocker, got %v", updated.Meta.Blockers)
	}

	// Mark child completed and ensure blocker is removed.
	writeTask(t, childFile, "developer", "medium", parentID, nil, true)
	tasks, err = parser.LoadTasks(tasksRoot)
	if err != nil {
		t.Fatalf("reload tasks: %v", err)
	}
	if _, err := UpdateBlockersFromChildren(tasks); err != nil {
		t.Fatalf("update blockers after completion: %v", err)
	}
	updated, err = parser.ParseFile(parentFile)
	if err != nil {
		t.Fatalf("parse parent after completion: %v", err)
	}
	for _, b := range updated.Meta.Blockers {
		if b == childID {
			t.Fatalf("did not expect child blocker after completion, got %v", updated.Meta.Blockers)
		}
	}
}

func writeTask(t *testing.T, path, role, priority, parent string, blockers []string, completed bool) {
	t.Helper()
	content := "---\n"
	content += "role: " + role + "\n"
	content += "priority: " + priority + "\n"
	content += "parent: " + parent + "\n"
	content += "blockers: ["
	for i, b := range blockers {
		if i > 0 {
			content += ", "
		}
		content += b
	}
	content += "]\n"
	content += "blocks: []\n"
	content += "date_created: 2026-01-27T00:00:00Z\n"
	content += "date_edited: 2026-01-27T00:00:00Z\n"
	content += "owner_approval: false\n"
	if completed {
		content += "completed: true\n"
	} else {
		content += "completed: false\n"
	}
	content += "---\n\n# Task\n"

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write task %s: %v", path, err)
	}
}
