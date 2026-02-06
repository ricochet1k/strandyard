package task

import (
	"fmt"
	"strings"
)

const (
	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
	StatusCancelled  = "cancelled"
	StatusDuplicate  = "duplicate"
)

// NormalizeStatus returns a canonical status string.
// Empty statuses remain empty (no default).
func NormalizeStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

// IsValidStatus returns true when the status is empty or one of the allowed values.
// Allowed values are: open, in_progress, done, cancelled, duplicate.
func IsValidStatus(status string) bool {
	if strings.TrimSpace(status) == "" {
		return true
	}
	switch NormalizeStatus(status) {
	case StatusOpen, StatusInProgress, StatusDone, StatusCancelled, StatusDuplicate:
		return true
	default:
		return false
	}
}

// AllowedStatusValues returns a slice of all allowed status values (excluding empty).
func AllowedStatusValues() []string {
	return []string{StatusOpen, StatusInProgress, StatusDone, StatusCancelled, StatusDuplicate}
}

// FormatStatusListForUser returns a human-readable list of allowed status values.
// Format: "open, in_progress, done, cancelled, or duplicate"
func FormatStatusListForUser() string {
	values := AllowedStatusValues()
	switch len(values) {
	case 0:
		return ""
	case 1:
		return values[0]
	case 2:
		return values[0] + " or " + values[1]
	default:
		return strings.Join(values[:len(values)-1], ", ") + ", or " + values[len(values)-1]
	}
}

// FormatStatusErrorMessage returns a user-friendly error message for invalid status values.
// Includes hints for common mistakes like "completed" or "pending".
func FormatStatusErrorMessage(status string) string {
	normalized := NormalizeStatus(status)
	base := fmt.Sprintf("invalid status %q: must be one of %s or empty", status, FormatStatusListForUser())

	// Provide helpful hints for common mistakes
	hint := getStatusHint(normalized)
	if hint != "" {
		base += "\n" + hint
	}

	return base
}

// getStatusHint returns a helpful hint for common mistakes
func getStatusHint(normalizedStatus string) string {
	switch normalizedStatus {
	case "completed":
		return "Did you mean 'done'? Use 'done' to mark a task as completed."
	case "pending", "waiting", "paused":
		return "Did you mean 'open' or 'in_progress'? Use 'open' for not yet started or 'in_progress' for active work."
	case "failed":
		return "Did you mean 'cancelled'? Use 'cancelled' for tasks that won't be completed."
	case "wontfix":
		return "Did you mean 'cancelled'? Use 'cancelled' for tasks that won't be completed."
	case "blocked":
		return "Use 'open' or 'in_progress' for task status. Blocking is managed via the 'blockers' field instead."
	}
	return ""
}

// IsActiveStatus returns true if a task status should be included in the free-list.
// Active statuses are "open" and "in_progress".
// Empty status defaults to "open" for backward compatibility.
func IsActiveStatus(status string) bool {
	if status == "" {
		return true // Default to active for backward compatibility
	}
	return status == StatusOpen || status == StatusInProgress
}

// IsOpen returns true if the task is in "open" status.
func (m *Metadata) IsOpen() bool {
	return m.Status == StatusOpen || m.Status == ""
}

// IsInProgress returns true if the task is in "in_progress" status.
func (m *Metadata) IsInProgress() bool {
	return m.Status == StatusInProgress
}

// IsDone returns true if the task is in "done" status.
func (m *Metadata) IsDone() bool {
	return m.Status == StatusDone
}

// IsCancelled returns true if the task is in "cancelled" status.
func (m *Metadata) IsCancelled() bool {
	return m.Status == StatusCancelled
}

// IsDuplicate returns true if the task is in "duplicate" status.
func (m *Metadata) IsDuplicate() bool {
	return m.Status == StatusDuplicate
}

// IsActive returns true if the task is in an active status (open or in_progress).
func (m *Metadata) IsActive() bool {
	return IsActiveStatus(m.Status)
}

// IsOpen returns true if the task is in "open" status.
func (t *Task) IsOpen() bool {
	return t.Meta.IsOpen()
}

// IsInProgress returns true if the task is in "in_progress" status.
func (t *Task) IsInProgress() bool {
	return t.Meta.IsInProgress()
}

// IsDone returns true if the task is in "done" status.
func (t *Task) IsDone() bool {
	return t.Meta.IsDone()
}

// IsCancelled returns true if the task is in "cancelled" status.
func (t *Task) IsCancelled() bool {
	return t.Meta.IsCancelled()
}

// IsDuplicate returns true if the task is in "duplicate" status.
func (t *Task) IsDuplicate() bool {
	return t.Meta.IsDuplicate()
}

// IsActive returns true if the task is in an active status (open or in_progress).
func (t *Task) IsActive() bool {
	return t.Meta.IsActive()
}
