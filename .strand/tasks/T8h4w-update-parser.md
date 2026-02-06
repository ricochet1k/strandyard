---
type: ""
role: developer
priority: ""
parent: E2k7x-metadata-format
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T14:15:00Z
owner_approval: false
completed: true
---

# Task Title

## Summary
Replace the current text-based parsing functions (`parseRole`, `parseBlockers`) with goldmark-frontmatter parsing to read task metadata from YAML frontmatter.

## Acceptance Criteria
- `ParseTaskMetadata` successfully reads YAML frontmatter from task files
- Handles missing or malformed frontmatter with clear error messages
- All metadata fields properly parsed into struct
- repair command works with new parser
- Tests pass

## Files
- pkg/metadata/parser.go (new)
- pkg/metadata/parser_test.go (new)
- cmd/validate.go
- cmd/next.go

## Example Format
```markdown
---
role: developer
parent: E2k7x-metadata-format
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
---Task content here...
```

## TODOs
- [ ] Create new `pkg/metadata/` package for frontmatter parsing
- [ ] Implement `ParseTaskMetadata(filepath string)` function using goldmark-frontmatter
- [ ] Define `TaskMetadata` struct with fields: Role, Parent, Blockers, DateCreated, DateEdited, OwnerApproval
- [ ] Update `cmd/validate.go` to use new parser instead of `parseRole` and `parseBlockers`
- [ ] Update `cmd/next.go` to use new parser
- [ ] Remove old parsing functions once migration complete
