---
type: fix
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-15T05:02:33.672855Z
date_edited: 2026-02-15T05:05:01.054252Z
owner_approval: false
completed: true
status: done
description: ""
---

# Investigate complete command root list update bug

## Summary
Investigate why `strand complete` does not update `root-tasks.md` as expected.

## Instructions
- Reproduce the root list update issue triggered by `strand complete`
- Identify the root cause in the completion/list regeneration flow
- Implement and verify a fix

## Acceptance Criteria
- Bug still exists
- Bug is fixed and verified locally
- Tests pass
- Build succeeds

## Completion Report
Identified that complete only updated free-tasks incrementally and skipped root-tasks updates. Added root list regeneration after successful incremental completion paths and added a regression test that verifies completed root tasks are removed from root-tasks.md. Verified with go build ./..., go test ./..., and strand repair.
