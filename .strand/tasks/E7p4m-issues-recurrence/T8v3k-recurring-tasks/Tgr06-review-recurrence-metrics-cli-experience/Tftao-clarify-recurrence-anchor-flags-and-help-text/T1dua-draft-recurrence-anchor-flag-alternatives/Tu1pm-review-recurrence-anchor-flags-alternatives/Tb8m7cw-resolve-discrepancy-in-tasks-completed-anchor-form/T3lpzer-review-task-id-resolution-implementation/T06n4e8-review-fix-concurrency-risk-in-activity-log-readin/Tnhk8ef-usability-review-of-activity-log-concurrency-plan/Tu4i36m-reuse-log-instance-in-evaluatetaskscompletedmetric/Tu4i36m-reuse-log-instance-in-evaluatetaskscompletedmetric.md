---
type: implement
role: developer
priority: medium
parent: Tnhk8ef-usability-review-of-activity-log-concurrency-plan
blockers: []
blocks: []
date_created: 2026-02-01T23:51:46.463593Z
date_edited: 2026-02-01T23:51:46.463593Z
owner_approval: false
completed: false
description: ""
---

# Reuse log instance in EvaluateTasksCompletedMetric

## Summary
Refactor EvaluateTasksCompletedMetric in pkg/task/recurrence.go to reuse the passed-in activity.Log if it's not nil. See design-docs/reuse-log-instance-plan.md for the implementation plan.

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
