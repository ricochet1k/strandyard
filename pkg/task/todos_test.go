package task

import (
	"fmt"
	"testing"
)

func TestParseTodoItems(t *testing.T) {
	primaryRole := testRoleName(t, "primary")
	secondaryRole := testRoleName(t, "secondary")
	content := fmt.Sprintf(`---
role: %s
priority: medium
---

# Example

## TODOs
Check this off one at a time.
1. [ ] (role: %s) Implement the behavior.
2. [x] (role: %s) Run tests.
3. [ ] Document it.

## Acceptance Criteria
- Done
`, primaryRole, primaryRole, secondaryRole)

	items, err := ParseTodoItems(content, "task.md")
	if err != nil {
		t.Fatalf("ParseTodoItems error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[0].Index != 1 || items[0].Checked || items[0].Role != primaryRole || items[0].Text != "Implement the behavior." {
		t.Fatalf("unexpected first item: %+v", items[0])
	}
	if items[1].Index != 2 || !items[1].Checked || items[1].Role != secondaryRole || items[1].Text != "Run tests." {
		t.Fatalf("unexpected second item: %+v", items[1])
	}
	if items[2].Index != 3 || items[2].Checked || items[2].Role != "" || items[2].Text != "Document it." {
		t.Fatalf("unexpected third item: %+v", items[2])
	}
}

func TestParseSubtaskItems(t *testing.T) {
	role := testRoleName(t, "subtask")
	content := fmt.Sprintf(`---
role: %s
priority: medium
---

# Example

## Tasks

- [ ] Manual item

- [ ] (subtask: T1aaa-first) First
- [x] (subtask: T2bbb-second) Second

## Acceptance Criteria
- Done
`, role)

	items, err := ParseSubtaskItems(content, "task.md")
	if err != nil {
		t.Fatalf("ParseSubtaskItems error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].ID != "T1aaa-first" || items[0].Checked || items[0].Title != "First" {
		t.Fatalf("unexpected first subtask: %+v", items[0])
	}
	if items[1].ID != "T2bbb-second" || !items[1].Checked || items[1].Title != "Second" {
		t.Fatalf("unexpected second subtask: %+v", items[1])
	}
}
