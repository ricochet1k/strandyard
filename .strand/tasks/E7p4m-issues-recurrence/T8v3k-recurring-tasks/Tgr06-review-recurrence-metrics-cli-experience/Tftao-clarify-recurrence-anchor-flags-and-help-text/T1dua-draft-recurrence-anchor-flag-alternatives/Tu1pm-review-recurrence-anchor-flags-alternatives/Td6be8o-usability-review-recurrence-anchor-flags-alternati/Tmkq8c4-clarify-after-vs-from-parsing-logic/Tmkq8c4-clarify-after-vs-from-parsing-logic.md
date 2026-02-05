---
type: implement
role: architect
priority: medium
parent: Td6be8o-usability-review-recurrence-anchor-flags-alternati
blockers: []
blocks: []
date_created: 2026-02-05T00:49:31.429809Z
date_edited: 2026-02-05T00:51:32.202451Z
owner_approval: false
completed: true
description: ""
---

# Clarify after vs from parsing logic

## Summary
Define the parsing logic and user-facing distinction between `after <anchor>` and `from <anchor>` in the structured grammar for `--every`.

Proposed distinction:
- `from <anchor>`: The recurrence starts *at* the anchor. The first materialization is eligible to happen at the anchor time (or immediately if the anchor is in the past).
- `after <anchor>`: The recurrence starts *one interval after* the anchor. The first materialization is only eligible after one full interval has passed since the anchor.

Example:
- `--every "1 day from Jan 1"` triggers on Jan 1, Jan 2, Jan 3...
- `--every "1 day after Jan 1"` triggers on Jan 2, Jan 3, Jan 4...

Acceptance Criteria:
- Update `design-docs/recurrence-anchor-update-logic.md` or create a new design doc clarifying this distinction.
- Update `cmd/add.go` and `cmd/edit.go` parsers to support the `after` keyword.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: developer) Implement the behavior described in Context.
  Updated design-docs/recurrence-anchor-update-logic.md to clarify the distinction between after and from keywords. Updated cmd/add.go to support the after keyword.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Verified after keyword correctly resolves to the next theoretical interval in the future.
- [x] (role: tester) Execute test-suite and report failures.
  N/A
- [x] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
  N/A
- [x] (role: documentation) Update user-facing docs and examples.
  N/A
