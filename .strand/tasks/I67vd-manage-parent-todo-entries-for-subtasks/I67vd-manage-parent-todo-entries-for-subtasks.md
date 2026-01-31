---
type: issue
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-01-29T22:26:35.930723Z
date_edited: 2026-01-29T22:26:35.930723Z
owner_approval: false
completed: false
---

# Manage parent TODO entries for subtasks

## Summary


## Summary
Ensure parent tasks keep a permanent TODO list of their subtasks. When a subtask is created, add a TODO entry in the parent task's `## Tasks` section; when the subtask is completed, check off the corresponding entry. The task library should own this section so it stays consistent and deterministic.

## Acceptance Criteria
- Parent task `## Tasks` includes entries for all subtasks
- Completed subtasks are checked off in the parent TODO list
- Updates are deterministic and require no manual edits

## Context
Provide relevant logs, links, and environment details.

## Impact
Describe severity and who/what is affected.

## Acceptance Criteria
- Define what fixes or mitigations are required.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## TODOs
- [ ] Add parent TODO entry when creating a subtask via `strand add --parent`
- [ ] Check off parent TODO entry when completing a subtask via `strand complete`
- [ ] Preserve non-subtask TODO items and keep deterministic ordering
