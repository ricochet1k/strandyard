---
type: task
role: developer
priority: medium
parent: T5gkbnm-clarify-strand-complete-todo-usage-and-role-valida
blockers: []
blocks: []
date_created: 2026-02-15T04:27:40.117955Z
date_edited: 2026-02-15T04:27:45.210325Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Investigate strand complete --todo indexing and role validation messaging

## Summary


## Summary
`strand complete --todo` accepts an absolute TODO index (all TODOs in file), but role validation errors print a renumbered list of only incomplete TODOs. This is confusing when earlier TODOs are already checked.

## Reproduction
1. Run `go run ./cmd/strand complete T6jry --todo 5 --role developer`.
2. Observe error:
   - `role validation failed: todo has role 'master-reviewer' but --role flag specifies 'developer'`
   - Followed by `Incomplete todos:` list numbered `1`, `2` only.
3. User passed `--todo 5`, but output implies available TODO numbers are `1` and `2`.

## Expected
Error/help text should clearly state how `--todo` is indexed and avoid ambiguous numbering in role validation output.

## Acceptance Criteria
- Repro is covered by automated test(s).
- Error output makes TODO index semantics explicit.
- `--help` and/or CLI docs clarify `--todo` indexing semantics.
- `go test ./...` passes.
- `go build ./...` succeeds.

## Instructions
Decide which task template would best fit this task and re-add it with that template and the same parent.
