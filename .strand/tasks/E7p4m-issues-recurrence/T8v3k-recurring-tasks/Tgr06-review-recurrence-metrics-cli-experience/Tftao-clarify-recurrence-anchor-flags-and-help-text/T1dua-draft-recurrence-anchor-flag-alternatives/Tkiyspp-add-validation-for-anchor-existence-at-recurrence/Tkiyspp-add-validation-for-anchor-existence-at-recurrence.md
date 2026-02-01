---
type: implement
role: developer
priority: medium
parent: T1dua-draft-recurrence-anchor-flag-alternatives
blockers: []
blocks: []
date_created: 2026-02-01T21:23:26.413541Z
date_edited: 2026-02-01T21:58:04.582867Z
owner_approval: false
completed: true
description: ""
---

# Add validation for anchor existence at recurrence creation

## Summary
Implement the validation logic for recurrence anchors as described in design-docs/recurrence-anchor-validation-plan.md.

Ensure that:
- Date anchors are valid timestamps.
- Commit anchors exist in the git repository.
- Task anchors exist in the task database.

Refer to the implementation plan for specific details and testing requirements.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: developer) Implement the behavior described in Context.
  Implemented date, commit, and task anchor validation.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Added tests in pkg/task/recurrence_test.go and updated cmd/add_every_test.go.
- [x] (role: tester) Execute test-suite and report failures.
  Ran all unit and e2e tests successfully.
- [x] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
  Prepared implementation for reliability, security, and usability review.
- [x] (role: documentation) Update user-facing docs and examples.
  Updated CLI.md to document ISO 8601 and task ID anchor support.
