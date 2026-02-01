# Implementation Plan: Reuse log instance in EvaluateTasksCompletedMetric

## Architecture Overview
The current implementation of `EvaluateTasksCompletedMetric` in `pkg/task/recurrence.go` takes an `*activity.Log` argument but ignores it for the actual query, opening its own instance instead. This is inefficient and potentially dangerous for concurrency if the passed-in instance holds a lock on the log file.

## Specific files to modify
- `pkg/task/recurrence.go`: Update `EvaluateTasksCompletedMetric` to use the provided log if it is not nil.

## Code structure/patterns to use

```go
func EvaluateTasksCompletedMetric(baseDir, anchor string, taskID string, log *activity.Log) (int, error) {
    // ... anchor resolution logic ...

    activeLog := log
    if activeLog == nil {
        var err error
        activeLog, err = activity.Open(baseDir)
        if err != nil {
            return 0, fmt.Errorf("failed to open activity log: %w", err)
        }
        defer activeLog.Close()
    }

    return activeLog.CountCompletionsSince(anchorTime)
}
```

## Integration points
- This function is used by `UpdateAnchor` and potentially other recurrence-related logic.
- The `activity.Log` struct (in `pkg/activity/log.go`) manages the file handle and mutex.

## Testing approach
- Update `TestEvaluateTasksCompletedMetric` or add a new test in `pkg/task/recurrence_test.go`.
- The new test should:
  1. Open an `activity.Log`.
  2. Pass it to `EvaluateTasksCompletedMetric`.
  3. Verify the result.
  4. Verify the log is still open/usable after the call.

## Decision rationale
- Reusing the log instance is necessary to respect the mutex and avoid redundant file I/O.
- Providing a fallback (`activity.Open`) ensures the function remains usable even if a log instance is not readily available.
