package task

import (
	"strings"
	"testing"
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
