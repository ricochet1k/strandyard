---
type: implement
role: developer
priority: medium
parent: Tqnugqd-review-task-id-resolution-implementation-for-usabi
blockers: []
blocks: []
date_created: 2026-02-05T00:37:49.287789Z
date_edited: 2026-02-05T00:37:49.287789Z
owner_approval: false
completed: false
description: ""
---

# Fix strand add failure when creating subtasks

## Summary
`strand add` fails when creating a task with a parent because it tries to load the newly created task from a fixed top-level path instead of its actual location under the parent directory. It also seems to expect `task.md` instead of `<task-id>.md`.

Acceptance Criteria:
- `strand add` successfully loads and displays the new task regardless of its location.
- Tests verify creating subtasks works without errors.

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
