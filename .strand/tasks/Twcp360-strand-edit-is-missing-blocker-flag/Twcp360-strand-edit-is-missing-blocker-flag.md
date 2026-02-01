---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-01T22:22:08.891654Z
date_edited: 2026-02-01T22:22:08.891654Z
owner_approval: false
completed: false
description: ""
---

# strand edit is missing --blocker flag

## Summary
## Summary
The `strand edit` command help text mentions it can edit blockers, but the `--blocker` flag is not registered and the implementation is marked as TODO in `cmd/edit.go`.

## Steps to Reproduce
1. Run `./strand edit --help`
2. Observe that `--blocker` is missing from the flags list.
3. Run `./strand edit <task-id> --blocker T1234`
4. Observe `Error: unknown flag: --blocker`
