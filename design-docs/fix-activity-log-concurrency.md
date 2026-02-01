# Implementation Plan: Fix Activity Log Concurrency Risk

## Problem
`activity.Log.ReadEntries` currently closes the write-only file handle and reopens the file for reading, then reopens for writing. This is not thread-safe and can lead to data loss or errors if multiple goroutines interact with the log.

## Proposed Changes

### 1. Refactor `activity.Log` struct
Add a `sync.RWMutex` to protect the log's state and ensure thread-safe access to the file.

```go
type Log struct {
    mu       sync.RWMutex
    filepath string
    file     *os.File
}
```

### 2. Update `WriteEntry`
Use `l.mu.Lock()` to ensure exclusive access during writes.

### 3. Refactor `ReadEntries`
Instead of closing and reopening the main file handle, `ReadEntries` should:
1. Use `l.mu.RLock()` to allow concurrent reads but block writes.
2. Open a separate, read-only file handle for the duration of the read operation.
3. Use a `defer` to close the temporary read handle.

This avoids the risky close/reopen cycle of the primary write handle.

### 4. Integration Points
- `CountCompletionsSince`, `CountCompletionsForTaskSince`, and `GetCompletionTimestampAtOffset` all call `ReadEntries`, so they will benefit from the improved safety.

## Testing Strategy
- Add a new test case in `pkg/activity/log_test.go` that performs concurrent reads and writes to verify stability under load.
- Ensure existing tests pass.

## Rationale
Opening a separate file handle for reading is the simplest way to support concurrent reading without disrupting the append-only write handle. The `RWMutex` ensures that we don't try to read while a write is in progress (which might be middle-of-line) and that multiple readers can proceed together.
