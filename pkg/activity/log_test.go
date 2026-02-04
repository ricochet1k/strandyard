package activity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ricochet1k/strandyard/pkg/task"
)

func TestOpen(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}
	defer log.Close()

	if log.filepath == "" {
		t.Error("expected filepath to be set")
	}

	if log.file == nil {
		t.Error("expected file to be open")
	}

	expectedPath := filepath.Join(tmpDir, defaultLogFilename)
	if log.filepath != expectedPath {
		t.Errorf("filepath mismatch: got %s, want %s", log.filepath, expectedPath)
	}
}

func TestWriteEntry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}
	defer log.Close()

	entry := Entry{
		Timestamp: time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC),
		TaskID:    "T3k7x-example",
		Type:      EventTaskCompleted,
		Report:    "Completed the task",
	}

	if err := log.WriteEntry(entry); err != nil {
		t.Fatalf("failed to write entry: %v", err)
	}

	if err := log.Close(); err != nil {
		t.Fatalf("failed to close log: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, defaultLogFilename))
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	var readEntry Entry
	if err := json.Unmarshal(data, &readEntry); err != nil {
		t.Fatalf("failed to unmarshal entry: %v", err)
	}

	if !cmp.Equal(entry, readEntry) {
		t.Errorf("entry mismatch (-want +got):\n%s", cmp.Diff(entry, readEntry))
	}
}

func TestWriteEntryDefaultTimestamp(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}
	defer log.Close()

	beforeWrite := time.Now().UTC()

	entry := Entry{
		TaskID: "T3k7x-example",
		Type:   EventTaskCompleted,
	}

	if err := log.WriteEntry(entry); err != nil {
		t.Fatalf("failed to write entry: %v", err)
	}

	if err := log.Close(); err != nil {
		t.Fatalf("failed to close log: %v", err)
	}

	afterWrite := time.Now().UTC()

	data, err := os.ReadFile(filepath.Join(tmpDir, defaultLogFilename))
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	var readEntry Entry
	if err := json.Unmarshal(data, &readEntry); err != nil {
		t.Fatalf("failed to unmarshal entry: %v", err)
	}

	if readEntry.Timestamp.IsZero() {
		t.Error("expected timestamp to be set")
	}

	if readEntry.Timestamp.Before(beforeWrite) {
		t.Error("timestamp is before write time")
	}

	if readEntry.Timestamp.After(afterWrite.Add(1 * time.Second)) {
		t.Error("timestamp is significantly after write time")
	}
}

