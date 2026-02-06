---
type: task
role: designer
priority: medium
parent: I8lew-allow-reapplying-templates-to-existing-tasks
blockers: []
blocks: []
date_created: 2026-02-06T04:48:43.012277Z
date_edited: 2026-02-06T04:48:43.012277Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Design: Allow reapplying templates to existing tasks

## Summary
## Summary
Design the mechanism for reapplying templates to existing tasks. 
Consider:
- How to detect changes in templates.
- How to merge template structure with existing task content without data loss.
- Command syntax (e.g., `strand edit --template <type>` or a new command `strand reapply`).
- Handling of frontmatter vs body.
- Dry run and interactive merge options.

## Acceptance Criteria
- Design document created in `design-docs/`.
- Alternatives considered and pros/cons listed.
- Implementation plan (epics/tasks) defined.

## Instructions
Decide which task template would best fit this task and re-add it with that template and the same parent.
