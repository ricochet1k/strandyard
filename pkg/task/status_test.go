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
