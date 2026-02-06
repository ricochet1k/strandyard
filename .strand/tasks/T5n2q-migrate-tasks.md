---
type: ""
role: developer
priority: ""
parent: E2k7x-metadata-format
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T14:16:00Z
owner_approval: false
completed: true
---

# Migrate Existing Tasks to New Format

## Summary
Convert all existing task files from simple field format to YAML frontmatter format.

## Acceptance Criteria
- All task files use YAML frontmatter
- `go run ./cmd/strand repair` passes with no errors
- No task data lost during migration
- Dates populated for all tasks

## Files
- tasks/D000001-review-design/task.md
- tasks/T000001-project-alpha/task.md
- Any other task files in tasks/

## TODOs
- [ ] Manually convert old tasks, or consider deleting them altogether.
- [ ] Convert D000001-review-design/task.md to new format
- [ ] Convert T000001-project-alpha/task.md to new format
- [ ] Convert any other existing task files
- [ ] Add date_created and date_edited fields (use file mtime for initial dates)
- [ ] Verify all migrated tasks repair successfully
