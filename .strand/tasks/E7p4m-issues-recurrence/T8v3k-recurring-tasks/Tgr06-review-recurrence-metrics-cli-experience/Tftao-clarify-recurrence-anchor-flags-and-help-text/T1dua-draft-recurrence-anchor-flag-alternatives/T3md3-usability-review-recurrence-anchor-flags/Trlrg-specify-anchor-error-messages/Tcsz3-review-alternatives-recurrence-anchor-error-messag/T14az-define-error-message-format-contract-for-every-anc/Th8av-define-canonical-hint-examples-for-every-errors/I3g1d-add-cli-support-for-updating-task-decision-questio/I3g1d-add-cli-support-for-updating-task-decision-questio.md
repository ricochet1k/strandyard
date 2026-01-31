---
type: issue
role: triage
priority: low
parent: Th8av-define-canonical-hint-examples-for-every-errors
blockers: []
blocks:
    - Th8av-define-canonical-hint-examples-for-every-errors
date_created: 2026-01-29T17:26:24.419964Z
date_edited: 2026-01-29T10:26:24.43021-07:00
owner_approval: false
completed: false
---

# Add CLI support for updating task decision/question sections

## Summary


## Summary
We manually updated a task section label from "Decisions Needed" to "Questions Needed". Provide a CLI command to update common task section headings to avoid manual edits.

## Repro
1. Observe tasks may include a "Decisions Needed" heading that should be "Questions Needed" per AGENTS policy.
2. There is no CLI command to update such headings across tasks.

## Acceptance Criteria
- New command (or existing command extension) can update section headings in task files without manual edits.
- Running the command updates master lists deterministically.

## Context
Provide relevant logs, links, and environment details.

## Impact
Describe severity and who/what is affected.

## Acceptance Criteria
- Define what fixes or mitigations are required.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.
