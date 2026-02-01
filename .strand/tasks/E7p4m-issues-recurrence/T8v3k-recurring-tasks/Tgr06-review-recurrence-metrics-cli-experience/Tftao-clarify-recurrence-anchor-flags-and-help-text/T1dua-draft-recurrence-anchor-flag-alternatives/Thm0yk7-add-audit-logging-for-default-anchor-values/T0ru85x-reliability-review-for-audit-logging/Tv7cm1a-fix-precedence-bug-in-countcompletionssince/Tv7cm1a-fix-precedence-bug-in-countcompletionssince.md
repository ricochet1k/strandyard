---
type: issue
role: developer
priority: medium
parent: T0ru85x-reliability-review-for-audit-logging
blockers: []
blocks: []
date_created: 2026-02-01T21:44:55.318909Z
date_edited: 2026-02-01T21:45:42.773416Z
owner_approval: false
completed: true
description: ""
---

# Fix precedence bug in CountCompletionsSince

## Summary
In pkg/activity/log.go, the logical condition in CountCompletionsSince has a precedence bug:
if entry.Type == EventTaskCompleted && entry.Timestamp.After(since) || entry.Timestamp.Equal(since)

It should be:
if entry.Type == EventTaskCompleted && (entry.Timestamp.After(since) || entry.Timestamp.Equal(since))

## Completion Report
Fixed precedence bug in CountCompletionsSince by adding parentheses around the timestamp comparison.
