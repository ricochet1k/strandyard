---
type: leaf
role: owner
priority: medium
parent: T8v3k-recurring-tasks
blockers: []
blocks:
    - T8v3k-recurring-tasks
    - Tcb90-document-recurrence-metrics-options
    - Tyvdv-extend-recurrence-schema-and-validation-for-new-me
date_created: 2026-01-28T17:32:15.028035Z
date_edited: 2026-02-01T09:27:56.970136Z
owner_approval: false
completed: false
description: ""
---

# Approve recurrence metrics design

## Context
Decisions for recurrence metrics have been approved and recorded in:
- [design-docs/recurrence-metrics.md](../../../../design-docs/recurrence-metrics.md) (Schema and metrics strategy)
- [design-docs/recurrence-anchor-flags-alternatives.md](../../../../design-docs/recurrence-anchor-flags-alternatives.md) (CLI repeatable --every flag)
- [design-docs/recurrence-anchor-error-messages-alternatives.md](../../../../design-docs/recurrence-anchor-error-messages-alternatives.md) (Standardized error output)
- [design-docs/recurrence-anchor-date-parsing.md](../../../../design-docs/recurrence-anchor-date-parsing.md) (Go `when` library selection)

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: owner) Implement the behavior described in Context.
  Approved recurrence metrics design: Option B (trigger array) for schema, date_completed for task tracking, repeatable --every flag for CLI, and github.com/olebedev/when for date parsing.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
- [ ] (subtask: Tl4cn) short description of subtask
