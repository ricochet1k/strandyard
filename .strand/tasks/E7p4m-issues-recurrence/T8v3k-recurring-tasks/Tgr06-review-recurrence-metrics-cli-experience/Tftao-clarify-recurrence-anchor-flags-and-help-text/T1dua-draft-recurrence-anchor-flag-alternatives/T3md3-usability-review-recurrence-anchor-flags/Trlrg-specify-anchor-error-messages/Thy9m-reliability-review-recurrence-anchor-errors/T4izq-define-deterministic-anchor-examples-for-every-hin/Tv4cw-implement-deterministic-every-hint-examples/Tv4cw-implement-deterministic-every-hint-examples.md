---
type: task
role: developer
priority: medium
parent: T4izq-define-deterministic-anchor-examples-for-every-hin
blockers: []
blocks:
    - T4izq-define-deterministic-anchor-examples-for-every-hin
date_created: 2026-01-29T19:24:46.390895Z
date_edited: 2026-02-01T15:56:45.27689Z
owner_approval: false
completed: true
description: ""
---

# Implement deterministic --every hint examples

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description


## Summary
Implement deterministic --every hint examples using the canonical anchors defined in design-docs/recurrence-anchor-hint-examples.md.

## Acceptance Criteria
- Hint lines match canonical examples exactly
- Tests cover date/time and commit anchor hints

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## Completion Report
Implemented deterministic --every hint examples with canonical examples and validation

- Added --every flag to add command with validation
- Implemented validateEvery function using canonical hint examples from design-docs/recurrence-anchor-hint-examples.md
- All hint examples match the deterministic specification exactly
- Added comprehensive tests covering valid and invalid --every flag inputs
- Implementation exits with code 2 for validation failures per design contract
