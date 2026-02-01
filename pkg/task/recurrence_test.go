package task

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ricochet1k/strandyard/pkg/activity"
)

// setupGitRepo initializes a git repository in a temporary directory
// and returns the path to the repository and a cleanup function.
// It can optionally initialize with an "unborn" HEAD or with an initial commit.
func setupGitRepo(t *testing.T, withInitialCommit bool) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "git-test-repo-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("failed to git init: %w", err)
	}

	if withInitialCommit {
		// Create initial commit
		filePath := filepath.Join(tmpDir, "README.md")
		err := os.WriteFile(filePath, []byte("Hello World"), 0o644)
		if err != nil {
			cleanup()
			return "", nil, fmt.Errorf("failed to write README.md: %w", err)
		}

		cmd = exec.Command("git", "add", "README.md")
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("failed to git add: %w", err)
		}

		cmd = exec.Command("git", "-c", "user.name=Test User", "-c", "user.email=test@example.com", "commit", "-m", "initial commit")
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			cleanup()
			return "", nil, fmt.Errorf("failed to git commit: %w", err)
		}
	}

	return tmpDir, cleanup, nil
}

func TestIsHeadValid(t *testing.T) {
	// Test case 1: Unborn HEAD (no commits)
	repoUnborn, cleanupUnborn, err := setupGitRepo(t, false)
	if err != nil {
		t.Fatalf("failed to setup unborn HEAD repo: %v", err)
	}
	defer cleanupUnborn()

	if isHeadValid(repoUnborn) {
		t.Errorf("isHeadValid returned true for unborn HEAD, expected false")
	}

	// Test case 2: Valid HEAD (initial commit)
	repoValid, cleanupValid, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup valid HEAD repo: %v", err)
	}
	defer cleanupValid()

	if !isHeadValid(repoValid) {
		t.Errorf("isHeadValid returned false for valid HEAD, expected true")
	}

	// Test case 3: Detached HEAD (pointing to a commit)
	// This is effectively covered by Test case 2, as detached HEAD still points to a valid commit
	// and `git rev-parse --verify HEAD` would return true.
}

func TestEvaluateGitMetric(t *testing.T) {
	// Test case 1: Unborn HEAD
	repoUnborn, cleanupUnborn, err := setupGitRepo(t, false)
	if err != nil {
		t.Fatalf("failed to setup unborn HEAD repo: %v", err)
	}
	defer cleanupUnborn()

	commits, err := EvaluateGitMetric(repoUnborn, "commits", "HEAD~0", "", nil)
	if err != nil {
		t.Errorf("EvaluateGitMetric for unborn HEAD commits returned an error: %v", err)
	}
	if commits != 0 {
		t.Errorf("EvaluateGitMetric for unborn HEAD commits returned %d, expected 0", commits)
	}

	lines, err := EvaluateGitMetric(repoUnborn, "lines_changed", "HEAD~0", "", nil)
	if err != nil {
		t.Errorf("EvaluateGitMetric for unborn HEAD lines_changed returned an error: %v", err)
	}
	if lines != 0 {
		t.Errorf("EvaluateGitMetric for unborn HEAD lines_changed returned %d, expected 0", lines)
	}

	// Test case 2: Valid HEAD with a few commits
	repoValid, cleanupValid, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup valid HEAD repo: %v", err)
	}
	defer cleanupValid()

	// Add more commits
	addCommit(t, repoValid, "file1.txt", "content1", "second commit")
	addCommit(t, repoValid, "file2.txt", "content2\nline2", "third commit")
	addCommit(t, repoValid, "file3.txt", "content3\nanotherline", "fourth commit") // 2 lines

	// Get first commit hash for anchor
	firstCommitHash, err := getCommitHash(t, repoValid, "HEAD~3")
	if err != nil {
		t.Fatalf("failed to get first commit hash: %v", err)
	}

	// Test commits metric
	commits, err = EvaluateGitMetric(repoValid, "commits", firstCommitHash, "", nil)
	if err != nil {
		t.Fatalf("EvaluateGitMetric for commits returned an error: %v", err)
	}
	if diff := cmp.Diff(3, commits); diff != "" {
		t.Errorf("EvaluateGitMetric for commits mismatch (-want +got):\n%s", diff)
	}

	// Test lines_changed metric
	// "HEAD~3..HEAD" should include 3 commits:
	// - second commit: file1.txt (1 addition)
	// - third commit: file2.txt (2 additions)
	// - fourth commit: file3.txt (2 additions)
	// Total additions: 1 + 2 + 2 = 5 lines
	lines, err = EvaluateGitMetric(repoValid, "lines_changed", firstCommitHash, "", nil)
	if err != nil {
		t.Fatalf("EvaluateGitMetric for lines_changed returned an error: %v", err)
	}
	if diff := cmp.Diff(5, lines); diff != "" {
		t.Errorf("EvaluateGitMetric for lines_changed mismatch (-want +got):\n%s", diff)
	}

	// Test case 3: Detached HEAD (pointing to a commit)
	// Create a detached HEAD in repoValid by checking out an old commit
	detachedRepo, cleanupDetached, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup detached HEAD repo: %v", err)
	}
	defer cleanupDetached()

	addCommit(t, detachedRepo, "fileD1.txt", "contentD1", "detached commit 1")
	addCommit(t, detachedRepo, "fileD2.txt", "contentD2", "detached commit 2")

	// Get the hash of the first commit (initial commit)
	initialCommitHash, err := getCommitHash(t, detachedRepo, "HEAD~2")
	if err != nil {
		t.Fatalf("failed to get initial commit hash: %v", err)
	}

	// Detach HEAD to the initial commit
	cmd := exec.Command("git", "checkout", initialCommitHash)
	cmd.Dir = detachedRepo
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to detach HEAD: %v", err)
	}

	// Now evaluate metric with detached HEAD
	detachedCommits, err := EvaluateGitMetric(detachedRepo, "commits", initialCommitHash, "", nil)
	if err != nil {
		t.Errorf("EvaluateGitMetric for detached HEAD commits returned an error: %v", err)
	}
	if detachedCommits != 0 {
		t.Errorf("EvaluateGitMetric for detached HEAD commits returned %d, expected 0", detachedCommits)
	}
	detachedLines, err := EvaluateGitMetric(detachedRepo, "lines_changed", initialCommitHash, "", nil)
	if err != nil {
		t.Errorf("EvaluateGitMetric for detached HEAD lines_changed returned an error: %v", err)
	}
	if detachedLines != 0 {
		t.Errorf("EvaluateGitMetric for detached HEAD lines_changed returned %d, expected 0", detachedLines)
	}

	// Test with a non-existent anchor (should now return 0, nil)
	commits, err = EvaluateGitMetric(repoValid, "commits", "nonexistenthash", "", nil)
	if err != nil {
		t.Errorf("EvaluateGitMetric for nonexistent anchor returned an error: %v", err)
	}
	if commits != 0 {
		t.Errorf("EvaluateGitMetric for nonexistent anchor returned %d, expected 0", commits)
	}

	// Test unsupported metric type
	_, err = EvaluateGitMetric(repoValid, "unsupported", "HEAD~1", "", nil)
	if err == nil {
		t.Errorf("EvaluateGitMetric for unsupported metric type expected an error, got nil")
	} else if !strings.Contains(err.Error(), "unsupported git metric type: unsupported") {
		t.Errorf("EvaluateGitMetric for unsupported metric type returned unexpected error: %v", err)
	}
}

