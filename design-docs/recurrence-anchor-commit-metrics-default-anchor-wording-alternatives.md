# Design Alternatives — Default anchor wording for commit metrics

## Summary
Clarify how documentation and hint examples describe the default anchor for commit-based recurrence metrics (`commits`, `lines_changed`) so users understand what an omitted anchor means. This document compares wording options for hints and docs while keeping examples deterministic.

## Project Principles
- Keep user-facing examples deterministic and stable for tests.
- Prefer wording that reduces ambiguity and improves recovery from errors.
- Align docs with actual behavior without adding repo-dependent content.

## HEAD availability (observed behavior)
- Commit-metric defaults rely on resolving `HEAD` (per `git rev-list --count <anchor>..HEAD` and `git diff --numstat <anchor>..HEAD`).
- In repos with an unborn `HEAD` (no commits), both commands fail with "ambiguous argument 'HEAD..HEAD'".
- In detached `HEAD` with a valid commit, both commands succeed; counts are `0` when anchor and `HEAD` are the same.
- Implication: any wording that implies a default `HEAD` anchor must note that a valid `HEAD` is required; otherwise the CLI should error and instruct the user to provide an explicit anchor or create an initial commit.

## Local verification steps
These steps mirror the git commands used by commit metrics in `design-docs/recurrence-metrics.md`.

Unborn `HEAD` (no commits):
```bash
mkdir -p /tmp/memmd-head-test-unborn
cd /tmp/memmd-head-test-unborn
git init
git rev-list --count HEAD..HEAD
git diff --numstat HEAD..HEAD
```
Expected: both commands error with "ambiguous argument 'HEAD..HEAD'".

Detached `HEAD` (valid commit):
```bash
mkdir -p /tmp/memmd-head-test-detached
cd /tmp/memmd-head-test-detached
git init
printf "test\n" > README.md
git add README.md
git -c user.name="memmd" -c user.email="memmd@example.com" commit -m "init"
git checkout --detach HEAD
git rev-list --count HEAD..HEAD
git diff --numstat HEAD..HEAD
```
Expected: commands succeed; rev-list returns `0` and diff output is empty.

## Alternatives

### Alternative A — Use "from now" wording for all metrics
- Description: Keep the existing "default is from now" phrasing for all metrics and avoid commit-specific wording in hints. If needed, add a brief note elsewhere in docs about commit metrics using `HEAD` internally.
- Assumptions: Users accept "now" as a generic anchor across metric types; the nuance that git metrics map to `HEAD` is optional.
- Pros:
  - Consistent language across all metrics.
  - Minimal doc churn; no hint string changes needed.
  - Keeps hint lines short.
- Cons:
  - "Now" is ambiguous for commit metrics and may imply time-based behavior.
  - Users may not realize that commit metrics use `HEAD` unless they read additional notes.
- Risks:
  - Increased support burden or confusion during error recovery.
- Rough effort estimate: 0.5–1 hour (if adding a small note).

### Alternative B — Use commit-specific default wording ("from HEAD")
- Description: For commit-based metrics, explicitly state the default anchor as `HEAD` in hints and docs. Keep "from now" wording for time-based metrics.
- Assumptions: Users benefit from explicit git terminology; "HEAD" is an acceptable default anchor reference.
- Pros:
  - Clear and actionable for commit metrics.
  - Aligns with git mental model and avoids "now" ambiguity.
  - Reinforces deterministic examples (no timestamps or hashes).
- Cons:
  - Introduces special-case wording per metric type.
  - Slightly longer hint strings or help text.
- Risks:
  - Inconsistent phrasing across metric types may look uneven in docs.
- Rough effort estimate: 1–2 hours (update hints/docs/tests to match).

### Alternative C — Define "now" as a metric-specific mapping
- Description: Keep "from now" as the primary wording but add a standardized sentence that defines "now" per metric: "For git metrics, 'now' means `HEAD`." Use this line consistently in docs and hint explanations.
- Assumptions: Users will read the clarification when presented; a single conceptual anchor is easier to teach.
- Pros:
  - Preserves a single “from now” concept across metrics.
  - Adds explicit mapping for commit metrics without rewriting examples.
  - Reduces ambiguity while keeping hints short.
- Cons:
  - Requires extra explanatory text in docs/hints.
  - Users may still miss the mapping if only seeing short examples.
- Risks:
  - Docs/hints could become verbose if repeated in multiple places.
- Rough effort estimate: 1–2 hours (add mapping sentence, ensure consistency).

## Decision
- Decision: Alternative B (use "from HEAD" for commit-based metrics).
- Default wording should not explicitly call out valid `HEAD` requirements.
- Behavior alignment: if `HEAD` is invalid or unborn, commit-based recurrence metrics are ignored and do not trigger recurring tasks.
- Rationale: explicit `HEAD` wording is the clearest and most deterministic for git metrics while keeping hints concise; treating invalid `HEAD` as a no-op avoids confusing errors and prevents accidental recurrences when the repo is not initialized.

## Post-decision cleanup
- Update hint examples and docs to use the chosen wording.
- Ensure any tests or snapshots with hint text are updated deterministically.
- Implement and test the "invalid HEAD is ignored" behavior for commit-based recurrence.

## Review Requests
- Request review from: `master-reviewer`, `reviewer-usability`, `reviewer-reliability`.
