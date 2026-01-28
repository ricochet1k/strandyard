package task

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestListFilteringAndSorting(t *testing.T) {
	t.Parallel()

	root := setupListFixture(t)

	cases := []struct {
		name string
		opts ListOptions
		want []string
	}{
		{
			name: "scope root",
			opts: ListOptions{Scope: "root"},
			want: []string{"E1a1a-epic", "T4a1a-completed", "T3a1a-blocked", "T5a1a-blocks", "T2a1a-free"},
		},
		{
			name: "scope free",
			opts: ListOptions{Scope: "free"},
			want: []string{"E1a1a-epic", "T4a1a-completed", "T1a1a-child", "T5a1a-blocks", "T2a1a-free"},
		},
		{
			name: "parent filter",
			opts: ListOptions{Scope: "all", Parent: "E1a1a-epic"},
			want: []string{"T1a1a-child"},
		},
		{
			name: "role filter",
			opts: ListOptions{Scope: "all", Role: "developer"},
			want: []string{"T4a1a-completed", "T1a1a-child", "T3a1a-blocked"},
		},
		{
			name: "priority filter",
			opts: ListOptions{Scope: "all", Priority: "high"},
			want: []string{"E1a1a-epic", "T4a1a-completed"},
		},
		{
			name: "blocked filter",
			opts: ListOptions{Scope: "all", Blocked: boolPtr(true)},
			want: []string{"T3a1a-blocked"},
		},
		{
			name: "blocks filter",
			opts: ListOptions{Scope: "all", Blocks: boolPtr(true)},
			want: []string{"T5a1a-blocks"},
		},
		{
			name: "owner approval filter",
			opts: ListOptions{Scope: "all", OwnerApproval: boolPtr(true)},
			want: []string{"T5a1a-blocks"},
		},
		{
			name: "path filter",
			opts: ListOptions{Scope: "all", Path: filepath.Join("tasks", "E1a1a-epic")},
			want: []string{"E1a1a-epic", "T1a1a-child"},
		},
		{
			name: "sort by id desc",
			opts: ListOptions{Scope: "all", Sort: "id", Order: "desc"},
			want: []string{"T5a1a-blocks", "T4a1a-completed", "T3a1a-blocked", "T2a1a-free", "T1a1a-child", "E1a1a-epic"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tasks, err := ListTasks(root, tc.opts)
			if err != nil {
				t.Fatalf("ListTasks failed: %v", err)
			}
			got := make([]string, 0, len(tasks))
			for _, task := range tasks {
				got = append(got, task.ID)
			}
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Fatalf("unexpected order\n got: %v\nwant: %v", got, tc.want)
			}
		})
	}
}

func TestFormatOutputs(t *testing.T) {
	root := setupListFixture(t)

	tasks, err := ListTasks(root, ListOptions{Scope: "all"})
	if err != nil {
		t.Fatalf("ListTasks failed: %v", err)
	}

	cases := []struct {
		name      string
		opts      ListOptions
		golden    string
		normalize bool
	}{
		{
			name:   "table",
			opts:   ListOptions{Format: "table"},
			golden: "testdata/list/table.txt",
		},
		{
			name:   "markdown",
			opts:   ListOptions{Format: "md"},
			golden: "testdata/list/markdown.md",
		},
		{
			name:      "json",
			opts:      ListOptions{Format: "json"},
			golden:    "testdata/list/list.json",
			normalize: true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			output, err := FormatList(tasks, tc.opts)
			if err != nil {
				t.Fatalf("FormatList failed: %v", err)
			}
			if tc.normalize {
				output = strings.ReplaceAll(output, filepath.ToSlash(root)+"/", "<ROOT>/")
			}
			assertGolden(t, tc.golden, output)
		})
	}
}

