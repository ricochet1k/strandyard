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

