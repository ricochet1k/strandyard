---
type: ""
role: developer
priority: ""
parent: E3q8p-next-command
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T14:20:00Z
owner_approval: false
completed: true
---

# Update Next to Default to First Free Task

## Summary
Remove the requirement for `--role` flag or `MEMMD_ROLE` env var. Make `next` command default to reading the first task from free-tasks.md.

## Acceptance Criteria
- `strand next` command runs successfully without flags
- Reads first task from tasks/free-tasks.md
- No error about missing role
- Still works with `--role` flag for optional filtering

## Files
- cmd/next.go

## TODOs
- [ ] Remove role requirement check from cmd/next.go
- [ ] Update logic to read free-tasks.md and select first task (no role filtering)
- [ ] Remove role filtering loops (first pass, second pass, fallback)
- [ ] Keep role flag as optional for future filtering if needed
- [ ] Update help text and error messages
- [ ] Test `strand next` works without any arguments
