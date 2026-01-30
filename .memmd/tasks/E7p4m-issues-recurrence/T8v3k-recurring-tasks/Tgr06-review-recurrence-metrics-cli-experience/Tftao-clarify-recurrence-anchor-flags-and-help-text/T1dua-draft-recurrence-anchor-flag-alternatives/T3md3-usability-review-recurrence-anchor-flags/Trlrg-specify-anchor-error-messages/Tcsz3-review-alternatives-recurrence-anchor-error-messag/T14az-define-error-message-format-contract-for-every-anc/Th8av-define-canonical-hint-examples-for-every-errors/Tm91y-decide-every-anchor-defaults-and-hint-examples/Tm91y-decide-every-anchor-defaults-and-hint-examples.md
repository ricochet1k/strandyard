---
type: task
role: owner
priority: medium
parent: Th8av-define-canonical-hint-examples-for-every-errors
blockers: []
blocks:
    - Th8av-define-canonical-hint-examples-for-every-errors
date_created: 2026-01-29T17:14:34.825269Z
date_edited: 2026-01-29T17:59:49.21692Z
owner_approval: false
completed: true
---

# Questions: --every anchor defaults and hint examples

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
## Summary
Clarify whether hint examples should omit `from <anchor>` and whether hints may include relative or human-friendly dates.

## Context
- design-docs/recurrence-anchor-error-messages-alternatives.md
- design-docs/recurrence-anchor-flags-alternatives.md

## Questions Needed
- Should `from <anchor>` be optional in hints (defaulting to "now" or another implicit anchor)?
- Are relative/human-friendly date formats allowed in hint examples, or must hints be strictly deterministic ISO 8601?
- If defaults are allowed, what are the exact default anchors per metric?

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
