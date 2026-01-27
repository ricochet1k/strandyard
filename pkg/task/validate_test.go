package task

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateMasterLists_FreeTasksPrioritySections(t *testing.T) {
	tmp := t.TempDir()
	rootsFile := filepath.Join(tmp, "root-tasks.md")
	freeFile := filepath.Join(tmp, "free-tasks.md")

	tasks := map[string]*Task{
		"T1aaa-high": {
			ID:       "T1aaa-high",
			FilePath: "tasks/T1aaa-high/T1aaa-high.md",
			Content:  "# High Task\n",
			Meta: Metadata{
				Priority: PriorityHigh,
			},
		},
		"T2bbb-default": {
			ID:       "T2bbb-default",
			FilePath: "tasks/T2bbb-default/T2bbb-default.md",
			Content:  "# Default Task\n",
			Meta:     Metadata{},
		},
		"T3ccc-low": {
			ID:       "T3ccc-low",
			FilePath: "tasks/T3ccc-low/T3ccc-low.md",
			Content:  "# Low Task\n",
			Meta: Metadata{
				Priority: PriorityLow,
			},
		},
	}

	if err := GenerateMasterLists(tasks, "tasks", rootsFile, freeFile); err != nil {
		t.Fatalf("GenerateMasterLists failed: %v", err)
	}

	got, err := os.ReadFile(freeFile)
	if err != nil {
		t.Fatalf("read free list: %v", err)
	}

	want := strings.Join([]string{
		"# Free tasks",
		"",
		"## High",
		"",
		"- [High Task](tasks/T1aaa-high/T1aaa-high.md)",
		"",
		"## Medium",
		"",
		"- [Default Task](tasks/T2bbb-default/T2bbb-default.md)",
		"",
		"## Low",
		"",
		"- [Low Task](tasks/T3ccc-low/T3ccc-low.md)",
		"",
		"",
	}, "\n")

	if string(got) != want {
		t.Fatalf("unexpected free list:\n--- got ---\n%s\n--- want ---\n%s", string(got), want)
	}
}
