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

# Some Task

## Summary
If task metadata doesn't have a role, extract role from the first TODO item in the task content.

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

## TODOs
- [ ] Parse task content for TODO items (e.g., `- [ ] (role: developer) Do something`)
- [ ] Extract role from first TODO if present
- [ ] Use format: `- [ ] (role: <role>) Task description`
- [ ] Fallback: if no role in metadata and no TODO role, print "Role: (none)"
- [ ] Add tests for TODO parsing
- [ ] (role: developer) Write the code
- [ ] (role: reviewer) Review the PR
