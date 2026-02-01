package activity

import (
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
