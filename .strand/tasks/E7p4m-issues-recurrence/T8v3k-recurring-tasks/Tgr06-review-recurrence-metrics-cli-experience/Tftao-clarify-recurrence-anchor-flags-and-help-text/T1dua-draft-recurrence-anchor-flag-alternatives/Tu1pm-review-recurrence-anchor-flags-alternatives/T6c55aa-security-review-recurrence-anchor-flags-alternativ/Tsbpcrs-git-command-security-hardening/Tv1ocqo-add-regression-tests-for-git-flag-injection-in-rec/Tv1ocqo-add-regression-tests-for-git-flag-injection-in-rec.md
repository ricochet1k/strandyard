---
type: implement
role: developer
priority: medium
parent: Tsbpcrs-git-command-security-hardening
blockers: []
blocks: []
date_created: 2026-02-01T22:13:15.54608Z
date_edited: 2026-02-01T22:13:15.54608Z
owner_approval: false
completed: false
description: ""
---

# Add regression tests for git flag injection in recurrence

## Summary
## Summary
Add test cases to pkg/task/recurrence_test.go that attempt to use flag-like strings as anchors and verify they are handled safely.

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
