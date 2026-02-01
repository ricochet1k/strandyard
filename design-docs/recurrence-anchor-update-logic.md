# Recurrence Anchor Update Logic

## Summary
Define how recurrence anchors are updated after a task is materialized to ensure deterministic behavior, avoid drift, and handle git-based metrics reliably.

## Principles
- **No-Drift**: For time-based recurrence, the next interval should be calculated from the *theoretical* due date, not the *actual* materialization time.
- **Auditability**: Every anchor resolution and update must be traceable.
- **Robustness**: Handle "overruns" (where a metric significantly exceeds the interval before the next evaluation).

## Update Logic by Metric Type

### 1. Time-based (days, weeks, months)
- **Anchor**: ISO 8601 timestamp.
- **Trigger**: `current_time >= anchor + interval`.
- **Update**: `new_anchor = old_anchor + interval`.
  - *Note*: If `current_time` is significantly past `old_anchor + interval` (e.g., multiple intervals missed), we may need to decide whether to "catch up" (create multiple tasks) or "skip ahead".
  - **Decision**: Skip ahead to the next theoretical interval after `current_time` to avoid flooding the task database, but log the skip in the activity log.

### 2. Git-based (commits)
- **Anchor**: Commit hash.
- **Trigger**: `git rev-list --count anchor..HEAD >= interval`.
- **Update**: `new_anchor = commit_at(old_anchor, interval)`.
  - The `new_anchor` is the specific commit hash that was exactly `interval` commits after the `old_anchor`.
  - Command: `git rev-list --reverse <old_anchor>..HEAD | sed -n '<interval>p'`.
  - This ensures that if 100 commits happened since the last anchor with an interval of 50, the next evaluation will correctly find another 50 commits ready.

### 3. Git-based (lines_changed)
- **Anchor**: Commit hash.
- **Trigger**: `git diff --numstat anchor..HEAD` (aggregated) `>= interval`.
- **Update**: `new_anchor = commit_crossing_threshold(old_anchor, interval)`.
  - The `new_anchor` is the *first* commit hash in the sequence where the cumulative lines changed since `old_anchor` met or exceeded the `interval`.
  - Implementation: Iterate through `git rev-list --reverse <old_anchor>..HEAD` and calculate cumulative delta.
  - This avoids drift and ensures all changes are eventually accounted for.

### 4. Tasks Completed
- **Anchor**: ISO 8601 timestamp of the last completion counted.
- **Trigger**: `activity_log.CountCompletionsSince(anchor) >= interval`.
- **Update**: `new_anchor = timestamp_of_completion_at(old_anchor, interval)`.
  - The `new_anchor` is the timestamp of the `interval`-th task completion since the `old_anchor`.

### Missing Git Anchors (Edge Case)
- **Scenario**: A commit hash anchor is missing from the repository (e.g., due to a force push).
- **Behavior**:
  - The evaluation should not fail.
  - Instead, it should resolve the current `HEAD` and update the task's `anchor` to this new hash.
  - A warning must be logged to the activity log: `recurrence_anchor_reset`.
  - The metric count for the current evaluation is treated as 0.

## Implementation Tasks (for Developers)
- [ ] Implement `GetCommitAtOffset(anchor, offset)` helper.
- [ ] Implement `GetCommitCrossingLinesThreshold(anchor, threshold)` helper.
- [ ] Update `MaterializeTask` logic to apply these update rules to the recurring task metadata.
- [ ] Ensure `activity.log` records these anchor updates.
