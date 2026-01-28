# Recurrence Metrics: Commits, Lines Changed, Tasks Completed

## Summary
Define recurrence metrics beyond time-based intervals so recurring tasks can trigger based on commit count, lines changed, or completed tasks. This document outlines schema options, data sourcing, validation rules, and follow-on implementation tasks.

## Goals
- Support commit-count, lines-changed, and tasks-completed metrics in recurrence definitions.
- Allow metrics to be used alone or combined with time-based intervals.
- Keep recurrence evaluation deterministic and auditable.

## Non-Goals
- Implement recurring task commands in this phase.
- Build a UI for configuring recurrence.

## Current State
- Recurrence support is planned in `tasks/E7p4m-issues-recurrence/T8v3k-recurring-tasks/T8v3k-recurring-tasks.md`.
- Only time-based intervals (days/weeks/months) and a `commits` unit are described.
- The CLI has no `recurring` command yet, so configuration is not currently possible.

## Data Sources
- **Commits**: `git rev-list --count <anchor>..HEAD` for commit counts.
- **Lines changed**: `git diff --numstat <anchor>..HEAD` aggregated into additions/deletions.
- **Tasks completed**: task metadata (requires a completion timestamp or a completion log).

## Schema Options
### Option A: Extend `recurrence_unit`
Introduce new units alongside the existing interval fields.
- `recurrence_interval`: integer
- `recurrence_unit`: `days | weeks | months | commits | lines_changed | tasks_completed`
- `recurrence_anchor`: ISO 8601 date or commit hash (reused for git-based metrics)
- `recurrence_next_due`: ISO 8601 date (computed for time-based units)
- **Pros**: Minimal new fields; aligns with current schema.
- **Cons**: Hard to represent combined metrics (time + commits) without multiple definitions.

### Option B: Add metric triggers array
Add a `recurrence_triggers` list with metric definitions that can be combined.
- `recurrence_triggers`:
  - `type`: `time | commits | lines_changed | tasks_completed`
  - `interval`: integer
  - `anchor`: ISO date or commit hash
- `recurrence_next_due`: computed when `type: time` is present
- **Pros**: Supports combinations and clearer validation per trigger.
- **Cons**: More schema surface area and migration work.

**Decision: deferred** — maintainer to choose between Options A and B before implementation.

## Tasks-Completed Metric Options
### Option A: Add `date_completed` to task metadata
- Store completion timestamp when `memmd complete` runs.
- Recurrence counts completed tasks where `date_completed > anchor`.
- **Pros**: Simple to compute and audit.
- **Cons**: Requires metadata migration for historical tasks.

### Option B: Track completion log in recurrence definition
- Store `recurrence_task_anchor` as a list of completed task IDs or a `last_completed_at` timestamp in the recurring definition.
- **Pros**: Avoids global schema changes.
- **Cons**: Re-implements completion tracking in recurrence logic.

**Decision: deferred** — pick a storage strategy after owner review.

## Validation Rules
- Require positive `recurrence_interval` for all triggers.
- Validate `recurrence_unit` (Option A) or `recurrence_triggers[].type` (Option B).
- Enforce presence of `recurrence_anchor` for git-based metrics.
- For tasks-completed metrics, require `date_completed` or equivalent anchor fields based on the chosen storage strategy.

## CLI/Template Impacts
- Add flags for metric selection once `recurring add` exists.
- Update templates to include the selected metric fields.
- Update `CLI.md` and examples to show metric-based recurrence.

## Integration Points
- `pkg/task` parser and validator for new recurrence fields.
- Recurrence materialization logic (when implemented) to query git and task metadata.
- `memmd complete` to capture completion timestamps if Option A for tasks-completed is selected.

## Testing Strategy
- Unit tests for schema validation and due-calculation per metric.
- Fixtures with deterministic git histories for commit/line metrics.
- Tests covering tasks-completed calculations once the storage approach is chosen.

## Open Questions
- Which schema option (Option A or B) should be adopted?
- Should lines-changed counts use additions only, deletions only, or total delta?
- Should task completion counting include issues, or only `type: task` entries?
