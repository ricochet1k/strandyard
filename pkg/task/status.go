package task

import "strings"

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
