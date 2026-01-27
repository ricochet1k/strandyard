package task

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFileWithKind(t *testing.T) {
	tmp := t.TempDir()
	taskID := "I1abc-example-issue"
	taskDir := filepath.Join(tmp, taskID)
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatalf("mkdir task dir: %v", err)
	}

	taskFile := filepath.Join(taskDir, taskID+".md")
	content := `---
kind: issue
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

	if parsed.Meta.Kind != "issue" {
		t.Fatalf("expected kind=issue, got %q", parsed.Meta.Kind)
	}
}
