---
type: ""
role: developer
priority: ""
parent: E6w3m-id-generation
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T14:31:00Z
owner_approval: false
completed: true
---

# Create Migration Tool for Existing Tasks

## Summary
Create a tool to migrate existing task IDs from sequential format (T000001) to new random base36 format (T3k7x).

## Acceptance Criteria
- Migration tool successfully renames all tasks
- All references updated correctly
- No broken links after migration
- Validate passes after migration
- Backup created for safety

## Files
- cmd/migrate.go (new) or scripts/migrate-ids.go
- All files in tasks/ (renamed/updated)

## TODOs
- [ ] Create migration script/command
- [ ] Scan all existing task directories
- [ ] Generate new ID for each task (preserve prefix, regenerate token, keep slug)
- [ ] Rename directories: `tasks/T000001-alpha/` → `tasks/T3k7x-alpha/`
- [ ] Rename task files to match new directory names
- [ ] Update all references to old IDs in:
- [ ] Parent fields in other tasks
- [ ] Blocker fields in other tasks
- [ ] Master lists (tasks/root-tasks.md, tasks/free-tasks.md)
- [ ] Documentation files
- [ ] Create backup before migration
- [ ] Provide dry-run option to preview changes
