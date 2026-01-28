package task

import (
	"github.com/yuin/goldmark/text"
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

	if err := os.WriteFile(filepath.Join(roleDir, "developer.md"), []byte("# Developer\n\nDeveloper role."), 0644); err != nil {
		t.Fatalf("failed to create developer role file: %v", err)
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

	tasks := map[string]*Task{
		"T1aaa-exists": {
			ID:       "T1aaa-exists",
			FilePath: "tasks/T1aaa-exists/T1aaa-exists.md",
			Content: `# Exists Task

This task exists.

See [T2bbb-missing](tasks/T2bbb-missing/T2bbb-missing.md) for details.
Also check [T1aaa-exists](tasks/T1aaa-exists/T1aaa-exists.md) for self-reference.
`,
			Meta: Metadata{
				Role: "developer",
			},
		},
		"T3ccc-no-links": {
			ID:       "T3ccc-no-links",
			FilePath: "tasks/T3ccc-no-links/T3ccc-no-links.md",
			Content: `# No Links Task

This task has no links to other tasks.
`,
			Meta: Metadata{
				Role: "developer",
			},
		},
	}

	// Parse the content to build ASTs
	for id, task := range tasks {
		reader := text.NewReader([]byte(task.Content))
		doc := parser.md.Parser().Parse(reader)
		task.Document = doc
		tasks[id] = task
	}

	v := NewValidator(tasks)
	errors := v.Validate()

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
	tmp := t.TempDir()
	rootsFile := filepath.Join(tmp, "root-tasks.md")
	freeFile := filepath.Join(tmp, "free-tasks.md")

	tasks := map[string]*Task{
		"T1aaa-high": {
			ID:       "T1aaa-high",
			FilePath: "tasks/T1aaa-high/T1aaa-high.md",
			Content:  "# High Task\n",
			Meta: Metadata{
				Priority: PriorityHigh,
			},
		},
		"T2bbb-default": {
			ID:       "T2bbb-default",
			FilePath: "tasks/T2bbb-default/T2bbb-default.md",
			Content:  "# Default Task\n",
			Meta:     Metadata{},
		},
		"T3ccc-low": {
			ID:       "T3ccc-low",
			FilePath: "tasks/T3ccc-low/T3ccc-low.md",
			Content:  "# Low Task\n",
			Meta: Metadata{
				Priority: PriorityLow,
			},
		},
	}

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