func TestWriteTaskCompletion(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}
	defer log.Close()

	taskID := "T3k7x-example"
	report := "Task completed successfully"

	if err := log.WriteTaskCompletion(taskID, report); err != nil {
		t.Fatalf("failed to write task completion: %v", err)
	}

	if err := log.Close(); err != nil {
		t.Fatalf("failed to close log: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, defaultLogFilename))
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	var readEntry Entry
	if err := json.Unmarshal(data, &readEntry); err != nil {
		t.Fatalf("failed to unmarshal entry: %v", err)
	}

	if readEntry.TaskID != taskID {
		t.Errorf("task_id mismatch: got %s, want %s", readEntry.TaskID, taskID)
	}

	if readEntry.Type != EventTaskCompleted {
		t.Errorf("type mismatch: got %s, want %s", readEntry.Type, EventTaskCompleted)
	}

	if readEntry.Report != report {
		t.Errorf("report mismatch: got %s, want %s", readEntry.Report, report)
	}
}

func TestWriteRecurrenceAnchorResolution(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}
	defer log.Close()

	taskID := "T3k7x-example"
	original := "HEAD"
	resolved := "abc1234567890"

	if err := log.WriteRecurrenceAnchorResolution(taskID, original, resolved); err != nil {
		t.Fatalf("failed to write recurrence anchor resolution: %v", err)
	}

	if err := log.Close(); err != nil {
		t.Fatalf("failed to close log: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, defaultLogFilename))
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	var readEntry Entry
	if err := json.Unmarshal(data, &readEntry); err != nil {
		t.Fatalf("failed to unmarshal entry: %v", err)
	}

	if readEntry.TaskID != taskID {
		t.Errorf("task_id mismatch: got %s, want %s", readEntry.TaskID, taskID)
	}

	if readEntry.Type != EventRecurrenceAnchorResolved {
		t.Errorf("type mismatch: got %s, want %s", readEntry.Type, EventRecurrenceAnchorResolved)
	}

	if readEntry.Metadata["original"] != original {
		t.Errorf("original anchor mismatch: got %s, want %s", readEntry.Metadata["original"], original)
	}

	if readEntry.Metadata["resolved"] != resolved {
		t.Errorf("resolved anchor mismatch: got %s, want %s", readEntry.Metadata["resolved"], resolved)
	}
}

func TestMultipleEntries(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}

	entries := []Entry{
		{TaskID: "T3k7x-first", Type: EventTaskCompleted, Report: "First task"},
		{TaskID: "E2k7x-second", Type: EventTaskCompleted, Report: "Second task"},
		{TaskID: "T8h4w-third", Type: EventTaskCompleted},
	}

	for _, entry := range entries {
		if err := log.WriteEntry(entry); err != nil {
			t.Fatalf("failed to write entry: %v", err)
		}
	}

	if err := log.Close(); err != nil {
		t.Fatalf("failed to close log: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, defaultLogFilename))
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	lines := parseJSONLines(data)
	if len(lines) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(lines))
	}

	for i, line := range lines {
		var readEntry Entry
		if err := json.Unmarshal(line, &readEntry); err != nil {
			t.Fatalf("failed to unmarshal entry %d: %v", i, err)
		}

		expected := entries[i]
		if readEntry.TaskID != expected.TaskID {
			t.Errorf("entry %d task_id mismatch: got %s, want %s", i, readEntry.TaskID, expected.TaskID)
		}
	}
}

func parseJSONLines(data []byte) [][]byte {
	lines := [][]byte{}
	current := []byte{}
	for _, b := range data {
		if b == '\n' {
			if len(current) > 0 {
				lines = append(lines, current)
				current = []byte{}
			}
		} else {
			current = append(current, b)
		}
	}
	if len(current) > 0 {
		lines = append(lines, current)
	}
	return lines
}

func TestReadEntries(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}

	entries := []Entry{
		{TaskID: "T3k7x-first", Type: EventTaskCompleted, Report: "First task"},
		{TaskID: "E2k7x-second", Type: EventTaskCompleted, Report: "Second task"},
		{TaskID: "T8h4w-third", Type: EventTaskCompleted},
	}

	for _, entry := range entries {
		if err := log.WriteEntry(entry); err != nil {
			t.Fatalf("failed to write entry: %v", err)
		}
	}

	readEntries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("failed to read entries: %v", err)
	}

	if len(readEntries) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(readEntries))
	}

	for i, entry := range readEntries {
		if entry.TaskID != entries[i].TaskID {
			t.Errorf("entry %d task_id mismatch: got %s, want %s", i, entry.TaskID, entries[i].TaskID)
		}
		if entry.Type != entries[i].Type {
			t.Errorf("entry %d type mismatch: got %s, want %s", i, entry.Type, entries[i].Type)
		}
	}
}

func TestReadEntriesEmptyLog(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}

	entries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("failed to read entries: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("expected 0 entries from empty log, got %d", len(entries))
	}
}

