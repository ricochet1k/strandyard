---
type: implement
role: tester
priority: high
parent: Ta92amw-reliability-status-field-completeness-in-free-list
blockers: []
blocks: []
date_created: 2026-02-05T22:06:25.683949Z
date_edited: 2026-02-05T22:06:25.683949Z
owner_approval: false
completed: false
description: ""
---

# Add integration tests for status field with free-list

## Summary
Create comprehensive integration tests for the free-list behavior with the new status field.

**Test coverage**:
1. Test that `strand next` only shows tasks with `open` or `in_progress` status
2. Test that `strand complete` updates status to `done` and removes from free-list
3. Test that tasks with `cancelled` status don't appear in free-list
4. Test that tasks with `duplicate` status don't appear in free-list
5. Test free-list is regenerated correctly after status changes
6. Test migration: tasks with old `completed: true` convert to `status: done` and are excluded from free-list

**Test files**:
- `pkg/task/free_list_test.go` - unit tests for free-list calculation
- Integration tests in main test suite

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [ ] (role: developer) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
