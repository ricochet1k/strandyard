---
type: task
role: developer
priority: low
parent: T3lpzer-review-task-id-resolution-implementation
blockers: []
blocks: []
date_created: 2026-02-01T23:38:33.596006Z
date_edited: 2026-02-01T23:38:33.596006Z
owner_approval: false
completed: false
description: ""
---

# New Task: Improve error messages for missing task ID anchors

## Description
## Summary
If a task ID is used as an anchor but is not found in the activity log, the system might fallback to date parsing and provide a confusing "invalid date format" error.

## Tasks
- [ ] Ensure that if an anchor looks like a task ID (e.g. starts with T, E, I) but is not found in the activity log, a specific error is returned.

## Acceptance Criteria
- User receives "task ID <id> not found in activity log" instead of "invalid date format" when a task ID anchor fails resolution.

Decide which task template would best fit this task and re-add it with that template and the same parent.
