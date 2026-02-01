---
type: task
role: architect
priority: medium
parent: Tnhk8ef-usability-review-of-activity-log-concurrency-plan
blockers: []
blocks: []
date_created: 2026-02-01T23:49:14.062985Z
date_edited: 2026-02-01T23:52:50.034332Z
owner_approval: false
completed: true
description: ""
---

# Description

## Description
The function `EvaluateTasksCompletedMetric` in `pkg/task/recurrence.go` takes an `*activity.Log` argument but currently ignores it and opens its own log instance. It should be refactored to reuse the passed-in log if it's not nil, to ensure it benefits from the concurrency protections (mutex) of that instance and to avoid redundant file handles and I/O.

Decide which task template would best fit this task and re-add it with that template and the same parent.

## Completion Report
Created implementation plan (design-docs/reuse-log-instance-plan.md) and broken down into implementation task Tu4i36m and review task T0fielf.
