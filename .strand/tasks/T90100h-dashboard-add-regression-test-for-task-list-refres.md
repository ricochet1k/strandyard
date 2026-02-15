---
type: implement
role: tester
priority: medium
parent: T39o95t-dashboard-switching-projects-between-projects-that
blockers: []
blocks: []
date_created: 2026-02-15T04:22:18.848063Z
date_edited: 2026-02-15T04:40:20.199738Z
owner_approval: false
completed: true
status: done
description: ""
---

# Dashboard: add regression test for task list refresh when switching projects

## Summary


## Acceptance Criteria
- Implementation matches the specification
- Tests cover the change and pass
- Build succeeds

## Completion Report
Added Playwright regression coverage for project switching: task table now verifies rows refresh to the selected project's tasks, and project-aware API mocks ensure task/role/template payloads switch correctly. Verified with CI=1 npm run test:e2e -- --grep 'switching projects refreshes task list' and npm run build in apps/dashboard.
