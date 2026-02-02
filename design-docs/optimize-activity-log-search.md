# Implementation Plan: Optimize Activity Log Search

## Problem
Currently, searching for the latest completion time of a task requires reading the entire activity log file into memory (via `ReadEntries`). As the log grows, this will become a performance bottleneck and memory pressure point.

## Proposed Changes

### 1. New Method: `GetLatestTaskCompletionTime`
Add a method to the `Log` struct that searches for the most recent `task_completed` event for a given `taskID` without necessarily reading the entire file.

```go
func (l *Log) GetLatestTaskCompletionTime(taskID string) (time.Time, error)
```

### 2. Backward Reading Strategy
Instead of forward scanning, the method will:
1. Use `l.mu.RLock()` to ensure consistency.
2. Open the log file for reading.
3. Seek to the end of the file.
4. Read blocks (e.g., 64KB) from the end moving backwards.
5. In each block, identify line boundaries (newline characters).
6. Parse each line (from most recent to oldest) as an `Entry`.
7. If an entry matches the `taskID` and has type `EventTaskCompleted`, return its timestamp.
8. If the start of the file is reached without a match, return a "not found" error.

### 3. Integration with Cache
- `GetLatestTaskCompletionTime` should *not* update or use the in-memory `entries` cache to avoid loading the whole file.
- However, if the cache is already populated and up-to-date (checked via `lastSize`), it might be faster to search the cache backwards first.
- Decision: For simplicity and to strictly meet the optimization goal, `GetLatestTaskCompletionTime` will always perform a backward file scan if the log is large, or just use the cache if it's already there.

### 4. Implementation Details: `pkg/activity/reverse_scanner.go`
Create a helper utility for backward line-by-line scanning to keep `log.go` clean.

## Testing Strategy
- Add unit tests in `pkg/activity/log_test.go`:
    - Search for a task completed recently (last line).
    - Search for a task completed long ago (middle of file).
    - Search for a task that was never completed.
    - Test with a very large log file (simulated).
- Verify that it works correctly with malformed lines (it should skip them).

## Rationale
Backward scanning is the most efficient way to find the "latest" event in an append-only log. It minimizes I/O and memory usage, especially when the desired event is recent.
