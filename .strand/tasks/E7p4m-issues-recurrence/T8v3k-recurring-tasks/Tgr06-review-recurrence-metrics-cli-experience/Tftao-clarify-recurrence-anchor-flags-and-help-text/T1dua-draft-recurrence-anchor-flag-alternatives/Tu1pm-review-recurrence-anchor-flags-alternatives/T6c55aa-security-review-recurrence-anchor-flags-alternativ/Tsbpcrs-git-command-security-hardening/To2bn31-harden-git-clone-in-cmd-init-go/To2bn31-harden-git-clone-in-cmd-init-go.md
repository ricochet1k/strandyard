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
date_edited: 2026-02-01T22:25:50.674738Z
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
- [x] (role: developer) Implement the behavior described in Context.
  Implemented -- separator in git clone call in cmd/init.go and added regression tests in cmd/init_test.go.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Added regression tests in cmd/init_test.go covering malicious and normal presets.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
