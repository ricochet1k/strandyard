package cmd

import (
	"strings"
	"testing"

	"github.com/ricochet1k/streamyard/pkg/task"
)

func TestRunListValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		opts    task.ListOptions
		wantErr bool
	}{
		{
			name:    "invalid scope",
			opts:    task.ListOptions{Scope: "nope"},
			wantErr: true,
		},
		{
			name:    "invalid priority",
			opts:    task.ListOptions{Scope: "all", Priority: "urgent"},
			wantErr: true,
		},
		{
			name:    "invalid sort",
			opts:    task.ListOptions{Scope: "all", Sort: "title"},
			wantErr: true,
		},
		{
			name:    "invalid order",
			opts:    task.ListOptions{Scope: "all", Order: "down"},
			wantErr: true,
		},
		{
			name:    "invalid format",
			opts:    task.ListOptions{Scope: "all", Format: "yaml"},
			wantErr: true,
		},
		{
			name:    "invalid group",
			opts:    task.ListOptions{Scope: "all", Group: "status"},
			wantErr: true,
		},
		{
			name:    "label unsupported",
			opts:    task.ListOptions{Scope: "all", Label: "bug"},
			wantErr: true,
		},
		{
			name:    "scope free with parent",
			opts:    task.ListOptions{Scope: "free", Parent: "T0000"},
			wantErr: true,
		},
		{
			name:    "scope free with group parent",
			opts:    task.ListOptions{Scope: "free", Group: "parent"},
			wantErr: true,
		},
		{
			name:    "parent with scope root",
			opts:    task.ListOptions{Scope: "root", Parent: "T0000"},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var output strings.Builder
			err := runList(&output, "tasks", tc.opts)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
