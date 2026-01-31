---
type: task
role: developer
priority: medium
parent: Ti6zj-taskdb-api-design-review
blockers: []
blocks:
    - Ti6zj-taskdb-api-design-review
date_created: 2026-01-31T17:19:06.527103Z
date_edited: 2026-01-31T17:19:06.551211Z
owner_approval: false
completed: false
---

# Design relationship management API

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
Design proper blocker/parent relationship APIs:
- Clarify responsibilities: what computes relationships vs. what sets them
- Proper naming for all relationship functions
- Resolve UpdateBlockersFromChildren / FixBlockerRelationships / SyncBlockersFromChildren confusion
- Define which should be public, private, or removed
- Ensure all operations maintain bidirectional integrity
- Document the mental model users should have

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [ ] (role: developer) Implement the behavior described in Context.
- [ ] (role: developer) Add unit and integration tests covering the main flows.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
