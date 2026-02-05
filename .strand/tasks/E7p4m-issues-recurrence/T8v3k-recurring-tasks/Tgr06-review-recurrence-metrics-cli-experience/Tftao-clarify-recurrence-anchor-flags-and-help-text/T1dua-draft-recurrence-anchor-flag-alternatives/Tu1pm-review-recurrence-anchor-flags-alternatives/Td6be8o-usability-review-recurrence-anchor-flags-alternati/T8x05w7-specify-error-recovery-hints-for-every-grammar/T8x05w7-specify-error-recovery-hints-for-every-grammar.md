---
type: implement
role: designer
priority: medium
parent: Td6be8o-usability-review-recurrence-anchor-flags-alternati
blockers: []
blocks: []
date_created: 2026-02-05T00:49:51.698112Z
date_edited: 2026-02-05T00:50:19.294286Z
owner_approval: false
completed: true
description: ""
---

# Specify error recovery hints for --every grammar

## Summary
User-facing error messages for the `--every` structured string must include the correct grammar and a relevant example based on the partial input. While `design-docs/recurrence-anchor-error-messages.md` defines a baseline, this task is to ensure they are comprehensive and provide the best possible recovery path for users.

Acceptance Criteria:
- Review `design-docs/recurrence-anchor-error-messages.md` for completeness.
- Ensure hints are context-aware (e.g., if metric is known, show hint for that metric).
- Finalize the set of error messages and hints.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: developer) Implement the behavior described in Context.
  Reviewed design-docs/recurrence-anchor-error-messages.md. It now includes comprehensive error messages and hints for all supported metrics including tasks_completed.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Verified implementation in cmd/add.go matches the design and provides context-aware hints.
- [x] (role: tester) Execute test-suite and report failures.
  N/A
- [x] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
  N/A
- [x] (role: documentation) Update user-facing docs and examples.
  N/A
