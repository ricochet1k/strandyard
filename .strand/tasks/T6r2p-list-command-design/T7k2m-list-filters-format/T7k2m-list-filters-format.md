---
type: ""
role: developer
priority: medium
parent: T6r2p-list-command-design
blockers: []
blocks: []
date_created: 2026-01-27T18:41:00Z
date_edited: 2026-01-27T17:04:24.917586-07:00
owner_approval: false
completed: true
---

# Implement list filtering, sorting, and formatting

## Context
See design doc: design-docs/list-command.md

## Acceptance Criteria
- Filters and sorting match the design doc in all listed combinations.
- Output formats are deterministic and stable for the same input.
- Tests cover the filter + sort matrix and format schemas.

## TODOs
- [ ] Add filtering helpers for role/priority/completed/blocked/blocks/owner-approval.
- [ ] Implement deterministic sorting with tie-breakers.
- [ ] Implement table, markdown, and JSON formatters with shared schema fields.
- [ ] Add golden tests for formatting outputs.
