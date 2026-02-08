---
type: issue
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-08T04:07:15.05556Z
date_edited: 2026-02-08T04:07:15.05556Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Repair mutates tasks after complete should be no-op

## Summary
## Summary
Running `strand repair` immediately after `strand complete` still reports repaired tasks. Per agent policy, `strand complete` should already update master lists and leave the task graph consistent.

## Repro Steps
1. `go run ./cmd/strand complete T06ubsf-consolidate-blocker-relationship-repair "<report>"`
2. `go run ./cmd/strand repair`

## Observed
- Output from complete:
  - `✓ Task T06ubsf marked as completed`
- Output from repair:
  - `repair: ok`
  - `Repaired 1 tasks`

## Expected
- `strand repair` should not need to modify any tasks immediately after `strand complete`.

## Affected Task IDs
- `T06ubsf-consolidate-blocker-relationship-repair`

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds
