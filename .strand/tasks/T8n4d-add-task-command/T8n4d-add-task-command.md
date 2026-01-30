---
type: ""
role: developer
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-27T21:28:47.007423-07:00
owner_approval: false
completed: true
---

# Add Task Subcommand

## Summary

Implement a CLI subcommand to create standard tasks with required metadata and deterministic IDs using templates and the existing filesystem conventions.

## Tasks

- [ ] Define required flags and defaults (role, parent, priority, blockers)
- [ ] Implement `strand task add` (or `strand add`) command skeleton
- [ ] Wire template rendering for standard tasks in `templates/`
- [ ] Ensure directory naming and ID generation are deterministic
- [ ] Validate created tasks via existing parser/validator

## Acceptance Criteria

- CLI command creates task directory and markdown file that pass validation
- Generated frontmatter adheres to required schema and ordering
- Example usage documented in the task body or CLI docs
