# Implementation Plan — Recurrence Anchor Audit Logging

## Summary
Implement audit logging for the resolution of default anchors ("now", "HEAD") in recurrence rules. This ensures that even when a rule uses a dynamic anchor, we have a record of what it resolved to at any point in time.

## Context
- `design-docs/recurrence-anchor-flags-alternatives.md`: Alternative D adopted optional anchors.
- `pkg/task/recurrence.go`: Logic for evaluating metrics.
- `pkg/activity/log.go`: Current activity logging implementation.

## Proposed Changes

### 1. Activity Log Enhancements
Update `pkg/activity/log.go` to support recurrence anchor resolution events:
- Add `EventRecurrenceAnchorResolved EventType = "recurrence_anchor_resolved"`.
- Update `Entry` struct to include a `Metadata map[string]string` field for storing the original and resolved values.
- Add `WriteRecurrenceAnchorResolution(taskID string, original, resolved string) error`.

### 2. Recurrence Evaluation Logging
Update `pkg/task/recurrence.go` (or the calling logic) to log resolutions:
- When `evaluateGitMetric` is called with "HEAD", log the resolved commit hash.
- When a time-based metric is evaluated with "now", log the resolved timestamp.
- These logs should only occur during evaluation/materialization, not during simple validation.

### 3. CLI Updates
Ensure `strand repair` and other commands that trigger recurrence evaluation have access to the activity log.

## Acceptance Criteria
- Running a command that evaluates a recurrence with a default anchor (e.g., `strand repair`) creates a `recurrence_anchor_resolved` entry in `.strand/activity.log`.
- The entry contains the original anchor string ("now" or "HEAD") and the resolved value.
- Existing tests pass.

## Architecture Decisions
- **Log Format**: Use the existing JSONL activity log.
- **Metadata**: Adding a generic `Metadata` map to `Entry` allows for future extensibility beyond just anchors.
- **Resolution Point**: Log at the time of evaluation to capture what was actually used for the calculation.

## Testing Strategy
- **Unit Tests**:
  - `pkg/activity`: Test writing and reading entries with metadata.
  - `pkg/task`: Mock git and time to verify `recurrence_anchor_resolved` events are triggered.
- **Integration Tests**:
  - Create a recurring task with a default anchor.
  - Run `strand repair`.
  - Verify `.strand/activity.log` contains the resolution entry.
