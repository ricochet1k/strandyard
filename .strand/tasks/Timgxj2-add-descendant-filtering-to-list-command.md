---
type: implement
role: developer
priority: medium
parent: Iquzx-add-parent-flag-to-list-command-to-list-only-commi
blockers: []
blocks: []
date_created: 2026-02-05T01:13:52.595941Z
date_edited: 2026-02-05T01:13:52.595941Z
owner_approval: false
completed: false
description: ""
---

# Add descendant filtering to list command

## Summary
Add a flag (e.g., `--descendants` or extend `--parent`) to the `list` command to list all tasks that are descendants of a given task ID, recursively.

Currently, `--children` only lists direct children.

Acceptance Criteria:
- New flag or extended flag allows recursive descendant listing.
- Works correctly with existing filters and output formats.

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
