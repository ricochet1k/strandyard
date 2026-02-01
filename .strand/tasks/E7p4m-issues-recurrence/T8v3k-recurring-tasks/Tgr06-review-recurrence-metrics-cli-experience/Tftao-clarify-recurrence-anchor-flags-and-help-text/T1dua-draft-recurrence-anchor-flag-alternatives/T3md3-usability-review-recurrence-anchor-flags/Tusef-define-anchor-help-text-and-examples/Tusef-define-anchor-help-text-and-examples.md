---
type: task
role: designer
priority: medium
parent: T3md3-usability-review-recurrence-anchor-flags
blockers:
    - Tcc9bjk-review-anchor-help-text-and-examples-alternatives
    - Tt3kj2u-usability-review-anchor-help-text-and-examples
blocks:
    - T3md3-usability-review-recurrence-anchor-flags
date_created: 2026-01-29T05:46:54.494094Z
date_edited: 2026-02-01T20:21:11.970741Z
owner_approval: false
completed: false
description: ""
---

# Define anchor help text and examples

## Context
- design-docs/anchor-help-text-and-examples-alternatives.md
- design-docs/recurrence-anchor-flags-alternatives.md
- design-docs/recurrence-anchor-error-messages.md
- CLI.md (recurring add section — needs update)

## Description
- Document per-unit anchor format mappings for `strand add --every`.
- Propose concise `--help` text and CLI.md snippet with examples for time units, git units, tasks_completed, and lines_changed.
- Include guidance for users who skim help output (short summary + example).

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] Draft alternatives document with per-unit anchor mappings
- [x] Propose `--help` text alternatives (compact, grouped, minimal)
- [x] Propose CLI.md examples and tables
- [x] Request review from master-reviewer and reviewer-usability
- [ ] Awaiting review feedback and owner decision

## Subtasks
- [ ] (subtask: Tcc9bjk) Description
- [ ] (subtask: Tt3kj2u) New Task: Usability review: anchor help text and examples
