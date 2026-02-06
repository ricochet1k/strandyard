---
type: issue
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:38:56.587564Z
date_edited: 2026-02-06T06:10:29.442179Z
owner_approval: false
completed: true
status: done
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

## Completion Report
Resolved by implementing status field support. Tasks can now be reopened by setting status to 'open' via 'strand edit --status open' or 'strand mark-in-progress'.
