---
type: task
role: developer
priority: medium
parent: Ti6zj-taskdb-api-design-review
blockers: []
blocks: []
date_created: 2026-01-31T17:18:48.599142Z
date_edited: 2026-02-15T05:12:33.908852Z
owner_approval: false
completed: true
status: done
description: ""
---

# Review blockers.go relationship management

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
Analyze pkg/task/blockers.go:
- Document UpdateBlockersFromChildren behavior and purpose
- Identify what it does well (original design)
- Note any issues with current implementation
- Understand the relationship computation logic
- Determine if this belongs in TaskDB or is a separate concern

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## Completion Report
Reviewed blocker relationship logic in pkg/task/blockers.go and documented behavior, strengths, and risks in design-docs/blockers-relationship-management-review.md. Captured boundary alternatives (TaskDB vs separate engine) with decision deferred to Owner, ran build/test/repair, and filed follow-up issue T09easy for reconciliation invariant gaps.

## Subtasks
- [x] (subtask: T09easy) Unify completion cleanup with blocker reconciliation invariants
  Reworked UpdateBlockersAfterCompletion to call ReconcileBlockerRelationships so completion cleanup matches canonical invariants, and extended tests to verify blocks/blockers are cleared after completion even with inconsistent edges. Ran go test ./..., go build ./..., and strand repair (0 tasks repaired).