func TestCountCompletionsSince(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}

	now := time.Now().UTC()
	yesterday := now.Add(-24 * time.Hour)

	entries := []Entry{
		{Timestamp: now.Add(-2 * time.Hour), TaskID: "T3k7x-recent", Type: EventTaskCompleted},
		{Timestamp: yesterday.Add(2 * time.Hour), TaskID: "E2k7x-yesterday", Type: EventTaskCompleted},
		{Timestamp: yesterday.Add(3 * time.Hour), TaskID: "T8h4w-older", Type: EventTaskCompleted},
	}

	for _, entry := range entries {
		if err := log.WriteEntry(entry); err != nil {
			t.Fatalf("failed to write entry: %v", err)
		}
	}

	since := now.Add(-3 * time.Hour)
	count, err := log.CountCompletionsSince(since)
	if err != nil {
		t.Fatalf("failed to count completions: %v", err)
	}

	expected := 1
	if count != expected {
		t.Errorf("expected %d completions since %v, got %d", expected, since, count)
	}
}

func TestCountCompletionsForTaskSince(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}

	now := time.Now().UTC()
	yesterday := now.Add(-24 * time.Hour)

	entries := []Entry{
		{Timestamp: now.Add(-2 * time.Hour), TaskID: "T3k7x-task1", Type: EventTaskCompleted},
		{Timestamp: now.Add(-1 * time.Hour), TaskID: "T3k7x-task1", Type: EventTaskCompleted},
		{Timestamp: now.Add(-2 * time.Hour), TaskID: "E2k7x-task2", Type: EventTaskCompleted},
		{Timestamp: yesterday.Add(2 * time.Hour), TaskID: "T3k7x-task1", Type: EventTaskCompleted},
	}

	for _, entry := range entries {
		if err := log.WriteEntry(entry); err != nil {
			t.Fatalf("failed to write entry: %v", err)
		}
	}

	since := now.Add(-3 * time.Hour)
	count, err := log.CountCompletionsForTaskSince("T3k7x-task1", since)
	if err != nil {
		t.Fatalf("failed to count completions: %v", err)
	}

	expected := 2
	if count != expected {
		t.Errorf("expected %d completions for T3k7x-task1 since %v, got %d", expected, since, count)
	}

	count2, err := log.CountCompletionsForTaskSince("E2k7x-task2", since)
	if err != nil {
		t.Fatalf("failed to count completions: %v", err)
	}

	expected2 := 1
	if count2 != expected2 {
		t.Errorf("expected %d completions for E2k7x-task2 since %v, got %d", expected2, since, count2)
	}
}

func TestCountCompletionsSinceIncludesBoundary(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}

	now := time.Now().UTC()
	since := now.Add(-1 * time.Hour)

	entry := Entry{
		Timestamp: since,
		TaskID:    "T3k7x-boundary",
		Type:      EventTaskCompleted,
	}

	if err := log.WriteEntry(entry); err != nil {
		t.Fatalf("failed to write entry: %v", err)
	}

	count, err := log.CountCompletionsSince(since)
	if err != nil {
		t.Fatalf("failed to count completions: %v", err)
	}

	if count != 1 {
		t.Errorf("expected boundary entry to be counted, got %d", count)
	}
}

func TestReadEntriesHandlesMalformedEntry(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}

	entry := Entry{
		Timestamp: time.Now().UTC(),
		TaskID:    "T3k7x-valid",
		Type:      EventTaskCompleted,
		Report:    "Valid entry",
	}

	if err := log.WriteEntry(entry); err != nil {
		t.Fatalf("failed to write entry: %v", err)
	}

	log.Close()

	activityLogPath := filepath.Join(tmpDir, defaultLogFilename)
	f, err := os.OpenFile(activityLogPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("failed to open log for appending malformed data: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(`{"invalid": "json", "missing fields"}` + "\n"); err != nil {
		t.Fatalf("failed to write malformed entry: %v", err)
	}

	log, err = Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to reopen log: %v", err)
	}

	entries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("expected no error for malformed entry (resilient parsing), got %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 valid entry, got %d", len(entries))
	}

	if entries[0].TaskID != "T3k7x-valid" {
		t.Errorf("expected valid entry TaskID 'T3k7x-valid', got %s", entries[0].TaskID)
	}
}

