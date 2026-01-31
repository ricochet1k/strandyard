---
type: ""
role: architect
priority: ""
parent: ""
blockers: []
blocks:
    - E5w8m-e2e-tests
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-31T04:41:33.149341Z
owner_approval: false
completed: true
---

# Implement YAML Frontmatter Metadata Format

## Summary
Replace the current simple field format (`Role: developer`) with YAML frontmatter using goldmark-frontmatter library. This provides a cleaner, more structured approach to task metadata that's easier to parse and extend.

## Context
**Owner Decision**: Use goldmark-frontmatter with YAML frontmatter for all task metadata including role, parent, blockers, owner approval, date created, date edited, and other helpful data.

**Current state**: Parser uses simple text parsing for `Role:`, `Blockers:`, etc. with regex and string matching.

**Target state**: All tasks use YAML frontmatter at the top of the file, parsed with goldmark-frontmatter library.

## Acceptance Criteria
- All task files use YAML frontmatter
- Parser reads metadata from frontmatter, not text parsing
- `repair` command works with new format
- All existing tasks migrated successfully
- Templates updated to new format

## References
- goldmark-frontmatter: https://github.com/abhinav/goldmark-frontmatter
- Original design doc specified goldmark for parsing/rendering

## TODOs
- [ ] [T3m9p-add-frontmatter-dep](T3m9p-add-frontmatter-dep/T3m9p-add-frontmatter-dep.md) - Add goldmark-frontmatter dependency
- [ ] [T8h4w-update-parser](T8h4w-update-parser/T8h4w-update-parser.md) - Update parser to read YAML frontmatter
- [ ] [T5n2q-migrate-tasks](T5n2q-migrate-tasks/T5n2q-migrate-tasks.md) - Migrate existing tasks to new format
- [ ] [T9x7k-update-templates](T9x7k-update-templates/T9x7k-update-templates.md) - Update templates to use YAML frontmatter

## Subtasks
- [x] (subtask: T3m9p) Add goldmark-frontmatter Dependency
- [x] (subtask: T5n2q) Migrate Existing Tasks to New Format
- [x] (subtask: T8h4w) Task Title
- [x] (subtask: T9x7k) {{ .Title }}
