---
type: implement
role: developer
priority: high
parent: Twqo3xp-status-field-migration-epic
blockers:
    - Tjqppdw-update-and-add-cli-commands-for-status-field
    - Tuowxgx-review-migration-and-testing-for-status-field
blocks: []
date_created: 2026-02-05T22:02:46.733426Z
date_edited: 2026-02-05T22:03:05.909698Z
owner_approval: false
completed: false
description: ""
---

# Migration and comprehensive testing for status field

## Summary
Implement backward compatibility migration and comprehensive testing for the status field changes.

**Implementation Plan**: See design-docs/status-field-migration.md (Phase 4-5)

**Specific changes**:
1. Task loading migration:
   - Handle `completed: true/false` to `status` conversion on load
   - Mark migrated tasks as dirty for rewrite with new status
   - Verify all existing tasks migrate correctly

2. Master list generation:
   - Update `tasks/free-tasks.md` generation to exclude non-active tasks
   - Verify `tasks/root-tasks.md` continues to work correctly

3. Activity logging:
   - Update or add event types for status transitions

4. Comprehensive testing:
   - Unit tests for status field parsing/serialization
   - Unit tests for migration logic
   - Unit tests for each status value
   - Integration tests for master list generation
   - E2E tests for all commands with new status
   - E2E tests for migration on old task files

5. Documentation updates if needed

**Acceptance Criteria**:
- All existing tasks migrate correctly
- Master lists generated correctly for new status values
- All tests pass
- No data loss during migration
- Old task files work with new code
- Activity logging tracks status changes
- E2E tests verify full workflow

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
