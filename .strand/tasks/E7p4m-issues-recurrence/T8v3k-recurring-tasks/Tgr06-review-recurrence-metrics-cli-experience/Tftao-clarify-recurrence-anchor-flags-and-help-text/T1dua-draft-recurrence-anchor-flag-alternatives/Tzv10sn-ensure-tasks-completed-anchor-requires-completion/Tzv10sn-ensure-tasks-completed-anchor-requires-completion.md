---
type: implement
role: developer
priority: high
parent: T1dua-draft-recurrence-anchor-flag-alternatives
blockers: []
blocks: []
date_created: 2026-02-01T20:27:58.689359Z
date_edited: 2026-02-01T20:27:58.689359Z
owner_approval: false
completed: false
description: ""
---

# Ensure tasks_completed anchor requires completion timestamp metadata

## Summary
When using the 'tasks_completed' anchor for recurring tasks, validate that tasks have completion timestamp metadata. This metadata is necessary to determine when a task was completed for recurrence scheduling purposes.

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
