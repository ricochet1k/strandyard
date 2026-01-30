---
type: task
role: documentation
priority: medium
parent: T4izq-define-deterministic-anchor-examples-for-every-hin
blockers:
    - Tqb9o-approve-canonical-every-hint-examples
    - Tsyeo-review-canonical-every-hint-examples-docs
blocks:
    - T4izq-define-deterministic-anchor-examples-for-every-hin
date_created: 2026-01-29T19:24:50.680451Z
date_edited: 2026-01-29T19:57:33.032262Z
owner_approval: false
completed: true
---

# Document canonical --every hint examples

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
## Summary
Update CLI documentation to reflect the canonical --every hint examples and anchor guidance.

## Tasks
- [ ] Update CLI.md or related docs with the canonical examples
- [ ] Ensure hint examples are consistent with design-docs/recurrence-anchor-hint-examples.md

## Acceptance Criteria
- Documentation references deterministic examples for date/time and commit anchors

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
