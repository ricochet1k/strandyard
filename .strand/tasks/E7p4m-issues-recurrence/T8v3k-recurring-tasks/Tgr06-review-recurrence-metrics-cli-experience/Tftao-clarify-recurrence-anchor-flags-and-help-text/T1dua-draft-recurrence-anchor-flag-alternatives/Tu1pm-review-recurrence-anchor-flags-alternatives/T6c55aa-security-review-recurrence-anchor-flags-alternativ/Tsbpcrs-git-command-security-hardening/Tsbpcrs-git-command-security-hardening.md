---
type: task
role: architect
priority: medium
parent: T6c55aa-security-review-recurrence-anchor-flags-alternativ
blockers:
    - T29wfxd-review-git-security-hardening
blocks: []
date_created: 2026-02-01T22:13:02.705583Z
date_edited: 2026-02-01T23:31:29.790034Z
owner_approval: false
completed: true
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
- [x] (subtask: T29wfxd) Description
- [x] (subtask: To2bn31) Harden git clone in cmd/init.go
- [x] (subtask: Tv1ocqo) Add regression tests for git flag injection in recurrence
- [x] (subtask: Tznw4mk) Harden git command calls in pkg/task/recurrence.go

## Completion Report
Epic complete. All git command hardening tasks implemented and verified. Hardened git clone in cmd/init.go and various git commands in pkg/task/recurrence.go. Comprehensive regression tests added and passed. Security and master reviews complete.
