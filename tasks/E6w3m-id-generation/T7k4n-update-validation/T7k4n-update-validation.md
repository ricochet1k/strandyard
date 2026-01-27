---
role: developer
parent: E6w3m-id-generation
blockers:
  - T2p8h-base36-generator
date_created: 2026-01-27
date_edited: 2026-01-27
---

# Update ID Validation Regex

## Summary

Update the ID validation regex in validate.go to accept the new 4-character base36 token format.

## Tasks

- [ ] Update regex in cmd/validate.go from `^[A-Z][0-9A-Za-z]{6}-[a-zA-Z0-9-]{1,}$` to accept 4-char tokens
- [ ] New regex should be: `^[A-Z][0-9a-z]{4}-[a-zA-Z0-9-]+$`
- [ ] Verify regex accepts valid IDs: `T3k7x-example`, `D9m2p-design`
- [ ] Verify regex rejects invalid IDs: `T123-bad`, `Tabcd-bad`, `T3K7X-bad` (uppercase token)
- [ ] Add test cases for regex validation
- [ ] Update error messages to reflect new format

## Acceptance Criteria

- Regex accepts format: `<PREFIX><4-lowercase-alphanumeric>-<slug>`
- PREFIX is single uppercase letter
- Token is 4 lowercase base36 characters (0-9, a-z)
- Slug is 1+ alphanumeric/hyphen characters
- Validation tests pass

## Files

- cmd/validate.go
- cmd/validate_test.go (when tests added)
