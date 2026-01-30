---
type: ""
role: developer
priority: ""
parent: E5w8m-e2e-tests
blockers:
    - T4p7k-test-framework-design
blocks:
    - E5w8m-e2e-tests
    - T7h5m-initial-e2e-tests
    - Tml0y-t9m4n-improved-task-templates
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-28T11:25:10.42807-07:00
owner_approval: false
completed: false
---

# Implement Sample Environment Setup

## Summary

Implement the test environment creation and management system that allows tests to spin up isolated environments with sample tasks, roles, and templates.

## Tasks

- [ ] Implement `TestEnv` struct with temp directory management
- [ ] Implement `CreateTask()` helper to create task files from structs
- [ ] Implement `CreateRole()` helper to create role files
- [ ] Implement `CreateTemplate()` helper to create template files
- [ ] Implement `RunCommand()` to execute CLI commands in test environment
- [ ] Implement `Cleanup()` to remove temp directories
- [ ] Add helpers for asserting file existence and content
- [ ] Add helpers for asserting command output

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
