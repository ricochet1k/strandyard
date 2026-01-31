package task

import (
	"fmt"
	"testing"
)

func TestParseTaskItems(t *testing.T) {
	primaryRole := testRoleName(t, "primary")
	secondaryRole := testRoleName(t, "secondary")
	content := fmt.Sprintf(`1. [ ] (role: %s) Implement the behavior.
2. [x] (role: %s) Run tests.
3. [ ] Document it.
- [ ] (subtask: T1aaa) Subtask 1
`, primaryRole, secondaryRole)

	items := ParseTaskItems(content)
	if len(items) != 4 {
		t.Fatalf("expected 4 items, got %d", len(items))
	}
	if items[0].Checked || items[0].Role != primaryRole || items[0].Text != "Implement the behavior." {
		t.Fatalf("unexpected first item: %+v", items[0])
	}
	if !items[1].Checked || items[1].Role != secondaryRole || items[1].Text != "Run tests." {
		t.Fatalf("unexpected second item: %+v", items[1])
	}
	if items[2].Checked || items[2].Role != "" || items[2].Text != "Document it." {
		t.Fatalf("unexpected third item: %+v", items[2])
	}
	if items[3].SubtaskID != "T1aaa" || items[3].Text != "Subtask 1" {
		t.Fatalf("unexpected fourth item: %+v", items[3])
	}
}

func TestFormatTaskItems(t *testing.T) {
	items := []TaskItem{
		{Checked: false, Role: "dev", Text: "Task 1"},
		{Checked: true, SubtaskID: "T1", Text: "Sub 1"},
	}

	todoStr := FormatTodoItems(items)
	expectedTodo := "1. [ ] (role: dev) Task 1\n2. [x] (subtask: T1) Sub 1"
	if todoStr != expectedTodo {
		t.Errorf("FormatTodoItems mismatch:\nGot: %q\nWant: %q", todoStr, expectedTodo)
	}

	subStr := FormatSubtaskItems(items)
	expectedSub := "- [ ] (role: dev) Task 1\n- [x] (subtask: T1) Sub 1"
	if subStr != expectedSub {
		t.Errorf("FormatSubtaskItems mismatch:\nGot: %q\nWant: %q", subStr, expectedSub)
	}
}
