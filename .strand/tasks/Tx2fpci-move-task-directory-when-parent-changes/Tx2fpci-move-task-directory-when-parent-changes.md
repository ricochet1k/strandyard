---
type: issue
role: developer
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:40:06.019695Z
date_edited: 2026-02-05T23:13:38.911626Z
owner_approval: false
completed: false
description: ""
---

# Move task directory when parent changes

## Summary
Automatically move the task directory when the `parent` field is updated.

## Description
Currently, `strand edit --parent` updates the metadata but leaves the directory in the old location, violating the "directory hierarchy mirrors lineage" convention.

## Requirements
- When `parent` changes (via edit), move the directory to the new parent's directory.
- Handle potential name collisions (though unlikely with IDs).
- Ensure `repair` also detects and fixes this drift if possible.
