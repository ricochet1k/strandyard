---
role: developer
parent: E5w8m-e2e-tests
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
priority: low
---

# Design E2E Test Framework

## Summary

Design the architecture and patterns for the e2e test framework, focusing on how to create isolated test environments and run CLI commands.

## Tasks

- [ ] Design test environment creation (temp directories, sample tasks, sample roles)
- [ ] Design CLI execution pattern (run commands in test environment)
- [ ] Design output validation patterns (stdout, stderr, files created, exit codes)
- [ ] Design cleanup mechanism (delete test environments after tests)
- [ ] Choose test fixtures approach (golden files, embedded resources, generated)
- [ ] Document test framework architecture
- [ ] Create example test skeleton showing pattern

## Acceptance Criteria

- Clear design for test environment lifecycle (setup, execute, validate, cleanup)
- Pattern for running CLI commands in isolated environments
- Pattern for asserting on outputs and side effects
- Documentation showing how to add new tests
- Example test demonstrating the pattern

## Files

- test/e2e/framework.go (design doc or initial implementation)
- test/e2e/README.md (framework documentation)
- test/e2e/example_test.go (example test)

## Example Test Pattern

```go
func TestValidateCommand(t *testing.T) {
    env := NewTestEnv(t)
    defer env.Cleanup()

    env.CreateTask("T3k7x-example", TaskOpts{...})
    env.CreateRole("developer")

    output, err := env.RunCommand("validate")

    assert.NoError(t, err)
    assert.Contains(t, output, "validate: ok")
    assert.FileExists(t, env.Path("tasks/free-tasks.md"))
}
```
