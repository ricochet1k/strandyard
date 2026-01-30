package task

import (
	"fmt"
	"strings"
	"testing"
)

func TestUpdateParentTodoEntriesPreservesManualItems(t *testing.T) {
	role := testRoleName(t, "parent")
	parentContent := fmt.Sprintf(`---
role: %s
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T00:00:00Z
owner_approval: false
completed: false
---

# Parent Task

## Tasks

- [ ] Manual item
- [x] (subtask: T2bbb-old) Old subtask entry

## Acceptance Criteria
- done
`, role)

	parent := &Task{
		ID:      "P1",
		Content: parentContent,
	}
	sub1 := &Task{ID: "T1aaa-first", Meta: Metadata{Parent: "P1", Completed: false}, Content: "# First\n"}
	sub2 := &Task{ID: "T2bbb-second", Meta: Metadata{Parent: "P1", Completed: true}, Content: "# Second\n"}

	tasks := map[string]*Task{
		"P1":           parent,
		"T1aaa-first":  sub1,
		"T2bbb-second": sub2,
	}

	changed, err := UpdateParentTodoEntries(tasks, "P1")
	if err != nil {
		t.Fatalf("UpdateParentTodoEntries error: %v", err)
	}
	if !changed {
		t.Fatalf("expected UpdateParentTodoEntries to report changes")
	}

	expected := "## Tasks\n\n- [ ] Manual item\n\n- [ ] (subtask: T1aaa-first) First\n- [x] (subtask: T2bbb-second) Second\n\n## Acceptance Criteria"
	if !strings.Contains(parent.Content, expected) {
		t.Fatalf("updated content missing expected tasks section\nExpected snippet:\n%s\nGot:\n%s", expected, parent.Content)
	}
}

func TestUpdateParentTodoEntriesInsertsTasksSection(t *testing.T) {
	role := testRoleName(t, "parent")
	parentContent := fmt.Sprintf(`---
role: %s
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T00:00:00Z
owner_approval: false
completed: false
---

# Parent Task

## Summary
Details.

## Acceptance Criteria
- done
`, role)

	parent := &Task{ID: "P1", Content: parentContent}
	sub := &Task{ID: "T1aaa-first", Meta: Metadata{Parent: "P1"}, Content: "# First\n"}

	tasks := map[string]*Task{
		"P1":          parent,
		"T1aaa-first": sub,
	}

	changed, err := UpdateParentTodoEntries(tasks, "P1")
	if err != nil {
		t.Fatalf("UpdateParentTodoEntries error: %v", err)
	}
	if !changed {
		t.Fatalf("expected UpdateParentTodoEntries to report changes")
	}

	tasksIndex := strings.Index(parent.Content, "## Tasks")
	acceptanceIndex := strings.Index(parent.Content, "## Acceptance Criteria")
	if tasksIndex == -1 || acceptanceIndex == -1 || tasksIndex > acceptanceIndex {
		t.Fatalf("expected ## Tasks section before Acceptance Criteria, got content:\n%s", parent.Content)
	}
	if !strings.Contains(parent.Content, "- [ ] (subtask: T1aaa-first) First") {
		t.Fatalf("expected subtask entry in inserted Tasks section\nGot:\n%s", parent.Content)
	}
}
