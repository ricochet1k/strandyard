---
type: issue
role: developer
priority: low
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:40:06.154877Z
date_edited: 2026-02-15T05:12:33.909204Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Rename task directory when title changes

## Summary


## Summary
Update the task slug (and directory/filename) when the title changes.

## Description
`strand edit --title` updates the H1 title in the file but the directory name (and thus the Task ID slug) remains unchanged.

## Requirements
- Option (or default behavior) to rename directory/file to match new title slug.
- Update all references (parent, blockers, blocks) in other files if the ID changes (slug change = ID change).
- Alternatively, provide a dedicated `strand rename` command.

## Subtasks
- [x] (subtask: Txebxxv) Move tasks to a flat directory structure.
  Moved all tasks to a flat directory structure in .strand/tasks/, updated code to support flat storage, and updated documentation.
