---
type: task
role: architect
priority: high
parent: T1dua-draft-recurrence-anchor-flag-alternatives
blockers: []
blocks: []
date_created: 2026-02-01T20:24:47.655537Z
date_edited: 2026-02-01T20:27:02.354895Z
owner_approval: false
completed: true
description: ""
---

# New Task: Add graceful error handling for invalid anchors during materialization

## Description
Design and plan how to add graceful error handling for invalid anchors during task materialization. This should include:
- What constitutes an invalid anchor
- How to detect invalid anchors
- How to gracefully handle them (error messages, fallback behavior, etc.)
- Implementation approach and file locations
- Testing strategy

As an architect, create a detailed implementation plan in design-docs/ that answers:
- Architecture overview of error handling
- Specific files to modify
- Code structure/patterns to use
- Integration points with existing materialization logic
- Testing approach
- Decision rationale and any remaining trade-offs

Do NOT write production code - that's the developer's job.

## Completion Report
Created detailed implementation plan in design-docs/graceful-error-handling-for-invalid-anchors.md. The design adopts fail-fast error handling (Option A) with clear error messages following the unified error format from recurrence-anchor-error-messages.md. Implementation involves modifying pkg/task/recurrence.go to return errors instead of silent no-op behavior, creating pkg/task/errors.go for unified error formatting, and creating pkg/task/logger.go for audit trail logging. Specific error messages are defined for invalid HEAD, unknown revisions, and git command errors, all with recovery hints. Testing strategy includes unit tests for all error scenarios and integration tests for materialization failures.
