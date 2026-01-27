---
role: developer
parent: E9m5w-validate-enhancements
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
---

# Add Task Link Validation

## Summary

Add validation to ensure all task references/links in task content point to existing tasks.

## Tasks

- [ ] Scan task content for task links (markdown links, direct references)
- [ ] Extract task IDs from links (format: `[text](path/to/T3k7x-task/file.md)`)
- [ ] Verify each referenced task ID exists in tasks directory
- [ ] Report broken links with clear error messages showing file and line number
- [ ] Handle both relative and absolute paths
- [ ] Add tests for link validation

## Acceptance Criteria

- Detects broken task links
- Reports file path and line number of broken link
- Validates links to task directories and task files
- Example error: `ERROR: Broken link in T3k7x-example/T3k7x-example.md:15: T9999-missing does not exist`

## Files

- cmd/validate.go
- pkg/metadata/linkchecker.go (new, optional)

## Link Formats to Validate

- Markdown links: `[Task Name](tasks/T3k7x-example/T3k7x-example.md)`
- Direct references: `T3k7x-example`
- Parent field: `parent: E2k7x-metadata-format`
- Blocker field: `blockers: [T5h7w-default-free-task]`