// Helper to add a commit to an existing git repository
func addCommit(t *testing.T, repoPath, fileName, content, commitMsg string) {
	filePath := filepath.Join(repoPath, fileName)
	err := os.WriteFile(filePath, []byte(content), 0o644)
	if err != nil {
		t.Fatalf("failed to write file %s: %v", fileName, err)
	}

	cmd := exec.Command("git", "add", fileName)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git add %s: %v", fileName, err)
	}

	// Ensure consistent committer info for reproducible tests
	cmd = exec.Command("git", "-c", "user.name=Test User", "-c", "user.email=test@example.com", "commit", "-m", commitMsg)
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to git commit: %v", err)
	}
}

// Helper to get commit hash for a given revision
func getCommitHash(t *testing.T, repoPath, revision string) (string, error) {
	cmd := exec.Command("git", "rev-parse", revision)
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get commit hash for revision %q: %w", revision, err)
	}
	return strings.TrimSpace(string(output)), nil
}

func TestEvaluateTasksCompletedMetric(t *testing.T) {
	// Create a temporary directory for the activity log
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test case 1: Empty activity log
	count, err := EvaluateTasksCompletedMetric(tmpDir, "Jan 28 2026 09:00 UTC", "", nil)
	if err != nil {
		t.Errorf("EvaluateTasksCompletedMetric for empty log returned an error: %v", err)
	}
	if count != 0 {
		t.Errorf("EvaluateTasksCompletedMetric for empty log returned %d, expected 0", count)
	}

	// Test case 2: Some completions
	// We need to directly write to the activity log for testing
	// The activity log is at tmpDir/activity.log
	activityLog, err := os.OpenFile(filepath.Join(tmpDir, "activity.log"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		t.Fatalf("failed to create activity log: %v", err)
	}
	defer activityLog.Close()

	// Write some completion entries with specific timestamps
	now := time.Now().UTC()
	entries := []string{
		fmt.Sprintf(`{"timestamp":"%s","task_id":"T1","type":"task_completed","report":"first"}`, now.Add(-48*time.Hour).Format(time.RFC3339)),
		fmt.Sprintf(`{"timestamp":"%s","task_id":"T2","type":"task_completed","report":"second"}`, now.Add(-24*time.Hour).Format(time.RFC3339)),
		fmt.Sprintf(`{"timestamp":"%s","task_id":"T3","type":"task_completed","report":"third"}`, now.Add(-2*time.Hour).Format(time.RFC3339)),
	}
	for _, entry := range entries {
		if _, err := activityLog.WriteString(entry + "\n"); err != nil {
			t.Fatalf("failed to write to activity log: %v", err)
		}
	}
	activityLog.Close()

	// Test with anchor 3 hours ago
	anchorTime := now.Add(-3 * time.Hour)
	anchorStr := anchorTime.Format("Jan 2 2006 15:04 MST")
	count, err = EvaluateTasksCompletedMetric(tmpDir, anchorStr, "", nil)
	if err != nil {
		t.Errorf("EvaluateTasksCompletedMetric returned an error: %v", err)
	}
	if count != 1 {
		t.Errorf("EvaluateTasksCompletedMetric returned %d, expected 1", count)
	}

	// Test with anchor 49 hours ago (all 3 completions)
	anchorTime = now.Add(-49 * time.Hour)
	anchorStr = anchorTime.Format("Jan 2 2006 15:04 MST")
	count, err = EvaluateTasksCompletedMetric(tmpDir, anchorStr, "", nil)
	if err != nil {
		t.Errorf("EvaluateTasksCompletedMetric returned an error: %v", err)
	}
	if count != 3 {
		t.Errorf("EvaluateTasksCompletedMetric returned %d, expected 3", count)
	}

	// Test case 3: Invalid date format
	_, err = EvaluateTasksCompletedMetric(tmpDir, "invalid date", "", nil)
	if err == nil {
		t.Errorf("EvaluateTasksCompletedMetric for invalid date expected an error, got nil")
	}
}

func TestRecurrenceAnchorResolutionLogging(t *testing.T) {
	// Setup git repo
	repo, cleanup, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup git repo: %v", err)
	}
	defer cleanup()

	// Setup activity log
	tmpDir, err := os.MkdirTemp("", "activity-log-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	log, err := activity.Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open activity log: %v", err)
	}
	defer log.Close()

	taskID := "T1234-test"

	// Test git HEAD resolution logging
	_, err = EvaluateGitMetric(repo, "commits", "HEAD", taskID, log)
	if err != nil {
		t.Fatalf("EvaluateGitMetric failed: %v", err)
	}

	// Test tasks_completed "now" resolution logging
	_, err = EvaluateTasksCompletedMetric(tmpDir, "now", taskID, log)
	if err != nil {
		t.Fatalf("EvaluateTasksCompletedMetric failed: %v", err)
	}

	// Verify entries in log
	entries, err := log.ReadEntries()
	if err != nil {
		t.Fatalf("failed to read entries: %v", err)
	}

	var gitResolved, timeResolved bool
	for _, entry := range entries {
		if entry.Type == activity.EventRecurrenceAnchorResolved {
			if entry.Metadata["original"] == "HEAD" {
				gitResolved = true
				if len(entry.Metadata["resolved"]) != 40 {
					t.Errorf("expected 40-char commit hash, got %q", entry.Metadata["resolved"])
				}
			}
			if entry.Metadata["original"] == "now" {
				timeResolved = true
				if entry.Metadata["resolved"] == "" {
					t.Errorf("expected resolved timestamp, got empty string")
				}
			}
		}
	}

	if !gitResolved {
		t.Errorf("expected git resolution entry not found")
	}
	if !timeResolved {
		t.Errorf("expected time resolution entry not found")
	}
}

