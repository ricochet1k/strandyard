---
type: design
role: developer
priority: high
parent: T1dua-draft-recurrence-anchor-flag-alternatives
blockers: []
blocks: []
date_created: 2026-02-01T20:27:58.689359Z
date_edited: 2026-02-01T20:55:03.530649Z
owner_approval: false
completed: false
description: ""
---

# Design activity log for tasks_completed metric

## Summary
The `tasks_completed` metric for recurring tasks should be based on an activity log rather than task completion timestamp metadata, since completed tasks may be deleted. Design an activity log system that records task completion events and can be queried for recurrence scheduling.

## Context
- Completed tasks may be deleted from the task database
- The `tasks_completed` metric needs to count completions regardless of whether the task file still exists
- An activity log (not yet designed) should record task completion events persistently
- This log should support querying for completion counts within a date range

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: developer) Design activity log schema (file format, fields, location)
  Updated task from 'implement' to 'design' type. The original task was based on an incorrect assumption that tasks_completed should validate completion timestamp metadata on tasks. Corrected approach: tasks_completed should use an activity log (not yet designed) since completed tasks may be deleted. Updated task title, summary, and TODOs to reflect the activity log design approach.
- [x] (role: developer) Implement activity log writer for task completion events
  Implemented activity log writer for task completion events. Created pkg/activity package with Log struct that writes JSONL entries to .strand/activity.log. Updated strand complete command to write to activity log when tasks are completed (both directly and via last todo completion).
- [x] (role: developer) Implement activity log reader for queries (e.g., count completions since date)
  Implemented activity log reader methods: ReadEntries(), CountCompletionsSince(), CountCompletionsForTaskSince(). All tests pass.
- [x] (role: developer) Update `strand complete` to write to activity log
  Verified activity log writing is implemented in complete.go for both direct completion and last todo completion. Added integration tests TestCompleteWritesToActivityLog and TestCompleteViaLastTodoWritesToActivityLog to verify the functionality.
- [x] (role: developer) Update recurrence evaluation to query activity log for `tasks_completed` metric
  Implemented evaluateTasksCompletedMetric function in pkg/task/recurrence.go that queries the activity log for task completion counts since a given anchor time. Added test TestEvaluateTasksCompletedMetric covering empty log, some completions, and invalid date format scenarios. All tests pass.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Added unit and integration tests for activity log main flows:
- [x] Added TestReadEntriesHandlesMalformedEntry in pkg/activity/log_test.go to verify resilience against malformed log entries
  Verified TestReadEntriesHandlesMalformedEntry test exists in pkg/activity/log_test.go and passes. The test correctly verifies that ReadEntries() returns an error when encountering malformed JSON in the activity log, ensuring resilience against corrupted entries.
- [ ] Added TestActivityLogIntegration in cmd/complete_test.go for end-to-end testing of the complete flow: creates multiple tasks, completes them, verifies activity log records all completions, and tests query functionality (CountCompletionsSince and CountCompletionsForTaskSince)
  All tests pass (go test ./...) and all builds pass (go build ./...).
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
