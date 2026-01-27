---
role: developer
parent: E2k7x-metadata-format
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
---

# Add goldmark-frontmatter Dependency

## Summary

Add the goldmark-frontmatter library to the project to enable parsing YAML frontmatter from markdown files.

## Tasks

- [ ] Run `go get github.com/abhinav/goldmark-frontmatter`
- [ ] Run `go get github.com/yuin/goldmark` (ensure latest version)
- [ ] Verify dependencies resolve correctly with `go mod tidy`
- [ ] Commit updated go.mod and go.sum

## Acceptance Criteria

- goldmark-frontmatter is listed in go.mod
- Project builds successfully: `go build ./...`
- Dependencies are properly resolved

## Files

- go.mod
- go.sum
