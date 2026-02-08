package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ricochet1k/strandyard/pkg/task"
)

func TestClaimMarksTaskInProgress(t *testing.T) {
	paths := setupTestProject(t, initOptions{ProjectName: "", StorageMode: storageLocal})
	roleName := testRoleName(t, "claim-cmd")
	if err := os.WriteFile(filepath.Join(paths.RolesDir, roleName+".md"), []byte("# "+roleName+"\n"), 0o644); err != nil {
		t.Fatalf("write role file: %v", err)
	}

	taskID := "T9a1a-claim-me"
	writeClaimTaskFile(t, paths.TasksDir, taskID, roleName)

	var out bytes.Buffer
	if err := runClaim(&out, taskID); err != nil {
		t.Fatalf("runClaim failed: %v", err)
	}
	if !strings.Contains(out.String(), "status set to in_progress") {
		t.Fatalf("expected success message, got: %s", out.String())
	}

	db := task.NewTaskDB(paths.TasksDir)
	if err := db.LoadAllIfEmpty(); err != nil {
		t.Fatalf("load tasks: %v", err)
	}
	tk, err := db.Get(taskID)
	if err != nil {
		t.Fatalf("get task: %v", err)
	}
	if tk.Meta.Status != task.StatusInProgress {
		t.Fatalf("expected status %q, got %q", task.StatusInProgress, tk.Meta.Status)
	}
}

func writeClaimTaskFile(t *testing.T, tasksDir, id, roleName string) {
	t.Helper()
	now := time.Date(2026, 2, 7, 12, 0, 0, 0, time.UTC).Format(time.RFC3339)
	content := strings.Join([]string{
		"---",
		"role: " + roleName,
		"priority: medium",
		"parent: \"\"",
		"blockers: []",
		"blocks: []",
		"date_created: " + now,
		"date_edited: " + now,
		"completed: false",
		"status: open",
		"---",
		"",
		"# " + id,
		"",
		"Task body.",
		"",
	}, "\n")

	if err := os.WriteFile(filepath.Join(tasksDir, id+".md"), []byte(content), 0o644); err != nil {
		t.Fatalf("write task file: %v", err)
	}
}
