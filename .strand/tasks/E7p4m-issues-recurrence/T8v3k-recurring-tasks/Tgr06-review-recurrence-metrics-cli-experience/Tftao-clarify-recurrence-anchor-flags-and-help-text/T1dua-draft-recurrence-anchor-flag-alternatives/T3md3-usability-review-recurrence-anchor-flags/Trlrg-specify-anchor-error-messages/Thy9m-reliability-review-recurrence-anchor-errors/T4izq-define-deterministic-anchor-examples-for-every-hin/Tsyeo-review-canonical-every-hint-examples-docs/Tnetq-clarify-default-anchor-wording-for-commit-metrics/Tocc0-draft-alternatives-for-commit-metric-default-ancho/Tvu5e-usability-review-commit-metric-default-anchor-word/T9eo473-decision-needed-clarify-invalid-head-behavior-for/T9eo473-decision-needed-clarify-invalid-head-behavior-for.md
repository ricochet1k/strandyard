---
type: task
role: owner
priority: high
parent: Tvu5e-usability-review-commit-metric-default-anchor-word
blockers: []
blocks: []
date_created: 2026-02-01T05:17:09.114657Z
date_edited: 2026-02-01T05:17:09.114657Z
owner_approval: false
completed: false
---

# Decision needed: Clarify invalid HEAD behavior for commit metrics

## Context
Provide links to relevant design documents, diagrams, and decision records.

## Description
The design document design-docs/recurrence-anchor-commit-metrics-default-anchor-wording-alternatives.md states in its HEAD availability section an implication that the CLI should error and instruct the user when HEAD is invalid or unborn. However, the Decision section states that commit-based recurrence metrics will be ignored and not trigger recurring tasks if HEAD is invalid or unborn.

This silent ignoring of recurrence metrics in cases of invalid HEAD can lead to a poor user experience, as users might expect tasks to recur but receive no feedback on why they aren't. This can be confusing and potentially lead to missed work or incorrect assumptions about the system's behavior.

**Recommendation:** Reconsider the behavior for invalid/unborn HEAD. Instead of silently ignoring, propose either:
1.  **Warning Message:** Log a clear warning message to the user, explaining that recurrence is being ignored due to an invalid HEAD and providing actionable guidance on how to resolve it (e.g., "recurrence for  metric ignored:  is unborn. To enable, make an initial commit or explicitly set an anchor.").
2.  **Explicit Error (contextual):** In situations where the user is actively querying or expecting recurrence, provide an explicit error with guidance. For background recurrence checks, a warning might be more appropriate.

Please clarify the desired behavior and update the design document and implementation accordingly.

## Escalation
Tasks are disposable. Use follow-up tasks for open questions/concerns. Record decisions and final rationale in design docs; do not edit this task to capture outcomes.

## Acceptance Criteria
- Clear, runnable steps to reproduce locally.
- Tests covering functionality and passing.
- Required reviews completed and blockers cleared.

## TODOs
Check this off one at a time with `strand complete <task_id> --role <my_given_role> --todo <num>` only if your Role matches the todo's role.
1. [ ] (role: developer) Implement the behavior described in Context.
2. [ ] (role: developer) Add unit and integration tests covering the main flows.
3. [ ] (role: tester) Execute test-suite and report failures.
4. [ ] (role: master-reviewer) Coordinate required reviews: `reviewer-reliability`, `reviewer-security`, `reviewer-usability`.
5. [ ] (role: documentation) Update user-facing docs and examples.

## Subtasks
Use subtasks for work that should be tracked separately or assigned to a different role. Use `strand add <type> "title" --parent <this_task_id> <<EOF description EOF`  to create subtasks.
