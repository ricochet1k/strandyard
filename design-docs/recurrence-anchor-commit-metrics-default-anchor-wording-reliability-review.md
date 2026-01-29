---
role: reviewer-reliability
priority: medium
---

# Review: Default anchor wording for commit metrics (Reliability)

## Artifacts
- design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-alternatives.md

## Scope
Evaluate reliability impacts of wording choices for commit-metric default anchors, including determinism and alignment with actual default behavior.
Implementation details and UI copy beyond the described hints/docs are out of scope.

## Review Focus
- Deterministic phrasing across environments
- Alignment with actual default anchor behavior for commit metrics
- Impact on snapshot/golden tests and documentation drift

## Findings
- Alternative A (keep "from now") preserves uniform phrasing but risks time-based interpretation for commit metrics, which can lead to inconsistent user expectations and harder-to-debug failures.
- Alternative B ("from HEAD") is the most deterministic and explicit for git contexts, but it depends on the presence of a valid HEAD reference.
- Alternative C (define "now" mapping to HEAD) can preserve a single concept, but it introduces a second sentence that must stay consistent across hints and docs to avoid drift.

## Reliability Considerations
- Prefer fixed, repo-agnostic tokens (e.g., `HEAD`) and avoid hashes/timestamps in hints to keep tests stable.
- Ensure any mapping sentence is centralized (single source of truth) so help text and docs do not diverge.
- Wording should not imply time-based evaluation for commit metrics; clarify that the anchor is a commit reference.

## Concerns captured as subtasks
- Confirm and document behavior when `HEAD` is missing or detached so wording does not overpromise (T0m0p-confirm-head-missing-behavior-for-commit-metric-de).

## Decision
- Decision: deferred.
