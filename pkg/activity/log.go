package activity

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultLogFilename = "activity.log"
)

// EventType represents the type of activity event
type EventType string

const (
	EventTaskCompleted            EventType = "task_completed"
	EventRecurrenceAnchorResolved EventType = "recurrence_anchor_resolved"
)

// Entry represents a single activity log entry
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	TaskID    string            `json:"task_id"`
	Type      EventType         `json:"type"`
	Report    string            `json:"report,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Log represents the activity log
type Log struct {
	mu       sync.RWMutex
	filepath string
	file     *os.File // write handle
	entries  []Entry  // cached entries
	lastSize int64    // last read size
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
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		err := l.file.Close()
		l.file = nil
		return err
	}
	return nil
}

// WriteEntry writes a new entry to the activity log
func (l *Log) WriteEntry(entry Entry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

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

	if err := l.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync log: %w", err)
	}

	// Invalidate cache by setting lastSize to -1 or just clear it.
	// We'll let ReadEntries handle re-reading.
	l.lastSize = -1

	return nil
}

// WriteTaskCompletion writes a task completion event to the activity log
func (l *Log) WriteTaskCompletion(taskID, report string) error {
	return l.WriteEntry(Entry{
		TaskID: taskID,
		Type:   EventTaskCompleted,
		Report: report,
	})
}

// WriteRecurrenceAnchorResolution writes a recurrence anchor resolution event to the activity log
func (l *Log) WriteRecurrenceAnchorResolution(taskID, original, resolved string) error {
	return l.WriteEntry(Entry{
		TaskID: taskID,
		Type:   EventRecurrenceAnchorResolved,
		Metadata: map[string]string{
			"original": original,
			"resolved": resolved,
		},
	})
}

// ReadEntries reads all entries from the activity log
func (l *Log) ReadEntries() ([]Entry, error) {
	l.mu.RLock()
	info, err := os.Stat(l.filepath)
	if err == nil && l.lastSize != -1 && info.Size() == l.lastSize {
		entries := make([]Entry, len(l.entries))
		copy(entries, l.entries)
		l.mu.RUnlock()
		return entries, nil
	}
	l.mu.RUnlock()

	l.mu.Lock()
	defer l.mu.Unlock()

	// Re-check size after acquiring write lock
	info, err = os.Stat(l.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat log: %w", err)
	}

	// If lastSize is -1 or file shrunk, re-read everything
	if l.lastSize == -1 || info.Size() < l.lastSize {
		l.entries = nil
		l.lastSize = 0
	}

	if info.Size() == l.lastSize {
		entries := make([]Entry, len(l.entries))
		copy(entries, l.entries)
		return entries, nil
	}

	file, err := os.Open(l.filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log for reading: %w", err)
	}
	defer file.Close()

	if _, err := file.Seek(l.lastSize, io.SeekStart); err != nil {
		return nil, fmt.Errorf("failed to seek to last position: %w", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Resilient parsing: skip malformed entries
			fmt.Fprintf(os.Stderr, "skipping malformed activity log entry: %v\n", err)
			continue
		}
		l.entries = append(l.entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log: %w", err)
	}

	l.lastSize = info.Size()

	entries := make([]Entry, len(l.entries))
	copy(entries, l.entries)
	return entries, nil
}

// GetLatestTaskCompletionTime returns the timestamp of the most recent completion of the given task.
func (l *Log) GetLatestTaskCompletionTime(taskID string) (time.Time, error) {
	l.mu.RLock()
	// Optimization: if cache is up to date, search it backwards
	info, err := os.Stat(l.filepath)
	if err == nil && l.lastSize != -1 && info.Size() == l.lastSize {
		for i := len(l.entries) - 1; i >= 0; i-- {
			entry := l.entries[i]
			if entry.TaskID == taskID && entry.Type == EventTaskCompleted {
				l.mu.RUnlock()
				return entry.Timestamp, nil
			}
		}
		l.mu.RUnlock()
		return time.Time{}, fmt.Errorf("task %s never completed (cached)", taskID)
	}
	l.mu.RUnlock()

	// Otherwise, scan the file backwards
	file, err := os.Open(l.filepath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to open log for reading: %w", err)
	}
	defer file.Close()

	scanner, err := NewReverseScanner(file)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to create reverse scanner: %w", err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			// Skip malformed entries
			continue
		}
		if entry.TaskID == taskID && entry.Type == EventTaskCompleted {
			return entry.Timestamp, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return time.Time{}, fmt.Errorf("error scanning log backwards: %w", err)
	}

	return time.Time{}, fmt.Errorf("task %s never completed", taskID)
}

// CountCompletionsSince counts task completion events since a given time
func (l *Log) CountCompletionsSince(since time.Time) (int, error) {
	entries, err := l.ReadEntries()
	if err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if entry.Type == EventTaskCompleted && (entry.Timestamp.After(since) || entry.Timestamp.Equal(since)) {
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

// GetCompletionTimestampAtOffset returns the timestamp of the 'offset'-th task completion since 'since'.
func (l *Log) GetCompletionTimestampAtOffset(since time.Time, offset int) (time.Time, error) {
	entries, err := l.ReadEntries()
	if err != nil {
		return time.Time{}, err
	}

	count := 0
	for _, entry := range entries {
		if entry.Type == EventTaskCompleted && (entry.Timestamp.After(since) || entry.Timestamp.Equal(since)) {
			count++
			if count == offset {
				return entry.Timestamp, nil
			}
		}
	}
	return time.Time{}, fmt.Errorf("offset %d not reached since %v (found %d)", offset, since, count)
}
