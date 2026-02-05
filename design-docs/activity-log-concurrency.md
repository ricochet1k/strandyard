# Activity Log Concurrency and Consistency Model

## Overview
The Activity Log (`activity.log`) is a newline-delimited JSON file stored in the project's `.strand/` directory. It tracks events such as task completions and recurrence anchor resolutions. This document defines the concurrency guarantees and consistency model for readers and writers of this log.

## Concurrency Guarantees

### Writing
- **Atomic Appends**: The log is opened with `os.O_APPEND`. On most modern operating systems (including Linux and macOS), writes to a file opened with `O_APPEND` are atomic if the data size is below the file system's block size (typically 4KB). Since our log entries are single JSON lines, they are usually well within this limit.
- **Process Safety**: Within a single Go process, the `Log` struct uses a `sync.RWMutex` to coordinate access between goroutines.
- **Durability**: Each write is followed by an explicit `file.Sync()` to ensure the data is flushed to the physical storage.

### Reading
- **Non-blocking Reads**: Readers use standard `os.Open` and do not block writers.
- **Incremental Reads**: The `Log` struct tracks `lastSize` and uses `os.Stat()` to detect new data. It only reads the delta since the last read, minimizing I/O.
- **Reverse Scanning**: For performance, some queries (like `GetLatestTaskCompletionTime`) scan the file from the end. This is done without locking the file, relying on the atomic nature of the appends.

## Consistency Model

### Within a Process
- **Read-Your-Own-Writes**: A process using the same `Log` instance will see its own writes immediately because `WriteEntry` invalidates the read cache (`lastSize = -1`).
- **Monotonic Reads**: Once an entry is read into the `Log.entries` cache, it is never removed (unless the file shrinks, which should not happen in normal operation).

### Across Processes
- **Eventual Consistency**: Separate CLI invocations will see each other's writes as soon as the OS flushes the data and `os.Stat()` reflects the new file size. Since we call `Sync()` after every write, this delay is minimal.
- **No Global Locking**: There is no global file lock (e.g., `flock`). This is intentional to keep the CLI fast and avoid deadlocks or stale lock files.
- **Resilience**: The parser skips malformed entries (e.g., partially written lines) and continues, ensuring that a crash during writing does not permanently break the log for readers.

## Trade-offs
- **ID Reuse**: We assume task IDs are unique enough that searching for the "latest" completion of a task ID in the log is sufficient. If a task ID is reused after being deleted, the log will contain history for both.
- **No Deletion**: The activity log is append-only and currently has no compaction or rotation mechanism. This ensures a complete audit trail but may grow over time.
