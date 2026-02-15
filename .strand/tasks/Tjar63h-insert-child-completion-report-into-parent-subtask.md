---
type: implement
role: developer
priority: medium
parent: Tfe1ssq-complete-command-should-also-insert-report-into-pa
blockers: []
blocks: []
date_created: 2026-02-15T04:52:53.295682Z
date_edited: 2026-02-15T05:12:39.1397Z
owner_approval: false
completed: true
status: done
description: ""
---

# Insert child completion report into parent subtask entry

## Summary
Confirmed bug: when `strand complete <child-id> "report"` completes a child task, the report is written to the child task body but not copied under the matching checked item in the parent task's `## Subtasks` section.

Repro steps:
1. Create parent + child task (child has `parent: <parent-id>`).
2. Run `strand repair` to ensure subtasks are materialized.
3. Run `strand complete <child-id> "child report line"`.
4. Observe parent `## Subtasks` shows only `- [x] (subtask: <short-id>) <title>` with no indented report lines.

Expected:
When a child is completed with a report, the parent subtask checkbox entry should include that report directly under the checked line so parent context contains completion details without opening each child task.

## Acceptance Criteria
- Implementation matches the specification
- Tests cover the change and pass
- Build succeeds

## Completion Report
Parent subtasks now copy child completion reports under checked subtask entries, including parsed completion-report sections from child body content. Added regression tests for parent subtask report propagation and complete command repaired-state behavior.
