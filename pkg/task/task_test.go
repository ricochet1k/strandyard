package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseFileWithType(t *testing.T) {
	tmp := t.TempDir()
	taskID := "T1abc-example"
	taskDir := filepath.Join(tmp, taskID)
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatalf("mkdir task dir: %v", err)
	}

	taskFile := filepath.Join(taskDir, taskID+".md")
	roleName := testRoleName(t, "role")
	typeName := testTypeName(t, "type")
	content := fmt.Sprintf(`---
type: %s
role: %s
priority: medium
parent:
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T00:00:00Z
owner_approval: false
completed: false
---

# Example Title
`, typeName, roleName)
	if err := os.WriteFile(taskFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write task file: %v", err)
	}

	parser := NewParser()
	parsed, err := parser.ParseFile(taskFile)
	if err != nil {
		t.Fatalf("parse file: %v", err)
	}

	if parsed.Meta.Type != typeName {
		t.Fatalf("expected type=%s, got %q", typeName, parsed.Meta.Type)
	}
}

func TestTaskSetTitleMarksDirtyAndUpdatesDateEdited(t *testing.T) {
	prevEdited := time.Date(2026, time.February, 1, 12, 0, 0, 0, time.UTC)
	task := Task{Meta: Metadata{DateEdited: prevEdited}}

	task.SetTitle("New Title")
	if task.TitleContent != "New Title" {
		t.Fatalf("expected title to be updated, got %q", task.TitleContent)
	}
	if !task.Dirty {
		t.Fatalf("expected task to be marked dirty")
	}
	if !task.Meta.DateEdited.After(prevEdited) {
		t.Fatalf("expected DateEdited to be updated, got %s", task.Meta.DateEdited.Format(time.RFC3339Nano))
	}

	updated := task.Meta.DateEdited
	task.SetTitle("New Title")
	if !task.Meta.DateEdited.Equal(updated) {
		t.Fatalf("expected DateEdited to remain unchanged on no-op title update")
	}
}

func TestTaskSetBodyStripsReservedSections(t *testing.T) {
	task := Task{}
	body := strings.TrimSpace(`
# Ignored Title
Intro text.

## TODOs
- [ ] One

## Progress
Did something.

## Subtasks
- [ ] Sub one

## Notes
Keep this.
`)

	task.SetBody(body)
	if !task.Dirty {
		t.Fatalf("expected task to be marked dirty")
	}
	if task.BodyContent != "## Notes\nKeep this." {
		t.Fatalf("unexpected body content: %q", task.BodyContent)
	}
}

func TestTaskMarkDirtyUpdatesDateOnce(t *testing.T) {
	prevEdited := time.Date(2026, time.February, 2, 12, 0, 0, 0, time.UTC)
	task := Task{Meta: Metadata{DateEdited: prevEdited}}

	task.MarkDirty()
	if !task.Dirty {
		t.Fatalf("expected task to be marked dirty")
	}
	if !task.Meta.DateEdited.After(prevEdited) {
		t.Fatalf("expected DateEdited to be updated, got %s", task.Meta.DateEdited.Format(time.RFC3339Nano))
	}
	updated := task.Meta.DateEdited

	task.MarkDirty()
	if !task.Meta.DateEdited.Equal(updated) {
		t.Fatalf("expected DateEdited to remain unchanged on subsequent MarkDirty")
	}
}

