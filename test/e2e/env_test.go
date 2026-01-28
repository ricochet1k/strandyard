package e2e

import (
	"testing"
)

func TestNewTestEnv(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	// Verify directories exist
	env.AssertFileExists("tasks")
	env.AssertFileExists("roles")
}

func TestCreateTask(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	env.CreateTask("T3k7x-example", TaskOpts{
		Role:     "developer",
		Parent:   "",
		Blockers: []string{},
		Priority: "high",
	})

	// Verify task file exists
	env.AssertFileExists("tasks/T3k7x-example/T3k7x-example.md")

	// Verify content
	env.AssertFileContains("tasks/T3k7x-example/T3k7x-example.md", "role: developer")
	env.AssertFileContains("tasks/T3k7x-example/T3k7x-example.md", "priority: high")
	env.AssertFileContains("tasks/T3k7x-example/T3k7x-example.md", "# T3k7x-example")
}

func TestCreateRole(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	env.CreateRole("developer")

	// Verify role file exists
	env.AssertFileExists("roles/developer.md")
	env.AssertFileContains("roles/developer.md", "# Developer")
}

func TestRunCommand(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	env.CreateTask("T3k7x-example", TaskOpts{
		Role:     "developer",
		Parent:   "",
		Blockers: []string{},
	})
	env.CreateRole("developer")

	output, err := env.RunCommand("repair")

	// Should succeed and output validation success
	if err != nil {
		t.Fatalf("Expected repair to succeed, got error: %v\nOutput: %s", err, output)
	}

	// Debug: print the output
	t.Logf("Repair output: %s", output)

	// Check for success indicators
	env.AssertFileContains("tasks/free-tasks.md", "T3k7x-example")
	env.AssertFileContains("tasks/root-tasks.md", "T3k7x-example")
}
