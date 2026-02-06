---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T07:29:16.741765Z
date_edited: 2026-02-06T07:29:16.741765Z
owner_approval: false
completed: false
status: ""
description: ""
---

# complete command does not update root-tasks.md

## Summary
## Summary
The 'strand complete' command updates 'free-tasks.md' but fails to remove completed tasks from 'root-tasks.md'.

## Steps to Reproduce
1. Find a root task that is currently open.
2. Run 'strand complete <task-id> "some report"'.
3. Check 'tasks/root-tasks.md'.

## Expected Result
The task should be removed from 'root-tasks.md'.

## Actual Result
The task remains in 'root-tasks.md' until a manual 'repair' is run.

## Acceptance Criteria
- 'strand complete' correctly updates both 'free-tasks.md' and 'root-tasks.md'.

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds
