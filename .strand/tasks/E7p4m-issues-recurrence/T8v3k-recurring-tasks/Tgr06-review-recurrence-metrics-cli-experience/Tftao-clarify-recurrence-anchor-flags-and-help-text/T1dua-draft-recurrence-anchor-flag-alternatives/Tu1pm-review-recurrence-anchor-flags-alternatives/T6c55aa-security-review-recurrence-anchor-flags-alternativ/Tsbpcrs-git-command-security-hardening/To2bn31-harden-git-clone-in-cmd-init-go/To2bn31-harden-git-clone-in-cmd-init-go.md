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
date_edited: 2026-02-01T22:33:45.505563Z
owner_approval: false
completed: true
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
- [x] (role: documentation) Update user-facing docs and examples.
  Updated CLI.md to clarify that --preset supports secure cloning of git repositories. Verified hardening manually by confirming that strand init --preset="--help" correctly identifies "--help" as a repository path rather than a flag.

## Subtasks
- [x] (subtask: T8qdyx5) Description
- [x] (subtask: Tm0na28) Description
- [x] (subtask: Tuk75nq) Description
