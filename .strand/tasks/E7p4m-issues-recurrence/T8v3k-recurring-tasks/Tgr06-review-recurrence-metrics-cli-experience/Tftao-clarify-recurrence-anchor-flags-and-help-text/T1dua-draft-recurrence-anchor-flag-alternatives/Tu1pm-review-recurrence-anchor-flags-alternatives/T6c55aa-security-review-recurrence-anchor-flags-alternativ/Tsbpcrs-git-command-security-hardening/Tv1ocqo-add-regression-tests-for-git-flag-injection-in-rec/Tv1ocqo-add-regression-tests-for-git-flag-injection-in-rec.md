---
type: implement
role: developer
priority: medium
parent: Tsbpcrs-git-command-security-hardening
blockers:
    - Thxloni-security-review-regression-tests-for-git-flag-inje
    - Tilakbz-reliability-review-regression-tests-for-git-flag-i
    - Tyfv8sz-usability-review-regression-tests-for-git-flag-inj
blocks:
    - T29wfxd-review-git-security-hardening
    - Twxmvkr-security-review-of-git-hardening
date_created: 2026-02-01T22:13:15.54608Z
date_edited: 2026-02-01T23:28:00.240373Z
owner_approval: false
completed: true
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
- [x] (role: developer) Implement the behavior described in Context.
  Added security-focused regression tests in pkg/task/recurrence_security_test.go covering flag injection in EvaluateGitMetric, ResolveGitHash, and GetCommitAtOffset.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Verified tests against hardened implementation.
- [x] (role: tester) Execute test-suite and report failures.
  Executed 'go test -v ./pkg/task -run FlagInjection' and verified all security tests pass.
- [x] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
  Delegated specialized reviews to reviewer-reliability (Tilakbz), reviewer-security (Thxloni), and reviewer-usability (Tyfv8sz).
- [x] (role: documentation) Update user-facing docs and examples.
  No user-facing documentation updates required for internal security tests.

## Subtasks
- [x] (subtask: Thxloni) Description
- [x] (subtask: Tilakbz) Description
- [x] (subtask: Tyfv8sz) Description
