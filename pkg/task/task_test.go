package task

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFileWithType(t *testing.T) {
	tmp := t.TempDir()
	taskID := "I1abc-example-issue"
	taskDir := filepath.Join(tmp, taskID)
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatalf("mkdir task dir: %v", err)
	}

	taskFile := filepath.Join(taskDir, taskID+".md")
	content := `---
type: issue
role: developer
priority: medium
parent:
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T00:00:00Z
owner_approval: false
completed: false
---

# Issue Title
`
	if err := os.WriteFile(taskFile, []byte(content), 0o644); err != nil {
		t.Fatalf("write task file: %v", err)
	}

	parser := NewParser()
	parsed, err := parser.ParseFile(taskFile)
	if err != nil {
		t.Fatalf("parse file: %v", err)
	}

	if parsed.Meta.Type != "issue" {
		t.Fatalf("expected type=issue, got %q", parsed.Meta.Type)
	}
}
