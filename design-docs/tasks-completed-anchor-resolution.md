# Implementation Plan — tasks_completed Anchor Resolution

## Summary
Resolve the discrepancy between documentation and implementation regarding `tasks_completed` anchors. The system will support both date/time anchors and task ID anchors for the `tasks_completed` metric by querying the activity log.

## Context
- `CLI.md` and `ValidateAnchor` claim task ID anchors are supported for `tasks_completed`.
- `EvaluateTasksCompletedMetric` and `UpdateAnchor` currently only support date anchors.
- `design-docs/recurrence-anchor-error-messages.md` and `design-docs/anchor-help-text-and-examples-alternatives.md` incorrectly state that `tasks_completed` requires a date anchor.

## Proposed Changes

### `pkg/activity`
- Add `GetLatestTaskCompletionTime(taskID string) (time.Time, error)` to the `Log` struct. This will search the activity log for the most recent `task_completed` event for the given task ID.

### `pkg/task`
- Update `EvaluateTasksCompletedMetric(baseDir, anchor string, taskID string, log *activity.Log) (int, error)`:
  - If `anchor` is not "now" or empty, try to resolve it as a task ID first using the activity log.
  - If it's a valid task ID found in the log, use its completion timestamp as the `anchorTime`.
  - Fallback to date parsing if it's not a task ID or not found in the log (or if date parsing succeeds).
- Update `UpdateAnchor(repoPath, baseDir, metric, currentAnchor string, interval int) (string, error)`:
  - Implement the TODO for `tasks_completed` task ID resolution.
  - Use `GetLatestTaskCompletionTime` to resolve task ID anchors to timestamps before calculating the next anchor.

### Documentation
- Update `design-docs/recurrence-anchor-error-messages.md`:
  - Update the "Ambiguous Anchor Type" section to reflect that `tasks_completed` accepts both date and task ID anchors.
  - Update the "Anchor Type by Metric" table.
- Update `design-docs/anchor-help-text-and-examples-alternatives.md`:
  - Update the "Per-Unit Anchor Mappings" table to include task ID as a valid anchor type for `tasks_completed`.

## Architecture Decisions
- **Activity Log as Source of Truth**: Task ID resolution for `tasks_completed` will rely solely on the activity log. This ensures consistency even if tasks are deleted from the filesystem.
- **Latest Completion**: If a task ID has multiple completion events in the log (unlikely for a specific task ID but possible if IDs are reused or for some types of tasks), the latest completion event will be used as the anchor.
- **Fallback to Date**: If an anchor string can be parsed as both a task ID and a date (unlikely given the `T...` prefix), task ID resolution takes precedence if found in the log.

## Verification Plan
- Unit tests in `pkg/activity/log_test.go` for `GetLatestTaskCompletionTime`.
- Unit tests in `pkg/task/recurrence_test.go` for `EvaluateTasksCompletedMetric` and `UpdateAnchor` with task ID anchors.
- Integration test in `cmd/add_every_test.go` or similar to verify `strand add --every "20 tasks_completed from T1234"` works as expected.
