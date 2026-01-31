package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateTaskLinks(t *testing.T) {
	// Create a parser for building test tasks
	parser := NewParser()

	// Create mock role files for testing
	tmpDir := t.TempDir()
	roleDir := filepath.Join(tmpDir, "roles")
	if err := os.MkdirAll(roleDir, 0755); err != nil {
		t.Fatalf("failed to create role dir: %v", err)
	}

	roleName := testRoleName(t, "validate")
	roleContent := fmt.Sprintf("# %s\n\nRole for validation.", strings.Title(roleName))
	if err := os.WriteFile(filepath.Join(roleDir, roleName+".md"), []byte(roleContent), 0o644); err != nil {
		t.Fatalf("failed to create role file: %v", err)
	}

	// Change working directory temporarily for role validation
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer os.Chdir(origWd)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}

	tasks := make(map[string]*Task)

	t1Content := fmt.Sprintf(`---
role: %s
---

# Exists Task

This task exists.

See [T2bbb-missing](tasks/T2bbb-missing/T2bbb-missing.md) for details.
Also check [T1aaa-exists](tasks/T1aaa-exists/T1aaa-exists.md) for self-reference.
`, roleName)
	t1, _ := parser.ParseString(t1Content, "T1aaa-exists")
	t1.FilePath = "tasks/T1aaa-exists/T1aaa-exists.md"
	tasks[t1.ID] = t1

	t3Content := fmt.Sprintf(`---
role: %s
---

# No Links Task

This task has no links to other tasks.
`, roleName)
	t3, _ := parser.ParseString(t3Content, "T3ccc-no-links")
	t3.FilePath = "tasks/T3ccc-no-links/T3ccc-no-links.md"
	tasks[t3.ID] = t3

	v := NewValidator(tasks)
	errors := v.ValidateAndRepair()

	// Should have 1 error for the missing T2bbb-missing task
	if len(errors) != 1 {
		t.Fatalf("expected 1 validation error, got %d", len(errors))
	}

	expectedError := ValidationError{
		TaskID:  "T1aaa-exists",
		File:    "tasks/T1aaa-exists/T1aaa-exists.md",
		Message: "broken link: task T2bbb-missing does not exist",
	}

	if errors[0].TaskID != expectedError.TaskID {
		t.Errorf("expected task ID %q, got %q", expectedError.TaskID, errors[0].TaskID)
	}
	if errors[0].File != expectedError.File {
		t.Errorf("expected file %q, got %q", expectedError.File, errors[0].File)
	}
	if errors[0].Message != expectedError.Message {
		t.Errorf("expected message %q, got %q", expectedError.Message, errors[0].Message)
	}
}

func TestExtractTaskIDFromPath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"tasks/T3k7x-example/T3k7x-example.md", "T3k7x-example"},
		{"../T5h7w-task/task.md", "T5h7w-task"},
		{"T2k9p-other", "T2k9p-other"},
		{"tasks/E4m8x-epic/E4m8x-epic.md", "E4m8x-epic"},
		{"normal-file.md", ""}, // no task ID pattern
		{"", ""},               // empty path
		{"../tasks/Z9999-nonexistent/file.md", "Z9999-nonexistent"},
	}

	for _, test := range tests {
		result := extractTaskIDFromPath(test.path)
		if result != test.expected {
			t.Errorf("extractTaskIDFromPath(%q) = %q, expected %q", test.path, result, test.expected)
		}
	}
}

func TestGenerateMasterLists_FreeTasksPrioritySections(t *testing.T) {
	parser := NewParser()
	tmp := t.TempDir()
	rootsFile := filepath.Join(tmp, "root-tasks.md")
	freeFile := filepath.Join(tmp, "free-tasks.md")

	tasks := make(map[string]*Task)

	t1, _ := parser.ParseString("# High Task\n", "T1aaa-high")
	t1.Meta.Priority = PriorityHigh
	t1.FilePath = "tasks/T1aaa-high/T1aaa-high.md"
	tasks[t1.ID] = t1

	t2, _ := parser.ParseString("# Default Task\n", "T2bbb-default")
	t2.FilePath = "tasks/T2bbb-default/T2bbb-default.md"
	tasks[t2.ID] = t2

	t3, _ := parser.ParseString("# Low Task\n", "T3ccc-low")
	t3.Meta.Priority = PriorityLow
	t3.FilePath = "tasks/T3ccc-low/T3ccc-low.md"
	tasks[t3.ID] = t3

	if err := GenerateMasterLists(tasks, "tasks", rootsFile, freeFile); err != nil {
		t.Fatalf("GenerateMasterLists failed: %v", err)
	}

	got, err := os.ReadFile(freeFile)
	if err != nil {
		t.Fatalf("read free list: %v", err)
	}

	want := strings.Join([]string{
		"# Free tasks",
		"",
		"## High",
		"",
		"- [High Task](tasks/T1aaa-high/T1aaa-high.md)",
		"",
		"## Medium",
		"",
		"- [Default Task](tasks/T2bbb-default/T2bbb-default.md)",
		"",
		"## Low",
		"",
		"- [Low Task](tasks/T3ccc-low/T3ccc-low.md)",
		"",
		"",
	}, "\n")

	if string(got) != want {
		t.Fatalf("unexpected free list:\n--- got ---\n%s\n--- want ---\n%s", string(got), want)
	}
}

