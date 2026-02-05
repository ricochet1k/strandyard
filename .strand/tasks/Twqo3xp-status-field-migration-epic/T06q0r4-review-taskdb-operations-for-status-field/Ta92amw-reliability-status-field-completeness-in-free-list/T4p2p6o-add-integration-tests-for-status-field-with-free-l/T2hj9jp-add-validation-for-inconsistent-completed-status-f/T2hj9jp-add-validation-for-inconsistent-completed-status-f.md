---
type: implement
role: architect
priority: high
parent: T4p2p6o-add-integration-tests-for-status-field-with-free-l
blockers: []
blocks: []
date_created: 2026-02-05T22:15:27.452666Z
date_edited: 2026-02-05T22:15:27.452666Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Add validation for inconsistent Completed/Status field states

## Summary
During the Status Field Migration, tasks transition from using a simple `completed: bool` field to a multi-state `status` field. This task ensures that validation catches and prevents inconsistent states where these two fields conflict with each other.

The validation should catch cases like:
- `Completed: true` with `Status: open` or `Status: in_progress` (should be `done`)
- `Completed: false` with `Status: done` (should be `open`, `in_progress`, `cancelled`, or `duplicate`)

## Context
The task database supports backward compatibility by maintaining both the old `completed` boolean field and the new `status` string field. During migration, it's critical that these two fields remain consistent. The `Validator` in `pkg/task/repair.go` should detect and report any inconsistencies so they can be fixed.

## Implementation Details
The validation logic already exists in `pkg/task/repair.go`:
- Function: `verifyCompletedStatusConsistency()` (lines 249-274)
- Test file: `pkg/task/repair_test.go` (lines 381-502)

**Validation Rules**:
1. If `Completed: true`, then `Status` must be `"done"` (or empty, which defaults to `open` in new tasks)
2. If `Completed: false`, then `Status` must NOT be `"done"`
3. Valid non-done statuses: `"open"`, `"in_progress"`, `"cancelled"`, `"duplicate"`

**Error Messages**:
- `"inconsistent state: Completed=true but Status="X" (should be 'done')"`
- `"inconsistent state: Completed=false but Status=done (should be 'open', 'in_progress', 'cancelled', or 'duplicate')"`

## Acceptance Criteria
- Validation catches all inconsistent Completed/Status combinations
- Tests verify 9 test cases (3 consistent, 3 inconsistent, 3 edge cases)
- All validation tests pass
- No false positives on valid state combinations
- Error messages are clear and actionable

## TODOs
- [x] (role: developer) Implement the behavior described in Context.
  Validation logic already implemented in `verifyCompletedStatusConsistency()` in pkg/task/repair.go (lines 249-274). Tests already exist in pkg/task/repair_test.go (lines 381-502) covering 9 test cases including consistent and inconsistent states.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Tests already exist in TestCompletedStatusConsistency in repair_test.go covering: Completed=true/Status=done (consistent), Completed=false/Status=open and in_progress (consistent), Completed=true with Status=open and in_progress (inconsistent), and Completed=false with Status=done (inconsistent).
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
