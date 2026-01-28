---
role: developer
parent: E5w8m-e2e-tests
blockers:
  - T2n9w-sample-environments
blocks:
  - E5w8m-e2e-tests
date_created: 2026-01-27
date_edited: 2026-01-27
priority: low
---

# Create Initial E2E Tests for Repair and Next

## Summary

Create the first set of end-to-end tests covering the repair and next commands with various scenarios.

## Tasks

- [ ] Create test suite for `repair` command:
  - Valid task structure
  - Invalid ID format
  - Missing role file
  - Broken parent links
  - Broken blocker links
  - Master list generation
- [ ] Create test suite for `next` command:
  - First free task selection
  - Empty free-tasks list
  - Role extraction from metadata
  - Role extraction from TODO
  - Output format validation
- [ ] Create test fixtures with known-good and known-bad task files
- [ ] Ensure all tests clean up properly
- [ ] Document test cases and expected behavior

## Acceptance Criteria

- All repair scenarios have tests
- All next scenarios have tests
- Tests pass with current implementation
- Tests fail appropriately when bugs introduced
- Tests run quickly (< 1 second total)
- Clear test names describing what they repair

## Files

- test/e2e/repair_test.go (new)
- test/e2e/next_test.go (new)
- test/e2e/testdata/ (fixtures)

## Example Test

```go
func TestRepair_ValidTasks(t *testing.T) {
    env := NewTestEnv(t)
    defer env.Cleanup()

    env.CreateTask("T3k7x-example", TaskOpts{
        Role: "developer",
        Parent: "",
        Blockers: []string{},
    })
    env.CreateRole("developer")

    output, err := env.RunCommand("repair")

    assert.NoError(t, err)
    assert.Contains(t, output, "repair: ok")
    env.AssertFileExists("tasks/root-tasks.md")
    env.AssertFileExists("tasks/free-tasks.md")
}

func TestRepair_InvalidID(t *testing.T) {
    env := NewTestEnv(t)
    defer env.Cleanup()

    // Create task with invalid ID (only 3 chars)
    env.CreateTaskRaw("T3kx-bad", "...")

    _, err := env.RunCommand("repair")

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "malformed ID")
}
```
