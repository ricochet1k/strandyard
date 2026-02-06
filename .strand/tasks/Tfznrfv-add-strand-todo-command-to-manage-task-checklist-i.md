---
type: issue
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T11:50:14.209489Z
date_edited: 2026-02-06T05:54:18.187301Z
owner_approval: false
completed: true
status: done
description: ""
---

# Add strand todo command to manage task checklist items

## Summary


## Summary
Add a command to add, remove, edit, and reorder TODO items within a task.

## Description
Currently, `--todo` flag on `complete` can only check off existing items. There is no CLI way to:
- Add new todo items
- Remove/delete todo items
- Edit/reword existing todos
- Reorder todos

## Requirements
- `strand todo add <task-id> "item text"` - append a todo
- `strand todo remove <task-id> <index>` - delete a todo
- `strand todo edit <task-id> <index> "new text"` - edit a todo
- `strand todo check <task-id> <index>` - check off a todo (alias for `complete --todo`)
- `strand todo uncheck <task-id> <index>` - uncheck a todo
- Alternatively, extend `strand edit` with `--todo-add`, `--todo-remove`, `--todo-edit` flags

## Acceptance Criteria
- All operations preserve checked/unchecked state of other items
- Index is 1-based (matching `complete --todo`)
- Provide clear error if task or index not found

## Completion Report
Implemented 'strand todo' command with add, remove, edit, check, uncheck, reorder, and list subcommands. Added corresponding methods to TaskDB and unit tests. Verified all operations work as expected and preserve task state.
