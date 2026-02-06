---
type: task
role: designer
priority: medium
parent: Ibv0j-create-delete-task-command
blockers: []
blocks: []
date_created: 2026-02-05T01:06:24.961814Z
date_edited: 2026-02-05T01:07:44.61865Z
owner_approval: false
completed: true
description: ""
---

# New Task: Design delete task command

## Description
Design a CLI command to safely delete tasks. Deletion should handle:
- Removing the task directory and its contents.
- Cleaning up references in parent tasks (TODO lists).
- Cleaning up references in other tasks (blockers/blocks).
- Updating master lists (root-tasks.md, free-tasks.md).
- Handling hierarchical deletion (should deleting a parent delete children?).

Deliverable: A design document in `design-docs/` outlining the command syntax, behavior, and safety considerations.

Decide which task template would best fit this task and re-add it with that template and the same parent.

## Completion Report
Triage complete. Re-added as designer task Tkoz0da under the same parent.
