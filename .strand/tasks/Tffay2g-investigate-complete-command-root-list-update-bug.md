---
type: fix
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-15T05:02:33.672855Z
date_edited: 2026-02-15T05:02:33.672855Z
owner_approval: false
completed: false
status: ""
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
