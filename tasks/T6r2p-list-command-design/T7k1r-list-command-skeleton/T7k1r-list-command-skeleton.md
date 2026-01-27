---
role: developer
priority: high
parent: T6r2p-list-command-design
blockers: []
blocks: []
date_created: 2026-01-27T18:40:00Z
date_edited: 2026-01-27T18:40:00Z
owner_approval: false
completed: false
---

# Implement list command skeleton

## Context
See design doc: design-docs/list-command.md

## Tasks
- [ ] Add `cmd/list.go` with Cobra wiring and flags.
- [ ] Add option parsing that maps to a `ListOptions` struct.
- [ ] Wire command to list scan/format pipeline (stub allowed if tests cover contract).

## Acceptance Criteria
- `memmd list --help` shows full flag set from the design doc.
- Command parses flags into a well-defined options struct.
- Non-zero exit codes for invalid flag combinations are in place.
