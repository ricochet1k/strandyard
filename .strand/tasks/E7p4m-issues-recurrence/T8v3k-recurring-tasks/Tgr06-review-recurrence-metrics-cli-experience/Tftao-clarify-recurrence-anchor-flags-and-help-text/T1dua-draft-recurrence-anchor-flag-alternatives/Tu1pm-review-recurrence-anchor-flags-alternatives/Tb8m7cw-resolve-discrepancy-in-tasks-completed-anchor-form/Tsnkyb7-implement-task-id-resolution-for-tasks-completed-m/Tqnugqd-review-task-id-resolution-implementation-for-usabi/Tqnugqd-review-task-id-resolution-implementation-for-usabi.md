---
type: review
role: reviewer-usability
priority: medium
parent: Tsnkyb7-implement-task-id-resolution-for-tasks-completed-m
blockers:
    - T8ficaa-support-short-task-id-resolution-in-getlatesttaskc
    - Thi6rb5-fix-strand-add-failure-when-creating-subtasks
    - Tjs0ulo-update-design-docs-to-include-task-id-anchors-for
    - Tqfbr5g-resolve-task-id-anchors-to-full-ids-during-recurre
    - Tsnyb3q-add-tasks-completed-hint-to-every-flag-validation
blocks: []
date_created: 2026-02-05T00:22:50.461374Z
date_edited: 2026-02-05T00:45:40.729182Z
owner_approval: false
completed: true
description: ""
---

# Description

Review the task ID resolution implementation for usability concerns. Examine user-facing documentation, examples, error messages, and help text. Verify that users can easily understand how to use task ID anchors in recurrence metrics.

Delegate concerns to the relevant role via subtasks.



## Completion Report
Usability review complete. Concerns identified: (1) design-docs/recurrence-anchor-error-messages.md and anchor-help-text-and-examples-alternatives.md are outdated. (2) cmd/add.go lacks a hint for tasks_completed validation failures. (3) Error messages for tasks_completed can be confusing. (4) Short task IDs are not supported in recurrence anchors. (5) strand add fails when creating subtasks. Five subtasks created (Tjs0ulo, Tsnyb3q, T8ficaa, Tqfbr5g, Thi6rb5) to address these.

## Subtasks
- [x] (subtask: T8ficaa) New Task: Support short task ID resolution in GetLatestTaskCompletionTime
- [ ] (subtask: Te65o1c) Resolve task ID anchors to full IDs during recurrence validation
- [x] (subtask: Tefkvy2) Update design docs to include task ID anchors for tasks_completed
- [x] (subtask: Thi6rb5) New Task: Fix strand add failure when creating subtasks
- [x] (subtask: Tjs0ulo) New Task: Update design docs to include task ID anchors for tasks_completed
- [x] (subtask: Tm0vg3a) Support short task ID resolution in GetLatestTaskCompletionTime
- [x] (subtask: Tm555wi) Fix strand add failure when creating subtasks
- [x] (subtask: Tqfbr5g) New Task: Resolve task ID anchors to full IDs during recurrence validation
- [ ] (subtask: Tsnyb3q) New Task: Add tasks_completed hint to --every flag validation in cmd/add.go
