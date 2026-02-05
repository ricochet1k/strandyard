---
type: issue
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:38:56.756967Z
date_edited: 2026-02-05T11:50:16.823304Z
owner_approval: false
completed: true
description: ""
---

# Add todo command to manage checklist items

## Summary


## Summary
Add a command to check off TODO items without marking the task complete.

## Description
`strand complete` fails if TODOs are unchecked, but there is no CLI way to check them.

## Requirements
- `strand todo check <task-id> <index>`
- `strand todo uncheck <task-id> <index>`
- `strand todo list <task-id>` (optional, maybe just show task)

## Completion Report
Replaced with better todo management issue (Tfznrfv)
