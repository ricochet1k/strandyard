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
date_edited: 2026-02-05T04:08:16.085823Z
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

  Approved recurrence metrics design: Option B (trigger array) for schema, date_completed for task tracking, repeatable --every flag for CLI, and github.com/olebedev/when for date parsing.
