---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-01-29T22:19:20.346242Z
date_edited: 2026-02-05T01:13:35.597115Z
owner_approval: false
completed: true
description: ""
---

# Investigate date_edited churn on no-op repair

## Summary
Observed task files where only `date_edited` changes after running `strand repair` or similar commands. Identify which command rewrites task files and prevent no-op rewrites when content/metadata is unchanged.

## Context
Provide relevant logs, links, and environment details.

## Impact
Describe severity and who/what is affected.

## Acceptance Criteria
- Define what fixes or mitigations are required.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Completion Report
Identified that SetBody was calling MarkDirty() without checking if the body actually changed. Fixed SetBody in pkg/task/task.go to only mark dirty if content changed. Also updated cmd/edit.go to only report an update and save if the task was actually modified.
