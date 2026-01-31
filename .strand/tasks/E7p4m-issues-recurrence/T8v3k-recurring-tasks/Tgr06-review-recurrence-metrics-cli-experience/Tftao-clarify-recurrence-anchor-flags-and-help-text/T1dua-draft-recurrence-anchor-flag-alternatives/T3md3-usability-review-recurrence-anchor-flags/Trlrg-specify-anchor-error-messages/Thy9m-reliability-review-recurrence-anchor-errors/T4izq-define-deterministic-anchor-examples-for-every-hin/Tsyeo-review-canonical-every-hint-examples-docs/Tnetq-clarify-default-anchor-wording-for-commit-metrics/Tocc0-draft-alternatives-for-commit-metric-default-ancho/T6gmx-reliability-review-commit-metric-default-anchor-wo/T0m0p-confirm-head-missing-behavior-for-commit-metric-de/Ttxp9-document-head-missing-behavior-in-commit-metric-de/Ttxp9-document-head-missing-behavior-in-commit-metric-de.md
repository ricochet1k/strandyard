---
type: task
role: documentation
priority: medium
parent: T0m0p-confirm-head-missing-behavior-for-commit-metric-de
blockers: []
blocks:
    - T0m0p-confirm-head-missing-behavior-for-commit-metric-de
date_created: 2026-01-29T22:15:40.319082Z
date_edited: 2026-01-29T22:18:52.660816Z
owner_approval: false
completed: true
---

# Document HEAD-missing behavior in commit-metric defaults

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description


## Summary
Capture the HEAD-missing/unborn behavior for commit-based recurrence metrics in user-facing docs once wording is finalized.

## Acceptance Criteria
- CLI/docs mention the HEAD requirement for commit-based default anchors.
- Wording aligns with the chosen default-anchor alternative.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [ ] Add a short note to CLI/docs explaining that commit-metric defaults rely on a valid HEAD.
- [ ] Provide a recovery hint for unborn repos (make initial commit or specify explicit anchor).
- [ ] (role: developer) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
