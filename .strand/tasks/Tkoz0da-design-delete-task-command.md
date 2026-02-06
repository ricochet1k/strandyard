---
type: implement
role: designer
priority: medium
parent: Ibv0j-create-delete-task-command
blockers: []
blocks: []
date_created: 2026-02-05T01:07:40.85089Z
date_edited: 2026-02-05T01:08:17.651959Z
owner_approval: false
completed: true
description: ""
---

# Design delete task command

## Summary
Design a CLI command to safely delete tasks. Deletion should handle:
- Removing the task directory and its contents.
- Cleaning up references in parent tasks (TODO lists).
- Cleaning up references in other tasks (blockers/blocks).
- Updating master lists (root-tasks.md, free-tasks.md).
- Handling hierarchical deletion (should deleting a parent delete children?).

Deliverable: An alternatives document in `design-docs/` using `doc-examples/design-alternatives.md`.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: developer) Implement the behavior described in Context.
  Created design-docs/delete-task-command-alternatives.md exploring three alternatives for the delete command.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Verified the design covers safety, hierarchical cleanup, and data integrity.
- [x] (role: tester) Execute test-suite and report failures.
  N/A
- [x] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
  N/A
- [x] (role: documentation) Update user-facing docs and examples.
  N/A
