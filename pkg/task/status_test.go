package task

import (
	"testing"
)

func TestIsValidStatus(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		expectValid bool
	}{
		{
			name:        "Valid: open",
			status:      "open",
			expectValid: true,
		},
		{
			name:        "Valid: in_progress",
			status:      "in_progress",
			expectValid: true,
		},
		{
			name:        "Valid: done",
			status:      "done",
			expectValid: true,
		},
		{
			name:        "Valid: cancelled",
			status:      "cancelled",
			expectValid: true,
		},
		{
			name:        "Valid: duplicate",
			status:      "duplicate",
			expectValid: true,
		},
		{
			name:        "Valid: empty string",
			status:      "",
			expectValid: true,
		},
		{
			name:        "Valid: whitespace only",
			status:      "   ",
			expectValid: true,
		},
		{
			name:        "Invalid: invalid_status",
			status:      "invalid_status",
			expectValid: false,
		},
		{
			name:        "Invalid: completed",
			status:      "completed",
			expectValid: false,
		},
		{
			name:        "Invalid: pending",
			status:      "pending",
			expectValid: false,
		},
		{
			name:        "Case insensitive: OPEN",
			status:      "OPEN",
			expectValid: true,
		},
		{
			name:        "Case insensitive: In_Progress",
			status:      "In_Progress",
			expectValid: true,
		},
		{
			name:        "Case insensitive: DONE",
			status:      "DONE",
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := IsValidStatus(tt.status)
			if valid != tt.expectValid {
				t.Errorf("IsValidStatus(%q) = %v, want %v", tt.status, valid, tt.expectValid)
			}
		})
	}
}

func TestNormalizeStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected string
	}{
		{
			name:     "lowercase: open",
			status:   "open",
			expected: "open",
		},
		{
			name:     "uppercase: OPEN",
			status:   "OPEN",
			expected: "open",
		},
		{
			name:     "mixed case: In_Progress",
			status:   "In_Progress",
			expected: "in_progress",
		},
		{
			name:     "whitespace: ' done '",
			status:   " done ",
			expected: "done",
		},
		{
			name:     "empty string",
			status:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			normalized := NormalizeStatus(tt.status)
			if normalized != tt.expected {
				t.Errorf("NormalizeStatus(%q) = %q, want %q", tt.status, normalized, tt.expected)
			}
		})
	}
}

func TestAllowedStatusValues(t *testing.T) {
	allowed := AllowedStatusValues()

	expectedValues := []string{
		StatusOpen,
		StatusInProgress,
		StatusDone,
		StatusCancelled,
		StatusDuplicate,
	}

	if len(allowed) != len(expectedValues) {
		t.Errorf("AllowedStatusValues() returned %d values, want %d", len(allowed), len(expectedValues))
	}

	// Check that all expected values are present
	for _, expected := range expectedValues {
		found := false
		for _, value := range allowed {
			if value == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("AllowedStatusValues() missing expected value: %q", expected)
		}
	}
}

func TestFormatStatusListForUser(t *testing.T) {
	expected := "open, in_progress, done, cancelled, or duplicate"
	result := FormatStatusListForUser()
	if result != expected {
		t.Errorf("FormatStatusListForUser() = %q, want %q", result, expected)
	}
}

func TestFormatStatusErrorMessage(t *testing.T) {
	tests := []struct {
		name    string
		status  string
		contain []string // strings that should be in the error message
	}{
		{
			name:   "invalid status without hint",
			status: "invalid_status",
			contain: []string{
				`invalid status "invalid_status"`,
				"open, in_progress, done, cancelled, or duplicate",
			},
		},
		{
			name:   "completed - should provide hint",
			status: "completed",
			contain: []string{
				`invalid status "completed"`,
				"Did you mean 'done'?",
				"mark a task as completed",
			},
		},
		{
			name:   "pending - should provide hint",
			status: "pending",
			contain: []string{
				`invalid status "pending"`,
				"Did you mean 'open' or 'in_progress'?",
			},
		},
		{
			name:   "blocked - should provide hint",
			status: "blocked",
			contain: []string{
				`invalid status "blocked"`,
				"blockers",
			},
		},
		{
			name:   "failed - should suggest cancelled",
			status: "failed",
			contain: []string{
				`invalid status "failed"`,
				"Did you mean 'cancelled'?",
			},
		},
		{
			name:   "wontfix - should suggest cancelled",
			status: "wontfix",
			contain: []string{
				`invalid status "wontfix"`,
				"Did you mean 'cancelled'?",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := FormatStatusErrorMessage(tt.status)
			for _, expected := range tt.contain {
				if !contains(msg, expected) {
					t.Errorf("FormatStatusErrorMessage(%q) missing %q\nGot: %q", tt.status, expected, msg)
				}
			}
		})
	}
}

// contains checks if needle is a substring of haystack
func contains(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
