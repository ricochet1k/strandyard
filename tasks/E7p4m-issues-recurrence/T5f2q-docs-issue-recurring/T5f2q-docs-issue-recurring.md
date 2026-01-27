---
role: documentation
priority: low
parent: E7p4m-issues-recurrence
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
completed: false
---

# Document Recurring Task Commands

## Summary

Update CLI documentation to describe the recurring task commands, including examples and metadata fields.

## Tasks

- [ ] Add command docs to `CLI.md`
- [ ] Provide examples for recurring task creation and materialization
- [ ] Document recurrence metadata schema and validation rules

## Implementation Plan

### Architecture overview

Document the new command surface area in `CLI.md` so it matches the actual CLI flags and metadata schema. Keep docs authoritative for the recurrence metadata and provide realistic examples aligned with templates and parser expectations.

### Files to modify

- `CLI.md` (new sections for `issue` and `recurring` commands)
- `templates/` (if docs need to reference template snippets)

### Approach

1. Add a “Recurring tasks” section describing:
   - `memmd recurring add` with required scheduling fields
   - `memmd recurring materialize` for generating due tasks
2. Include copy-pasteable examples showing frontmatter results and file layout.
3. Document validation rules (required fields, accepted units, determinism) and note how `validate` treats recurring definitions vs. materialized tasks.

### Integration points

- Align examples with the templates and parser schema defined in `pkg/task`.
- Ensure command names and flags match `cmd/` definitions.

### Testing approach

- Run `memmd validate` after doc updates if docs include embedded examples or references that affect templates.
- Spot-check examples against actual CLI output once commands exist.

### Alternatives considered

- **Separate docs file for issues/recurrence**: rejected to keep CLI usage centralized in `CLI.md`.

## Acceptance Criteria

- `CLI.md` includes new command descriptions and examples
- Docs reflect the accepted metadata schema
- Examples are consistent with actual command flags
