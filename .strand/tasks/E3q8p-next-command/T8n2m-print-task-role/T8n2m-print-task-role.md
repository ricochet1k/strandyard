---
role: developer
parent: E3q8p-next-command
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27T14:20:00Z
completed: true
---

# Update Next to Print Task's Role

## Summary

Change `next` command output to print the task's role (from metadata) instead of the role doc from roles/ directory.

## Tasks

- [ ] Remove role doc loading and printing (lines reading `roles/<role>.md`)
- [ ] Extract role from task metadata (YAML frontmatter)
- [ ] Print role in simple format (e.g., "Role: developer")
- [ ] Remove separator line ("---") between role doc and task
- [ ] Simplify output to just: role + task content

## Acceptance Criteria

- Output shows task's role from metadata
- Does not read or print role doc from roles/ directory
- Clean, minimal output
- Example output:
  ```
  Role: developer

  # Task Title
  Task content...
  ```

## Files

- cmd/next.go
