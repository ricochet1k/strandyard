---
type: implement
role: developer
priority: medium
parent: Tsbpcrs-git-command-security-hardening
blockers: []
blocks:
    - T29wfxd-review-git-security-hardening
    - Twxmvkr-security-review-of-git-hardening
date_created: 2026-02-01T22:13:09.503114Z
date_edited: 2026-02-01T22:13:09.503114Z
owner_approval: false
completed: false
description: ""
---

# Harden git clone in cmd/init.go

## Summary
## Summary
Use -- separator in git clone for the preset argument to prevent flag injection.

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
