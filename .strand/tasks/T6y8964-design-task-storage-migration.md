---
type: implement
role: designer
priority: medium
parent: Iff91-add-cli-support-for-task-storage-migration
blockers: []
blocks: []
date_created: 2026-02-05T01:09:32.253096Z
date_edited: 2026-02-05T01:09:56.823936Z
owner_approval: false
completed: true
description: ""
---

# Design task storage migration

## Summary
Design a CLI command or process to migrate task storage between "global" and "local" modes.
The migration should handle:
- Moving `tasks/`, `roles/`, and `templates/` to the new location (e.g., from top-level to `.strand/`).
- Updating the global configuration mapping.
- Handling Git tracking: if migrating to local storage, should the files be `git add`-ed? If migrating from local, should they be removed from Git?
- Updating master lists paths.

Deliverable: An alternatives document in `design-docs/` using `doc-examples/design-alternatives.md`.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: developer) Implement the behavior described in Context.
  Created design-docs/task-storage-migration-alternatives.md exploring three alternatives for storage migration.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Verified the design covers file moves, config updates, and Git tracking.
- [x] (role: tester) Execute test-suite and report failures.
  N/A
- [x] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
  N/A
- [x] (role: documentation) Update user-facing docs and examples.
  N/A
