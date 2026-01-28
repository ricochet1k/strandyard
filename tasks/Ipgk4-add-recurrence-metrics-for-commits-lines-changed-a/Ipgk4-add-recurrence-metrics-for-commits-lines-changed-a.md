---
type: issue
role: triage
priority: medium
parent: ""
blockers: []
blocks: []
date_created: 2026-01-28T13:13:38.441725Z
date_edited: 2026-01-28T17:23:40.728217Z
owner_approval: false
completed: true
---

# Add recurrence metrics for commits, lines changed, and tasks completed

## Summary
Recurring tasks currently only support time-based recurrence. We need to add recurrence metrics based on commits, lines changed, and tasks completed so recurring tasks can trigger using those metrics instead of (or alongside) time passing.

## Steps to Reproduce
1. Configure a recurring task using the existing recurrence options.
2. Attempt to use commits, lines changed, or tasks completed as a recurrence trigger.
3. Observe that only time-based recurrence options are available.

## Expected Result
Recurrence configuration allows selecting commits, lines changed, or tasks completed as recurrence metrics, either alone or combined with time-based recurrence.

## Actual Result
Recurrence configuration only supports time passing; commit, line-change, and task-completed metrics are unavailable.

## Acceptance Criteria
- Recurrence configuration supports commit-count, lines-changed, and tasks-completed metrics.
- Metrics can be used instead of or in addition to time-based recurrence.
- User can verify the metric selection via the CLI/UI where recurrence is configured.

## Triage Notes
- `go run . recurring` returns "unknown command \"recurring\"", so no recurrence configuration is available yet.
- `tasks/E7p4m-issues-recurrence/T8v3k-recurring-tasks/T8v3k-recurring-tasks.md` only specifies time-based units plus `commits`; it does not cover line changes or tasks completed.

## Next Steps
- `T968i-design-recurrence-metrics-for-commits-lines-change` to define the schema and metric sourcing for commits, line changes, and completed tasks.
