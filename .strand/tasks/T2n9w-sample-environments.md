---
type: ""
role: developer
priority: ""
parent: E5w8m-e2e-tests
blockers: []
blocks:
    - E5w8m-e2e-tests
    - T7h5m-initial-e2e-tests
    - Tml0y-t9m4n-improved-task-templates
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-02-01T04:27:50.386299Z
owner_approval: false
completed: true
---

# Implement Sample Environment Setup

## Summary
Implement the test environment creation and management system that allows tests to spin up isolated environments with sample tasks, roles, and templates.

## Acceptance Criteria
- `NewTestEnv(t)` creates isolated temp directory
- Helper methods make it easy to set up test scenarios
- `RunCommand()` executes commands with proper working directory
- `Cleanup()` reliably removes test directories
- Can create complex test scenarios with minimal boilerplate

## Files
- test/e2e/env.go (new)
- test/e2e/env_test.go (new)
- test/e2e/helpers.go (new)

## Example Usage
```go
env := NewTestEnv(t)
defer env.Cleanup()

env.CreateTask("T3k7x-example", TaskOpts{
    Role: "developer",
    Blockers: []string{},
})

env.CreateRole("developer")

output, err := env.RunCommand("repair", "--path", "tasks")
env.AssertFileExists("tasks/free-tasks.md")
```

## TODOs
- [x] Implement `TestEnv` struct with temp directory management
- [x] Implement `CreateTask()` helper to create task files from structs
- [x] Implement `CreateRole()` helper to create role files
- [x] Implement `CreateTemplate()` helper to create template files
- [x] Implement `RunCommand()` to execute CLI commands in test environment
- [x] Implement `Cleanup()` to remove temp directories
- [x] Add helpers for asserting file existence and content
- [x] Add helpers for asserting command output
