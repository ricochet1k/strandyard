---
type: issue
role: reviewer-usability
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-01T20:16:22.968702Z
date_edited: 2026-02-15T05:12:33.908815Z
owner_approval: false
completed: true
status: done
description: ""
---

# Clarify missing anchor default behavior in recurrence

## Summary


## Completion Report
Usability review complete. Concerns: default anchor behavior for anchor-less tasks_completed rules is unclear in user-facing docs/help; filed follow-up T66i4y6 to clarify defaults and examples.

## Subtasks
- [x] (subtask: T66i4y6) Document default anchor when recurrence omits explicit anchor
  Clarified recurrence default-anchor docs and help text across metric families, explicitly documenting tasks_completed defaulting to the current UTC evaluation timestamp when no anchor is provided. Added anchor-less tasks_completed examples in command help and CLI docs and verified with go test ./... and go build ./....
