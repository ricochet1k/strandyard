---
type: task
role: developer
priority: medium
parent: T4izq-define-deterministic-anchor-examples-for-every-hin
blockers:
    - Tm2sq-review-canonical-every-hint-examples-implementatio
    - Tqb9o-approve-canonical-every-hint-examples
blocks:
    - T4izq-define-deterministic-anchor-examples-for-every-hin
date_created: 2026-01-29T19:24:46.390895Z
date_edited: 2026-01-29T12:24:46.400087-07:00
owner_approval: false
completed: false
---

# Implement deterministic --every hint examples

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description


## Summary
Implement deterministic --every hint examples using the canonical anchors defined in design-docs/recurrence-anchor-hint-examples.md.

## Acceptance Criteria
- Hint lines match canonical examples exactly
- Tests cover date/time and commit anchor hints

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [ ] Wire canonical hint examples into --every parsing errors
- [ ] Ensure hint strings are deterministic and stable across runs
- [ ] Add unit/integration tests that assert canonical hint examples
- [ ] (role: developer) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