func TestReadEntriesConcurrency(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "concurrency-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}
	defer log.Close()

	// Number of concurrent operations
	n := 100
	var wg sync.WaitGroup
	wg.Add(n * 2)

	errors := make(chan error, n*2)

	for i := 0; i < n; i++ {
		// Concurrent writers
		go func(id int) {
			defer wg.Done()
			err := log.WriteTaskCompletion(fmt.Sprintf("task-%d", id), "completed")
			if err != nil {
				errors <- fmt.Errorf("writer %d failed: %w", id, err)
			}
		}(i)

		// Concurrent readers
		go func(id int) {
			defer wg.Done()
			_, err := log.ReadEntries()
			if err != nil {
				errors <- fmt.Errorf("reader %d failed: %w", id, err)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	errCount := 0
	for err := range errors {
		errCount++
		t.Logf("Error: %v", err)
	}

	if errCount > 0 {
		t.Errorf("Encountered %d errors during concurrent operations", errCount)
	}

	// Final read to check consistency
	entries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("final read failed: %v", err)
	}
	if len(entries) != n {
		t.Errorf("expected %d entries, got %d", n, len(entries))
	}
}

func TestGetLatestTaskCompletionTime(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "latest-time-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}
	defer log.Close()

	taskID := "T3k7x-target"
	t1 := time.Date(2026, 2, 1, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 2, 1, 11, 0, 0, 0, time.UTC)
	t3 := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)

	log.WriteEntry(Entry{Timestamp: t1, TaskID: taskID, Type: EventTaskCompleted})
	log.WriteEntry(Entry{Timestamp: t2, TaskID: "other", Type: EventTaskCompleted})
	log.WriteEntry(Entry{Timestamp: t3, TaskID: taskID, Type: EventTaskCompleted})

	// Test 1: Finding latest (last entry)
	got, err := log.GetLatestTaskCompletionTime(taskID)
	if err != nil {
		// Log file content for debugging
		data, _ := os.ReadFile(filepath.Join(tmpDir, defaultLogFilename))
		t.Logf("Log content:\n%s", string(data))
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(t3) {
		t.Errorf("expected %v, got %v", t3, got)
	}

	// Test 2: Finding latest when it's not the last entry
	got, err = log.GetLatestTaskCompletionTime("other")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.Equal(t2) {
		t.Errorf("expected %v, got %v", t2, got)
	}

	// Test 3: Task never completed
	_, err = log.GetLatestTaskCompletionTime("non-existent")
	if err == nil {
		t.Error("expected error for non-existent task, got nil")
	}

	// Test 4: Verify it works after clearing cache
	log.lastSize = -1
	got, err = log.GetLatestTaskCompletionTime(taskID)
	if err != nil {
		t.Fatalf("unexpected error after cache clear: %v", err)
	}
	if !got.Equal(t3) {
		t.Errorf("expected %v, got %v", t3, got)
	}
}

// Comprehensive error recovery tests for corrupted activity logs

func TestErrorRecovery_TruncatedLogFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "truncated-log-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a log with some entries
	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}

	entries := []Entry{
		{TaskID: "T1aaa-first", Type: EventTaskCompleted, Report: "First task"},
		{TaskID: "T2bbb-second", Type: EventTaskCompleted, Report: "Second task"},
		{TaskID: "T3ccc-third", Type: EventTaskCompleted, Report: "Third task"},
	}

	for _, entry := range entries {
		if err := log.WriteEntry(entry); err != nil {
			t.Fatalf("failed to write entry: %v", err)
		}
	}

	if err := log.Close(); err != nil {
		t.Fatalf("failed to close log: %v", err)
	}

	// Truncate the file in the middle of a JSON line
	activityLogPath := filepath.Join(tmpDir, defaultLogFilename)
	data, err := os.ReadFile(activityLogPath)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	// Find a good place to truncate (middle of the file)
	truncatePos := len(data) / 2
	if err := os.WriteFile(activityLogPath, data[:truncatePos], 0o644); err != nil {
		t.Fatalf("failed to truncate file: %v", err)
	}

	// Reopen and verify recovery
	log, err = Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to reopen truncated log: %v", err)
	}
	defer log.Close()

	recoveredEntries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("expected no error reading truncated log, got %v", err)
	}

	// Should recover at least one complete entry
	if len(recoveredEntries) == 0 {
		t.Error("expected at least one recovered entry from truncated log")
	}

	// Verify recovered entries are valid
	for _, entry := range recoveredEntries {
		if entry.TaskID == "" {
			t.Error("recovered entry has empty TaskID")
		}
		if entry.Type == "" {
			t.Error("recovered entry has empty Type")
		}
	}
}

