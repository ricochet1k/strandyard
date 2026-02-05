---
type: implement
role: developer
priority: high
parent: Twqo3xp-status-field-migration-epic
blockers:
    - Tzk35d7-update-metadata-struct-and-add-status-helpers
blocks:
    - Tjqppdw-update-and-add-cli-commands-for-status-field
date_created: 2026-02-05T22:02:36.862406Z
date_edited: 2026-02-05T22:04:34.267776Z
owner_approval: false
completed: false
description: ""
---

# Update TaskDB operations for status field

## Summary
Update the TaskDB in `pkg/task/taskdb.go` to support the new status field operations.

**Implementation Plan**: See design-docs/status-field-migration.md (Phase 2)

**Specific changes**:
1. Update `SetCompleted()` to also set status field
2. Add `SetStatus(taskID, status)` method with validation
3. Update `CompleteTask()` to set `status: done`
4. Add `CancelTask(taskID, reason)` to set `status: cancelled`
5. Add `MarkDuplicate(taskID, duplicateOf)` to set `status: duplicate`
6. Add `MarkInProgress(taskID)` to set `status: in_progress`
7. Update `UpdateBlockersAfterCompletion()` to check status
8. Update filtering logic to respect status field:
   - `GetIncompleteTodos()` - exclude non-active tasks
   - `CalculateIncrementalFreeListUpdate()` - consider status
   - Task validation - only open/in_progress in free-list

**Acceptance Criteria**:
- All new methods work correctly
- Status is persisted to disk
- Filtering respects status values
- No regression in existing TaskDB operations
- Unit tests for each new method
- Integration tests for status transitions

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
