---
type: implement
role: developer
priority: high
parent: Twqo3xp-status-field-migration-epic
blockers:
    - Tccehse-review-update-metadata-struct-for-status-field
blocks:
    - Ta3bynh-update-taskdb-operations-for-status-field
date_created: 2026-02-05T22:02:32.076427Z
date_edited: 2026-02-05T22:02:55.687073Z
owner_approval: false
completed: false
description: ""
---

# Update Metadata struct and add status helpers

## Summary
Update the task data model in `pkg/task/task.go` to support the new status field.

**Implementation Plan**: See design-docs/status-field-migration.md (Phase 1)

**Specific changes**:
1. Add `Status string` field to Metadata struct with yaml tag
2. Keep `Completed bool` temporarily for backward compatibility
3. Add helper methods:
   - `IsOpen()` - returns true if status is "open"
   - `IsDone()` - returns true if status is "done"  
   - `IsActive()` - returns true if status is "open" or "in_progress"
   - `SetStatus(status string)` - validates and sets status
   - `GetStatus()` - returns current status with default handling
4. Add status constants: StatusOpen, StatusInProgress, StatusDone, StatusDuplicate, StatusCancelled
5. Update migration logic to convert `completed: bool` to `status` on load

**Acceptance Criteria**:
- All status values parse and serialize correctly
- Helper methods return correct values
- Backward compatibility with existing completed: bool tasks
- Unit tests cover all status values
- No regression in existing task loading/parsing

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
