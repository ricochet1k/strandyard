---
type: task
role: architect
priority: medium
parent: T6c55aa-security-review-recurrence-anchor-flags-alternativ
blockers:
    - T29wfxd-review-git-security-hardening
    - To2bn31-harden-git-clone-in-cmd-init-go
    - Tv1ocqo-add-regression-tests-for-git-flag-injection-in-rec
    - Tznw4mk-harden-git-command-calls-in-pkg-task-recurrence-go
blocks: []
date_created: 2026-02-01T22:13:02.705583Z
date_edited: 2026-02-01T22:13:19.326977Z
owner_approval: false
completed: false
description: ""
---

# New Task: Git Command Security Hardening

## Description


## Summary
Epic for hardening all git command executions against flag injection.
Links to: [Implementation Plan — Git Command Security Hardening](design-docs/git-command-security-hardening.md)

## Tracks
- **Hardening**: Implementation of -- separators and validation.
- **Verification**: Regression tests and manual verification.

Decide which task template would best fit this task and re-add it with that template and the same parent.

## Subtasks
- [ ] (subtask: T29wfxd) Description
- [ ] (subtask: To2bn31) Harden git clone in cmd/init.go
- [ ] (subtask: Tv1ocqo) Add regression tests for git flag injection in recurrence
- [ ] (subtask: Tznw4mk) Harden git command calls in pkg/task/recurrence.go
