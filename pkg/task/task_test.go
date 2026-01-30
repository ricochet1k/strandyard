package task

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
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
