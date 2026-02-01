---
type: task
role: owner
priority: medium
parent: Tio6w-review-alternatives-commit-metric-default-anchor-w
blockers:
    - Tuu18-update-docs-and-hints-to-say-from-head
    - Tv5ra-implement-ignore-invalid-head-for-commit-metrics
blocks:
    - Tio6w-review-alternatives-commit-metric-default-anchor-w
date_created: 2026-01-29T22:24:02.846165Z
date_edited: 2026-02-01T04:22:23.490153Z
owner_approval: false
completed: true
---

# Decide default anchor wording for commit metrics (A/B/C)

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description


## Summary
Choose the default anchor wording for commit-based metrics in docs and hint examples.

## Context
- design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-alternatives.md
- design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-reliability-review.md
- design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-review.md

## Decision
- Select Alternative A, B, or C from the alternatives doc.
- Confirm whether the wording must explicitly mention the need for a valid HEAD or initial commit.

## Acceptance Criteria
- Decision recorded in the alternatives doc under the Decision section.
- Follow-up tasks created for doc/help text updates as needed.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [ ] (role: developer) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
- [x] (subtask: Tuu18) Update docs and hints to say from HEAD
- [ ] (subtask: Tv5ra) Implement ignore invalid HEAD for commit metrics
