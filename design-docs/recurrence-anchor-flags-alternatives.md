# Design Alternatives — Recurrence Anchor Flags and Help Text

## Summary
Clarify how `memmd recurring add` captures the recurrence anchor across time-based and git-based metrics. This document focuses on the CLI flags and help text, aligned with `design-docs/recurrence-metrics.md` and the `recurring add` section in `CLI.md`.

## Context
- `design-docs/recurrence-metrics.md`
- `design-docs/recurrence-anchor-flags-alternatives.md`
- `CLI.md` (recurring add section)

## Project Principles
- Keep CLI usage unambiguous and self-documenting.
- Preserve deterministic behavior and auditability in recurrence definitions.
- Prefer minimal flags unless clarity or validation requires specificity.
- Avoid breaking existing command usage unless there is a clear benefit.

## Alternatives

### Alternative A — Single `--anchor` with unit-specific help text
- Description: Keep a single `--anchor` flag. Validation and help text are driven by `--unit`. The help text explicitly states what format is required for each unit (including tasks-completed anchors).
- Assumptions: `--unit` remains required; `--anchor` is required for all recurrence types; future units include `commits`, `lines_changed`, `tasks_completed` as defined in `recurrence_triggers`.
- Pros:
  - Minimal CLI surface area; aligns with existing `CLI.md` usage.
  - Easy to maintain and document.
  - Backward compatible with current flags.
- Cons:
  - Users may still be unsure which anchor format applies without reading help text.
  - Ambiguous for future mixed-trigger definitions (if `recurrence_triggers` expands).
- Risks:
  - Misconfigured anchors if users skim help output.
- Effort: Small (help text + validation tweaks).

### Alternative B — Split anchors by type (`--anchor-date`, `--anchor-commit`)
- Description: Replace the single `--anchor` with explicit flags per anchor type. Require `--anchor-date` for time-based units and `--anchor-commit` for git-based units. Optionally keep `--anchor` as a deprecated alias for one release cycle.
- Assumptions: Most recurrence types can be clearly mapped to either date or commit anchors; tasks-completed uses date anchors.
- Pros:
  - Removes ambiguity; help text is self-evident.
  - Validation can be strict with clearer error messages.
- Cons:
  - Adds flags and migration complexity.
  - Harder to evolve to mixed-trigger definitions if they are added later.
- Risks:
  - Breaking change unless a deprecation period is supported.
- Effort: Medium (flag changes, migration path, docs updates).

### Alternative C — `--anchor` plus explicit `--anchor-type`
- Description: Keep `--anchor` but add `--anchor-type` with values like `date` or `commit`. Require it when `--unit` is a git-based metric, or allow `auto` with validation fallback.
- Assumptions: Help text can guide users to set `--anchor-type` when ambiguous; tasks-completed uses `--anchor-type date`.
- Pros:
  - More explicit than Alternative A without adding multiple anchor flags.
  - Easier future expansion if anchors diversify.
- Cons:
  - Adds another flag to explain; potential confusion between `--unit` and `--anchor-type`.
  - Still requires help text to connect `--unit` and `--anchor-type`.
- Risks:
  - Users may select inconsistent `--unit`/`--anchor-type` pairs.
- Effort: Medium (new flag + validation + docs).

### Alternative D — Single `--every` with structured string parsing
- Description: Replace `--unit`/`--anchor` with a single repeatable `--every` flag that accepts a small structured string: `<amount> <metric> [from <anchor>]` (e.g., `--every "5 days from 2026-01-28T00:00:00Z"`). Allow multiple `--every` flags for mixed metrics.
- Assumptions: The parser enforces a strict grammar (not free-form NLP). The anchor is optional and defaults to "now" or the last completion depending on metric.
- Pros:
  - Single flag with clear, compact syntax; supports mixed metrics by repeating the flag.
  - Avoids multiple anchor flags while still allowing per-metric anchors.
- Cons:
  - Requires quoting in shells; more complex parsing and validation.
  - Harder to generate precise, beginner-friendly help text.
- Risks:
  - Higher chance of user formatting errors; more complex error messages.
  - Future extensions could complicate the grammar.
- Effort: Medium to Large (parser + validation + docs + tests).

## Help Text Sketches (for evaluation)
- Alternative A example:
  - `--anchor`: Anchor for recurrence. For time units (days/weeks/months): ISO 8601 timestamp. For git units (commits/lines_changed): commit hash. For tasks_completed: ISO 8601 timestamp of last completion.
- Alternative B example:
  - `--anchor-date`: ISO 8601 timestamp (required for time-based recurrence).
  - `--anchor-commit`: Commit hash (required for git-based recurrence).
- Alternative C example:
  - `--anchor`: Anchor value; format depends on `--anchor-type`.
  - `--anchor-type`: Anchor type: date|commit (required for git-based units).
- Alternative D example:
  - `--every`: Repeatable recurrence rule. Format: `<amount> <metric> [from <anchor>]` (e.g., `"5 days from 2026-01-28T00:00:00Z"`, `"50 commits from a1b2c3d4"`).

## Decision
- Decision: Alternative D. Use repeatable `--every` with a strict structured grammar and per-metric anchors via `from <anchor>`.

## Recommendation
- Recommendation: Alternative D to keep the CLI surface to a single flag while supporting mixed metrics via repeated `--every` rules.

## Review Requests
- Request review from: `master-reviewer`, `reviewer-usability`, `reviewer-reliability`.
