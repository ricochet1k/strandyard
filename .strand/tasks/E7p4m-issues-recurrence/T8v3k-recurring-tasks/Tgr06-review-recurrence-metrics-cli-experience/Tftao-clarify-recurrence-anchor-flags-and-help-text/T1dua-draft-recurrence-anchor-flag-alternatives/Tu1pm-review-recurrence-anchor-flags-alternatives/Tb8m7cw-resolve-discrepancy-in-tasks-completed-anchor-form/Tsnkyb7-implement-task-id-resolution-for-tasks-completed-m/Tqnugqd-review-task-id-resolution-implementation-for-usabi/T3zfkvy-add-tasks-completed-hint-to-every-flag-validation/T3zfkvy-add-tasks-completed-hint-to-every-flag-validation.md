---
type: implement
role: developer
priority: medium
parent: Tqnugqd-review-task-id-resolution-implementation-for-usabi
blockers: []
blocks: []
date_created: 2026-02-05T00:48:15.025652Z
date_edited: 2026-02-05T00:48:15.025652Z
owner_approval: false
completed: false
description: ""
---

# Add tasks_completed hint to --every flag validation in cmd/add.go

## Summary
The `validateEvery` function in `cmd/add.go` provides helpful hints when anchor validation fails for most metrics, but it is missing a specific hint for the `tasks_completed` metric. A hint should be added to show valid anchor formats (date or task ID) for `tasks_completed`.

Acceptance Criteria:
- Validation failure for `tasks_completed` anchor displays a hint like: `hint: --every "20 tasks_completed from T1a1a"`.

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
