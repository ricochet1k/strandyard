---
role: developer
parent: E7p4m-issues-recurrence
blockers: []
date_created: 2026-01-27
date_edited: 2026-01-27
completed: false
---

# Add Recurring Task Support

## Summary

Implement recurring task definitions (e.g., clean up AGENTS.md every N days or commits) and a CLI command to create and materialize them.

## Tasks

- [ ] Define recurrence metadata schema (interval type, interval value, anchor date/commit)
- [ ] Decide where recurrence definitions live (task frontmatter vs. separate registry)
- [ ] Add CLI subcommand to add recurring task definitions
- [ ] Add CLI subcommand to materialize due recurring tasks into normal task directories
- [ ] Ensure deterministic ordering and IDs for generated tasks
- [ ] Update validation rules to check recurrence definitions

## Acceptance Criteria

- Recurring tasks can be created via CLI with explicit interval settings
- Due recurring tasks can be materialized into normal tasks without manual edits
- Generated tasks appear in `free-tasks.md` when unblocked
- Validation reports malformed recurrence metadata