func TestErrorRecovery_MalformedJSONEntryTimestamps(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "malformed-timestamp-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create activity log with malformed timestamps
	activityLogPath := filepath.Join(tmpDir, defaultLogFilename)
	malformedContent := `{"timestamp":"not-a-valid-timestamp","task_id":"T1k7x-bad-time","type":"task_completed","report":"Bad timestamp"}
{"timestamp":"2026-02-02T12:00:00Z","task_id":"T2k7x-good-time","type":"task_completed","report":"Good timestamp"}
{"timestamp":"","task_id":"T3k7x-empty-time","type":"task_completed","report":"Empty timestamp"}
{"timestamp":1234567890,"task_id":"T4k7x-number-time","type":"task_completed","report":"Number timestamp"}
`

	if err := os.WriteFile(activityLogPath, []byte(malformedContent), 0o644); err != nil {
		t.Fatalf("failed to write malformed activity log: %v", err)
	}

	// Test that the log recovers gracefully
	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open malformed activity log: %v", err)
	}
	defer log.Close()

	entries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("expected no error reading malformed timestamp log, got %v", err)
	}

	// Should recover entries despite malformed timestamps
	if len(entries) == 0 {
		t.Error("expected at least one recovered entry from malformed timestamp log")
	}

	// Verify that valid entries are properly parsed
	foundValidEntry := false
	for _, entry := range entries {
		if entry.TaskID == "T2k7x-good-time" {
			foundValidEntry = true
			if entry.Type != EventTaskCompleted {
				t.Errorf("expected type %q, got %q", EventTaskCompleted, entry.Type)
			}
		}
	}

	if !foundValidEntry {
		t.Error("expected to find valid entry with good timestamp")
	}
}

