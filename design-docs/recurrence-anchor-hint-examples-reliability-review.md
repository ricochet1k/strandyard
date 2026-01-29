---
role: reviewer-reliability
priority: medium
---

# Review: Canonical --every hint examples (Reliability)

## Artifacts
- design-docs/recurrence-anchor-hint-examples.md
- design-docs/recurrence-anchor-error-messages-alternatives.md
- design-docs/recurrence-anchor-error-messages-reliability-review.md

## Scope
Review deterministic hint examples for `memmd recurring add --every` anchor parsing.
Implementation details for parsing or formatting are out of scope.

## Review Focus
- Determinism and stability of example strings for automation
- Suitability for error output contracts and test fixtures
- Avoidance of locale/time-dependent formatting

## Findings
- Examples are stable, time-independent constants that avoid runtime values.
- Date anchors are explicitly UTC and human-friendly; ISO 8601 is limited to validation contexts.
- Commit anchor uses `HEAD` with a fixed placeholder for hash-only tests, avoiding repo-dependent values.

## Concerns captured as subtasks
- None.
