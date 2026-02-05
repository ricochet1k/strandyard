---
type: issue
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:38:57.07846Z
date_edited: 2026-02-05T04:38:57.07846Z
owner_approval: false
completed: false
description: ""
---

# Normalize short IDs to full IDs in relationships

## Summary
## Summary
Auto-normalize short IDs to full IDs when adding/editing relationships.

## Description
The validation logic currently requires full IDs in `parent` and `blockers` fields. If a user manually edits a file and uses a short ID, `repair` or `next` might fail.

## Requirements
- When parsing, if a short ID is found in a relationship field, resolve it to the full ID.
- When saving/repairing, rewrite the file with the full ID to maintain consistency.