func TestCalculateIncrementalFreeListUpdate(t *testing.T) {
	tasks := map[string]*Task{
		"T1completed": {
			ID: "T1completed",
			Meta: Metadata{
				Priority:  "high",
				Blockers:  []string{},
				Completed: false,
			},
			FilePath: "tasks/T1completed.md",
		},
		"T2blocked": {
			ID: "T2blocked",
			Meta: Metadata{
				Priority:  "medium",
				Blockers:  []string{"T1completed"},
				Completed: false,
			},
			FilePath: "tasks/T2blocked.md",
		},
		"T3blocked": {
			ID: "T3blocked",
			Meta: Metadata{
				Priority:  "low",
				Blockers:  []string{"T1completed", "T4other"},
				Completed: false,
			},
			FilePath: "tasks/T3blocked.md",
		},
		"T4other": {
			ID: "T4other",
			Meta: Metadata{
				Priority:  "medium",
				Blockers:  []string{},
				Completed: true,
			},
			FilePath: "tasks/T4other.md",
		},
	}

	update, err := CalculateIncrementalFreeListUpdate(tasks, "T1completed")
	if err != nil {
		t.Fatalf("CalculateIncrementalFreeListUpdate failed: %v", err)
	}

	// Check that T1completed is removed
	if !containsString(update.RemoveTaskIDs, "T1completed") {
		t.Errorf("Expected T1completed to be in RemoveTaskIDs")
	}

	// Check that T2blocked is added (now unblocked)
	if len(update.AddTasks) != 2 {
		t.Fatalf("Expected 2 tasks to be added, got %d", len(update.AddTasks))
	}

	addedIDs := make(map[string]bool)
	for _, task := range update.AddTasks {
		addedIDs[task.ID] = true
	}

	if !addedIDs["T2blocked"] {
		t.Errorf("Expected T2blocked to be added")
	}
	if !addedIDs["T3blocked"] {
		t.Errorf("Expected T3blocked to be added (T4other is already completed)")
	}
}

func TestUpdateFreeListIncrementally(t *testing.T) {
	tempDir := t.TempDir()

	// Create an initial free-tasks.md file
	freeFile := filepath.Join(tempDir, "free.md")
	initialContent := `# Free tasks

## High

- [T1 Title](tasks/T1.md)

## Medium

- [T2 Title](tasks/T2.md)

## Low

- [T3 Title](tasks/T3.md)
`

	if err := os.WriteFile(freeFile, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to write initial free tasks file: %v", err)
	}

	// Create test tasks
	tasks := map[string]*Task{
		"T1": {
			ID: "T1",
			Meta: Metadata{
				Priority:  "high",
				Blockers:  []string{},
				Completed: false,
			},
			FilePath:     "tasks/T1.md",
			TitleContent: "T1 Title",
		},
		"T2": {
			ID: "T2",
			Meta: Metadata{
				Priority:  "medium",
				Blockers:  []string{},
				Completed: false,
			},
			FilePath:     "tasks/T2.md",
			TitleContent: "T2 Title",
		},
		"T4": {
			ID: "T4",
			Meta: Metadata{
				Priority:  "low",
				Blockers:  []string{},
				Completed: false,
			},
			FilePath:     "tasks/T4.md",
			TitleContent: "T4 Title",
		},
	}

	// Create update: remove T1, add T4
	update := IncrementalFreeListUpdate{
		RemoveTaskIDs: []string{"T1"},
		AddTasks:      []*Task{tasks["T4"]},
	}

	if err := UpdateFreeListIncrementally(tasks, freeFile, update); err != nil {
		t.Fatalf("UpdateFreeListIncrementally failed: %v", err)
	}

	// Read updated content
	updatedContent, err := os.ReadFile(freeFile)
	if err != nil {
		t.Fatalf("Failed to read updated free tasks file: %v", err)
	}

	content := string(updatedContent)

	// Check that T1 is removed
	if strings.Contains(content, "T1 Title") {
		t.Errorf("T1 Title should have been removed")
	}

	// Check that T4 is added in Low section
	lowSection := extractSection(content, "Low")
	if !strings.Contains(lowSection, "T4 Title") {
		t.Errorf("T4 Title not found in Low section")
	}

	// Check that T2 is still there
	mediumSection := extractSection(content, "Medium")
	if !strings.Contains(mediumSection, "T2 Title") {
		t.Errorf("T2 Title not found in Medium section")
	}
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func extractSection(content, sectionName string) string {
	start := strings.Index(content, "## "+sectionName)
	if start == -1 {
		return ""
	}
	start += len("## " + sectionName)

	end := strings.Index(content[start:], "## ")
	if end == -1 {
		return content[start:]
	}
	return content[start : start+end]
}
