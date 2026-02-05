---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-01-29T22:19:58.257252Z
date_edited: 2026-02-05T01:08:34.213217Z
owner_approval: false
completed: true
description: ""
---

# Add CLI support for editing task summary/body

## Summary
Manual edits to task markdown are sometimes needed to fix mistakes, but policy requires CLI support. Add a command or flag to update task summary/body safely without hand-editing task files.

## Context
Provide relevant logs, links, and environment details.

## Impact
Describe severity and who/what is affected.

## Acceptance Criteria
- Define what fixes or mitigations are required.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Completion Report
Issue resolved. The 'edit' command already supports updating the task body via stdin. Users can use heredocs or pipes to provide a new body, which will update all sections except for the protected ones like TODOs and Subtasks.
