# Status Field Migration Design

## Overview
Replace the simple `completed: bool` field with a multi-state `status` field that captures the full task lifecycle. This enables better visibility into task states and eliminates the need for a separate delete command.

## Status Values
- `open` - Task is ready to work on (default for new tasks)
- `in_progress` - Task is actively being worked on
- `done` - Task successfully completed
- `duplicate` - Task is a duplicate of another task
- `cancelled` - Task explicitly cancelled or won't-fix
- `blocked` - (Optional) Task waiting on dependencies (implicit in `blockers` array, may not need explicit status)

## Current State
**Data Model** (`pkg/task/task.go`):
- `Completed bool` in `Metadata` struct (line 23)

**Usage Points**:
1. `pkg/task/taskdb.go`:
   - `SetCompleted(taskID, completed bool)` - lines 276-293
   - `CompleteTask(taskID, report)` - line 600
   - `CompleteTodo(...)` - sets `Meta.Completed = true` at line 564
   - `UpdateBlockersAfterCompletion()` - processes completed tasks

2. `cmd/complete.go`:
   - Checks `t.Meta.Completed` at line 94
   - Sets completion state at line 105 via `db.CompleteTask()`
   - Command documentation references "completed: true"

3. `cmd/next.go`:
   - Filters for free tasks (tasks not completed and with no blockers)
   - Free-list generation respects completed status

4. `cmd/list.go`:
   - `--completed` flag to show only completed tasks (line 23)

5. Migration references:
   - Expects transformation: `completed: true` → `status: done`
   - Expects transformation: `completed: false` → `status: open`

## Implementation Strategy

### Phase 1: Data Model Updates
1. Update `Metadata` struct in `pkg/task/task.go`:
   - Add `Status string` field with yaml tag
   - Keep `Completed bool` field temporarily for backward compatibility
   - Add helper functions: `IsOpen()`, `IsDone()`, `GetStatus()`, `SetStatus()`

2. Add status validation:
   - Valid values: "open", "in_progress", "done", "duplicate", "cancelled"
   - Default: "open" for new tasks
   - Migration logic: convert boolean on load if needed

### Phase 2: TaskDB Updates
1. Update `pkg/task/taskdb.go`:
   - Modify `SetCompleted()` to also set status
   - Add `SetStatus(taskID, status)` method
   - Update `CompleteTask()` to set `status: done`
   - Add `CancelTask(taskID, reason)` to set `status: cancelled`
   - Add `MarkDuplicate(taskID, duplicateOf)` to set `status: duplicate`
   - Add `MarkInProgress(taskID)` to set `status: in_progress`
   - Update `UpdateBlockersAfterCompletion()` to check status instead of completed boolean

2. Update filtering logic:
   - `GetIncompleteTodos()` - check if task status allows incomplete todos (not done/duplicate/cancelled)
   - `CalculateIncrementalFreeListUpdate()` - consider status field
   - Validation rules should reflect that only `open` or `in_progress` tasks appear in free-list

### Phase 3: CLI Command Updates
1. `cmd/complete.go`:
   - Update to set `status: done` instead of `completed: true`
   - Update command help text
   - Validation: prevent completing tasks with status `duplicate` or `cancelled`

2. `cmd/next.go`:
   - Filter to show only `open` or `in_progress` tasks (not done/duplicate/cancelled)
   - Update free-list filtering logic

3. `cmd/list.go`:
   - Update `--completed` flag logic to check `status == done`
   - Consider adding `--status` flag for more granular filtering

4. New commands:
   - `strand cancel <task-id> "reason"` - sets `status: cancelled`
   - `strand mark-duplicate <task-id> <duplicate-of>` - sets `status: duplicate`
   - `strand mark-in-progress <task-id>` - sets `status: in_progress`

### Phase 4: Migration & Backward Compatibility
1. Task loading:
   - When reading a task file, if `completed: true/false` exists and `status` is missing:
     - Load both fields
     - Convert: `completed: true` → `status: done`
     - Convert: `completed: false` → `status: open` (or keep if explicitly open)
   - Mark task as dirty so it rewrites with new status field

2. Master list generation:
   - `tasks/free-tasks.md` - regenerate to exclude non-open/in-progress tasks
   - `tasks/root-tasks.md` - no change needed, includes all tasks

3. Activity logging:
   - Update event types or add new ones for different status transitions

### Phase 5: Testing
1. Unit tests:
   - Test status field parsing and serialization
   - Test migration from completed boolean to status
   - Test each status value behavior
   - Test filtering logic with new statuses

2. Integration tests:
   - `strand complete` sets correct status
   - `strand next` filters correctly
   - `strand list --completed` works with new status
   - New commands create correct status values

3. E2E tests:
   - Verify master lists are regenerated correctly
   - Verify old task files migrate on first load

## Files to Modify
- **pkg/task/task.go**: Add `Status` field and helpers
- **pkg/task/taskdb.go**: Update all completion logic and add new status methods
- **cmd/complete.go**: Update to use status field
- **cmd/next.go**: Update filtering to respect status
- **cmd/list.go**: Update --completed flag behavior
- **cmd/root.go**: Add new commands (cancel, mark-duplicate, mark-in-progress)
- **pkg/task/free_list.go**: Update free-list generation for new status values
- **pkg/task/repair.go**: Update validation for status field
- **Tests**: Update all existing tests and add new ones

## Acceptance Criteria
- All existing `completed: true/false` tasks migrate to `status: done/open` on first load
- Tasks with status other than `open`/`in_progress` don't appear in next/free-list
- All four new status values work correctly
- `strand complete` sets `status: done`
- New `strand cancel`, `strand mark-duplicate`, `strand mark-in-progress` commands work
- No regression in existing functionality (complete, next, list commands)
- Master lists are correctly regenerated with new status values

## Decision: Implementation Approach
**Decision: Implement in phases, starting with data model and core TaskDB changes, followed by CLI updates.**

This approach minimizes risk by:
1. Centralizing the status logic in the data model and TaskDB
2. Maintaining backward compatibility during migration
3. Allowing CLI changes to be tested independently
4. Enabling incremental rollout to users

## Dependencies
- None - this is a self-contained feature
- All changes are internal to the task data model and CLI
- No external API or library changes required
