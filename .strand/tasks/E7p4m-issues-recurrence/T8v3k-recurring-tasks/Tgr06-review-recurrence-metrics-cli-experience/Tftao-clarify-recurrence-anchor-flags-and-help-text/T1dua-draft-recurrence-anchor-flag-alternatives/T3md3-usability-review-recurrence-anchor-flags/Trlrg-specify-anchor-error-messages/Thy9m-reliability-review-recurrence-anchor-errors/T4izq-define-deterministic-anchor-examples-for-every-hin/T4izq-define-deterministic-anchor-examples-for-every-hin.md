---
type: task
role: architect
priority: medium
parent: Thy9m-reliability-review-recurrence-anchor-errors
blockers:
    - Tm2sq-review-canonical-every-hint-examples-implementatio
    - Tm6qi-document-canonical-every-hint-examples
    - Tqb9o-approve-canonical-every-hint-examples
    - Tsyeo-review-canonical-every-hint-examples-docs
    - Tv4cw-implement-deterministic-every-hint-examples
blocks:
    - Thy9m-reliability-review-recurrence-anchor-errors
date_created: 2026-01-29T19:20:50.103845Z
date_edited: 2026-01-31T17:29:31.078586Z
owner_approval: false
completed: true
---

# Define deterministic anchor examples for --every hints

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description

## Summary
Define canonical, deterministic anchor examples for each anchor type in `--every` hint lines.

## Details
- Specify fixed examples for date/time anchors (human-friendly + optional ISO 8601).
- Specify fixed examples for commit anchors (e.g., `HEAD` or a placeholder hash) that do not vary per run.
- Ensure examples are stable for tests and automation (no current-time rendering).

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## Subtasks
- [x] (subtask: Tm2sq) Review canonical --every hint examples (implementation)
- [x] (subtask: Tm6qi) Document canonical --every hint examples
- [x] (subtask: Tqb9o) Approve canonical --every hint examples
- [x] (subtask: Tsyeo) Review canonical --every hint examples (docs)
- [ ] (subtask: Tv4cw) Implement deterministic --every hint examples