func TestValidateAnchor(t *testing.T) {
	// Setup git repo
	repo, cleanup, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup git repo: %v", err)
	}
	defer cleanup()

	headHash, err := getCommitHash(t, repo, "HEAD")
	if err != nil {
		t.Fatalf("failed to get head hash: %v", err)
	}

	// Setup tasks map
	tasks := map[string]*Task{
		"T1234-existing": {ID: "T1234-existing"},
	}

	tests := []struct {
		name    string
		metric  string
		anchor  string
		isValid bool
	}{
		{"valid date ISO", "days", "2026-01-28T09:00:00Z", true},
		{"valid date human", "days", "Jan 28 2026 09:00 UTC", true},
		{"invalid date", "days", "invalid-date", false},
		{"valid commit HEAD", "commits", "HEAD", true},
		{"valid commit hash", "commits", headHash, true},
		{"invalid commit hash", "commits", "0123456789abcdef0123456789abcdef01234567", false},
		{"valid task anchor", "tasks_completed", "T1234-existing", true},
		{"invalid task anchor", "tasks_completed", "T9999-missing", false},
		{"valid tasks_completed date", "tasks_completed", "Jan 28 2026 09:00 UTC", true},
		{"valid tasks_completed now", "tasks_completed", "now", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAnchor(tt.metric, tt.anchor, repo, tasks)
			if tt.isValid {
				if err != nil {
					t.Errorf("ValidateAnchor(%s, %s) expected no error but got %v", tt.metric, tt.anchor, err)
				}
			} else {
				if err == nil {
					t.Errorf("ValidateAnchor(%s, %s) expected error but got nil", tt.metric, tt.anchor)
				}
			}
		})
	}
}
