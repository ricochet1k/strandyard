---
type: task
role: designer
priority: medium
parent: T3md3-usability-review-recurrence-anchor-flags
blockers: []
blocks:
    - T3md3-usability-review-recurrence-anchor-flags
date_created: 2026-01-29T05:46:54.494094Z
date_edited: 2026-01-28T22:46:54.504258-07:00
owner_approval: false
completed: false
---

# Define anchor help text and examples

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
- Document per-unit anchor format mappings for `strand recurring add`.
- Propose concise `--help` text and CLI.md snippet with examples for time units, git units, tasks_completed, and lines_changed.
- Include guidance for users who skim help output (short summary + example).

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
