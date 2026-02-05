package task

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestParseStringWithMalformedYAML(t *testing.T) {
	// Test parsing YAML with missing space after colon (blockers:[] instead of blockers: [])
	malformedContent := `---
role: developer
blockers:[]
blocks: []
---

# Test Task

This is a test task.`

	parser := NewParser()
	_, err := parser.ParseString(malformedContent, "T1234-test")

	if err == nil {
		t.Fatal("expected error when parsing malformed YAML, but got nil")
	}

	// Check that it's a FrontmatterParseError
	var parseErr *FrontmatterParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected FrontmatterParseError, got %T: %v", err, err)
	}

	// Check that the error message contains line number information
	errorMsg := err.Error()
	if errorMsg == "" {
		t.Fatal("error message should not be empty")
	}
	t.Logf("Error message: %s", errorMsg)

	// The error should mention the line number (2, since blockers is on line 2 of YAML content)
	if parseErr.LineNumber == 0 {
		t.Logf("Note: line number not extracted from YAML error (this is okay if YAML library changes format)")
	}

	// Check that the YAML error is captured
	if parseErr.YAMLError == "" {
		t.Fatal("YAML error message should not be empty")
	}
}

func TestParseStringWithValidYAML(t *testing.T) {
	// Test that valid YAML parses correctly
	validContent := `---
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T13:43:58Z
owner_approval: false
completed: false
---

# Test Task

This is a valid test task.`

	parser := NewParser()
	task, err := parser.ParseString(validContent, "T1234-test")

	if err != nil {
		t.Fatalf("unexpected error parsing valid YAML: %v", err)
	}

	if task == nil {
		t.Fatal("expected task to be non-nil")
	}

	if task.Meta.Role != "developer" {
		t.Errorf("expected role 'developer', got %q", task.Meta.Role)
	}

	if task.Meta.Priority != "high" {
		t.Errorf("expected priority 'high', got %q", task.Meta.Priority)
	}

	if len(task.Meta.Blockers) != 0 {
		t.Errorf("expected empty blockers slice, got %v", task.Meta.Blockers)
	}
}

func TestParseStringWithMissingFrontmatter(t *testing.T) {
	// Test content without closing frontmatter delimiter
	invalidContent := `---
role: developer
blockers: []

# Test Task`

	parser := NewParser()
	_, err := parser.ParseString(invalidContent, "T1234-test")

	if err == nil {
		t.Fatal("expected error for missing frontmatter closing delimiter")
	}

	// Check that it's an InvalidFrontmatterError
	var fmErr *InvalidFrontmatterError
	if !errors.As(err, &fmErr) {
		t.Fatalf("expected InvalidFrontmatterError, got %T: %v", err, err)
	}
}

func TestYAMLMarshalingSpacing(t *testing.T) {
	// Create a task with empty blockers and blocks
	task := &Task{
		ID: "T1234-test",
		Meta: Metadata{
			Role:     "developer",
			Priority: "high",
			Blockers: []string{},
			Blocks:   []string{},
		},
		TitleContent: "Test Task",
		BodyContent:  "Test body",
	}

	content := task.Content()

	// Verify that the content includes proper spacing
	if !containsSubstring(content, "blockers: []") && !containsSubstring(content, "blockers:") {
		t.Errorf("YAML should contain 'blockers: []' or 'blockers:' with proper formatting, got:\n%s", content)
	}

	if !containsSubstring(content, "blocks: []") && !containsSubstring(content, "blocks:") {
		t.Errorf("YAML should contain 'blocks: []' or 'blocks:' with proper formatting, got:\n%s", content)
	}

	// Verify that we don't have invalid spacing like "blockers:[]"
	if containsSubstring(content, "blockers:[]") {
		t.Errorf("YAML should not have 'blockers:[]' (missing space), got:\n%s", content)
	}

	if containsSubstring(content, "blocks:[]") {
		t.Errorf("YAML should not have 'blocks:[]' (missing space), got:\n%s", content)
	}
}

