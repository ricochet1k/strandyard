---
role: developer
parent: E4n7k-update-docs
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27T14:25:00Z
completed: true
---

# Update Design Docs to Match Flat Template Structure

## Summary

Update design-docs/commands-design.md to reflect the Owner's decision to keep flat template structure (templates/ instead of templates/task-templates/).

## Tasks

- [ ] Update references to `templates/task-templates/` → `templates/`
- [ ] Update references to `templates/doc-templates/` → `doc-examples/` (or remove if not applicable)
- [ ] Update ID format specification from 6-char to 4-char random base36
- [ ] Remove or update any mentions of sequential IDs
- [ ] Update task metadata examples to show YAML frontmatter
- [ ] Update command examples to reflect new behavior
- [ ] Ensure consistency across entire document

## Acceptance Criteria

- design-docs/commands-design.md accurately reflects Owner decisions
- Template paths correct (flat structure)
- ID format correct (4-char base36)
- Metadata format shows YAML frontmatter
- No contradictions with actual implementation

## Files

- design-docs/commands-design.md
