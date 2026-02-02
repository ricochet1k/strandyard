---
type: task
role: architect
priority: medium
parent: Tpr12oz-reliability-review-of-activity-log-concurrency-pla
blockers: []
blocks: []
date_created: 2026-02-02T00:02:21.874142Z
date_edited: 2026-02-02T00:06:27.071876Z
owner_approval: false
completed: true
description: ""
---

# New Task: Address reliability concerns from activity log concurrency review

## Description


## Summary
The reliability review of the activity log concurrency plan identified several concerns that should be addressed in the implementation or design document.

## Acceptance Criteria
- Implementation handles file opening errors gracefully.
- Write failures do not corrupt the `Log` state.
- Design doc or code comments address the file handle limit trade-offs.

Decide which task template would best fit this task and re-add it with that template and the same parent.

## TODOs
- [x] (role: architect) Investigate potential for file handle exhaustion if `ReadEntries` is called frequently under high concurrency.
  I investigated the potential for file handle exhaustion during concurrent calls to ReadEntries.
  Findings:
  * The current implementation is not thread-safe and causes race conditions (crashes) when multiple goroutines attempt to close/re-open the log file.
  * Opening a new file handle for every read (as previously planned) could hit the OS file handle limit (ulimit -n) under high concurrency.
  * I updated design-docs/fix-activity-log-concurrency.md to include a caching strategy that uses an RWMutex and only re-reads the file when it has changed on disk.
  * I added a section to the design doc addressing handle exhaustion trade-offs and graceful error handling.
- [x] (role: architect) Ensure robust error handling when opening the separate read-only file handle (e.g., handling "too many open files" or permission issues).
  Updated design doc with robust error handling strategies for file opening failures.
- [x] (role: architect) Verify that `WriteEntry` handles write failures (like disk full) gracefully and doesn't leave the file handle in a bad state.
  Addressed write failure resilience in design doc, relying on single write calls and existing malformed entry detection.
- [x] (role: architect) Consider if a persistent read handle (with its own mutex) would be more efficient/reliable than opening a new one every time, if concurrency is high.
  Proposed caching strategy in design doc as a more efficient alternative to persistent read handles.
