---
role: developer
parent: E9m5w-validate-enhancements
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
---

# Add Blocker Status Validation

## Summary

Validate that blocker relationships are bidirectional and consistent, and that free-tasks.md accurately reflects tasks with no blockers.

## Tasks

- [ ] Validate bidirectional blocker relationships:
  - If task A blocks task B, then task B should list A in blockers
  - If task A has blocker B, then task B should list A in blocks
- [ ] Validate free-tasks.md only contains tasks with empty blockers array
- [ ] Validate tasks in free-tasks.md actually exist
- [ ] Report inconsistent blocker relationships with clear errors
- [ ] Suggest fixes for broken relationships

## Acceptance Criteria

- Detects missing bidirectional blocker links
- Detects tasks in free-tasks.md that have blockers
- Detects tasks not in free-tasks.md that should be (no blockers)
- Example errors:
  - `ERROR: Task T3k7x has blocker T5h7w, but T5h7w doesn't list T3k7x in blocks`
  - `ERROR: Task T8n2m is in free-tasks.md but has blocker T6p4k`
  - `ERROR: Task T2h9m has no blockers but is not in free-tasks.md`

## Files

- cmd/validate.go

## Bidirectional Validation Example

Task A:
```yaml
blockers:
  - T5h7w-example
```

Task B (T5h7w-example):
```yaml
blocks:
  - Ta3k7x-taskname  # Should list Task A
```
