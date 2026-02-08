---
type: issue
role: triage
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T00:08:55.661003Z
date_edited: 2026-02-08T04:06:37.994519Z
owner_approval: false
completed: true
status: done
description: ""
---

# Allow reapplying templates to existing tasks

## Summary
When templates evolve (new sections, structure, or fields), there is no way to push those changes to tasks that were already created. Add a command or flag that lets us re-apply a template to an existing task so it can inherit the latest structure without losing existing content, optionally opening a merge-style prompt or supporting dry runs. This keeps documentation consistent when the template author makes improvements.

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds

## Completion Report
Confirmed the feature request and created a design task (Tk259ie) for the designer role.

## Subtasks
- [ ] (subtask: Tk259ie) Design: Allow reapplying templates to existing tasks
