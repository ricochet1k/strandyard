---
type: implement
role: developer
priority: medium
parent: Tsbpcrs-git-command-security-hardening
blockers: []
blocks:
    - T29wfxd-review-git-security-hardening
    - Twxmvkr-security-review-of-git-hardening
date_created: 2026-02-01T22:13:06.506931Z
date_edited: 2026-02-01T22:13:06.506931Z
owner_approval: false
completed: false
description: ""
---

# Harden git command calls in pkg/task/recurrence.go

## Summary
## Summary
Inject -- separator before positional arguments in git rev-list, diff, rev-parse, and show to prevent flag injection.

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
