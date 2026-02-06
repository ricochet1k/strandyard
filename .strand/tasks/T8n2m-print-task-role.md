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

# Task Title

## Summary
Change `next` command output to print the task's role (from metadata) instead of the role doc from roles/ directory.

## Acceptance Criteria
- Output shows task's role from metadata
- Does not read or print role doc from roles/ directory
- Clean, minimal output
- Example output:
  ```
  Role: developerTask content...
  ```

## Files
- cmd/next.go

## TODOs
- [ ] Remove role doc loading and printing (lines reading `roles/<role>.md`)
- [ ] Extract role from task metadata (YAML frontmatter)
- [ ] Print role in simple format (e.g., "Role: developer")
- [ ] Remove separator line ("---") between role doc and task
- [ ] Simplify output to just: role + task content
