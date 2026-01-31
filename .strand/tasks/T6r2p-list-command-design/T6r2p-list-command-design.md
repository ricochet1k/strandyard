---
type: ""
role: architect
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-31T04:41:33.149358Z
owner_approval: false
completed: true
---

# Design List Command

## Summary
Define the CLI contract and output formats for a `strand list` command, including filters, sorting, and groupings aligned with master lists and task metadata.

## Design Document
- design-docs/list-command.md

## Acceptance Criteria
- Command proposal includes flags, defaults, and output schema
- Deterministic ordering rules are specified
- Implementation notes identify which packages/files should change

## TODOs
- [x] Define supported list scopes (all tasks, root tasks, free tasks, by parent)
- [x] Specify filters (role, priority, completed, blockers, labels if present)
- [x] Decide output formats (table, markdown, JSON) and default sort order
- [x] Identify data sources (scan vs. master lists) and determinism requirements
- [x] Define CLI flags and subcommands to keep the interface stable
- [ ] tasks/T6r2p-list-command-design/T7k1r-list-command-skeleton/T7k1r-list-command-skeleton.md
- [ ] tasks/T6r2p-list-command-design/T7k2m-list-filters-format/T7k2m-list-filters-format.md
- [ ] tasks/T6r2p-list-command-design/T7k3d-list-docs/T7k3d-list-docs.md

## Subtasks
- [x] (subtask: T7k1r) Implement list command skeleton
- [x] (subtask: T7k2m) Implement list filtering, sorting, and formatting
- [x] (subtask: T7k3d) Document list command
- [x] (subtask: T9p4r) Add these aliases:
