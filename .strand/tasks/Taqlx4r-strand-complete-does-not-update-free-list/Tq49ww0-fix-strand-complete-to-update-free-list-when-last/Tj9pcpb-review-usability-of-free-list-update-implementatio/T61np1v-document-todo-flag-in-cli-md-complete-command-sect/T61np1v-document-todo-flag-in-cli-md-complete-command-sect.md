---
type: task
role: documentation
priority: medium
parent: Tj9pcpb-review-usability-of-free-list-update-implementatio
blockers: []
blocks: []
date_created: 2026-02-05T12:10:09.964915Z
date_edited: 2026-02-05T12:10:09.964915Z
owner_approval: false
completed: false
description: ""
---

# New Task: Document --todo flag in CLI.md complete command section

## Description
## Context
The `strand complete` command supports a `--todo` flag to mark individual TODO items as completed, but this is not documented in CLI.md.

## Current Documentation Gap
The CLI.md file only shows basic usage:
```bash
strand complete <task-id> [report]
```

It doesn't mention:
- `--todo` flag for checking off individual todo items
- `--role` flag for role validation
- Behavior when completing the last TODO item (task auto-completes)

## Acceptance Criteria
- `--todo` flag is documented with examples
- `--role` flag is documented with examples
- Behavior of auto-completing a task via last TODO is explained
- Free-list update behavior is documented
- Example shows completing a todo item

Decide which task template would best fit this task and re-add it with that template and the same parent.
