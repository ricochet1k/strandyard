---
type: fix
role: developer
priority: medium
parent: T5gkbnm-clarify-strand-complete-todo-usage-and-role-valida
blockers: []
blocks: []
date_created: 2026-02-15T04:54:00.452326Z
date_edited: 2026-02-15T04:54:00.452326Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Investigate strand complete --todo indexing and role validation messaging

## Summary
`strand complete --todo` currently accepts an absolute TODO index across all TODOs in the task file, but role validation errors currently print a renumbered list of incomplete TODOs only.

This creates confusing output when a user passes a higher absolute index (for example `--todo 5`) and then sees an `Incomplete todos:` list numbered `1`, `2`.

## Reproduction
1. Run `go run ./cmd/strand complete T6jry --todo 5 --role developer`.
2. Observe role mismatch error and an incomplete TODO list renumbered from 1.

## Expected
Clarify index semantics in error/help/docs so users understand `--todo` is absolute and not the incomplete-list position.

## Acceptance Criteria
- Repro is covered by automated test(s).
- Error output makes TODO index semantics explicit.
- `--help` and/or CLI docs clarify `--todo` indexing semantics.
- `go test ./...` passes.
- `go build ./...` succeeds.

## Acceptance Criteria
- Bug still exists
- Bug is fixed and verified locally
- Tests pass
- Build succeeds
