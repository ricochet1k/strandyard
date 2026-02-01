---
type: implement
role: developer
priority: high
parent: T3lpzer-review-task-id-resolution-implementation
blockers: [T06n4e8-review-fix-concurrency-risk-in-activity-log-readin]
blocks: []
date_created: 2026-02-01T23:41:26.678109Z
date_edited: 2026-02-01T23:41:26.678109Z
owner_approval: false
completed: false
description: ""
---

# Fix concurrency risk in activity log reading

## Summary
Refactor pkg/activity/log.go to use sync.RWMutex and avoid closing the write handle during reads.

See design-docs/fix-activity-log-concurrency.md for details.

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
