package cmd

import (
	"testing"
)

func TestValidateEvery(t *testing.T) {
	// Test that validateEvery handles various inputs correctly
	tests := []struct {
		name    string
		every   []string
		isValid bool
	}{
		{"empty every slice", []string{}, true},
		{"valid default days", []string{"10 days"}, true},
		{"valid commits with HEAD", []string{"50 commits from HEAD"}, true},
		{"valid lines_changed with HEAD", []string{"500 lines_changed from HEAD"}, true},
		{"valid tasks_completed", []string{"20 tasks_completed"}, true},
		{"invalid format - missing metric", []string{"10"}, false},
		{"invalid amount - non-integer", []string{"ten days"}, false},
		{"unsupported metric", []string{"10 hours"}, false},
		{"invalid commit anchor", []string{"50 commits from invalid"}, false},
		{"invalid date anchor", []string{"10 days from invalid date"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEvery(tt.every)

			if tt.isValid {
				if err != nil {
					t.Errorf("validateEvery(%v) expected no error but got %v", tt.every, err)
				}
			} else {
				if err == nil {
					t.Errorf("validateEvery(%v) expected error but got nil", tt.every)
				}
			}
		})
	}
}
