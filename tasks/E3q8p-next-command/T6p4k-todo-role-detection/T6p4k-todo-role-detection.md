---
role: developer
parent: E3q8p-next-command
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27T14:20:00Z
completed: true
---

# Add TODO-Based Role Detection

## Summary

If task metadata doesn't have a role, extract role from the first TODO item in the task content.

## Tasks

- [ ] Parse task content for TODO items (e.g., `- [ ] (role: developer) Do something`)
- [ ] Extract role from first TODO if present
- [ ] Use format: `- [ ] (role: <role>) Task description`
- [ ] Fallback: if no role in metadata and no TODO role, print "Role: (none)"
- [ ] Add tests for TODO parsing

## Acceptance Criteria

- If task has role in metadata, use that
- If no metadata role, check first TODO for `(role: xxx)` pattern
- Parse TODOs correctly with regex or simple string matching
- Gracefully handle tasks with no role information
- Example TODO format: `- [ ] (role: developer) Implement the feature`

## Files

- cmd/next.go
- pkg/metadata/parser.go (add TODO parsing function)

## Example

Task with TODO role:
```markdown
---
role:
---

# Some Task

## Tasks
- [ ] (role: developer) Write the code
- [ ] (role: reviewer) Review the PR
```

Should output: `Role: developer`
