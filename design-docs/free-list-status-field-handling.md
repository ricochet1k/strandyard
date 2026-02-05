# Free-List Status Field Handling Design

## Overview
The free-list is a critical mechanism for identifying actionable tasks (no blockers, ready to work on). As part of the Status Field Migration, we need to ensure the free-list correctly handles the new `status` field and excludes tasks with non-active statuses.

## Current State
- Free-list is generated in `pkg/task/free_list.go`
- Current logic checks `!task.Meta.Completed` to determine if a task is eligible
- Free-list calculation happens in `CalculateIncrementalFreeListUpdate()`

## Design Decision: Status Field Eligibility
**Status values eligible for free-list**: `open`, `in_progress`
**Status values excluded from free-list**: `done`, `cancelled`, `duplicate`

### Rationale
- `open`: Task is ready to be worked on
- `in_progress`: Task is actively being worked on
- `done`: Task is complete (shouldn't be in free-list)
- `cancelled`: Task will not be completed (shouldn't be in free-list)
- `duplicate`: Task is redundant (shouldn't be in free-list)

## Implementation Strategy

### Phase 1: Free-List Calculation Update
**File**: `pkg/task/free_list.go`

Add helper function:
```go
func IsActiveStatus(status string) bool {
    return status == "open" || status == "in_progress"
}
```

Update `CalculateIncrementalFreeListUpdate()`:
- Replace `!task.Meta.Completed` checks with `IsActiveStatus(task.Meta.Status)`
- Ensure backward compatibility: if `Status` is empty, default to `open`

### Phase 2: Validation Rules
**File**: `pkg/task/repair.go`

Add validation checks:
- All tasks in `free-tasks.md` must have `status: open` or `status: in_progress`
- All tasks in `free-tasks.md` must have empty `blockers` array
- `repair` command should fix free-list by recalculating from all tasks

### Phase 3: Testing
**Files**: `pkg/task/free_list_test.go`, integration tests

Test cases:
1. Tasks with `done` status excluded from free-list
2. Tasks with `cancelled` status excluded from free-list
3. Tasks with `duplicate` status excluded from free-list
4. Tasks with `open` status included in free-list
5. Tasks with `in_progress` status included in free-list
6. Free-list regeneration after status changes
7. Backward compatibility with old `completed: bool` format

## Impact Analysis

### Affected Commands
- `strand next` - filters tasks from free-list
- `strand complete` - sets status to `done` (removes from free-list)
- `strand list --completed` - filters based on status
- `strand repair` - regenerates free-list with status logic

### Integration Points
- Task loading/parsing - ensure Status field is populated
- Master list generation - `tasks/free-tasks.md` must reflect status field
- Blocker logic - only interact with active-status tasks
- Recurrence handling - may need to check status for "still needs doing"

## Edge Cases to Handle
1. Tasks with undefined/empty `status` field - treat as `open`
2. Migration from old `completed: bool` - convert on first load
3. Tasks transitioning between statuses - free-list should update immediately
4. Concurrent status updates - ensure free-list remains consistent

## Acceptance Criteria
- [ ] Free-list correctly excludes tasks with non-active statuses
- [ ] Free-list includes `open` and `in_progress` tasks
- [ ] Validation catches inconsistencies
- [ ] All tests pass
- [ ] No regression in existing free-list behavior
- [ ] `strand next` only shows actionable tasks
