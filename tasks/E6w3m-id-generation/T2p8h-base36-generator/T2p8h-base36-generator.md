---
role: developer
parent: E6w3m-id-generation
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27T14:30:00Z
completed: true
---

# Implement 4-Char Base36 Token Generator

## Summary

Create a cryptographically secure random ID generator that produces 4-character base36 tokens for task IDs.

## Tasks

- [ ] Create `pkg/idgen/` package
- [ ] Implement `GenerateToken()` function using `crypto/rand`
- [ ] Convert random bytes to base36 encoding (0-9, a-z)
- [ ] Ensure exactly 4 characters output
- [ ] Add collision detection (check existing task IDs)
- [ ] Implement `GenerateID(prefix, title string)` that combines prefix + token + slug
- [ ] Add slugify function to convert title to slug
- [ ] Write comprehensive tests

## Acceptance Criteria

- `GenerateToken()` returns 4-char base36 string
- Uses crypto/rand for cryptographic security
- Slugify converts titles like "Implement Parser" → "implement-parser"
- GenerateID produces IDs like "T3k7x-implement-parser"
- Tests verify randomness and format
- No collisions in generated IDs

## Files

- pkg/idgen/generator.go (new)
- pkg/idgen/generator_test.go (new)

## Example Usage

```go
id := idgen.GenerateID("T", "Implement Parser")
// Output: "T3k7x-implement-parser"
```
