package task

import "strings"

const (
	PriorityHigh   = "high"
	PriorityMedium = "medium"
	PriorityLow    = "low"
)

// NormalizePriority returns a canonical priority string.
// Empty priorities default to "medium".
func NormalizePriority(priority string) string {
	p := strings.ToLower(strings.TrimSpace(priority))
	if p == "" {
		return PriorityMedium
	}
	return p
}

// PriorityRank returns a sortable rank where lower is higher priority.
func PriorityRank(priority string) int {
	switch NormalizePriority(priority) {
	case PriorityHigh:
		return 0
	case PriorityMedium:
		return 1
	case PriorityLow:
		return 2
	default:
		return 3
	}
}

// IsValidPriority returns true when the priority is empty or one of the known values.
func IsValidPriority(priority string) bool {
	if strings.TrimSpace(priority) == "" {
		return true
	}
	switch NormalizePriority(priority) {
	case PriorityHigh, PriorityMedium, PriorityLow:
		return true
	default:
		return false
	}
}
