# Canonical Hint Examples for --every Anchors

## Summary
Define deterministic, reusable hint examples for `memmd recurring add --every` anchor parsing errors so tests and automation remain stable.

## Context
- design-docs/recurrence-anchor-error-messages-alternatives.md
- design-docs/recurrence-anchor-error-messages-reliability-review.md
- design-docs/recurrence-anchor-flags-alternatives.md
- design-docs/recurrence-metrics.md

## Decision
Approved: prefer anchor-less examples that use the default "from now" behavior for hint lines. Include a small number of explicit anchor examples to demonstrate supported syntax. These examples are used in error hint lines and tests.

## Canonical Anchor Examples

### Date/time anchors (time-based + tasks-completed metrics)
- Human-friendly: `Jan 28 2026 09:00 UTC`
- ISO 8601 (optional secondary reference): `2026-01-28T09:00:00Z`

Use the human-friendly anchor in explicit anchor examples; include the ISO 8601 variant only in documentation or tests that specifically validate ISO parsing.

### Commit anchors (commit + lines-changed metrics)
- Canonical anchor: `HEAD`

Use `HEAD` in explicit anchor examples. If a test requires a fixed hash token, use a placeholder like `0123456789abcdef` (do not use `git rev-parse` output in hint examples).

## Canonical --every Examples by Metric

These are the stable examples to embed in hint lines. Most examples omit the anchor to rely on the default "from now" behavior.

- `days`/`weeks`/`months` (default anchor): `--every "10 days"`
- `commits` (default anchor): `--every "50 commits"`
- `lines_changed` (default anchor): `--every "500 lines_changed"`
- `tasks_completed` (default anchor): `--every "20 tasks_completed"`

Explicit anchor examples (use sparingly in hints/tests):
- `days` (date anchor): `--every "10 days from Jan 28 2026 09:00 UTC"`
- `commits` (commit anchor): `--every "50 commits from HEAD"`

Notes:
- Keep amounts and metrics stable across tests to reduce churn.
- Prefer human-friendly anchors when you do specify one; ISO 8601 is acceptable for validation examples.

## Implementation Notes
- Hint lines should prefer the default anchor (no explicit anchor) to reflect typical usage.
- When an explicit anchor is included, use deterministic strings (no timestamps, no locale-dependent formatting).

## Local Verification Steps
1. Use the canonical examples above when adding tests for invalid anchors.
2. Assert that hint lines match exactly, including the anchor example.
