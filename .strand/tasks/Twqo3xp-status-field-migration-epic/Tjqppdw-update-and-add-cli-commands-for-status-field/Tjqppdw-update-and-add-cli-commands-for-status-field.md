---
type: implement
role: developer
priority: high
parent: Twqo3xp-status-field-migration-epic
blockers:
    - Ta3bynh-update-taskdb-operations-for-status-field
    - Tqa5qvs-review-cli-commands-for-status-field
blocks:
    - Tg05aq8-migration-and-comprehensive-testing-for-status-fie
date_created: 2026-02-05T22:02:41.772886Z
date_edited: 2026-02-05T22:03:02.57131Z
owner_approval: false
completed: false
description: ""
---

# Update and add CLI commands for status field

## Summary
Update and create CLI commands to support the new status field.

**Implementation Plan**: See design-docs/status-field-migration.md (Phase 3)

**Specific changes**:
1. Update `cmd/complete.go`:
   - Set `status: done` instead of `completed: true`
   - Update command help text
   - Prevent completing tasks with status duplicate/cancelled

2. Update `cmd/next.go`:
   - Filter to show only open/in_progress tasks
   - Update free-list filtering logic

3. Update `cmd/list.go`:
   - Update `--completed` flag to check `status == done`
   - Consider adding `--status` flag for granular filtering

4. Create new commands:
   - `strand cancel <task-id> "reason"` - sets `status: cancelled`
   - `strand mark-duplicate <task-id> <duplicate-of>` - sets `status: duplicate`
   - `strand mark-in-progress <task-id>` - sets `status: in_progress`

5. Update command help/documentation

**Acceptance Criteria**:
- All commands work with new status field
- strand complete sets correct status
- strand next filters correctly
- strand cancel/mark-duplicate/mark-in-progress work
- No regression in existing commands
- All new commands have help text
- Integration tests for each command

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [ ] (role: developer) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
