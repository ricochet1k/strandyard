---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T23:54:31Z
date_edited: 2026-01-28T05:09:15.190291Z
owner_approval: false
completed: true
---

# Add progress recording command

## Summary
There is no way to record incremental progress on tasks. Users can only mark entire tasks as completed via `strand complete`, but cannot track progress on individual TODO items within tasks.

## Steps to Reproduce
1. Create a task with multiple TODO items
2. Complete some TODO items manually by editing the markdown file
3. Notice there's no CLI command to mark individual TODOs as complete
4. No way to see overall progress percentage for a task

## Expected Result
A `strand progress` command should exist that can:
- Mark individual TODO items as complete/incomplete
- Show progress percentage for tasks
- List completed vs remaining TODO items
- Update date_edited when progress is recorded

## Actual Result
No progress recording capability exists. Users must manually edit markdown files to track TODO completion.

## Acceptance Criteria
- `strand progress <task-id>` shows current progress (X/Y completed)
- `strand progress <task-id> --item <number>` marks specific TODO as complete
- `strand progress <task-id> --item <number> --undo` marks TODO as incomplete
- Progress updates update date_edited timestamp
- Progress command works with both - [ ] and - [x] checkbox formats
