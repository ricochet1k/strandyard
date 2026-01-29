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
  - Missing anchor: "missing anchor: expected 'from <anchor>' after metric 'commits'"
  - Recovery hint: "Example: --every \"10 commits from a1b2c3d4\""
  - Malformed date anchor: "invalid anchor '2026-13-01': expected ISO 8601 timestamp after 'from'"
  - Recovery hint: "Use 2026-01-01T00:00:00Z"
  - Malformed commit anchor: "invalid anchor '2026-01-01': expected commit hash after 'from'"
  - Recovery hint: "Use 'git rev-parse HEAD'"
  - Unit/anchor mismatch: "metric 'days' expects a date anchor, got commit hash"
  - Recovery hint: "Use ISO 8601 timestamp for days/weeks/months metrics"
  - Ambiguous anchor type: "metric 'tasks_completed' requires date anchor"
  - Recovery hint: "Use completion timestamp: 2026-01-28T00:00:00Z"

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
  - "invalid --every value: missing anchor after 'commits'"
  - "hint: --every \"10 commits from a1b2c3d4\""

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
- Decision: Deferred to Owner.

## Review Requests
- Request review from: `reviewer` (master), `reviewer-usability`, `reviewer-reliability`.
