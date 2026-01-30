package e2e

import (
	"strings"
	"testing"
)

func TestNewTestEnv(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()

	// Verify directories exist
	env.AssertFileExists(".strand/tasks")
	env.AssertFileExists(".strand/roles")
}

func TestCreateTask(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()
	roleName := testRoleName(t, "task")

	env.CreateTask("T3k7x-example", TaskOpts{
		Role:     roleName,
		Parent:   "",
		Blockers: []string{},
		Priority: "high",
	})

	// Verify task file exists
	env.AssertFileExists(".strand/tasks/T3k7x-example/T3k7x-example.md")

	// Verify content
	env.AssertFileContains(".strand/tasks/T3k7x-example/T3k7x-example.md", "role: "+roleName)
	env.AssertFileContains(".strand/tasks/T3k7x-example/T3k7x-example.md", "priority: high")
	env.AssertFileContains(".strand/tasks/T3k7x-example/T3k7x-example.md", "# T3k7x-example")
}

func TestCreateRole(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()
	roleName := testRoleName(t, "role")

	env.CreateRole(roleName)

	// Verify role file exists
	env.AssertFileExists(".strand/roles/" + roleName + ".md")
	env.AssertFileContains(".strand/roles/"+roleName+".md", "# "+strings.Title(roleName))
}

func TestRunCommand(t *testing.T) {
	env := NewTestEnv(t)
	defer env.Cleanup()
	roleName := testRoleName(t, "role")

	env.CreateTask("T3k7x-example", TaskOpts{
		Role:     roleName,
		Parent:   "",
		Blockers: []string{},
	})
	env.CreateRole(roleName)

	output, err := env.RunCommand("repair")

	// Should succeed and output validation success
	if err != nil {
		t.Fatalf("Expected repair to succeed, got error: %v\nOutput: %s", err, output)
	}

	// Debug: print the output
	t.Logf("Repair output: %s", output)

	// Check for success indicators
	env.AssertFileContains(".strand/tasks/free-tasks.md", "T3k7x-example")
	env.AssertFileContains(".strand/tasks/root-tasks.md", "T3k7x-example")
}
