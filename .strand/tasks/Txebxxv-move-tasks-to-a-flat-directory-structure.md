---
type: implement
role: developer
priority: high
parent: T3c1yfi-rename-task-directory-when-title-changes
blockers: []
blocks: []
date_created: 2026-02-06T04:36:51.030861Z
date_edited: 2026-02-06T05:12:29.026379Z
owner_approval: false
completed: true
status: done
description: ""
---

# Move tasks to a flat directory structure.

## Summary
Storing tasks in a nested directory structure means the entire structure has to be searched in order to find a task by ID, there is no good shortcut. Just drop the subdirectory requirement and store everything in the same directory.

## Acceptance Criteria
- Implementation matches the specification
- Tests cover the change and pass
- Build succeeds

## Completion Report
Moved all tasks to a flat directory structure in .strand/tasks/, updated code to support flat storage, and updated documentation.