func TestErrorRecovery_MissingRequiredFields(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "missing-fields-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create activity log with entries missing required fields
	activityLogPath := filepath.Join(tmpDir, defaultLogFilename)
	missingFieldsContent := `{"timestamp":"2026-02-02T12:00:00Z","type":"task_completed","report":"Missing task_id"}
{"timestamp":"2026-02-02T13:00:00Z","task_id":"T2k7x-good","type":"task_completed","report":"Complete valid entry"}
{"timestamp":"2026-02-02T14:00:00Z","task_id":"T3k7x-missing-type","report":"Missing type"}
{"task_id":"T4k7x-missing-timestamp","type":"task_completed","report":"Missing timestamp"}
{"type":"task_completed"}  // Missing all required fields
{"timestamp":"2026-02-02T15:00:00Z","task_id":"T5k7x-another-good","type":"task_completed","report":"Another valid entry"}
`

	if err := os.WriteFile(activityLogPath, []byte(missingFieldsContent), 0o644); err != nil {
		t.Fatalf("failed to write activity log with missing fields: %v", err)
	}

	// Test that the log recovers gracefully
	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open activity log with missing fields: %v", err)
	}
	defer log.Close()

	entries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("expected no error reading log with missing fields, got %v", err)
	}

	// Should recover the valid entries
	validTaskIDs := []string{"T2k7x-good", "T5k7x-another-good"}
	foundCount := 0

	for _, entry := range entries {
		for _, validID := range validTaskIDs {
			if entry.TaskID == validID {
				foundCount++
				if entry.Type != EventTaskCompleted {
					t.Errorf("entry %s: expected type %q, got %q", entry.TaskID, EventTaskCompleted, entry.Type)
				}
				break
			}
		}
	}

	if foundCount != len(validTaskIDs) {
		t.Errorf("expected to recover %d valid entries, found %d", len(validTaskIDs), foundCount)
	}

	// Create a task with invalid YAML frontmatter
	taskDir := filepath.Join(tmpDir, "tasks", "T1aaa-invalid")
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatalf("failed to create task dir: %v", err)
	}

	invalidTaskContent := `---
role: developer
priority: high
parent: [invalid yaml array format
blockers: "invalid yaml string instead of array"
date_created: 2026-02-02T00:00:00Z
---

# Invalid Task

This task has invalid YAML frontmatter.
`

	taskFile := filepath.Join(taskDir, "T1aaa-invalid.md")
	if err := os.WriteFile(taskFile, []byte(invalidTaskContent), 0o644); err != nil {
		t.Fatalf("failed to write invalid task file: %v", err)
	}

	// Test that the system can handle this gracefully
	parser := task.NewParser()
	_, err = parser.ParseFile(taskFile)
	if err == nil {
		t.Error("expected error parsing invalid YAML frontmatter")
	}

	// The error should be informative and not crash
	if !strings.Contains(err.Error(), "yaml") && !strings.Contains(err.Error(), "frontmatter") {
		t.Errorf("expected YAML-related error, got: %v", err)
	}
}

func TestErrorRecovery_MalformedTaskIDs(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "malformed-id-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name        string
		taskID      string
		expectError bool
	}{
		{"valid_id", "T1k7x-valid", false},
		{"too_short", "T1k7x", true}, // only 4 chars
		{"no_prefix", "1k7x-example", true},
		{"invalid_chars", "T1@7x-invalid", true},
		{"uppercase_only", "t1k7x-lowercase", true}, // prefix must be uppercase
		{"missing_dash", "T1k7xexample", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			content := fmt.Sprintf(`---
role: developer
priority: medium
---

# Test Task

Content for %s
`, tc.taskID)

			parser := task.NewParser()
			_, err := parser.ParseString(content, tc.taskID)

			if tc.expectError && err == nil {
				t.Errorf("expected error for task ID %q, got none", tc.taskID)
			}
			if !tc.expectError && err != nil {
				t.Errorf("expected no error for task ID %q, got %v", tc.taskID, err)
			}
		})
	}
}

