---
type: task
role: owner
priority: medium
parent: T14az-define-error-message-format-contract-for-every-anc
blockers:
    - I3g1d-add-cli-support-for-updating-task-decision-questio
    - Tm91y-decide-every-anchor-defaults-and-hint-examples
    - Tmxs6-evaluate-date-parsing-library-for-every-anchors
blocks:
    - T14az-define-error-message-format-contract-for-every-anc
date_created: 2026-01-29T16:55:20.213085Z
date_edited: 2026-01-31T04:41:33.150446Z
owner_approval: false
completed: false
---

# Define canonical hint examples for --every errors

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description


## Summary
Select canonical example anchors and full `--every` examples for each unit/metric to use in hint lines.

## Acceptance Criteria
- One canonical example per unit/metric is specified.
- Examples avoid non-deterministic content.

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
- [ ] (subtask: I3g1d) Add CLI support for updating task decision/question sections
- [x] (subtask: Tm91y) Questions: --every anchor defaults and hint examples
- [x] (subtask: Tmxs6) Evaluate date parsing library for --every anchors
