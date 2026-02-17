---
type: implement
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-15T04:44:45.012858Z
date_edited: 2026-02-17T21:37:48.992817Z
owner_approval: false
completed: false
status: ""
description: ""
---

# web: route task create/update through TaskDB and refresh master lists

## Summary


## Problem
The web task endpoints still bypass TaskDB behavior in update paths and do not refresh deterministic master lists after mutations.

## Repro Evidence
- `pkg/web/handlers.go` PATCH writes `t.Meta.Blockers` and `t.Meta.Blocks` directly instead of using TaskDB relationship methods.
- `createTask` persists tasks and relationships but never regenerates `tasks/root-tasks.md` and `tasks/free-tasks.md`.

## Required Changes
1. Update web PATCH handling to diff blocker/blocks updates and apply them via TaskDB APIs (`AddBlocker`/`RemoveBlocker`, `AddBlocked`/`RemoveBlocked`) with ID resolution/validation.
2. Ensure task create/update operations refresh master lists (incremental or `GenerateMasterLists`) so root/free lists stay in sync.
3. Add/adjust web tests to cover blocker link consistency and master list regeneration for create/update flows.
4. Verify with `go test ./...` and `go build ./...`.

## Acceptance Criteria
- Implementation matches the specification
- Tests cover the change and pass
- Build succeeds
