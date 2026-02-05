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
