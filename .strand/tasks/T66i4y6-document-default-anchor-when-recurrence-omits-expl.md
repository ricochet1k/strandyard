---
type: issue
role: designer
priority: medium
parent: T3xpql7-clarify-missing-anchor-default-behavior-in-recurre
blockers: []
blocks: []
date_created: 2026-02-15T04:23:36.594526Z
date_edited: 2026-02-15T04:29:14.890137Z
owner_approval: false
completed: true
status: done
description: ""
---

# Document default anchor when recurrence omits explicit anchor

## Summary


## Summary
Clarify the user-facing default anchor behavior when `--every` is provided without `from/after`, especially for `tasks_completed`.

## Problem
Current docs and help text describe defaults for time-based and git-based metrics, but do not clearly state what happens for `tasks_completed` when no anchor is provided. This leaves users unsure whether recurrence starts from now, from latest completion, or from task creation.

## Acceptance Criteria
- CLI docs explicitly state the default anchor for each metric family, including `tasks_completed`.
- At least one example shows anchor-less `tasks_completed` with plain-language behavior.
- Wording avoids ambiguous phrases like "from now" without defining the effective anchor timestamp.

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds

## Completion Report
Clarified recurrence default-anchor docs and help text across metric families, explicitly documenting tasks_completed defaulting to the current UTC evaluation timestamp when no anchor is provided. Added anchor-less tasks_completed examples in command help and CLI docs and verified with go test ./... and go build ./....
