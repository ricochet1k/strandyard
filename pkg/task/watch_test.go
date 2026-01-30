package task

import (
	"path/filepath"
	"testing"
)

func TestIsTaskFilePath(t *testing.T) {
	root := filepath.Join("tasks", "T1234-example")
	cases := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "task.md",
			path: filepath.Join(root, "task.md"),
			want: true,
		},
		{
			name: "README.md",
			path: filepath.Join(root, "README.md"),
			want: true,
		},
		{
			name: "named by directory",
			path: filepath.Join(root, "T1234-example.md"),
			want: true,
		},
		{
			name: "root list",
			path: filepath.Join("tasks", "root-tasks.md"),
			want: false,
		},
		{
			name: "other markdown",
			path: filepath.Join(root, "notes.md"),
			want: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if got := isTaskFilePath(tc.path); got != tc.want {
				t.Fatalf("isTaskFilePath(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}
