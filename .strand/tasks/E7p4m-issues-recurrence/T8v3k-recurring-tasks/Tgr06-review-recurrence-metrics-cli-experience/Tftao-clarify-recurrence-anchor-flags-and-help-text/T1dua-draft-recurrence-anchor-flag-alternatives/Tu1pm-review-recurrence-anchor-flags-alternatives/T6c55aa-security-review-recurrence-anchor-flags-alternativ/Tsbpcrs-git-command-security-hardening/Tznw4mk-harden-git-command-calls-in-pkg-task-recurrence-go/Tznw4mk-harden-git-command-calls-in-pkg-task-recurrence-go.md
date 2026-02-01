---
type: implement
role: developer
priority: medium
parent: Tsbpcrs-git-command-security-hardening
blockers:
    - T3p1fkv-security-review-harden-git-command-calls-in-pkg-ta
    - Torj7f1-reliability-review-harden-git-command-calls-in-pkg
    - Tu52hy0-usability-review-harden-git-command-calls-in-pkg-t
blocks:
    - T29wfxd-review-git-security-hardening
    - Twxmvkr-security-review-of-git-hardening
date_created: 2026-02-01T22:13:06.506931Z
date_edited: 2026-02-01T23:24:32.939103Z
owner_approval: false
completed: true
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
- [x] (role: developer) Implement the behavior described in Context.
  Hardened git command calls in pkg/task/recurrence.go using --end-of-options for rev-parse and resolving anchors to hashes for other commands to prevent flag injection.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Verified main flows with existing and new tests.
- [x] (role: tester) Execute test-suite and report failures.
  Executed 'go test ./...' and verified all tests pass, including new security tests.
- [x] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
  Delegated specialized reviews to reviewer-reliability (Torj7f1), reviewer-security (T3p1fkv), and reviewer-usability (Tu52hy0).
- [x] (role: documentation) Update user-facing docs and examples.
  No user-facing documentation updates required for this internal hardening. Verification confirmed hardening is transparent.

## Subtasks
- [x] (subtask: T3p1fkv) Description
- [x] (subtask: Torj7f1) Description
- [x] (subtask: Tu52hy0) Description
