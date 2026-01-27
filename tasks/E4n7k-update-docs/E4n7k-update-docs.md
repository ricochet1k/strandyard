---
role: architect
parent:
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27T14:25:00Z
completed: true
---

# Update Documentation

## Summary

Update all documentation to reflect Owner decisions: flat template structure, YAML frontmatter format, 4-char random IDs, and proper document organization.

## Context

**Owner Decisions**:
- Template organization: Use current flat structure (templates/leaf.md), update docs to match
- Move design-alternatives-review.md out of doc-examples (it's not an example)
- Update AGENTS.md to reflect YAML frontmatter format
- Update design-docs to match implementation

## Subtasks

1. [T8p5w-design-docs-templates](T8p5w-design-docs-templates/T8p5w-design-docs-templates.md) - Update design-docs to match flat template structure
2. [T6m3h-agents-frontmatter](T6m3h-agents-frontmatter/T6m3h-agents-frontmatter.md) - Update AGENTS.md to reflect YAML frontmatter
3. [T9k2n-move-alternatives-doc](T9k2n-move-alternatives-doc/T9k2n-move-alternatives-doc.md) - Move design-alternatives-review.md to design-docs

## Acceptance Criteria

- All documentation accurate and up-to-date
- design-alternatives-review.md in correct location
- AGENTS.md shows YAML frontmatter format
- design-docs/commands-design.md reflects current decisions
- No contradictions between docs and implementation

## Files

- design-docs/commands-design.md
- AGENTS.md
- doc-examples/design-alternatives-review.md (move)
