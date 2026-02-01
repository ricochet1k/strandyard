---
type: task
role: designer
priority: medium
parent: T3md3-usability-review-recurrence-anchor-flags
blockers:
    - Ti8ig-usability-review-recurrence-anchor-errors
blocks:
    - T3md3-usability-review-recurrence-anchor-flags
date_created: 2026-01-29T05:46:59.618092Z
date_edited: 2026-02-01T20:18:56.987463Z
owner_approval: false
completed: true
description: ""
---

# Specify anchor error messages

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
- Define user-facing error messages for missing or malformed anchors.
- Cover unit/anchor mismatches and ambiguity (when unit implies commit vs date).
- Provide example recovery hints for each error case.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## Subtasks
- [x] (subtask: Tcsz3) Review alternatives: recurrence anchor error messages
- [x] (subtask: Thy9m) Reliability review: recurrence anchor errors
- [x] (subtask: Ti8ig) Usability review: recurrence anchor errors

## Completion Report
Design complete: Adopted Alternative B unified error format with structured reason + hint line. Specified error messages for all failure modes (malformed anchors, unit/anchor mismatches, ambiguity, invalid metrics/amounts, missing anchors). Defined canonical examples from hint-examples doc. Documented anchor type mapping, implementation notes, and local verification steps. Open concerns captured: CLI.md needs --every flag update, exit code convention alignment.
