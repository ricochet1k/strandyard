---
type: leaf
role: reviewer-reliability
priority: medium
parent: T8v3k-recurring-tasks
blockers: []
blocks:
    - T8v3k-recurring-tasks
    - Tyvdv-extend-recurrence-schema-and-validation-for-new-me
date_created: 2026-01-28T17:32:19.22925Z
date_edited: 2026-02-05T00:59:50.487827Z
owner_approval: false
completed: true
description: ""
---

# Review recurrence metrics schema

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## Subtasks
- [ ] (subtask: Thnhh) short description of subtask

## Completion Report
Reliability review complete. Verdict: Approved. The recurrence metrics schema using 'every: []string' with a structured grammar balances simplicity and flexibility. It correctly supports multiple triggers per task. The implementation ensures reliability through strict validation in the CLI (add/edit) and by resolving anchors (short IDs, Git hashes) to their canonical forms before storage. This avoids ambiguity and makes task files self-contained. The 'after' vs 'from' logic prevents drift in time-based recurrence.
