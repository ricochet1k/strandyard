---
type: implement
role: developer
priority: medium
parent: Td6be8o-usability-review-recurrence-anchor-flags-alternati
blockers: []
blocks: []
date_created: 2026-02-05T00:50:42.300718Z
date_edited: 2026-02-05T00:50:42.300718Z
owner_approval: false
completed: false
description: ""
---

# Include unit list and examples in --every help text

## Summary
Update the help text for the `--every` flag in `cmd/add.go` and `cmd/edit.go` to include a full list of valid units and provide scannable examples for each category (time, git, activity).

Acceptance Criteria:
- Help text lists all valid units: days, weeks, months, commits, lines_changed, tasks_completed.
- Help text provides examples like: "10 days", "50 commits from HEAD", "20 tasks_completed from T1a1a".

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