func TestTaskContentIncludesSections(t *testing.T) {
	created := time.Date(2026, time.February, 3, 9, 0, 0, 0, time.UTC)
	edited := time.Date(2026, time.February, 3, 10, 0, 0, 0, time.UTC)
	task := Task{
		ID:           "T1abc-example",
		Meta:         Metadata{Role: "developer", Priority: "medium", DateCreated: created, DateEdited: edited},
		TitleContent: "Example Title",
		BodyContent:  "Body text.",
		TodoItems: []TaskItem{
			{Checked: false, Role: "developer", Text: "Do work"},
		},
		SubsItems: []TaskItem{
			{Checked: true, SubtaskID: "T2def-sub", Text: "Subtask"},
		},
		ProgressContent: "Progress update.",
		OtherContent:    "Footer block.",
	}

	content := task.Content()
	if !strings.Contains(content, "---\n") {
		t.Fatalf("expected frontmatter to be present")
	}
	if !strings.Contains(content, "# Example Title") {
		t.Fatalf("expected title to be present")
	}
	if !strings.Contains(content, "Body text.") {
		t.Fatalf("expected body to be present")
	}
	if !strings.Contains(content, "## TODOs") {
		t.Fatalf("expected TODOs section to be present")
	}
	if !strings.Contains(content, "## Subtasks") {
		t.Fatalf("expected Subtasks section to be present")
	}
	if !strings.Contains(content, "## Progress") {
		t.Fatalf("expected Progress section to be present")
	}
	if !strings.Contains(content, "Footer block.") {
		t.Fatalf("expected OtherContent to be present")
	}

	idxTitle := strings.Index(content, "# Example Title")
	idxTodo := strings.Index(content, "## TODOs")
	idxSubtasks := strings.Index(content, "## Subtasks")
	idxProgress := strings.Index(content, "## Progress")
	if idxTitle == -1 || idxTodo == -1 || idxSubtasks == -1 || idxProgress == -1 {
		t.Fatalf("expected all sections to appear in content")
	}
	if !(idxTitle < idxTodo && idxTodo < idxSubtasks && idxSubtasks < idxProgress) {
		t.Fatalf("expected sections in order: title, todos, subtasks, progress")
	}
}

func TestTaskGetEffectiveRolePrefersTodoRole(t *testing.T) {
	task := Task{
		Meta: Metadata{Role: "developer"},
		TodoItems: []TaskItem{
			{Checked: true, Role: "owner", Text: "Checked role"},
			{Checked: false, Role: "designer", Text: "Active role"},
		},
	}

	if got := task.GetEffectiveRole(); got != "designer" {
		t.Fatalf("expected effective role to be %q, got %q", "designer", got)
	}
}

func TestWriteDirtyTasksWritesOnlyDirty(t *testing.T) {
	tmp := t.TempDir()

	dirtyPath := filepath.Join(tmp, "dirty.md")
	cleanPath := filepath.Join(tmp, "clean.md")
	if err := os.WriteFile(cleanPath, []byte("keep\n"), 0o644); err != nil {
		t.Fatalf("write clean file: %v", err)
	}

	dirty := &Task{FilePath: dirtyPath}
	dirty.MarkDirty()

	clean := &Task{FilePath: cleanPath}

	updated, err := WriteDirtyTasks(map[string]*Task{
		"dirty": dirty,
		"clean": clean,
	})
	if err != nil {
		t.Fatalf("WriteDirtyTasks: %v", err)
	}
	if updated != 1 {
		t.Fatalf("expected 1 task updated, got %d", updated)
	}
	if _, err := os.Stat(dirtyPath); err != nil {
		t.Fatalf("expected dirty task to be written: %v", err)
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		t.Fatalf("read clean file: %v", err)
	}
	if string(content) != "keep\n" {
		t.Fatalf("expected clean file to remain unchanged, got %q", string(content))
	}
}

func TestWriteAllTasksWritesAllAndClearsDirty(t *testing.T) {
	tmp := t.TempDir()

	firstPath := filepath.Join(tmp, "first.md")
	secondPath := filepath.Join(tmp, "second.md")

	first := &Task{FilePath: firstPath}
	second := &Task{FilePath: secondPath}
	second.MarkDirty()

	updated, err := WriteAllTasks(map[string]*Task{
		"first":  first,
		"second": second,
	})
	if err != nil {
		t.Fatalf("WriteAllTasks: %v", err)
	}
	if updated != 2 {
		t.Fatalf("expected 2 tasks updated, got %d", updated)
	}
	if first.Dirty || second.Dirty {
		t.Fatalf("expected Dirty to be false after WriteAllTasks")
	}

	if _, err := os.Stat(firstPath); err != nil {
		t.Fatalf("expected first task to be written: %v", err)
	}
	if _, err := os.Stat(secondPath); err != nil {
		t.Fatalf("expected second task to be written: %v", err)
	}
}
