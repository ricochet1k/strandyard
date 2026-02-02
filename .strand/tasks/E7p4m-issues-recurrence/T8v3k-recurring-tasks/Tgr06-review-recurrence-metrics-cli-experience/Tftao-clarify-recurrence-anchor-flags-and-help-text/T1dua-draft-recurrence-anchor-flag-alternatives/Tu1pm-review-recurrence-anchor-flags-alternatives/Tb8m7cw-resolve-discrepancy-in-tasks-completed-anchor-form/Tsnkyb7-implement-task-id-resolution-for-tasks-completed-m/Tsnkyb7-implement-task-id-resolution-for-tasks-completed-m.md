---
type: task
role: developer
priority: medium
parent: Tb8m7cw-resolve-discrepancy-in-tasks-completed-anchor-form
blockers: []
blocks: []
date_created: 2026-02-01T23:34:10.537422Z
date_edited: 2026-02-02T01:19:51.06607Z
owner_approval: false
completed: false
description: ""
---

# Implement task ID resolution for tasks_completed metric

## Summary
Implement task ID resolution for the `tasks_completed` recurrence metric as described in `design-docs/tasks-completed-anchor-resolution.md`.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: developer) Implement `GetLatestTaskCompletionTime` in `pkg/activity/log.go`.
  Implemented GetLatestTaskCompletionTime in pkg/activity/log.go with comprehensive tests
- [x] (role: developer) Update `EvaluateTasksCompletedMetric` in `pkg/task/recurrence.go` to support task ID anchors.
  Updated EvaluateTasksCompletedMetric to try resolving anchors as task IDs first, then fallback to date parsing
- [x] (role: developer) Update `UpdateAnchor` in `pkg/task/recurrence.go` to support task ID anchors.
  Updated UpdateAnchor to resolve task ID anchors using GetLatestTaskCompletionTime before calculating next anchor
- [x] (role: developer) Add unit and integration tests covering the main flows.
  Added comprehensive unit tests in recurrence_test.go for EvaluateTasksCompletedMetric and UpdateAnchor with task ID anchors. All tests pass.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples in `design-docs/` as specified in the plan.