func TestErrorRecovery_CircularReferences(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "circular-ref-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create tasks with circular references
	taskDir := filepath.Join(tmpDir, "tasks")
	if err := os.MkdirAll(taskDir, 0o755); err != nil {
		t.Fatalf("failed to create tasks dir: %v", err)
	}

	// Task A blocks B
	taskAContent := `---
role: developer
priority: medium
blockers: []
blocks: ["T2bbb-circular-b"]
---

# Task A

Blocks task B
`

	taskADir := filepath.Join(taskDir, "T1aaa-circular-a")
	if err := os.MkdirAll(taskADir, 0o755); err != nil {
		t.Fatalf("failed to create task A dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskADir, "T1aaa-circular-a.md"), []byte(taskAContent), 0o644); err != nil {
		t.Fatalf("failed to write task A file: %v", err)
	}

	// Task B blocks A (creating circular dependency)
	taskBContent := `---
role: developer
priority: medium
blockers: []
blocks: ["T1aaa-circular-a"]
---

# Task B

Blocks task A (circular)
`

	taskBDir := filepath.Join(taskDir, "T2bbb-circular-b")
	if err := os.MkdirAll(taskBDir, 0o755); err != nil {
		t.Fatalf("failed to create task B dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(taskBDir, "T2bbb-circular-b.md"), []byte(taskBContent), 0o644); err != nil {
		t.Fatalf("failed to write task B file: %v", err)
	}

	// Test that the validator detects circular references
	parser := task.NewParser()
	tasks := make(map[string]*task.Task)

	taskA, _ := parser.ParseFile(filepath.Join(taskADir, "T1aaa-circular-a.md"))
	tasks[taskA.ID] = taskA

	taskB, _ := parser.ParseFile(filepath.Join(taskBDir, "T2bbb-circular-b.md"))
	tasks[taskB.ID] = taskB

	v := task.NewValidator(tasks)
	errors := v.ValidateAndRepair()

	// Should detect circular reference
	foundCircular := false
	for _, err := range errors {
		if strings.Contains(err.Message, "circular") || strings.Contains(err.Message, "cycle") {
			foundCircular = true
			break
		}
	}

	if !foundCircular {
		t.Error("expected validation to detect circular reference")
	}
}

func TestErrorRecovery_MissingParentChildRelationships(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "parent-child-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a child task without a parent
	childDir := filepath.Join(tmpDir, "tasks", "T1aaa-orphan")
	if err := os.MkdirAll(childDir, 0o755); err != nil {
		t.Fatalf("failed to create child task dir: %v", err)
	}

	childContent := `---
role: developer
priority: medium
parent: "E1aaa-missing-parent"
---

# Orphan Task

This task references a non-existent parent.
`

	childFile := filepath.Join(childDir, "T1aaa-orphan.md")
	if err := os.WriteFile(childFile, []byte(childContent), 0o644); err != nil {
		t.Fatalf("failed to write child task file: %v", err)
	}

	// Test that the system detects missing parent
	parser := task.NewParser()
	ttask, err := parser.ParseFile(childFile)
	if err != nil {
		t.Fatalf("failed to parse child task: %v", err)
	}

	// Create a validator with only the child task
	tasks := map[string]*task.Task{ttask.ID: ttask}
	v := task.NewValidator(tasks)
	errors := v.ValidateAndRepair()

	// Should detect missing parent
	foundMissingParent := false
	for _, err := range errors {
		if strings.Contains(err.Message, "parent") && strings.Contains(err.Message, "missing") {
			foundMissingParent = true
			break
		}
	}

	if !foundMissingParent {
		t.Error("expected validation to detect missing parent")
	}
}

func TestErrorRecovery_CorruptedJSONLines(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "corrupted-json-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an activity log with corrupted JSON lines
	activityLogPath := filepath.Join(tmpDir, defaultLogFilename)
	corruptedContent := `{"timestamp":"2026-02-02T12:00:00Z","task_id":"T1k7x-valid","type":"task_completed","report":"Valid entry"}
{"invalid": "json", "missing: "comma", "unclosed: "quote"
{"timestamp":"2026-02-02T13:00:00Z","task_id":"T2k7x-valid","type":"task_completed","report":"Another valid entry"}
{"malformed json without quotes}
{"timestamp":"2026-02-02T14:00:00Z","task_id":"T3k7x-valid","type":"task_completed","report":"Third valid entry"}
`

	if err := os.WriteFile(activityLogPath, []byte(corruptedContent), 0o644); err != nil {
		t.Fatalf("failed to write corrupted activity log: %v", err)
	}

	// Test that the log recovers gracefully
	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open corrupted activity log: %v", err)
	}
	defer log.Close()

	entries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("expected no error reading corrupted log, got %v", err)
	}

	// Should recover the 3 valid entries
	if len(entries) != 3 {
		t.Errorf("expected 3 valid entries, got %d", len(entries))
	}

	// Verify recovered entries are valid
	validTaskIDs := []string{"T1k7x-valid", "T2k7x-valid", "T3k7x-valid"}
	for i, entry := range entries {
		if entry.TaskID != validTaskIDs[i] {
			t.Errorf("entry %d: expected TaskID %q, got %q", i, validTaskIDs[i], entry.TaskID)
		}
		if entry.Type != EventTaskCompleted {
			t.Errorf("entry %d: expected type %q, got %q", i, EventTaskCompleted, entry.Type)
		}
	}
}

