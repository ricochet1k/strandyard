package e2e

import (
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TaskOpts represents options for creating a task
type TaskOpts struct {
	Role     string
	Parent   string
	Blockers []string
	Priority string
}

// CreateTask creates a task directory and markdown file
func (e *TestEnv) CreateTask(taskID string, opts TaskOpts) {
	taskDir := filepath.Join(e.tasksDir, taskID)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		e.t.Fatalf("Failed to create task dir %s: %v", taskDir, err)
	}

	taskFile := filepath.Join(taskDir, taskID+".md")
	content := e.generateTaskContent(taskID, opts)

	if err := os.WriteFile(taskFile, []byte(content), 0644); err != nil {
		e.t.Fatalf("Failed to write task file %s: %v", taskFile, err)
	}
}

// CreateRole creates a role markdown file
func (e *TestEnv) CreateRole(roleName string) {
	roleDir := filepath.Join(e.baseDir, "roles")
	roleFile := filepath.Join(roleDir, roleName+".md")
	content := fmt.Sprintf("# %s\n\nTODO: Add role description.", strings.Title(roleName))

	if err := os.WriteFile(roleFile, []byte(content), 0644); err != nil {
		e.t.Fatalf("Failed to write role file %s: %v", roleFile, err)
	}
}

// CreateTemplate creates a template file
func (e *TestEnv) CreateTemplate(templateName string, content string) {
	templateDir := filepath.Join(e.baseDir, "templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		e.t.Fatalf("Failed to create templates dir: %v", err)
	}

	templateFile := filepath.Join(templateDir, templateName)
	if err := os.WriteFile(templateFile, []byte(content), 0644); err != nil {
		e.t.Fatalf("Failed to write template file %s: %v", templateFile, err)
	}
}

// RunCommand executes a CLI command in test environment
func (e *TestEnv) RunCommand(args ...string) (string, error) {
	if strandBinary == "" {
		return "", fmt.Errorf("strand binary not built (check TestMain)")
	}

	cmd := exec.Command(strandBinary, args...)
	cmd.Dir = e.rootDir

	// Pass through the isolated config dir
	cmd.Env = append(os.Environ(), "STRAND_CONFIG_DIR="+e.rootDir)

	output, err := cmd.CombinedOutput()
	return string(output), err
}

// AssertFileExists asserts that a file exists in test environment
func (e *TestEnv) AssertFileExists(relPath string) {
	fullPath := e.Path(relPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		e.t.Fatalf("Expected file %s to exist, but it doesn't", relPath)
	}
}

// AssertFileContains asserts that a file contains specific content
func (e *TestEnv) AssertFileContains(relPath, expectedContent string) {
	fullPath := e.Path(relPath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		e.t.Fatalf("Failed to read file %s: %v", relPath, err)
	}

	if !strings.Contains(string(content), expectedContent) {
		e.t.Fatalf("File %s does not contain expected content:\nExpected: %s\nActual:\n%s",
			relPath, expectedContent, string(content))
	}
}

// generateTaskContent generates task markdown content
func (e *TestEnv) generateTaskContent(taskID string, opts TaskOpts) string {
	blockersLine := "blockers: []"
	if len(opts.Blockers) > 0 {
		blockerList := make([]string, len(opts.Blockers))
		for i, b := range opts.Blockers {
			blockerList[i] = fmt.Sprintf("  - %s", b)
		}
		blockersLine = fmt.Sprintf("blockers:\n%s", strings.Join(blockerList, "\n"))
	}

	priority := opts.Priority
	if priority == "" {
		priority = "medium"
	}

	return fmt.Sprintf(`---
role: %s
parent: %s
priority: %s
%s
date_created: 2026-01-27
date_edited: 2026-01-27
---

# %s

## Summary

TODO: Add task summary.

## Tasks

- [ ] Add task description
- [ ] Implement acceptance criteria

## Acceptance Criteria

- Task is complete
- All requirements met
`,
		opts.Role, opts.Parent, priority, blockersLine, taskID)
}

// CreateTaskRaw creates a task with custom content
func (e *TestEnv) CreateTaskRaw(taskID, content string) {
	taskDir := filepath.Join(e.tasksDir, taskID)
	if err := os.MkdirAll(taskDir, 0755); err != nil {
		e.t.Fatalf("Failed to create task dir %s: %v", taskDir, err)
	}

	taskFile := filepath.Join(taskDir, taskID+".md")
	if err := os.WriteFile(taskFile, []byte(content), 0644); err != nil {
		e.t.Fatalf("Failed to write task file %s: %v", taskFile, err)
	}
}

// containsFlag checks if a flag is already present in args
func testToken(parts ...string) string {
	h := fnv.New32a()
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
	}
	return fmt.Sprintf("%08x", h.Sum32())[:6]
}

func testRoleName(t *testing.T, suffix string) string {
	name := strings.TrimSpace(t.Name())
	return "role-" + testToken("role", name, suffix)
}
