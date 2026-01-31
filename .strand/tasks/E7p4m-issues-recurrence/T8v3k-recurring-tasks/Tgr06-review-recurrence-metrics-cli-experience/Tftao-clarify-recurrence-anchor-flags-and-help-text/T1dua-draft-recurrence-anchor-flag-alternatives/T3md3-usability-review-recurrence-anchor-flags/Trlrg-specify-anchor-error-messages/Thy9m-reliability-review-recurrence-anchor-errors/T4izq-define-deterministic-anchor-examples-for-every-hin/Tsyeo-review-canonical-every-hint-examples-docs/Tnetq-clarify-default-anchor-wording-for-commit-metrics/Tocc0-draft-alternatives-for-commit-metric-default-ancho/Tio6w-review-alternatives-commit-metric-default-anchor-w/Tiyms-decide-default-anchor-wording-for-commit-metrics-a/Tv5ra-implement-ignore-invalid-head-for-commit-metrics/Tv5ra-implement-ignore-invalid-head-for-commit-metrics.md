---
type: task
role: developer
priority: medium
parent: Tiyms-decide-default-anchor-wording-for-commit-metrics-a
blockers: []
blocks:
    - Tiyms-decide-default-anchor-wording-for-commit-metrics-a
date_created: 2026-01-29T22:32:28.748147Z
date_edited: 2026-01-29T15:32:28.757968-07:00
owner_approval: false
completed: false
---

# Implement ignore invalid HEAD for commit metrics

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
Update commit-based recurrence metrics to treat invalid or unborn HEAD as a no-op (skip triggering recurring tasks) rather than erroring. Align behavior with owner decision in design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-alternatives.md. Add unit/integration tests covering unborn HEAD and detached HEAD cases.

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
