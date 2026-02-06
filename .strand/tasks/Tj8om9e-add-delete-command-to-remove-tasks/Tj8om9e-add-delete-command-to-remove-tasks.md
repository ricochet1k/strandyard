---
type: issue
role: architect
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T04:38:56.456015Z
date_edited: 2026-02-06T04:31:22.881677Z
owner_approval: false
completed: true
status: ""
description: ""
---

# Consider: Remove delete command, use status states instead

## Summary
Reconsider the need for a hard delete command in light of status states.

## Description
With the introduction of proper status states (`open`, `done`, `duplicate`, `cancelled`), there may be no need for a hard delete:
- Tasks marked `duplicate` are clearly identified as redundant
- Tasks marked `cancelled` are clearly identified as won't-do
- Marking a task with these statuses rather than deleting preserves history and relationships

## Questions
- Should we support permanent deletion, or rely on status states?
- If deletion is needed, should it only work for tasks with status `duplicate` or `cancelled`?
- Should there be an archive command to hide completed/obsolete tasks?

## Acceptance Criteria
Decision recorded in design doc or AGENTS.md about whether hard delete is needed.
