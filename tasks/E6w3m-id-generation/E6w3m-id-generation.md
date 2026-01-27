---
role: architect
parent:
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27T14:31:00Z
completed: true
---

# Update ID Generation System

## Summary

Update task ID generation to use 4-character random base36 tokens instead of the current sequential numbering or 6-character format.

## Context

**Owner Decision**: Task IDs should use 4 random base36 characters for the token portion.

**Current state**: Tasks use sequential format like `T000001-project-alpha` or design spec calls for 6-char random tokens.

**Target state**: Task IDs use format `<prefix><4-char-base36>-<slug>`, e.g., `T3k7x-implement-parser`, `D9m2p-api-design`.

## Subtasks

1. [T2p8h-base36-generator](T2p8h-base36-generator/T2p8h-base36-generator.md) - Implement 4-char base36 token generator
2. [T7k4n-update-validation](T7k4n-update-validation/T7k4n-update-validation.md) - Update ID validation regex
3. [T4m9x-migration-tool](T4m9x-migration-tool/T4m9x-migration-tool.md) - Create migration tool for existing tasks

## Acceptance Criteria

- ID generator creates 4-char random base36 tokens using crypto/rand
- Validation regex accepts new format
- Existing tasks can be migrated to new ID format
- No ID collisions (36^4 = 1.6M possible IDs should be sufficient)

## References

- Base36 uses characters: 0-9, a-z (36 total)
- 4 characters = 36^4 = 1,679,616 possible combinations