func setupListFixture(t *testing.T) string {
	t.Helper()

	root := filepath.Join(t.TempDir(), "tasks")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("mkdir tasks root: %v", err)
	}

	writeListTask(t, root, "E1a1a-epic", taskFixture{
		Role:        "architect",
		Priority:    "high",
		Parent:      "",
		Blockers:    nil,
		Blocks:      nil,
		Completed:   false,
		Owner:       false,
		DateCreated: "2026-01-01T00:00:00Z",
		DateEdited:  "2026-01-02T00:00:00Z",
		Title:       "Epic",
	})

	writeListTask(t, root, "T1a1a-child", taskFixture{
		Role:        "developer",
		Priority:    "medium",
		Parent:      "E1a1a-epic",
		Blockers:    nil,
		Blocks:      nil,
		Completed:   false,
		Owner:       false,
		DateCreated: "2026-01-03T00:00:00Z",
		DateEdited:  "2026-01-04T00:00:00Z",
		Title:       "Child Task",
		DirParent:   "E1a1a-epic",
	})

	writeListTask(t, root, "T2a1a-free", taskFixture{
		Role:        "designer",
		Priority:    "low",
		Parent:      "",
		Blockers:    nil,
		Blocks:      nil,
		Completed:   false,
		Owner:       false,
		DateCreated: "2026-01-05T00:00:00Z",
		DateEdited:  "2026-01-06T00:00:00Z",
		Title:       "Free Task",
	})

	writeListTask(t, root, "T3a1a-blocked", taskFixture{
		Role:        "developer",
		Priority:    "medium",
		Parent:      "",
		Blockers:    []string{"T9x9x-blocker"},
		Blocks:      nil,
		Completed:   false,
		Owner:       false,
		DateCreated: "2026-01-07T00:00:00Z",
		DateEdited:  "2026-01-08T00:00:00Z",
		Title:       "Blocked Task",
	})

	writeListTask(t, root, "T4a1a-completed", taskFixture{
		Role:        "developer",
		Priority:    "high",
		Parent:      "",
		Blockers:    nil,
		Blocks:      nil,
		Completed:   true,
		Owner:       false,
		DateCreated: "2026-01-09T00:00:00Z",
		DateEdited:  "2026-01-10T00:00:00Z",
		Title:       "Completed Task",
	})

	writeListTask(t, root, "T5a1a-blocks", taskFixture{
		Role:        "reviewer",
		Priority:    "medium",
		Parent:      "",
		Blockers:    nil,
		Blocks:      []string{"T2a1a-free"},
		Completed:   false,
		Owner:       true,
		DateCreated: "2026-01-11T00:00:00Z",
		DateEdited:  "2026-01-12T00:00:00Z",
		Title:       "Blocks Task",
	})

	return root
}

type taskFixture struct {
	Role        string
	Priority    string
	Parent      string
	Blockers    []string
	Blocks      []string
	Completed   bool
	Owner       bool
	DateCreated string
	DateEdited  string
	Title       string
	DirParent   string
}

func writeListTask(t *testing.T, root, id string, fixture taskFixture) {
	t.Helper()

	dir := root
	if fixture.DirParent != "" {
		dir = filepath.Join(dir, fixture.DirParent)
	}
	dir = filepath.Join(dir, id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir task dir: %v", err)
	}

	filePath := filepath.Join(dir, id+".md")
	created, err := time.Parse(time.RFC3339, fixture.DateCreated)
	if err != nil {
		t.Fatalf("parse date created: %v", err)
	}
	edited, err := time.Parse(time.RFC3339, fixture.DateEdited)
	if err != nil {
		t.Fatalf("parse date edited: %v", err)
	}

	content := strings.Join([]string{
		"---",
		"role: " + fixture.Role,
		"priority: " + fixture.Priority,
		"parent: " + fixture.Parent,
		"blockers: " + formatYAMLList(fixture.Blockers),
		"blocks: " + formatYAMLList(fixture.Blocks),
		"date_created: " + created.Format(time.RFC3339),
		"date_edited: " + edited.Format(time.RFC3339),
		"owner_approval: " + formatBool(fixture.Owner),
		"completed: " + formatBool(fixture.Completed),
		"---",
		"",
		"# " + fixture.Title,
		"",
		"Body",
		"",
	}, "\n")

	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write task file: %v", err)
	}
}

func formatYAMLList(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	return "[\"" + strings.Join(values, "\", \"") + "\"]"
}

func formatBool(value bool) string {
	if value {
		return "true"
	}
	return "false"
}

func boolPtr(value bool) *bool {
	v := value
	return &v
}

func assertGolden(t *testing.T, path, got string) {
	t.Helper()

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir golden dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden: %v", err)
	}
	if strings.TrimSpace(string(want)) != strings.TrimSpace(got) {
		t.Fatalf("golden mismatch\n--- got ---\n%s\n--- want ---\n%s", got, string(want))
	}
}
