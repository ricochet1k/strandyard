---
type: implement
role: developer
priority: high
parent: T6wwk1p-review-status-field-validation-for-usability
blockers: []
blocks: []
date_created: 2026-02-05T22:25:40.992166Z
date_edited: 2026-02-05T22:27:37.625731Z
owner_approval: false
completed: false
status: ""
description: ""
---

# Improve status field validation error messages

## Summary
The current error messages for status field validation use Go slice formatting (e.g., "[open in_progress done cancelled duplicate]"), which may confuse end users. Improve the error messages to clearly communicate which status values are allowed and provide helpful hints for common mistakes.

## Requirements
- Error messages should be human-readable and not expose Go internals
- Should list allowed values in a user-friendly format
- Should provide hints for common mistakes (e.g., if user tries "completed" instead of "done")
- Error messages should be consistent with other field validation messages

## Design doc
See design-docs/status-field-validation-error-messages.md for detailed specification.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: developer) Implement the behavior described in Context.
  Implemented user-friendly error message formatting and hints. Created FormatStatusListForUser() and FormatStatusErrorMessage() functions with hints for common mistakes like completed->done, pending->open/in_progress, failed/wontfix->cancelled, blocked->use blockers field.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Added comprehensive unit tests for FormatStatusListForUser() and FormatStatusErrorMessage() covering all hint scenarios. Added integration tests in repair_test.go for invalid status validation. Updated existing test expectations for new error message format. All 100+ tests pass.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
