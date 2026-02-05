---
type: implement
role: developer
priority: medium
parent: Tqnugqd-review-task-id-resolution-implementation-for-usabi
blockers: []
blocks: []
date_created: 2026-02-05T00:45:36.55135Z
date_edited: 2026-02-05T00:45:36.55135Z
owner_approval: false
completed: false
description: ""
---

# Resolve task ID anchors to full IDs during recurrence validation

## Summary
During recurrence validation in `cmd/add.go` and `cmd/edit.go`, task ID anchors for the `tasks_completed` metric should be resolved to their full canonical IDs. This resolved ID should then be stored in the task metadata instead of the short ID provided by the user.

This ensures that task files are self-contained and resilient to future changes in the task database (like adding ambiguous short IDs).

Acceptance Criteria:
- `strand add --every "20 tasks_completed from <short-id>"` stores the full ID in the task file.
- `strand edit --every ...` also resolves and stores full IDs.
- Validation fails if the short ID is ambiguous.

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
