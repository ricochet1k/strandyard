---
type: task
role: documentation
priority: medium
parent: T0m0p-confirm-head-missing-behavior-for-commit-metric-de
blockers: []
blocks:
    - T0m0p-confirm-head-missing-behavior-for-commit-metric-de
date_created: 2026-01-29T22:15:40.319082Z
date_edited: 2026-01-29T15:15:40.330632-07:00
owner_approval: false
completed: false
---

# Document HEAD-missing behavior in commit-metric defaults

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
## Summary
Capture the HEAD-missing/unborn behavior for commit-based recurrence metrics in user-facing docs once wording is finalized.

## Tasks
- [ ] Add a short note to CLI/docs explaining that commit-metric defaults rely on a valid HEAD.
- [ ] Provide a recovery hint for unborn repos (make initial commit or specify explicit anchor).

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
Check this off one at a time with `memmd complete <task_id> --role <my_given_role> --todo <num>` only if your Role matches the todo's role.
1. [ ] (role: developer) Implement the behavior described in Context.
2. [ ] (role: developer) Add unit and integration tests covering the main flows.
3. [ ] (role: tester) Execute test-suite and report failures.
4. [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
5. [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
Use subtasks for work that should be tracked separately or assigned to a different role. Use `memmd add <type> "title" --parent <this_task_id> <<EOF description EOF`  to create subtasks.
