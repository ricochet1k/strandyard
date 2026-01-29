# Design Alternatives — Recurrence Anchor Error Messages

## Summary
Compare approaches for user-facing error messages when parsing recurrence anchors in `memmd recurring add --every`, with an emphasis on clarity, recoverability, and deterministic behavior.

## Context
- design-docs/recurrence-anchor-flags-alternatives.md
- design-docs/recurrence-metrics.md
- CLI.md (recurring add section)

## Project Principles
- Keep CLI usage unambiguous and self-documenting.
- Make error messages actionable with recovery hints and examples.
- Preserve deterministic behavior and auditability in recurrence definitions.
- Avoid breaking existing command usage unless there is a clear benefit.

## Alternatives

### Alternative A — Metric-specific errors with inline recovery hints
- Description: Emit specific error messages for each grammar failure, including the metric, expected anchor type, and a concrete example hint.
- Assumptions: Users prefer precise guidance over generic errors; tests can assert exact strings.
- Pros:
  - Highly actionable; the user sees the fix immediately.
  - Keeps the CLI self-documenting without extra flags.
  - Enables precise test coverage for each grammar branch.
- Cons:
  - Requires maintaining message strings as grammar evolves.
  - More conditional logic in validation paths.
- Risks:
  - Small grammar changes can cascade into many expected-string updates.
- Rough effort estimate: Medium.
- Example messages:
  - Recovery hint: "Example: --every \"10 commits\" (defaults to now)"
  - Malformed date anchor: "invalid anchor '2026-13-01': expected ISO 8601 timestamp after 'from'"
  - Recovery hint: "Use \"Jan 28 2026 09:00\""
  - Malformed commit anchor: "invalid anchor '2026-01-01': expected commit hash after 'from'"
  - Recovery hint: "Use 'git rev-parse HEAD'"
  - Unit/anchor mismatch: "metric 'days' expects a date anchor, got commit hash"
  - Recovery hint: "Use a date like \"Jan 28 2026 09:00\" for days/weeks/months metrics"
  - Ambiguous anchor type: "metric 'tasks_completed' requires date anchor"
  - Recovery hint: "Use completion time like \"Jan 28 2026 09:00\""

### Alternative B — Unified error format with structured reason + hint line
- Description: Use a consistent error prefix and structured reason, followed by a separate hint line (e.g., `hint:`) with a minimal example.
- Assumptions: Consistent formatting is more valuable than per-metric phrasing.
- Pros:
  - Easier to parse and test; reduces string churn.
  - Standardized output for help text and potential future tooling.
- Cons:
  - Less tailored to specific metrics unless the reason string is rich.
  - Slightly less friendly for first-time users.
- Risks:
  - Overly generic messages could increase support churn.
- Rough effort estimate: Low to Medium.
- Example messages:
  - "memmd: error: invalid --every value: expected date anchor after 'from', got '2026-13-01'"
  - "hint: --every \"10 days from Jan 28 2026 09:00\""

### Alternative C — Default anchors with warning-only errors
- Description: If the anchor is missing, default to `now` for time-based metrics and `HEAD` for commit-based metrics, and only error on malformed anchors; emit a warning when defaults are applied.
- Assumptions: Users are comfortable with implicit defaults; warnings are acceptable in non-error output.
- Pros:
  - Reduces friction for common cases.
  - Keeps the command working even with incomplete input.
- Cons:
  - Less explicit; may surprise users and weaken auditability.
  - Warnings may be ignored in scripts.
- Risks:
  - Silent behavior changes could be interpreted as breaking for strict workflows.
- Rough effort estimate: Medium.

## Decision
- Decision: Alternative B with a fixed error output contract.
- Missing anchors are explicitly allowed and preferred for common usage.
- Use "from now" for immediate run + recur; use "after now" to schedule the first run at the next interval after now.

## Output Contract
- Error prefix: `memmd: error: ` (stable, single-line prefix for the primary error line).
- Hint line: optional second line prefixed with `hint: `.
- Output channel: errors and hints emit to stderr only; no stdout output on failure.
- Exit codes:
  - `2` for parse/validation failures related to `--every`.
  - `1` for other runtime errors.
  - `0` on success.

## Defaults and Hint Examples
- `from <anchor>` is optional; if omitted, the anchor defaults to `now` and is treated as "from now".
- Default anchor is `now` for all metrics, interpreted relative to what the metric measures (for example, `HEAD` for commit-based metrics).
- `after now` means the first run occurs at the next interval after the current time; `from now` triggers an immediate run and then recurs.
- Hint examples should prefer human-friendly dates (for example, "Jan 28 2026 09:00") and may include ISO 8601 as a secondary reference if needed.

## Review Requests
- Request review from: `reviewer` (master), `reviewer-usability`, `reviewer-reliability`.
