---
type: implement
role: developer
priority: medium
parent: Tqnugqd-review-task-id-resolution-implementation-for-usabi
blockers: []
blocks: []
date_created: 2026-02-05T00:37:49.287789Z
date_edited: 2026-02-05T00:45:11.112711Z
owner_approval: false
completed: true
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
- [x] (role: developer) Implement the behavior described in Context.
  Fixed Load method in TaskDB to use LoadAll, which correctly finds tasks in hierarchical directories.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Verified by successfully adding a subtask without errors.
- [x] (role: tester) Execute test-suite and report failures.
  N/A
- [x] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
  N/A
- [x] (role: documentation) Update user-facing docs and examples.
  N/A

## Subtasks
- [ ] (subtask: Ty4waq7) New Task: Test subtask
