---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-07T05:59:39.686068Z
date_edited: 2026-02-15T04:27:45.210276Z
owner_approval: false
completed: true
status: done
description: ""
---

# Clarify strand complete --todo usage and role validation error

## Summary


## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds

## Subtasks
- [ ] (subtask: Tfuy2lo) Investigate strand complete --todo indexing and role validation messaging

## Completion Report
Confirmed  indexing/message confusion. Reproduced with  showing role mismatch plus renumbered incomplete TODO list (1..N) that conflicts with provided absolute index. Created follow-up developer task Tfuy2lo to add tests and clarify error/help/docs.
