---
role: developer
parent: E9m5w-validate-enhancements
blockers:
  - T3k8p-link-validation
  - T7w4n-blocker-validation
date_created: 2026-01-27
date_edited: 2026-01-27
---

# Ensure Single Right Way Validation

## Summary

Remove any optional validation modes and ensure there's one right way to structure tasks. Keep validation simple and strict.

## Tasks

- [ ] Remove any lenient/strict mode options (if they exist or were planned)
- [ ] Ensure validation always fails fast on first error category
- [ ] Provide clear, actionable error messages
- [ ] No warnings - only errors that must be fixed
- [ ] Document the canonical task format in AGENTS.md
- [ ] Update error messages to be specific and helpful

## Acceptance Criteria

- No `--strict` or `--lenient` flags
- Validation either passes completely or fails with errors
- Each error message clearly states what's wrong and how to fix it
- No ambiguity in what constitutes a valid task
- Documentation clearly describes the one right way

## Files

- cmd/validate.go
- AGENTS.md (update canonical format)

## Error Message Examples

Good error messages:
- `ERROR: Task T3k7x-example missing required frontmatter field 'role'`
- `ERROR: Invalid ID format 'T123-bad' in tasks/T123-bad/ - must be <PREFIX><4-lowercase-alphanumeric>-<slug>`

Bad error messages:
- `WARNING: Task might have issues`
- `ERROR: Invalid format`
