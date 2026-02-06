---
type: issue
role: triage
priority: medium
parent: ""
blockers:
    - Twv7603-implement-role-opt-in-filtering-for-strand-next
blocks: []
date_created: 2026-01-29T22:19:55.48898Z
date_edited: 2026-02-06T07:36:26.408091Z
owner_approval: false
completed: true
status: done
description: ""
---

# Support role opt-in filtering for strand next

## Summary
Add role metadata (likely frontmatter) to let `strand next` skip roles marked as opt-in by default unless explicitly requested (for example via `strand next --role <role>`). Define the role tag, update parsing, and ensure owner tasks can be run in order when requested.

## Context
Provide relevant logs, links, and environment details.

## Impact
Describe severity and who/what is affected.

## Acceptance Criteria
- Define what fixes or mitigations are required.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Subtasks
- [ ] (subtask: Twv7603) Implement role opt-in filtering for strand next

## Completion Report
Confirmed that 'owner' role is hardcoded in cmd/next.go and there is no way to make other roles opt-in. Created implementation task Twv7603 to add 'opt_in' metadata to roles and update 'strand next' to use it.
