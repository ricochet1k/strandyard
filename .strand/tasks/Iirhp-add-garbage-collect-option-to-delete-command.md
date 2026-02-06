---
type: issue
role: triage
priority: high
parent: ""
blockers:
    - Tvbdhas-design-gc-command-for-status-based-garbage-collect
blocks: []
date_created: 2026-02-06T00:08:30.986484Z
date_edited: 2026-02-06T04:53:49.056196Z
owner_approval: false
completed: true
status: done
description: ""
---

# Add garbage-collect option to delete command

## Summary
The delete command currently removes a single task, but there is no built-in way to sweep completed or aged tasks after a cleanup window. Add a garbage-collect mode or `--old-completed` flag that can target completed tasks or tasks older than a configured age, optionally respecting projects or parents. This will help keep `.strand/tasks` manageable without manually deleting files.

## Acceptance Criteria
- Issue still exists
- Issue is fixed and verified locally
- Tests pass
- Build succeeds

## Subtasks
- [ ] (subtask: Tvbdhas) Design gc command for status-based garbage collection

## Completion Report
Confirmed that the 'delete' command is missing and that the project previously decided to use status states instead of hard deletes. However, the need for mass garbage collection is valid. Created a design task Tvbdhas for an architect to design a 'gc' command that aligns with the status-based system.
