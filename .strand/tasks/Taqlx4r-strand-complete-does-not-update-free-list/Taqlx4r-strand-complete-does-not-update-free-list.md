---
type: issue
role: triage
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-04T23:21:02.366797Z
date_edited: 2026-02-05T04:07:06.454724Z
owner_approval: false
completed: false
description: ""
---

# strand complete does not update free-list

## Summary


## Description
When a task is marked as completed using `strand complete`, it remains in the "free" list used by `strand next` until a manual `strand repair` is run.

## Expected Behavior
Any command that modifies task completion status (like `complete`) should automatically update the free-list index.

## Technical Details
The implementation of the complete command should ensure it uses TaskDB APIs that handle free-list maintenance internally, or explicitly trigger a free-list update upon status changes to ensure the cache remains consistent without requiring manual repairs.
