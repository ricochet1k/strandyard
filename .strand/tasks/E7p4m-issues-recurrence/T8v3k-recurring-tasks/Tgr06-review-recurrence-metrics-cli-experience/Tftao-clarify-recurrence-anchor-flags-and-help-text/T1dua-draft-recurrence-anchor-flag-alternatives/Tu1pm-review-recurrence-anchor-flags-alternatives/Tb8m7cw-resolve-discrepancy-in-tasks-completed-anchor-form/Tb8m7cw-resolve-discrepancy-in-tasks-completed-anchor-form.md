---
type: issue
role: architect
priority: medium
parent: Tu1pm-review-recurrence-anchor-flags-alternatives
blockers:
    - T3lpzer-review-task-id-resolution-implementation
    - Tsnkyb7-implement-task-id-resolution-for-tasks-completed-m
blocks: []
date_created: 2026-02-01T22:00:36.683234Z
date_edited: 2026-02-01T23:36:24.402671Z
owner_approval: false
completed: true
description: ""
---

# Resolve discrepancy in tasks_completed anchor format documentation

## Summary
There is a discrepancy between the documentation (`CLI.md`, `recurrence_metrics.md`) and the implementation (`pkg/task/recurrence.go`) regarding `tasks_completed` anchors. The documentation claims task IDs are supported as anchors, but the implementation only supports dates.

I have created an implementation plan in `design-docs/tasks-completed-anchor-resolution.md` to fully support task ID anchors by querying the activity log. This approach is more robust and aligns with the existing documentation and validation logic.

I have created two subtasks:
1. `T3lpzer`: Review the implementation plan (assigned to `reviewer-reliability`).
2. `Tsnkyb7`: Implement the changes (assigned to `developer`).

## Subtasks
- [ ] (subtask: T3lpzer) Review task ID resolution implementation plan
- [ ] (subtask: Tsnkyb7) Implement task ID resolution for tasks_completed metric

## Completion Report
Created implementation plan in design-docs/tasks-completed-anchor-resolution.md and created subtasks for implementation and review.
