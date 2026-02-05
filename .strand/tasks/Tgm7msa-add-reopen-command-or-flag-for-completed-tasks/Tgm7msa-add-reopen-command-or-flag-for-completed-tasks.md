---
type: issue
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:38:56.587564Z
date_edited: 2026-02-05T11:50:35.259503Z
owner_approval: false
completed: false
description: ""
---

# Simplify with status states: no separate reopen needed

## Summary
With status states, reopening a task is just setting `status: open` instead of `status: done`.

## Description
Once status states are implemented (Tyhxhb7), there is no need for a dedicated reopen command:
- Set `status: open` to reopen a task
- Set `status: in_progress` if someone is actively working on it
- Use `strand edit --status <value>` to change status

## Rationale
Status states provide a cleaner abstraction than reverting a boolean flag. This task can be completed (and removed from backlog) once Tyhxhb7 is implemented.
