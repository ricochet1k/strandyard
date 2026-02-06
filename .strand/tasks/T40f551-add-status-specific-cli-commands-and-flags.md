---
type: implement
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T05:05:35.90684Z
date_edited: 2026-02-06T05:05:35.90684Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Add status-specific CLI commands and flags

## Summary
## Summary
Update CLI commands to expose status field functionality to users.

## Deliverables
- Add strand cancel <task-id> [reason] command.
- Add strand mark-duplicate <task-id> <duplicate-of> command.
- Add strand mark-in-progress <task-id> command.
- Add --status flag to strand edit.
- Add --status flag to strand list for filtering.
- Update strand list columns to include status by default or as an option.

## Acceptance Criteria
- Implementation matches the specification
- Tests cover the change and pass
- Build succeeds
