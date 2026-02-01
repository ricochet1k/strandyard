---
type: implement
role: developer
priority: medium
parent: Tsbpcrs-git-command-security-hardening
blockers:
    - T8qdyx5-reliability-review-harden-git-clone-in-cmd-init-go
    - Tm0na28-security-review-harden-git-clone-in-cmd-init-go
    - Tuk75nq-usability-review-harden-git-clone-in-cmd-init-go
blocks:
    - T29wfxd-review-git-security-hardening
    - Twxmvkr-security-review-of-git-hardening
date_created: 2026-02-01T22:13:09.503114Z
date_edited: 2026-02-01T22:31:03.053473Z
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
- [x] (role: tester) Execute test-suite and report failures.
  Executed 'go test ./...' and verified that TestInitWithMaliciousPreset specifically tests for flag injection. All tests passed.
- [x] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
  Delegated specialized reviews to reviewer-reliability (T8qdyx5), reviewer-security (Tm0na28), and reviewer-usability (Tuk75nq).
- [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
- [x] (subtask: T8qdyx5) Description
- [ ] (subtask: Tm0na28) Description
- [ ] (subtask: Tuk75nq) Description
