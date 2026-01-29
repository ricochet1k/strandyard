---
type: task
role: architect
priority: medium
parent: Thy9m-reliability-review-recurrence-anchor-errors
blockers: []
blocks:
    - Thy9m-reliability-review-recurrence-anchor-errors
date_created: 2026-01-29T19:20:50.103845Z
date_edited: 2026-01-29T12:20:50.114069-07:00
owner_approval: false
completed: false
---

# Define deterministic anchor examples for --every hints

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
## Summary
Define canonical, deterministic anchor examples for each anchor type in `--every` hint lines.

## Details
- Specify fixed examples for date/time anchors (human-friendly + optional ISO 8601).
- Specify fixed examples for commit anchors (e.g., `HEAD` or a placeholder hash) that do not vary per run.
- Ensure examples are stable for tests and automation (no current-time rendering).

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
Check this off one at a time with `memmd complete <task_id> --role <my_given_role> --todo <num>` only if your Role matches the todo's role.
1. [ ] (role: developer) Implement the behavior described in Context.
2. [ ] (role: developer) Add unit and integration tests covering the main flows.
3. [ ] (role: tester) Execute test-suite and report failures.
4. [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
5. [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
Use subtasks for work that should be tracked separately or assigned to a different role. Use `memmd add <type> "title" --parent <this_task_id> <<EOF description EOF`  to create subtasks.
