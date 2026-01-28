package e2e

import (
	"testing"
)

func TestRepair_ValidTasks(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	env.CreateTask("T3k7x-example", TaskOpts{
		Role:     "developer",
		Parent:   "",
		Blockers: []string{},
	})
	env.CreateRole("developer")

	output, err := env.RunCommand("repair")

	if err != nil {
		t.Fatalf("Expected repair to succeed, got error: %v", err)
	}

	// Check for success indicators
	env.AssertFileExists("tasks/free-tasks.md")
	env.AssertFileExists("tasks/root-tasks.md")
	env.AssertFileContains("tasks/free-tasks.md", "T3k7x-example")
	env.AssertFileContains("tasks/root-tasks.md", "T3k7x-example")

	// Check output contains success message
	if !contains(output, "repair: ok") {
		t.Errorf("Expected output to contain 'repair: ok', got: %s", output)
	}
}

func TestRepair_InvalidID(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	// Create task with invalid ID (only 3 chars in token)
	env.CreateTaskRaw("T3kx-bad", `---
role: developer
parent: ""
blockers: []
---
# Bad Task
`)

	output, err := env.RunCommand("repair")

	if err == nil {
		t.Fatal("Expected repair to fail with invalid ID")
	}

	// Check error message contains ID validation error
	if !contains(output, "malformed ID") && !contains(err.Error(), "malformed ID") {
		t.Errorf("Expected output or error to contain 'malformed ID', got: output=%s error=%v", output, err)
	}
}

func TestRepair_MissingRole(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	env.CreateTask("T3k7x-example", TaskOpts{
		Role:     "nonexistent",
		Parent:   "",
		Blockers: []string{},
	})
	// Don't create the role file

	output, err := env.RunCommand("repair")

	if err == nil {
		t.Fatal("Expected repair to fail with missing role")
	}

	// Check error message contains role validation error
	if !contains(output, "role file") && !contains(err.Error(), "role file") {
		t.Errorf("Expected output or error to contain 'role file', got: output=%s error=%v", output, err)
	}
}

func TestNext_FreeTask(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	env.CreateTask("T3k7x-example", TaskOpts{
		Role:     "developer",
		Parent:   "",
		Blockers: []string{},
	})
	env.CreateRole("developer")

	// First run repair to generate free-tasks.md
	_, err := env.RunCommand("repair")
	if err != nil {
		t.Fatalf("Expected repair to succeed: %v", err)
	}

	output, err := env.RunCommand("next")

	if err != nil {
		t.Fatalf("Expected next to succeed, got error: %v. Output: %s", err, output)
	}

	// Check output contains task information
	if !contains(output, "T3k7x-example") {
		t.Errorf("Expected output to contain task ID, got: %s", output)
	}
	if !contains(output, "Role: developer") {
		t.Errorf("Expected output to contain role, got: %s", output)
	}
}

func TestNext_NoFreeTasks(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	// Create two tasks that block each other (so none are free)
	env.CreateTask("T5h7w-blocker", TaskOpts{
		Role:     "developer",
		Parent:   "",
		Blockers: []string{"T3k7x-blocked"},
	})
	env.CreateTask("T3k7x-blocked", TaskOpts{
		Role:     "developer",
		Parent:   "",
		Blockers: []string{"T5h7w-blocker"},
	})
	env.CreateRole("developer")

	output, err := env.RunCommand("next")

	if err != nil {
		t.Fatalf("Expected next to succeed even with no free tasks, got error: %v. Output: %s", err, output)
	}

	// Check output indicates no free tasks
	if !contains(output, "No free tasks found") {
		t.Errorf("Expected output to contain 'No free tasks found', got: %s", output)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
