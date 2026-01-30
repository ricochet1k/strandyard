---
type: ""
role: developer
priority: high
parent: E7p4m-issues-recurrence
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T16:35:33.237513-07:00
owner_approval: false
completed: true
---

# Add Issue Subcommand

## Summary

Add a CLI subcommand to create issue-style tasks (non-recurring) with required metadata and a default template.

## Tasks

- [ ] Define issue metadata fields and defaults (role, parent, blockers, labels if needed)
- [ ] Add `memmd issue add` (or `memmd add --issue`) command skeleton
- [ ] Create issue task template in `templates/`
- [ ] Ensure new issues are created with deterministic IDs and directory layout
- [ ] Validate created issues via existing parser/validator

## Implementation Plan

### Architecture overview

Introduce an issue-flavored task creation flow that reuses the existing task creation pipeline and templates. Issues are regular tasks with a specific template and a small metadata extension (e.g., `type: issue` or `issue: true`) so validation and listing logic stays consistent. The command should write a task directory + markdown file that passes the existing parser/validator.

### Files to modify

- `cmd/` (new issue subcommand or flags on existing add/new command)
- `pkg/task/` (extend metadata schema/types for issue fields)
- `templates/` (new `issue.md` template, or extend existing task template to accept kind)
- `CLI.md` (usage examples and flags; coordinated with docs task)

### Approach

1. **Metadata schema**: add a small, deterministic extension to frontmatter. Suggested minimal field: `type: issue` (string) or `issue: true` (bool). Prefer a single field to avoid drift; ensure sorted field output on write.
2. **Command shape**: add `memmd issue add` (preferred for clarity) or `memmd add --kind issue`. Use the existing task creation helper (if any) to avoid duplicated frontmatter logic.
3. **Template**: add `templates/issue.md` with issue-specific headings (Summary/Steps/Acceptance Criteria). Ensure template excludes ID/parent (derived from filesystem) per conventions.
4. **ID and directory layout**: reuse task ID generator and directory creation rules. If issues need a different prefix, define it explicitly and update validation to accept it.
5. **Validation**: ensure the validator treats issue tasks as standard tasks with additional optional metadata; enforce schema if `type: issue` is present.

### Integration points

- Task creation helpers (wherever `new`/`add` currently writes frontmatter).
- Validation in `pkg/task` that checks required metadata and schema.
- Template loader used by CLI.

### Testing approach

- Add/extend unit tests in `pkg/task` to parse frontmatter with `type: issue` and ensure it round-trips.
- CLI tests (if present) to assert directory + file creation with deterministic output.
- Golden file/template test for `templates/issue.md` if template rendering is tested.

### Alternatives considered

- **Separate “issues” registry folder**: rejected because it splits the task database and complicates validation and master list generation.
- **No schema field, only template difference**: rejected because validation cannot distinguish issues for future reporting or filters.

## Acceptance Criteria

- CLI command creates a task directory and markdown file for an issue
- Generated task conforms to frontmatter conventions and directory naming rules
- `memmd repair` passes after adding an issue
- Example usage documented in task body
