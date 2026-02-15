---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-07T06:57:44.764315Z
date_edited: 2026-02-15T04:44:48.319671Z
owner_approval: false
completed: true
status: done
description: ""
---

# web: task create/update bypass TaskDB and skip master list refresh

## Summary
Web API endpoints in `pkg/web/handlers.go` update task metadata directly and do not keep root/free master lists in sync:

- `handleTaskGetOrUpdate` PATCH path sets `t.Meta.Blockers/Blocks/Completed` directly, bypassing TaskDB relationship management and validation.
- `handleTaskCreate` writes new task files + blocker links without regenerating `tasks/root-tasks.md` or `tasks/free-tasks.md`.

This means the dashboard can leave relationships and master lists stale until a manual `strand repair` runs.

Steps to reproduce:
1) Start the web server and create a task with a blocker from the dashboard.
2) Inspect the blocker’s `blocks` list and `tasks/free-tasks.md` — they do not update unless repair is run.

Expected:
- Web endpoints should use TaskDB relationship methods (AddBlocker/RemoveBlocker/SetCompleted/SetParent).
- Master lists should be updated incrementally or via GenerateMasterLists after web writes.

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds

## Completion Report
Confirmed from code-path inspection that web PATCH still writes blockers/blocks directly and web create/update flows do not regenerate root/free master lists. Routed to developer with follow-up task Tvnko7r-web-route-task-create-update-through-taskdb-and-re to implement TaskDB-based relationship updates, master list refresh, and tests.
