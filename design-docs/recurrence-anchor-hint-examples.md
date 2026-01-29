# Canonical Hint Examples for --every Anchors

## Summary
Define deterministic, reusable hint examples for `memmd recurring add --every` anchor parsing errors so tests and automation remain stable.

## Context
- design-docs/recurrence-anchor-error-messages-alternatives.md
- design-docs/recurrence-anchor-error-messages-reliability-review.md
- design-docs/recurrence-anchor-flags-alternatives.md
- design-docs/recurrence-metrics.md

## Decision
Adopt a fixed set of anchor examples that never depend on the current time, system locale, or repository state. These examples are used in error hint lines and tests.

## Canonical Anchor Examples

### Date/time anchors (time-based + tasks-completed metrics)
- Human-friendly: `Jan 28 2026 09:00 UTC`
- ISO 8601 (optional secondary reference): `2026-01-28T09:00:00Z`

Use the human-friendly anchor in hint lines; include the ISO 8601 variant only in documentation or tests that specifically validate ISO parsing.

### Commit anchors (commit + lines-changed metrics)
- Canonical anchor: `HEAD`

Use `HEAD` in hint lines to avoid run-to-run variability. If a test requires a fixed hash token, use a placeholder like `0123456789abcdef` (do not use `git rev-parse` output in hint examples).

## Canonical --every Examples by Metric

These are the stable examples to embed in hint lines:

- `days`/`weeks`/`months` (date anchor): `--every "10 days from Jan 28 2026 09:00 UTC"`
- `commits` (commit anchor): `--every "50 commits from HEAD"`
- `lines_changed` (commit anchor): `--every "500 lines_changed from HEAD"`
- `tasks_completed` (date anchor): `--every "20 tasks_completed from Jan 28 2026 09:00 UTC"`

Notes:
- Do not use `now` in hint examples.
- Keep amounts and metrics stable across tests to reduce churn.
- Prefer human-friendly anchors in hints; ISO 8601 is acceptable for validation examples.

## Implementation Notes
- Hint lines should be deterministic strings (no timestamps, no locale-dependent formatting).
- Examples must not be generated from runtime values.

## Local Verification Steps
1. Use the canonical examples above when adding tests for invalid anchors.
2. Assert that hint lines match exactly, including the anchor example.
