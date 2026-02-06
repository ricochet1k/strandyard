---
type: implement
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-06T05:05:35.635463Z
date_edited: 2026-02-06T05:05:35.635463Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Implement missing Task and Metadata status helpers

## Summary
## Summary
Add helper methods to Metadata and Task structs to easily check and set status.

## Deliverables
- IsOpen(), IsDone(), IsInProgress(), IsCancelled(), IsDuplicate() on Metadata and/or Task.
- SetStatus(string) helper that handles normalization and potentially sets the Completed boolean for backward compatibility.
- Move IsActiveStatus from free_list.go to status.go.

## Acceptance Criteria
- Implementation matches the specification
- Tests cover the change and pass
- Build succeeds
