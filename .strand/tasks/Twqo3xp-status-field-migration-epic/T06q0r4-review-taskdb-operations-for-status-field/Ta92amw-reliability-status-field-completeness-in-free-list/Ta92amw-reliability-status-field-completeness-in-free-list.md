---
type: task
role: architect
priority: high
parent: T06q0r4-review-taskdb-operations-for-status-field
blockers:
    - T4p2p6o-add-integration-tests-for-status-field-with-free-l
    - Tblt03z-update-free-list-generation-to-check-status-field
    - Tg9s49n-add-free-list-status-validation-rules
blocks: []
date_created: 2026-02-05T22:06:06.5552Z
date_edited: 2026-02-05T22:06:50.024883Z
owner_approval: false
completed: true
description: ""
---

# New Task: Reliability: Status field completeness in free-list calculation

## Description


## Overview
Ensure the free-list calculation (`CalculateIncrementalFreeListUpdate()`) correctly handles the new status field and validates that only `open` or `in_progress` tasks appear in the free-list.

## Context
As part of the Status Field Migration epic, the `completed: bool` field is being replaced with multi-state `status` values. The free-list is a critical mechanism that lists tasks with no blockers (ready to work on), and it must exclude tasks with non-open statuses.

## Planning Tasks
1. Review free-list generation logic in `pkg/task/free_list.go`
2. Define validation rules for free-list inclusion based on status field
3. Identify all code paths that read/update the free-list
4. Design test strategy for status field integration with free-list
5. Create child implementation tasks for developers

## Deliverables
- Updated `free_list.go` to check status field instead of completed bool
- Validation rules documented in code comments
- Child tasks created for each implementation concern
- Design doc updated with free-list handling strategy

Decide which task template would best fit this task and re-add it with that template and the same parent.

## Subtasks
- [ ] (subtask: T4p2p6o) Add integration tests for status field with free-list
- [ ] (subtask: Tblt03z) Update free-list generation to check status field
- [ ] (subtask: Tg9s49n) Add free-list status validation rules

## Completion Report
Broke down free-list status field reliability requirement into 3 implementable child tasks:

1. Tblt03z: Update free-list generation to check status field - change free-list calculation to exclude non-active statuses
2. Tg9s49n: Add free-list status validation rules - ensure free-list consistency via validation in repair command
3. T4p2p6o: Add integration tests - comprehensive test coverage for status field with free-list

Created design document (free-list-status-field-handling.md) with:
- Status value eligibility rules (open/in_progress = included, done/cancelled/duplicate = excluded)
- Implementation strategy across three phases
- Impact analysis on affected commands (next, complete, list, repair)
- Edge cases and acceptance criteria

All child tasks are properly scoped for developers and testers to implement independently.
