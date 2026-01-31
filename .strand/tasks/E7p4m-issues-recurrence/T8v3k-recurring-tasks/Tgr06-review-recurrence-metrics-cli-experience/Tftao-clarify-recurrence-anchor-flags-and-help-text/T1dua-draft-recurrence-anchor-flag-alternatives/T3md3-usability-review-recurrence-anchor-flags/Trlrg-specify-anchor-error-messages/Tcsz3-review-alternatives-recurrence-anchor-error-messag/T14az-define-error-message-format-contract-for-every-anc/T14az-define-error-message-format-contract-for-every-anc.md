---
type: task
role: reviewer-reliability
priority: medium
parent: Tcsz3-review-alternatives-recurrence-anchor-error-messag
blockers:
    - Tegcz-decide-every-error-output-contract-prefix-stderr-e
    - Th8av-define-canonical-hint-examples-for-every-errors
blocks:
    - Tcsz3-review-alternatives-recurrence-anchor-error-messag
date_created: 2026-01-29T15:28:01.821675Z
date_edited: 2026-01-31T17:29:31.078577Z
owner_approval: false
completed: true
---

# Define error message format contract for --every anchor parsing

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description


## Escalation
If new concerns or decisions arise, create follow-up tasks instead of editing this task.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [ ] (role: developer) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
- [x] (subtask: Tegcz) Decide --every error output contract (prefix, stderr, exit code)
- [ ] (subtask: Th8av) Define canonical hint examples for --every errors
