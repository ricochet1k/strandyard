---
type: issue
role: architect
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-02-05T11:50:12.275258Z
date_edited: 2026-02-05T22:03:15.823716Z
owner_approval: false
completed: true
description: ""
---

# Replace completed boolean with status field (open/in_progress/done/duplicate/cancelled)

## Summary


## Summary
Replace the `completed: true/false` boolean with a multi-state `status` field.

## Description
A simple boolean doesn't capture the full lifecycle of a task. Tasks can be successfully completed, duplicated, or explicitly cancelled—each should be tracked distinctly.

## Proposed Status Values
- `open` - ready to work on (default)
- `in_progress` - actively being worked on
- `done` - successfully completed
- `duplicate` - marked as duplicate of another task
- `cancelled` - explicitly cancelled/won't fix
- `blocked` - optionally, waiting on dependencies (implicit in blockers array)

## Benefits
- Better reporting and project visibility
- Eliminates the need for a separate `delete` command (use `duplicate` or `cancelled`)
- Clearer intent for stakeholders

## Migration
- Existing `completed: true` → `status: done`
- Existing `completed: false` → `status: open`
- `strand next` should respect status (only return `open` or `in_progress` tasks)
- `strand complete` should set status to `done` and prompt for status if not done
- New `strand cancel` command to set status to `cancelled`
- New `strand mark-duplicate` command to set status to `duplicate`

## Completion Report
Completed architectural planning for multi-state status field migration.

**Design Document**: design-docs/status-field-migration.md - Comprehensive implementation plan covering all 5 phases

**Epic Created**: Twqo3xp-status-field-migration-epic - Organizes all implementation work

**Implementation Tasks** (4 total):
1. Tzk35d7-update-metadata-struct-and-add-status-helpers - Data model updates
2. Ta3bynh-update-taskdb-operations-for-status-field - TaskDB operations
3. Tjqppdw-update-and-add-cli-commands-for-status-field - CLI command updates
4. Tg05aq8-migration-and-comprehensive-testing-for-status-fie - Migration and testing

**Review Tasks** (4 total):
- Tccehse-review-update-metadata-struct-for-status-field (reviewer-reliability)
- T06q0r4-review-taskdb-operations-for-status-field (reviewer-reliability)
- Tqa5qvs-review-cli-commands-for-status-field (reviewer-usability)
- Tuowxgx-review-migration-and-testing-for-status-field (tester)

**Status Values**: open, in_progress, done, duplicate, cancelled

**Key Implementation Details**:
- Backward compatibility with existing completed: bool field
- Automatic migration on task load
- Master list regeneration with new status values
- New CLI commands: cancel, mark-duplicate, mark-in-progress
- Comprehensive testing at unit, integration, and E2E levels
