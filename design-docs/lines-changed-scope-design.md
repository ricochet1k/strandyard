# Design — lines_changed Scope and Default

## Summary
Define how the `lines_changed` recurrence metric accounts for changes in the repository. This includes whether to count additions, deletions, or both (total delta), and how to optionally scope the metric to specific file paths or patterns.

## Context
- `design-docs/recurrence-metrics.md`
- `pkg/task/recurrence.go` (current implementation: `additions + deletions`)

## Project Principles
- Keep recurrence evaluation deterministic and auditable.
- Prefer minimal flags unless clarity or validation requires specificity.
- Avoid breaking existing command usage.

## Alternatives

### Alternative A — Total Delta (Current Default)
- **Description**: Count both additions and deletions as "lines changed".
- **Pros**: Simplest to understand; reflects general activity level.
- **Cons**: Can be misleading for large refactors (high count but little new functionality).
- **Default**: Yes.

### Alternative B — Distinct Units (`lines_added`, `lines_deleted`)
- **Description**: Instead of a generic `lines_changed`, provide specific units for additions and deletions.
- **Pros**: Very clear intent (e.g., trigger every 1000 lines of growth); no new grammar needed.
- **Cons**: Doesn't support "total activity" as easily without combining triggers.
- **Effort**: Small (add new units to `recurrence.go`).

### Alternative C — Scope Flag/Parameter
- **Description**: Add an optional path pattern to the `lines_changed` trigger.
- **Syntax**: `--every "500 lines_changed [from <anchor>] [in <pattern>]"`
- **Pros**: Extremely flexible; allow monitoring specific components (e.g., `pkg/activity`).
- **Cons**: Adds complexity to the parser and metadata.
- **Effort**: Medium.

## Recommendation
**Adopt Alternative A as the default behavior** (total delta) for the `lines_changed` unit, but **also implement Alternative B** by adding `lines_added` and `lines_deleted` as distinct units. This provides clarity and flexibility without complicating the core grammar.

Wait on Alternative C until a specific use case for path-scoping is identified.

## Decision
Decision: deferred to Owner.
