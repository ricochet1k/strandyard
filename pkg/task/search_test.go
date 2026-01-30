package task

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSearchTasksMatchesTitleBodyTodosAndSkipsSubtasks(t *testing.T) {
	root := filepath.Join(t.TempDir(), "tasks")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("mkdir tasks root: %v", err)
	}

	role := testRoleName(t, "search")
	writeSearchTask(t, root, role, "T1a1a-child", "Child Task", strings.Join([]string{
		"## Summary",
		"Child summary",
		"",
		"## Tasks",
		"- [ ] Implement child thing",
	}, "\n"))

	writeSearchTask(t, root, role, "T2a1a-parent", "Parent Task", strings.Join([]string{
		"## Summary",
		"Parent summary",
		"",
		"## Tasks",
		"- [ ] (subtask: T1a1a-child) Child Task",
		"- [ ] Write docs",
	}, "\n"))

	cases := []struct {
		name  string
		query string
		want  []string
	}{
		{
			name:  "matches title",
			query: "Parent Task",
			want:  []string{"T2a1a-parent"},
		},
		{
			name:  "matches description",
			query: "Parent summary",
			want:  []string{"T2a1a-parent"},
		},
		{
			name:  "matches todos",
			query: "WRITE DOCS",
			want:  []string{"T2a1a-parent"},
		},
		{
			name:  "skips subtask names",
			query: "Child Task",
			want:  []string{"T1a1a-child"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tasks, err := SearchTasks(root, SearchOptions{Query: tc.query})
			if err != nil {
				t.Fatalf("SearchTasks failed: %v", err)
			}
			got := make([]string, 0, len(tasks))
			for _, task := range tasks {
				got = append(got, task.ID)
			}
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Fatalf("unexpected results\n got: %v\nwant: %v", got, tc.want)
			}
		})
	}
}

func TestSearchTasksRequiresQuery(t *testing.T) {
	root := filepath.Join(t.TempDir(), "tasks")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("mkdir tasks root: %v", err)
	}

	_, err := SearchTasks(root, SearchOptions{Query: " "})
	if err == nil {
		t.Fatalf("expected error for empty query")
	}
}

func writeSearchTask(t *testing.T, root, role, id, title, body string) {
	t.Helper()

	dir := filepath.Join(root, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir task dir: %v", err)
	}

	filePath := filepath.Join(dir, id+".md")
	content := strings.Join([]string{
		"---",
		"role: " + role,
		"priority: medium",
		"parent: ",
		"blockers: []",
		"blocks: []",
		"date_created: 2026-01-01T00:00:00Z",
		"date_edited: 2026-01-01T00:00:00Z",
		"owner_approval: false",
		"completed: false",
		"---",
		"",
		"# " + title,
		"",
		body,
		"",
	}, "\n")

	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write task file: %v", err)
	}
}
