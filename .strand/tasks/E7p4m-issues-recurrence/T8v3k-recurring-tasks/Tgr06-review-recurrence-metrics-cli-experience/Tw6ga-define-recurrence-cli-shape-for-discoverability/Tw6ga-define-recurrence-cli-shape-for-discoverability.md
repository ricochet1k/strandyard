---
type: task
role: designer
priority: medium
parent: Tgr06-review-recurrence-metrics-cli-experience
blockers: []
blocks:
    - Tgr06-review-recurrence-metrics-cli-experience
date_created: 2026-01-28T19:01:01.684537Z
date_edited: 2026-02-05T00:59:16.172779Z
owner_approval: false
completed: true
description: ""
---

# Define recurrence CLI shape for discoverability

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## Subtasks
- [ ] (subtask: Tw6ga) short description of subtask

## Completion Report
Decision: Adopted Alternative D (--every flag with structured string parsing). This provides a single, discoverable entry point for all recurrence triggers while maintaining a clean CLI surface. Implemented in cmd/add.go and cmd/edit.go with scannable help text and examples.
