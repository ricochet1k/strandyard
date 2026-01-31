---
type: task
role: developer
priority: medium
parent: Tnetq-clarify-default-anchor-wording-for-commit-metrics
blockers: []
blocks:
    - Tnetq-clarify-default-anchor-wording-for-commit-metrics
date_created: 2026-01-29T20:04:09.214075Z
date_edited: 2026-01-29T13:04:09.223055-07:00
owner_approval: false
completed: false
---

# Add validation for non-ASCII quotes in task bodies

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description


## Summary
Add validation in task creation/repair to flag non-ASCII quotes in task bodies so manual edits are avoidable.

## Acceptance Criteria
- strand add warns or fails on smart quotes in body text.
- strand repair surfaces non-ASCII quotes when present.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [ ] Identify where body text is normalized during strand add/repair.
- [ ] Decide whether to reject or warn on non-ASCII quotes.
- [ ] Add validation and tests covering smart quote input.
- [ ] (role: developer) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
