---
type: task
role: designer
priority: medium
parent: Tnetq-clarify-default-anchor-wording-for-commit-metrics
blockers:
    - T6gmx-reliability-review-commit-metric-default-anchor-wo
    - Tfxi6-usability-review-commit-metric-default-anchor-word
    - Tio6w-review-alternatives-commit-metric-default-anchor-w
    - Tvu5e-usability-review-commit-metric-default-anchor-word
blocks:
    - Tnetq-clarify-default-anchor-wording-for-commit-metrics
date_created: 2026-01-29T19:59:33.996756Z
date_edited: 2026-01-29T19:22:43.455133-07:00
owner_approval: false
completed: false
---

# Draft alternatives for commit-metric default anchor wording

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
## Summary
Draft alternatives and decision points for how docs and hint examples describe the default anchor for commit-based recurrence metrics.

## Tasks

- [ ] Write alternatives doc: design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-alternatives.md
- [ ] Request reviews from master reviewer and relevant specialized reviewers
- [ ] Capture owner decision and update design docs after selection

- [x] (subtask: T6gmx-reliability-review-commit-metric-default-anchor-wo) Reliability review: commit-metric default anchor wording
- [x] (subtask: Tfxi6-usability-review-commit-metric-default-anchor-word) Usability review: commit-metric default anchor wording
- [x] (subtask: Tio6w-review-alternatives-commit-metric-default-anchor-w) Review alternatives: commit-metric default anchor wording
- [ ] (subtask: Tvu5e-usability-review-commit-metric-default-anchor-word) Usability review: commit-metric default anchor wording

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
