---
type: ""
role: developer
priority: ""
parent: E2k7x-metadata-format
blockers: []
blocks: []
date_created: 2026-01-27T14:15:00Z
date_edited: 2026-01-27T14:17:00Z
owner_approval: false
completed: true
---

# {{ .Title }}

## Summary
Update task templates to use YAML frontmatter format instead of simple field format.

## Acceptance Criteria
- Template uses YAML frontmatter format
- Template variables properly expand when creating new tasks
- Created tasks repair successfully
- No markdown heading metadata in body

## Files
- templates/leaf.md

## Example Template Format
```markdown
---
---

## Context
...
```

## TODOs
- [ ] Update templates/leaf.md to use YAML frontmatter header
- [ ] Remove `## Role`, `## Track` markdown headings from template body
- [ ] (role: designer) Add frontmatter with template variables: example, example.
- [ ] Include date_created and date_edited in frontmatter (should be set to current date when task created)
- [ ] Test template expansion works correctly with new format
