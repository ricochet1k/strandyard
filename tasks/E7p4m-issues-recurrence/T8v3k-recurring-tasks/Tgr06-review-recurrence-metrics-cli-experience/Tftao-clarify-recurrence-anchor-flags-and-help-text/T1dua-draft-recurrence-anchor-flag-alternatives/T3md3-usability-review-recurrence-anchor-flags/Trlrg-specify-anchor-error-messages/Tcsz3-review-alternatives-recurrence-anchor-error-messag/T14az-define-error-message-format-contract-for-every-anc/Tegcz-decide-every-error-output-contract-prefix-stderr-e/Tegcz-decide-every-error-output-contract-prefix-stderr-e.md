---
type: task
role: owner
priority: medium
parent: T14az-define-error-message-format-contract-for-every-anc
blockers: []
blocks:
    - T14az-define-error-message-format-contract-for-every-anc
date_created: 2026-01-29T16:55:16.051856Z
date_edited: 2026-01-29T17:05:51.143718Z
owner_approval: false
completed: true
---

# Decide --every error output contract (prefix, stderr, exit code)

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
## Summary
Define the user-visible error output contract for `--every` anchor parsing, including the exact prefix wording, stderr vs stdout behavior, and exit code mapping for parse/validation failures.

## Acceptance Criteria
- Error prefix string is specified and stable.
- stderr/stdout behavior is documented.
- Exit code mapping is documented.

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
