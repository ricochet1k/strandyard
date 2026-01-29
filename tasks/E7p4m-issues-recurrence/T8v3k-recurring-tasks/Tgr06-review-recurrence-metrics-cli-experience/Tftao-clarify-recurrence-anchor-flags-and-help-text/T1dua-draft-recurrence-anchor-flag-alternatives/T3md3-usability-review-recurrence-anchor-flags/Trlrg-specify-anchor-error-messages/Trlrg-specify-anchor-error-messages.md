---
type: task
role: designer
priority: medium
parent: T3md3-usability-review-recurrence-anchor-flags
blockers:
    - Tcsz3-review-alternatives-recurrence-anchor-error-messag
    - Thy9m-reliability-review-recurrence-anchor-errors
    - Ti8ig-usability-review-recurrence-anchor-errors
blocks:
    - T3md3-usability-review-recurrence-anchor-flags
date_created: 2026-01-29T05:46:59.618092Z
date_edited: 2026-01-29T13:01:12.794384-07:00
owner_approval: false
completed: false
---

# Specify anchor error messages

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
- Define user-facing error messages for missing or malformed anchors.
- Cover unit/anchor mismatches and ambiguity (when unit implies commit vs date).
- Provide example recovery hints for each error case.

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
