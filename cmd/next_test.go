package cmd

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ricochet1k/strandyard/pkg/task"
)

func TestNextClaimSkipsClaimedTask(t *testing.T) {
	paths := setupTestProject(t, initOptions{ProjectName: "", StorageMode: storageLocal})
	roleName := testRoleName(t, "claim")
	writeRoleFile(t, filepath.Join(paths.RolesDir, roleName+".md"), roleName)

	now := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)
	taskA := "T1a1a-alpha"
	taskB := "T2a1a-beta"
	writeNextTaskFile(t, paths.TasksDir, taskA, roleName, task.StatusOpen, now)
	writeNextTaskFile(t, paths.TasksDir, taskB, roleName, task.StatusOpen, now)

	if err := runRepair(io.Discard, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
		t.Fatalf("runRepair failed: %v", err)
	}

	var claimedOutput bytes.Buffer
	if err := runNextWithOptions(&claimedOutput, "", "", nextOptions{
		Claim:        true,
		ClaimTimeout: time.Hour,
		Now:          func() time.Time { return now },
	}); err != nil {
		t.Fatalf("runNextWithOptions claim failed: %v", err)
	}
	if !strings.Contains(claimedOutput.String(), taskA) {
		t.Fatalf("expected claimed output to include %s, got: %s", taskA, claimedOutput.String())
	}

	db := task.NewTaskDB(paths.TasksDir)
	if err := db.LoadAllIfEmpty(); err != nil {
		t.Fatalf("failed to load tasks: %v", err)
	}
	claimedTask, err := db.Get(taskA)
	if err != nil {
		t.Fatalf("failed to get claimed task: %v", err)
	}
	if claimedTask.Meta.Status != task.StatusInProgress {
		t.Fatalf("expected claimed task status %q, got %q", task.StatusInProgress, claimedTask.Meta.Status)
	}

	var nextOutput bytes.Buffer
	if err := runNextWithOptions(&nextOutput, "", "", nextOptions{
		ClaimTimeout: time.Hour,
		Now:          func() time.Time { return now.Add(5 * time.Minute) },
	}); err != nil {
		t.Fatalf("runNextWithOptions second call failed: %v", err)
	}
	if !strings.Contains(nextOutput.String(), taskB) {
		t.Fatalf("expected second output to include %s, got: %s", taskB, nextOutput.String())
	}
	if strings.Contains(nextOutput.String(), taskA) {
		t.Fatalf("expected second output to exclude claimed task %s, got: %s", taskA, nextOutput.String())
	}
}

func TestNextReopensExpiredInProgressTask(t *testing.T) {
	paths := setupTestProject(t, initOptions{ProjectName: "", StorageMode: storageLocal})
	roleName := testRoleName(t, "expire")
	writeRoleFile(t, filepath.Join(paths.RolesDir, roleName+".md"), roleName)

	now := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)
	taskID := "T3a1a-expired"
	writeNextTaskFile(t, paths.TasksDir, taskID, roleName, task.StatusInProgress, now.Add(-2*time.Hour))

	if err := runRepair(io.Discard, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
		t.Fatalf("runRepair failed: %v", err)
	}

	var output bytes.Buffer
	if err := runNextWithOptions(&output, "", "", nextOptions{
		ClaimTimeout: time.Hour,
		Now:          func() time.Time { return now },
	}); err != nil {
		t.Fatalf("runNextWithOptions failed: %v", err)
	}
	if !strings.Contains(output.String(), taskID) {
		t.Fatalf("expected output to include reopened task %s, got: %s", taskID, output.String())
	}

	db := task.NewTaskDB(paths.TasksDir)
	if err := db.LoadAllIfEmpty(); err != nil {
		t.Fatalf("failed to load tasks: %v", err)
	}
	reopenedTask, err := db.Get(taskID)
	if err != nil {
		t.Fatalf("failed to get task: %v", err)
	}
	if reopenedTask.Meta.Status != task.StatusOpen {
		t.Fatalf("expected reopened task status %q, got %q", task.StatusOpen, reopenedTask.Meta.Status)
	}
}

func writeRoleFile(t *testing.T, path, roleName string) {
	t.Helper()
	content := "# " + roleName + "\n\nrole description\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write role file: %v", err)
	}
}

func writeNextTaskFile(t *testing.T, tasksDir, id, roleName, status string, edited time.Time) {
	t.Helper()
	content := "---\n" +
		"role: " + roleName + "\n" +
		"priority: medium\n" +
		"parent: \"\"\n" +
		"blockers: []\n" +
		"blocks: []\n" +
		"date_created: " + edited.Format(time.RFC3339) + "\n" +
		"date_edited: " + edited.Format(time.RFC3339) + "\n" +
		"completed: false\n"
	if status != "" {
		content += "status: " + status + "\n"
	}
	content += "---\n\n# " + id + "\n\nTask body.\n"

	if err := os.WriteFile(filepath.Join(tasksDir, id+".md"), []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write task file: %v", err)
	}
}
