---
type: implement
role: architect
priority: high
parent: Tq49ww0-fix-strand-complete-to-update-free-list-when-last
blockers: []
blocks: []
date_created: 2026-02-05T21:45:59.620259Z
date_edited: 2026-02-05T21:57:32.643414Z
owner_approval: false
completed: false
description: ""
---

# Fix atomicity issue in free-list update

## Summary


## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
- [x] (role: developer) Implement the behavior described in Context.
  Fixed atomicity issue by moving CalculateIncrementalFreeListUpdate to after all task state changes. Now calculates based on final state instead of stale state.
- [x] (role: developer) Add unit and integration tests covering the main flows if they don't already exist.
  Added unit test (TestAtomicityOfFreeListCalculation) that verifies the atomicity fix for free-list updates. The test ensures CalculateIncrementalFreeListUpdate correctly identifies newly-unblocked tasks by calculating based on final state rather than stale state. Tests cover single blocker scenarios and multiple blocker edge cases.
- [ ] (role: tester) Execute test-suite and report failures.
- [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
- [ ] (role: documentation) Update user-facing docs and examples.
