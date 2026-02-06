---
type: fix
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T04:38:06.794888Z
date_edited: 2026-02-06T15:43:27.097811Z
owner_approval: false
completed: true
status: done
description: ""
---

# dashboard: When +Add Task is clicked, it should not have parent task filled in

## Summary


## Acceptance Criteria
- Bug still exists
- Bug is fixed and verified locally
- Tests pass
- Build succeeds

## Completion Report
Fixed parent pre-fill bug in dashboard. When clicking '+ Add Task', the parent field is now correctly empty instead of being pre-filled with the currently selected task. Only '+ Add Subtask' button pre-fills the parent. Added addTaskParent signal and openAddTaskModal function to properly track the intended parent task.
