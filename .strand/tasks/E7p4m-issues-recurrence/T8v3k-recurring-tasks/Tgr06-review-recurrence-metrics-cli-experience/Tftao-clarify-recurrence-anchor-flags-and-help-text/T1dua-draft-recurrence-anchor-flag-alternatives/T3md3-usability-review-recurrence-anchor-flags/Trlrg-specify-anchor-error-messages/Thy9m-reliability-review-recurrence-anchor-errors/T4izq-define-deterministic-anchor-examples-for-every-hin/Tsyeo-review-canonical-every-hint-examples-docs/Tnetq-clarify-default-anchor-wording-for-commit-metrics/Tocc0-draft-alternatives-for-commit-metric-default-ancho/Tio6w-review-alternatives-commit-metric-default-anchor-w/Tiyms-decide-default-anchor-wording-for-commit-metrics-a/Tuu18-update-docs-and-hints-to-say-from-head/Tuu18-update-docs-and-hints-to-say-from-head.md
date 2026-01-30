---
type: task
role: developer
priority: medium
parent: Tiyms-decide-default-anchor-wording-for-commit-metrics-a
blockers: []
blocks:
    - Tiyms-decide-default-anchor-wording-for-commit-metrics-a
date_created: 2026-01-29T22:32:32.821004Z
date_edited: 2026-01-29T15:32:32.830203-07:00
owner_approval: false
completed: false
---

# Update docs and hints to say from HEAD

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
Update user-facing docs and hint examples to use "from HEAD" for commit-based metrics and keep time-based metrics wording unchanged. Do not add explicit "valid HEAD" requirements in copy. Keep examples deterministic and update snapshots/tests as needed.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
Check this off one at a time with `strand complete <task_id> --role <my_given_role> --todo <num>` only if your Role matches the todo's role.
1. [ ] (role: developer) Implement the behavior described in Context.
2. [ ] (role: developer) Add unit and integration tests covering the main flows.
3. [ ] (role: tester) Execute test-suite and report failures.
4. [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
5. [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
Use subtasks for work that should be tracked separately or assigned to a different role. Use `strand add <type> "title" --parent <this_task_id> <<EOF description EOF`  to create subtasks.
