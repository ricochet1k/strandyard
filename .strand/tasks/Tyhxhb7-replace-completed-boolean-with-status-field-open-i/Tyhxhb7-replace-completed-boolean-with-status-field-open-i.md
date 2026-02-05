---
type: issue
role: architect
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T11:50:12.275258Z
date_edited: 2026-02-05T11:50:12.275258Z
owner_approval: false
completed: false
description: ""
---

# Replace completed boolean with status field (open/in_progress/done/duplicate/cancelled)

## Summary
## Summary
Replace the `completed: true/false` boolean with a multi-state `status` field.

## Description
A simple boolean doesn't capture the full lifecycle of a task. Tasks can be successfully completed, duplicated, or explicitly cancelled—each should be tracked distinctly.

## Proposed Status Values
- `open` - ready to work on (default)
- `in_progress` - actively being worked on
- `done` - successfully completed
- `duplicate` - marked as duplicate of another task
- `cancelled` - explicitly cancelled/won't fix
- `blocked` - optionally, waiting on dependencies (implicit in blockers array)

## Benefits
- Better reporting and project visibility
- Eliminates the need for a separate `delete` command (use `duplicate` or `cancelled`)
- Clearer intent for stakeholders

## Migration
- Existing `completed: true` → `status: done`
- Existing `completed: false` → `status: open`
- `strand next` should respect status (only return `open` or `in_progress` tasks)
- `strand complete` should set status to `done` and prompt for status if not done
- New `strand cancel` command to set status to `cancelled`
- New `strand mark-duplicate` command to set status to `duplicate`
