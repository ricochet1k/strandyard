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
- Document per-unit anchor format mappings for `memmd recurring add`.
- Propose concise `--help` text and CLI.md snippet with examples for time units, git units, tasks_completed, and lines_changed.
- Include guidance for users who skim help output (short summary + example).

## Escalation
If new concerns or decisions arise, create follow-up tasks instead of editing this task.

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
