---
type: task
role: designer
priority: medium
parent: Tsyeo-review-canonical-every-hint-examples-docs
blockers: []
blocks:
    - Tsyeo-review-canonical-every-hint-examples-docs
date_created: 2026-01-29T19:54:15.970672Z
date_edited: 2026-02-01T09:43:49.512189Z
owner_approval: false
completed: true
description: ""
---

# Clarify default anchor wording for commit metrics

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description


## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## Subtasks
- [x] (subtask: Iy8y6) Flag non-ASCII quotes in task bodies
- [x] (subtask: Tocc0) Draft alternatives for commit-metric default anchor wording

## Completion Report
Decision: Alternative B (use 'from HEAD' for commit-based metrics). Behavior alignment: if HEAD is invalid or unborn, commit-based recurrence metrics are ignored.
