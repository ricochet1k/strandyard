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

### 3. Refactor `ReadEntries` with Caching
Instead of closing and reopening the main file handle, `ReadEntries` should:
1. Use `l.mu.RLock()` to allow concurrent reads but block writes.
2. Check if the file on disk has changed (by comparing current size with `lastSize` or using `lastModified`).
3. If unchanged, return the cached `entries`.
4. If changed:
    - Open a separate, read-only file handle.
    - Read ONLY the new entries (starting from `lastSize`) if possible, or re-read the whole file if necessary (simpler to start).
    - Update the cache and `lastSize`.
5. Use a `defer` to close the temporary read handle.

This avoids the risky close/reopen cycle and significantly reduces file I/O and handle count under high concurrency.

### 4. Address File Handle Exhaustion
Under extremely high concurrency, even temporary read handles could hit the OS limit (`ulimit -n`). 
- **Graceful Error Handling**: Ensure that `os.Open` failures (especially `EMFILE`) are handled and returned as clear errors to the caller.
- **Caching**: The caching strategy above minimizes the window where handles are open and avoids opening them at all for many calls.
- **Optional Semaphore**: If exhaustion remains a concern, a internal semaphore could be used to limit concurrent read operations to a safe number (e.g., 64).

### 6. Robust Error Handling and Write Resilience
- **File Opening Errors**: All `os.Open` and `os.OpenFile` calls must check for errors. If a read-only handle cannot be opened (e.g., `EMFILE`), `ReadEntries` should return a wrapped error allowing the caller to retry or fail gracefully.
- **Write Failure Resilience**: 
    - `WriteEntry` should continue to use a single `Write` call for the marshaled JSON plus newline to maximize the chance of atomic line writes.
    - If a write fails (e.g., disk full), the error is returned immediately. Since we use `O_APPEND`, the file pointer remains at the end for the next successful write.
    - We rely on the existing malformed entry detection in `ReadEntries` to skip or report any partial writes that might have occurred during a crash or disk-full event.
    - `Sync()` is called after every write to ensure durability, especially important for a log.

### 7. Integration Points
- `CountCompletionsSince`, `CountCompletionsForTaskSince`, and `GetCompletionTimestampAtOffset` all call `ReadEntries`, so they will benefit from the improved safety.

## Testing Strategy
- Add a new test case in `pkg/activity/log_test.go` that performs concurrent reads and writes to verify stability under load.
- Ensure existing tests pass.

## Rationale
Opening a separate file handle for reading is the simplest way to support concurrent reading without disrupting the append-only write handle. The `RWMutex` ensures that we don't try to read while a write is in progress (which might be middle-of-line) and that multiple readers can proceed together.
