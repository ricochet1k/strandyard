package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	roleFile := filepath.Join(e.Root(), "roles", roleName+".md")
	content := fmt.Sprintf("# %s\n\nTODO: Add role description.", strings.Title(roleName))

	if err := os.WriteFile(roleFile, []byte(content), 0644); err != nil {
		e.t.Fatalf("Failed to write role file %s: %v", roleFile, err)
	}
}

// CreateTemplate creates a template file
func (e *TestEnv) CreateTemplate(templateName string, content string) {
	templateDir := filepath.Join(e.Root(), "templates")
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
	// Build binary first to avoid go.mod issues
	binaryPath := "/tmp/memmd-test"
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %v", err)
	}
	repoRoot := filepath.Clean(filepath.Join(wd, "../.."))
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = repoRoot
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build failed: %v\nOutput: %s", err, string(output))
	}

	// Make binary executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return "", fmt.Errorf("failed to make binary executable: %v", err)
	}

	var cmd *exec.Cmd

	// Special handling for 'next' command which has hardcoded paths
	if len(args) > 0 && args[0] == "next" {
		// For 'next', run from test environment so hardcoded paths work
		cmd = exec.Command(binaryPath, args...)
		cmd.Dir = e.rootDir
	} else {
		// For other commands, run from project root with custom paths
		allArgs := args
		command := ""
		if len(args) > 0 {
			command = args[0]
		}

		// Commands that support --path flag
		if command == "repair" {
			if !containsFlag(args, "--path") {
				allArgs = append(allArgs, "--path", e.tasksDir)
			}
			if !containsFlag(args, "--roots") {
				allArgs = append(allArgs, "--roots", e.Path("tasks/root-tasks.md"))
			}
			if !containsFlag(args, "--free") {
				allArgs = append(allArgs, "--free", e.Path("tasks/free-tasks.md"))
			}
		}

		cmd = exec.Command(binaryPath, allArgs...)
		cmd.Dir = repoRoot
	}

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
func containsFlag(args []string, flag string) bool {
	for _, arg := range args {
		if arg == flag {
			return true
		}
	}
	return false
}
