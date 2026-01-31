---
type: task
role: owner
priority: medium
parent: Thy9m-reliability-review-recurrence-anchor-errors
blockers: []
blocks:
    - Thy9m-reliability-review-recurrence-anchor-errors
date_created: 2026-01-29T19:20:55.261026Z
date_edited: 2026-01-29T12:20:55.271287-07:00
owner_approval: false
completed: false
---

# Align exit code conventions for --every failures

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description


## Summary
Confirm and document CLI-wide exit code conventions for `--every` parse/validation failures.

## Details
- Decide whether exit code `2` is reserved for parse/validation errors across CLI commands.
- Document the contract for `--every` alongside other error cases (parse vs runtime).
- Ensure automation guidance remains stable for scripts and CI.

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
