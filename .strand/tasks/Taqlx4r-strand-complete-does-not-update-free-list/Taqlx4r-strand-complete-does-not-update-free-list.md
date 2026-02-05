---
type: issue
role: triage
priority: high
parent: ""
blockers:
    - Tq49ww0-fix-strand-complete-to-update-free-list-when-last
blocks: []
date_created: 2026-02-04T23:21:02.366797Z
date_edited: 2026-02-05T12:01:26.291278Z
owner_approval: false
completed: true
description: ""
---

# strand complete does not update free-list

## Summary


## Description
When a task is marked as completed using `strand complete`, it remains in the "free" list used by `strand next` until a manual `strand repair` is run.

## Expected Behavior
Any command that modifies task completion status (like `complete`) should automatically update the free-list index.

## Technical Details
The implementation of the complete command should ensure it uses TaskDB APIs that handle free-list maintenance internally, or explicitly trigger a free-list update upon status changes to ensure the cache remains consistent without requiring manual repairs.

## Subtasks
- [ ] (subtask: Tq49ww0) Fix strand complete to update free-list when last todo completes task

## Completion Report
Confirmed issue: when completing the last TODO item in a task, the free-list is not updated. Root cause found in cmd/complete.go:runCompleteTodo() at line 238 - the function returns without updating free-list when result.TaskCompleted is true. Created follow-up implement task (Tq49ww0) to fix this issue.
