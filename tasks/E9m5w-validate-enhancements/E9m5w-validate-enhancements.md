---
role: architect
parent:
blockers: []
blocks:
  - E5w8m-e2e-tests
date_created: 2026-01-27
date_edited: 2026-01-27
completed: true
---

# Enhance Repair Command

## Summary

Improve the repair command to add task link validation, blocker status validation, and maintain simplicity with one right way of doing things.

## Context

**Owner Decision**: Default to simplicity and one right way of doing things where possible. Add validation for task links and blocker status.

**Current state**: Repairs role files exist, parent links valid, malformed IDs, regenerates master lists.

**Target state**: Also repairs task links are valid, blocker relationships are correct, and maintains simple, strict validation with no optional modes.

## Subtasks

1. [T3k8p-link-validation](T3k8p-link-validation/T3k8p-link-validation.md) - Add task link validation
2. [T7w4n-blocker-validation](T7w4n-blocker-validation/T7w4n-blocker-validation.md) - Add blocker status validation
3. [T2h9m-simplify-validation](T2h9m-simplify-validation/T2h9m-simplify-validation.md) - Ensure single right way validation

## Acceptance Criteria

- Repairs all task links point to existing tasks
- Repairs blocker relationships are bidirectional and consistent
- Repairs free-tasks.md only contains tasks with no blockers
- One validation mode (no --strict/--lenient options)
- Clear, actionable error messages
- Fails fast on any validation error

## References

- Current implementation: cmd/validate.go
