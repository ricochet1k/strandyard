---
type: implement
role: developer
priority: high
parent: Ta92amw-reliability-status-field-completeness-in-free-list
blockers: []
blocks: []
date_created: 2026-02-05T22:06:18.340899Z
date_edited: 2026-02-05T22:06:18.340899Z
owner_approval: false
completed: false
description: ""
---

# Update free-list generation to check status field

## Summary
Update the free-list calculation to exclude tasks with non-active statuses (anything other than `open` or `in_progress`).

**File**: `pkg/task/free_list.go`

**Changes**:
1. In `CalculateIncrementalFreeListUpdate()`, update the predicate to check `task.Meta.Status` instead of `!task.Meta.Completed`
2. Verify tasks are in `open` or `in_progress` state, not `done`, `cancelled`, or `duplicate`
3. Add helper function to check if status is "active" (open or in_progress)

**Tests**:
- Verify tasks with status `done` don't appear in free-list
- Verify tasks with status `cancelled` don't appear in free-list
- Verify tasks with status `duplicate` don't appear in free-list
- Verify tasks with `open` and `in_progress` statuses DO appear in free-list

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
