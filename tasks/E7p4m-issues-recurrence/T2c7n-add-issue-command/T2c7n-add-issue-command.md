---
role: developer
parent: E7p4m-issues-recurrence
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
completed: false
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

## Acceptance Criteria

- CLI command creates a task directory and markdown file for an issue
- Generated task conforms to frontmatter conventions and directory naming rules
- `memmd validate` passes after adding an issue
- Example usage documented in task body
