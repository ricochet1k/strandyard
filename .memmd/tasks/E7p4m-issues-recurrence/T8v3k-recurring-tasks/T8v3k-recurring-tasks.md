---
type: ""
role: developer
priority: low
parent: E7p4m-issues-recurrence
blockers:
    - T968i-design-recurrence-metrics-for-commits-lines-change
    - Tcb90-document-recurrence-metrics-options
    - Tgr06-review-recurrence-metrics-cli-experience
    - Thnhh-review-recurrence-metrics-schema
    - Tl4cn-approve-recurrence-metrics-design
    - Tyvdv-extend-recurrence-schema-and-validation-for-new-me
blocks:
    - E7p4m-issues-recurrence
    - Iquw5-create-recurring-review-task-plan
date_created: 2026-01-27T00:00:00Z
date_edited: 2026-01-29T19:22:43.455133-07:00
owner_approval: false
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

- [x] (subtask: T968i-design-recurrence-metrics-for-commits-lines-change) Design recurrence metrics for commits, lines changed, and tasks completed
- [ ] (subtask: Tcb90-document-recurrence-metrics-options) Document recurrence metrics options
- [x] (subtask: Tgr06-review-recurrence-metrics-cli-experience) Review recurrence metrics CLI experience
- [ ] (subtask: Thnhh-review-recurrence-metrics-schema) Review recurrence metrics schema
- [ ] (subtask: Tl4cn-approve-recurrence-metrics-design) Approve recurrence metrics design
- [ ] (subtask: Tyvdv-extend-recurrence-schema-and-validation-for-new-me) Extend recurrence schema and validation for new metrics

## Implementation Plan

### Architecture overview

Model recurrence as a first-class task definition that can generate concrete tasks on demand. Recurrence definitions should live in the filesystem-backed task tree and be parsed by the same `pkg/task` parser, but marked with a `type: recurring` (or `recurring: true`) and scheduling metadata. A separate command materializes due instances into normal tasks, placing them in deterministic directories and updating the definition’s “last run” or “next due” metadata.

### Files to modify

- `pkg/task/` (metadata schema, validation, scheduling fields)
- `cmd/` (new `recurring add` and `recurring materialize` commands or subcommands)
- `templates/` (new `recurring.md` template for definitions)
- `CLI.md` (usage and schema; coordinated with docs task)

### Scheduling model

- **Required fields**: `recurrence_interval` (int), `recurrence_unit` (`days`, `weeks`, `months`, `commits`), `recurrence_anchor` (ISO date or commit hash), `recurrence_next_due` (ISO date; optional if computed).
- **Optional fields**: `recurrence_max_instances`, `recurrence_timezone`, `recurrence_last_run`.
- Treat commit-based intervals as a future enhancement if commit data is not currently tracked; gate via validation rules.

### Approach

1. **Definition location**: store recurrence definitions as tasks under a deterministic folder (e.g., `tasks/recurring/<ID>/` or under the parent epic), so they are discoverable by `scan`.
2. **Parser + validator**: extend `pkg/task` to parse recurrence fields, repair required fields, and compute `next_due` deterministically.
3. **Materialization**: implement `memmd recurring materialize` that:
   - scans recurrence definitions
   - filters those due as of “now”
   - generates concrete task directories with a reference to the source (e.g., `parent` or `source_recurring_id`)
   - updates `recurrence_last_run` / `recurrence_next_due` in the definition file
4. **Determinism**: generated task IDs should be based on the definition ID + date stamp to avoid collisions and ensure stable ordering.
5. **Validation integration**: ensure `repair` checks recurrence metadata and that generated tasks conform to normal task rules.

### Integration points

- `cmd/validate.go` should understand recurring definitions and skip treating them as normal tasks for master list inclusion, unless configured.
- `tasks/free-tasks.md` generation should include only materialized tasks (not recurring definitions).

### Testing approach

- Unit tests for parsing/validation of recurrence fields.
- Tests for due calculation across units (days/weeks/months).
- CLI integration tests (if present) that materialize tasks and update the definition deterministically.

### Alternatives considered

- **External registry (JSON/YAML) for recurrence**: rejected to keep single-source-of-truth in task markdown.
- **Cron-style strings only**: rejected as too opaque for deterministic validation and user ergonomics.

## Acceptance Criteria

- Recurring tasks can be created via CLI with explicit interval settings
- Due recurring tasks can be materialized into normal tasks without manual edits
- Generated tasks appear in `free-tasks.md` when unblocked
- Validation reports malformed recurrence metadata
