---
type: task
role: architect
priority: medium
parent: Iirhp-add-garbage-collect-option-to-delete-command
blockers: []
blocks: []
date_created: 2026-02-06T04:53:45.169725Z
date_edited: 2026-02-06T04:53:45.169725Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Design gc command for status-based garbage collection

## Summary
The 'delete' command is currently missing, and there is a project decision (Tj8om9e) to use status states (done, cancelled, duplicate) instead of hard deletes. 

However, there is still a need for garbage collection to keep the .strand/tasks directory manageable.

Design a 'gc' command that:
- Targets tasks with status 'done', 'cancelled', or 'duplicate'.
- Supports an '--age' flag (e.g. 30 days).
- Safely removes task directories and cleans up references.
- Reconciles with the decision to avoid hard deletes for individual tasks.

## Instructions
Decide which task template would best fit this task and re-add it with that template and the same parent.
