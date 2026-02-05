---
type: issue
role: triage
priority: medium
parent: ""
blockers:
    - Tnb5ir6-design-task-storage-migration
blocks: []
date_created: 2026-01-30T02:22:43.426796Z
date_edited: 2026-02-05T01:10:03.891355Z
owner_approval: false
completed: true
description: ""
---

# Add CLI support for task storage migration

## Summary
We manually moved tasks/roles/templates into .strand for local storage initialization. Add a first-class command to migrate existing repos without manual moves (including updating git-tracked paths and master lists).

## Context
Provide relevant logs, links, and environment details.

## Impact
Describe severity and who/what is affected.

## Acceptance Criteria
- Define what fixes or mitigations are required.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Completion Report
Confirmed that manual migration was used and a CLI command is needed. Created a design task Tnb5ir6 to define the migration process and handle Git-related concerns.

## Subtasks
- [x] (subtask: T6y8964) Design task storage migration
- [x] (subtask: Tnb5ir6) New Task: Design task storage migration