func TestErrorRecovery_EmptyLogFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "empty-log-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an empty activity log file
	activityLogPath := filepath.Join(tmpDir, defaultLogFilename)
	if err := os.WriteFile(activityLogPath, []byte{}, 0o644); err != nil {
		t.Fatalf("failed to create empty activity log: %v", err)
	}

	// Test that the system handles empty logs gracefully
	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open empty activity log: %v", err)
	}
	defer log.Close()

	entries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("expected no error reading empty log, got %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("expected 0 entries from empty log, got %d", len(entries))
	}

	// Test that we can still write new entries
	newEntry := Entry{
		TaskID: "T1k7x-new",
		Type:   EventTaskCompleted,
		Report: "New entry in empty log",
	}

	if err := log.WriteEntry(newEntry); err != nil {
		t.Fatalf("failed to write entry to previously empty log: %v", err)
	}

	entries, err = log.ReadEntries()
	if err != nil {
		t.Fatalf("failed to read entries after writing to empty log: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 entry after writing to empty log, got %d", len(entries))
	}
}

func TestErrorRecovery_ActivityLogPermissions(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("skipping permissions test when running as root")
	}

	tmpDir, err := os.MkdirTemp("", "permissions-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a log file with no write permissions
	activityLogPath := filepath.Join(tmpDir, defaultLogFilename)
	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open activity log: %v", err)
	}
	log.Close()

	// Remove write permissions
	if err := os.Chmod(activityLogPath, 0o444); err != nil {
		t.Fatalf("failed to remove write permissions: %v", err)
	}

	// Test that write operations fail gracefully
	log, err = Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to reopen activity log: %v", err)
	}
	defer log.Close()

	entry := Entry{
		TaskID: "T1k7x-test",
		Type:   EventTaskCompleted,
		Report: "Test entry",
	}

	err = log.WriteEntry(entry)
	if err == nil {
		t.Error("expected error writing to read-only log file")
	}

	// Restore permissions for cleanup
	if err := os.Chmod(activityLogPath, 0o644); err != nil {
		t.Logf("failed to restore permissions: %v", err)
	}
}

func TestErrorRecovery_RapidConcurrentCorruption(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "concurrent-corruption-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open log: %v", err)
	}
	defer log.Close()

	// Simulate rapid concurrent writes that could cause corruption
	n := 50
	var wg sync.WaitGroup
	errors := make(chan error, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			entry := Entry{
				TaskID: fmt.Sprintf("T%04d-concurrent", id),
				Type:   EventTaskCompleted,
				Report: fmt.Sprintf("Concurrent entry %d", id),
			}
			if err := log.WriteEntry(entry); err != nil {
				errors <- fmt.Errorf("writer %d failed: %w", id, err)
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Count any errors
	errCount := 0
	for err := range errors {
		errCount++
		t.Logf("Error during concurrent writes: %v", err)
	}

	// Verify log integrity after concurrent writes
	entries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("failed to read entries after concurrent writes: %v", err)
	}

	// Should have all successful writes minus any errors
	expectedEntries := n - errCount
	if len(entries) != expectedEntries {
		t.Errorf("expected %d entries after concurrent writes, got %d", expectedEntries, len(entries))
	}

	// Verify all entries are valid
	taskIDs := make(map[string]bool)
	for _, entry := range entries {
		if entry.TaskID == "" || entry.Type == "" {
			t.Error("found invalid entry after concurrent writes")
		}
		if taskIDs[entry.TaskID] {
			t.Errorf("found duplicate task ID: %s", entry.TaskID)
		}
		taskIDs[entry.TaskID] = true
	}
}
