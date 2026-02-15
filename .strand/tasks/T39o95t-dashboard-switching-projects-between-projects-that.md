---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-04T23:23:39.42562Z
date_edited: 2026-02-15T04:40:20.20304Z
owner_approval: false
completed: true
status: done
description: ""
---

# Dashboard: switching projects between projects that have issues doesn't correctly update the list of tasks

This might have something to do with MyTransitionGroup. CSS styles/classes might not be correct for the transition group, they are subtle.



## Deliverables
Make sure tests check this.

## Completion Report
Could not reproduce stale task-list behavior on current branch. Verified project-switch e2e scenario passes, then created follow-up implement task T90100h to add explicit regression coverage that task rows refresh when changing projects.

## Subtasks
- [x] (subtask: T90100h) Dashboard: add regression test for task list refresh when switching projects
