---
type: implement
role: developer
priority: high
parent: Ta92amw-reliability-status-field-completeness-in-free-list
blockers: []
blocks: []
date_created: 2026-02-05T22:06:22.251958Z
date_edited: 2026-02-05T22:06:22.251958Z
owner_approval: false
completed: false
description: ""
---

# Add free-list status validation rules

## Summary
Define and implement validation rules to ensure the free-list is accurate when status field is used.

**File**: `pkg/task/repair.go` and `pkg/task/free_list.go`

**Changes**:
1. Add validation in `repair.go` to check that all tasks in free-tasks.md have `status: open` or `status: in_progress`
2. Add validation that no tasks in free-tasks.md have `blockers` array
3. Create warning or error if a task with status `done`/`cancelled`/`duplicate` appears in free-tasks.md
4. Update master list regeneration logic to respect status field

**Tests**:
- Validation catches tasks with wrong status in free-list
- Validation allows open and in_progress tasks
- repair command fixes free-list when status values are incorrect

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
