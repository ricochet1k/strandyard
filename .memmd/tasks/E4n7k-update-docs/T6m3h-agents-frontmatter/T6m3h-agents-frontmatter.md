---
role: developer
parent: E4n7k-update-docs
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27T14:25:00Z
completed: true
---

# Update AGENTS.md to Reflect YAML Frontmatter

## Summary

Update AGENTS.md to show the canonical task format using YAML frontmatter instead of simple field format.

## Tasks

- [ ] Update "Data model and filesystem conventions" section
- [ ] Replace example showing `Role: developer` with YAML frontmatter example
- [ ] Update "Example task.md layout" section with frontmatter format
- [ ] List all frontmatter fields: role, parent, blockers, blocks, date_created, date_edited, owner_approval
- [ ] Update ID format specification to 4-char base36
- [ ] Update template references to flat structure
- [ ] Ensure parsing rules section mentions goldmark-frontmatter

## Acceptance Criteria

- AGENTS.md shows YAML frontmatter as canonical format
- All examples use frontmatter, not simple fields
- Frontmatter fields documented clearly
- ID format shows 4-char base36
- Authoritative and unambiguous

## Files

- AGENTS.md

## Example Format to Show

```markdown
---
role: developer
parent: E2k7x-metadata-format
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
---

# Title: Initialize project skeleton

## Description

Add initial Go module, scaffold cobra CLI, and commit.
```