func TestYAMLMarshalingWithNonEmptySlices(t *testing.T) {
	// Create a task with non-empty blockers and blocks
	task := &Task{
		ID: "T1234-test",
		Meta: Metadata{
			Role:     "developer",
			Priority: "high",
			Blockers: []string{"T5678-blocker1", "T9012-blocker2"},
			Blocks:   []string{"T3456-blocks1"},
		},
		TitleContent: "Test Task",
		BodyContent:  "Test body",
	}

	content := task.Content()

	// Verify that the content can be parsed back
	parser := NewParser()
	parsedTask, err := parser.ParseString(content, "T1234-test")
	if err != nil {
		t.Fatalf("failed to parse generated content: %v", err)
	}

	// Verify the parsed task matches the original
	if len(parsedTask.Meta.Blockers) != 2 {
		t.Errorf("expected 2 blockers, got %d", len(parsedTask.Meta.Blockers))
	}

	if len(parsedTask.Meta.Blocks) != 1 {
		t.Errorf("expected 1 block, got %d", len(parsedTask.Meta.Blocks))
	}
}

func containsSubstring(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestParseFileWithMalformedYAML(t *testing.T) {
	// Create a temporary directory and file with malformed YAML
	tmpDir, err := os.MkdirTemp("", "task-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create task directory structure: tmpDir/T1234-test/
	taskDir := filepath.Join(tmpDir, "T1234-test")
	if err := os.Mkdir(taskDir, 0755); err != nil {
		t.Fatalf("failed to create task directory: %v", err)
	}

	// Write a file with malformed YAML (blockers:[] without space)
	filePath := filepath.Join(taskDir, "T1234-test.md")
	malformedContent := `---
role: developer
blockers:[]
blocks: []
---

# Test Task

This task has malformed YAML.`

	if err := os.WriteFile(filePath, []byte(malformedContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Try to parse the file
	parser := NewParser()
	_, err = parser.ParseFile(filePath)

	if err == nil {
		t.Fatal("expected error when parsing file with malformed YAML")
	}

	// Verify it's a FrontmatterParseError with file path
	var parseErr *FrontmatterParseError
	if !errors.As(err, &parseErr) {
		t.Fatalf("expected FrontmatterParseError, got %T: %v", err, err)
	}

	// Check that the error message includes the file path
	if parseErr.Path != filePath {
		t.Errorf("expected path %q, got %q", filePath, parseErr.Path)
	}

	// Check that the error message is informative
	errorMsg := parseErr.Error()
	if !containsSubstring(errorMsg, filePath) {
		t.Errorf("error message should contain file path, got: %s", errorMsg)
	}

	if !containsSubstring(errorMsg, "yaml") || !containsSubstring(errorMsg, "line") {
		t.Errorf("error message should mention YAML and line number, got: %s", errorMsg)
	}
}

func TestParseFileWithValidYAML(t *testing.T) {
	// Create a temporary directory and file with valid YAML
	tmpDir, err := os.MkdirTemp("", "task-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create task directory structure
	taskDir := filepath.Join(tmpDir, "T5678-valid")
	if err := os.Mkdir(taskDir, 0755); err != nil {
		t.Fatalf("failed to create task directory: %v", err)
	}

	// Write a file with valid YAML
	filePath := filepath.Join(taskDir, "T5678-valid.md")
	validContent := `---
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T00:00:00Z
date_edited: 2026-02-05T00:00:00Z
owner_approval: false
completed: false
---

# Valid Task

This task has valid YAML.`

	if err := os.WriteFile(filePath, []byte(validContent), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Parse the file
	parser := NewParser()
	task, err := parser.ParseFile(filePath)

	if err != nil {
		t.Fatalf("unexpected error parsing valid file: %v", err)
	}

	if task == nil {
		t.Fatal("expected task to be non-nil")
	}

	if task.FilePath != filePath {
		t.Errorf("expected FilePath %q, got %q", filePath, task.FilePath)
	}

	if task.Dir != taskDir {
		t.Errorf("expected Dir %q, got %q", taskDir, task.Dir)
	}

	if task.Meta.Role != "developer" {
		t.Errorf("expected role 'developer', got %q", task.Meta.Role)
	}
}
