---
role: developer
parent: E2k7x-metadata-format
blockers: []
date_created: 2026-01-27T14:15:00Z
date_edited: 2026-01-27T14:17:00Z
completed: true
---

# Update Templates to Use YAML Frontmatter

## Summary

Update task templates to use YAML frontmatter format instead of simple field format.

## Tasks

- [ ] Update templates/leaf.md to use YAML frontmatter header
- [ ] Remove `## Role`, `## Track` markdown headings from template body
- [ ] Add frontmatter with template variables: `{{ .Role }}`, `{{ .Parent }}`, etc.
- [ ] Include date_created and date_edited in frontmatter (should be set to current date when task created)
- [ ] Test template expansion works correctly with new format

## Acceptance Criteria

- Template uses YAML frontmatter format
- Template variables properly expand when creating new tasks
- Created tasks validate successfully
- No markdown heading metadata in body

## Files

- templates/leaf.md

## Example Template Format

```markdown
---
role: {{ .Role }}
parent: {{ .Parent }}
blockers: []
date_created: {{ .DateCreated }}
date_edited: {{ .DateEdited }}
---

# {{ .Title }}

## Context
...
```
