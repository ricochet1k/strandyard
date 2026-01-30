---
type: task
role: owner
priority: medium
parent: Tio6w-review-alternatives-commit-metric-default-anchor-w
blockers:
    - Tuu18-update-docs-and-hints-to-say-from-head
    - Tv5ra-implement-ignore-invalid-head-for-commit-metrics
blocks:
    - Tio6w-review-alternatives-commit-metric-default-anchor-w
date_created: 2026-01-29T22:24:02.846165Z
date_edited: 2026-01-30T02:22:10.222134Z
owner_approval: false
completed: true
---

# Decide default anchor wording for commit metrics (A/B/C)

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
## Summary
Choose the default anchor wording for commit-based metrics in docs and hint examples.

## Context
- design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-alternatives.md
- design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-reliability-review.md
- design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-review.md

## Decision
- Select Alternative A, B, or C from the alternatives doc.
- Confirm whether the wording must explicitly mention the need for a valid HEAD or initial commit.

## Tasks

- [ ] (subtask: Tuu18-update-docs-and-hints-to-say-from-head) Update docs and hints to say from HEAD
- [ ] (subtask: Tv5ra-implement-ignore-invalid-head-for-commit-metrics) Implement ignore invalid HEAD for commit metrics

## Acceptance Criteria
- Decision recorded in the alternatives doc under the Decision section.
- Follow-up tasks created for doc/help text updates as needed.

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
