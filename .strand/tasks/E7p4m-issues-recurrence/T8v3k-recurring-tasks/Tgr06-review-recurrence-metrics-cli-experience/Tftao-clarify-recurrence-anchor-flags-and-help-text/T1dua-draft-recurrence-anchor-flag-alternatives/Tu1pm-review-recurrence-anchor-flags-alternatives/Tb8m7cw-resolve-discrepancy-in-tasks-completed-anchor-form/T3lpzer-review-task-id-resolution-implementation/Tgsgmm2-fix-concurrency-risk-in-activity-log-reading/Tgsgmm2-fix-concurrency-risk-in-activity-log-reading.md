---
type: task
role: architect
priority: high
parent: T3lpzer-review-task-id-resolution-implementation
blockers: []
blocks: []
date_created: 2026-02-01T23:38:22.115795Z
date_edited: 2026-02-01T23:38:22.115795Z
owner_approval: false
completed: false
description: ""
---

# New Task: Fix concurrency risk in activity log reading

## Description
## Summary
`activity.Log.ReadEntries` currently closes the write-only file handle and reopens the file for reading, then reopens for writing. This is not thread-safe and can lead to data loss or errors if multiple processes (or the same process with multiple goroutines) interact with the log.

## Tasks
- [ ] Refactor `activity.Log` to support safe concurrent reads and writes.
- [ ] Ensure read operations do not interfere with the active write handle.

## Acceptance Criteria
- Activity log operations are thread-safe.
- No risk of data loss during log rotation or reopening.

Decide which task template would best fit this task and re-add it with that template and the same parent.
