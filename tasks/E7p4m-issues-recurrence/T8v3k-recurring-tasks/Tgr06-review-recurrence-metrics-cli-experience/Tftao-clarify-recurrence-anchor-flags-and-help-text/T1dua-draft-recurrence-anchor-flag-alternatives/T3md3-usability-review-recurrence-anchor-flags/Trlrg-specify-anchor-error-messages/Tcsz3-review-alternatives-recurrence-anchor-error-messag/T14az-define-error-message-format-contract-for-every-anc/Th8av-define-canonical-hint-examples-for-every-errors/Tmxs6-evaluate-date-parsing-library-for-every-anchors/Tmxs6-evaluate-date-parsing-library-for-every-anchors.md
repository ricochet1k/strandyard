---
type: task
role: developer
priority: medium
parent: Th8av-define-canonical-hint-examples-for-every-errors
blockers: []
blocks:
    - Th8av-define-canonical-hint-examples-for-every-errors
date_created: 2026-01-29T17:14:39.176361Z
date_edited: 2026-01-29T10:14:39.186089-07:00
owner_approval: false
completed: false
---

# Evaluate date parsing library for --every anchors

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
## Summary
Research and recommend a date parsing library that supports flexible, human-friendly inputs and is suitable for deterministic CLI parsing.

## Context
Potential requirement: allow relative date expressions in `--every` anchors.

## Acceptance Criteria
- Compare at least 2 candidate libraries with pros/cons (Go compatibility, license, determinism, locale behavior).
- Recommend one library and document constraints for CLI usage.

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
