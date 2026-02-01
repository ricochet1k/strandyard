package activity

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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
