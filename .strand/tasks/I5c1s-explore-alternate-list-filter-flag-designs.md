---
type: issue
role: triage
priority: high
parent: ""
blockers: []
blocks: []
date_created: 2026-01-28T01:45:19Z
date_edited: 2026-01-27T18:53:24.840439-07:00
owner_approval: false
completed: true
---

# Explore alternate list filter flag designs

## Summary
The current list command exposes multiple boolean filter flags (completed/blocked/blocks/owner-approval) and relies on flag values like `--completed=false` for negation. The design doc specifies boolean flags but doesn’t evaluate alternative flag patterns. We need to explore a clearer filter syntax before finalizing CLI docs and defaults.

## Steps to Reproduce
1. Review `design-docs/list-command.md` (Filters section) to see the current boolean flag contract.
2. Review `cmd/list.go` to confirm the current list flags and their help text.
3. Compare with alternate patterns (single `--filter` flag, `--status` enum, `--blocked/--unblocked`, or `--completed/--incomplete`).

## Expected Result
A documented, consistent filter flag pattern with explicit true/false and negation semantics, with pros/cons recorded.

## Actual Result
Multiple boolean flags are defined, and negation relies on `--flag=false` (not obvious in help), with no consolidated design decision on alternative flag patterns.

## Acceptance Criteria
- Clear steps to reproduce are documented.
- Expected vs actual behavior is explained.
- Follow-on design task opened to evaluate alternative filter flag designs.
