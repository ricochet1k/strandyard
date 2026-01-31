---
type: task
role: designer
priority: medium
parent: T3md3-usability-review-recurrence-anchor-flags
blockers:
    - Tcsz3-review-alternatives-recurrence-anchor-error-messag
    - Thy9m-reliability-review-recurrence-anchor-errors
    - Ti8ig-usability-review-recurrence-anchor-errors
blocks:
    - T3md3-usability-review-recurrence-anchor-flags
date_created: 2026-01-29T05:46:59.618092Z
date_edited: 2026-01-31T04:41:33.150443Z
owner_approval: false
completed: false
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

## TODOs
- [ ] (role: developer) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
- [x] (subtask: Tcsz3) Review alternatives: recurrence anchor error messages
- [x] (subtask: Thy9m) Reliability review: recurrence anchor errors
- [ ] (subtask: Ti8ig) Usability review: recurrence anchor errors
