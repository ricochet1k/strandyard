---
type: issue
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:38:57.21779Z
date_edited: 2026-02-06T04:32:10.0118Z
owner_approval: false
completed: true
status: ""
description: ""
---

# Add --blocks flag to edit command

## Summary


## Summary
Allow editing the 'blocks' relationship via CLI.

## Description
`strand edit` supports `--blocker` (dependencies) but not `--blocks` (reverse dependencies).

## Requirements
- Add `--blocks` flag to `strand edit`.
- Update the referenced tasks to add the current task to their `blockers` list.
