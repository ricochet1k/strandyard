package activity

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultLogFilename = "activity.log"
)

// EventType represents the type of activity event
type EventType string

const (
	EventTaskCompleted EventType = "task_completed"
)

// Entry represents a single activity log entry
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	TaskID    string    `json:"task_id"`
	Type      EventType `json:"type"`
	Report    string    `json:"report,omitempty"`
}

// Log represents the activity log
type Log struct {
	filepath string
	file     *os.File
}

// Open opens the activity log for appending
func Open(logDir string) (*Log, error) {
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	fp := filepath.Join(logDir, defaultLogFilename)
	file, err := os.OpenFile(fp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open activity log: %w", err)
	}

	return &Log{
		filepath: fp,
		file:     file,
	}, nil
}

// Close closes the activity log file
func (l *Log) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// WriteEntry writes a new entry to the activity log
func (l *Log) WriteEntry(entry Entry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	if _, err := l.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	return l.file.Sync()
}

// WriteTaskCompletion writes a task completion event to the activity log
func (l *Log) WriteTaskCompletion(taskID, report string) error {
	return l.WriteEntry(Entry{
		TaskID: taskID,
		Type:   EventTaskCompleted,
		Report: report,
	})
}

// ReadEntries reads all entries from the activity log
func (l *Log) ReadEntries() ([]Entry, error) {
	if err := l.file.Close(); err != nil {
		return nil, fmt.Errorf("failed to close log for reading: %w", err)
	}

	file, err := os.Open(l.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log for reading: %w", err)
	}
	defer file.Close()

	var entries []Entry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal entry: %w", err)
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log: %w", err)
	}

	l.file, err = os.OpenFile(l.filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to reopen log for writing: %w", err)
	}

	return entries, nil
}

// CountCompletionsSince counts task completion events since a given time
func (l *Log) CountCompletionsSince(since time.Time) (int, error) {
	entries, err := l.ReadEntries()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if entry.Type == EventTaskCompleted && entry.Timestamp.After(since) || entry.Timestamp.Equal(since) {
			count++
		}
	}
	return count, nil
}

// CountCompletionsForTaskSince counts completion events for a specific task since a given time
func (l *Log) CountCompletionsForTaskSince(taskID string, since time.Time) (int, error) {
	entries, err := l.ReadEntries()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if entry.Type == EventTaskCompleted && entry.TaskID == taskID && (entry.Timestamp.After(since) || entry.Timestamp.Equal(since)) {
			count++
		}
	}
	return count, nil
}
