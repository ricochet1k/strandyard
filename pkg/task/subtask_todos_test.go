package task

import (
	"strings"
	"testing"
	"time"
)

func TestUpdateParentTodoEntriesPreservesManualItems(t *testing.T) {
	parser := NewParser()
	parentContent := `---
role: developer
parent: ""
---

# Parent Task

## TODOs

1. [ ] Manual item
2. [x] (subtask: T2bbb-old) Old subtask entry

## Acceptance Criteria
- done
`

	parent, _ := parser.ParseString(parentContent, "P1")
	sub1, _ := parser.ParseString("# First\n", "T1aaa")
	sub1.Meta.Parent = "P1"
	sub2, _ := parser.ParseString("# Second\n", "T2bbb")
	sub2.Meta.Parent = "P1"
	sub2.Meta.Completed = true

	tasks := map[string]*Task{
		"P1":    parent,
		"T1aaa": sub1,
		"T2bbb": sub2,
	}

	changed, err := UpdateParentTodoEntries(tasks, "P1")
	if err != nil {
		t.Fatalf("UpdateParentTodoEntries error: %v", err)
	}
	if !changed {
		t.Fatalf("expected UpdateParentTodoEntries to report changes")
	}

	content := parent.Content()
	if !strings.Contains(content, "- [ ] Manual item") {
		t.Errorf("manual item missing:\n%s", content)
	}
	if !strings.Contains(content, "- [ ] (subtask: T1aaa) First") {
		t.Errorf("sub1 missing:\n%s", content)
	}
	if !strings.Contains(content, "- [x] (subtask: T2bbb) Second") {
		t.Errorf("sub2 missing:\n%s", content)
	}
}

func TestUpdateParentTodoEntriesInsertsSubtasksSection(t *testing.T) {
	parser := NewParser()
	parentContent := `---
role: developer
parent: ""
---

# Parent Task

## Summary
Details.

## Acceptance Criteria
- done
`

	parent, _ := parser.ParseString(parentContent, "P1")
	sub, _ := parser.ParseString("# First\n", "T1aaa")
	sub.Meta.Parent = "P1"

	tasks := map[string]*Task{
		"P1":    parent,
		"T1aaa": sub,
	}

	changed, err := UpdateParentTodoEntries(tasks, "P1")
	if err != nil {
		t.Fatalf("UpdateParentTodoEntries error: %v", err)
	}
	if !changed {
		t.Fatalf("expected UpdateParentTodoEntries to report changes")
	}

	content := parent.Content()
	if !strings.Contains(content, "## Subtasks") {
		t.Errorf("Subtasks section missing:\n%s", content)
	}
	if !strings.Contains(content, "- [ ] (subtask: T1aaa) First") {
		t.Errorf("subtask entry missing:\n%s", content)
	}
}

func TestUpdateParentTodoEntriesUsesCreationOrder(t *testing.T) {
	parser := NewParser()
	parent, _ := parser.ParseString("# Parent\n", "P1")

	late, _ := parser.ParseString("# Late\n", "T2bbb")
	late.Meta.Parent = "P1"
	late.Meta.DateCreated = time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)

	early, _ := parser.ParseString("# Early\n", "T1aaa")
	early.Meta.Parent = "P1"
	early.Meta.DateCreated = time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC)

	tasks := map[string]*Task{
		"P1":    parent,
		"T2bbb": late,
		"T1aaa": early,
	}

	if _, err := UpdateParentTodoEntries(tasks, "P1"); err != nil {
		t.Fatalf("UpdateParentTodoEntries error: %v", err)
	}

	if len(parent.SubsItems) != 2 {
		t.Fatalf("expected 2 subtasks, got %d", len(parent.SubsItems))
	}
	if parent.SubsItems[0].Text != "Early" || parent.SubsItems[1].Text != "Late" {
		t.Fatalf("unexpected subtask order: %#v", parent.SubsItems)
	}
}

func TestUpdateParentTodoEntriesPreservesExistingSubtaskOrder(t *testing.T) {
	parser := NewParser()
	parentContent := `# Parent

## Subtasks
- [ ] (subtask: T2bbb) Second
- [ ] (subtask: T1aaa) First
`
	parent, _ := parser.ParseString(parentContent, "P1")

	first, _ := parser.ParseString("# First\n", "T1aaa")
	first.Meta.Parent = "P1"
	first.Meta.DateCreated = time.Date(2026, 2, 1, 9, 0, 0, 0, time.UTC)

	second, _ := parser.ParseString("# Second\n", "T2bbb")
	second.Meta.Parent = "P1"
	second.Meta.DateCreated = time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)

	tasks := map[string]*Task{
		"P1":    parent,
		"T1aaa": first,
		"T2bbb": second,
	}

	if _, err := UpdateParentTodoEntries(tasks, "P1"); err != nil {
		t.Fatalf("UpdateParentTodoEntries error: %v", err)
	}

	if len(parent.SubsItems) != 2 {
		t.Fatalf("expected 2 subtasks, got %d", len(parent.SubsItems))
	}
	if parent.SubsItems[0].SubtaskID != "T2bbb" || parent.SubsItems[1].SubtaskID != "T1aaa" {
		t.Fatalf("expected existing order to be preserved, got: %#v", parent.SubsItems)
	}
}
