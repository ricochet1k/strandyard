# E2E Test Framework

## Overview

This framework provides isolated test environments for testing the StrandYard CLI commands end-to-end.

## Architecture

### TestEnv

The `TestEnv` struct manages a temporary directory with the following structure:

```
/tmp/strand-test-*/
├── tasks/          # Task files
├── roles/          # Role files
└── templates/      # Template files (optional)
```

### Core Components

- **Test Creation**: `NewTestEnv(t)` creates an isolated environment with cleanup
- **Task Management**: `CreateTask()`, `CreateTaskRaw()` for setting up test tasks
- **Role Management**: `CreateRole()` for creating role files
- **Command Execution**: `RunCommand()` for executing CLI commands
- **Assertions**: `AssertFileExists()`, `AssertFileContains()` for validation

## Usage Pattern

```go
func TestRepairCommand(t *testing.T) {
    env := NewTestEnv(t)
    defer env.Cleanup()

    // Setup test data
    env.CreateTask("T3k7x-example", TaskOpts{
        Role:     "developer",
        Parent:   "",
        Blockers: []string{},
    })
    env.CreateRole("developer")

    // Execute command
    _, err := env.RunCommand("repair")

    // Assert results
    if err != nil {
        t.Fatalf("Expected repair to succeed, got error: %v", err)
    }
    env.AssertFileExists("tasks/free-tasks.md")
    env.AssertFileContains("tasks/free-tasks.md", "T3k7x-example")
}
```

## Best Practices

1. **Always use defer env.Cleanup()** immediately after NewTestEnv
2. **Use descriptive task IDs** that follow the format `<PREFIX><4-char>-<slug>`
3. **Create role files** when testing role-dependent functionality
4. **Test both success and failure scenarios** for robust coverage
5. **Use TaskOpts struct** for consistent task creation
6. **Repair both command output and file side effects**

## File Structure

- `env.go` - TestEnv struct and basic functionality
- `helpers.go` - Task/role creation and command execution helpers
- `env_test.go` - Basic framework tests
- Additional test files for specific CLI commands

## Integration with CI

The framework runs in temporary directories, making it safe for CI environments:
- No state pollution between tests
- Automatic cleanup prevents disk space issues
- Works in any environment where Go can build the project