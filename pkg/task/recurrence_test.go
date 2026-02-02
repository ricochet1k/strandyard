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

func TestGetCommitAtOffset(t *testing.T) {
	repo, cleanup, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup git repo: %v", err)
	}
	defer cleanup()

	// Add 5 more commits
	for i := 1; i <= 5; i++ {
		addCommit(t, repo, fmt.Sprintf("file%d.txt", i), "content", fmt.Sprintf("commit %d", i))
	}

	initialHash, err := getCommitHash(t, repo, "HEAD~5")
	if err != nil {
		t.Fatalf("failed to get anchor hash: %v", err)
	}

	// Test offset 3
	hash3, err := GetCommitAtOffset(repo, initialHash, 3)
	if err != nil {
		t.Fatalf("GetCommitAtOffset failed: %v", err)
	}
	expectedHash3, _ := getCommitHash(t, repo, "HEAD~2")
	if hash3 != expectedHash3 {
		t.Errorf("GetCommitAtOffset(3) = %s, want %s", hash3, expectedHash3)
	}

	// Test offset 6 (exceeds)
	_, err = GetCommitAtOffset(repo, initialHash, 6)
	if err == nil {
		t.Errorf("GetCommitAtOffset(6) expected error, got nil")
	}
}

func TestGetCommitCrossingLinesThreshold(t *testing.T) {
	repo, cleanup, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup git repo: %v", err)
	}
	defer cleanup()

	// Add commits with known line changes
	addCommit(t, repo, "f1.txt", "1\n2\n3", "c1") // 3 lines
	addCommit(t, repo, "f2.txt", "4\n5", "c2")    // 2 lines (total 5)
	addCommit(t, repo, "f3.txt", "6", "c3")       // 1 line (total 6)

	initialHash, err := getCommitHash(t, repo, "HEAD~3")
	if err != nil {
		t.Fatalf("failed to get anchor hash: %v", err)
	}

	// Test threshold 4 (should be c2)
	hash, err := GetCommitCrossingLinesThreshold(repo, initialHash, 4)
	if err != nil {
		t.Fatalf("GetCommitCrossingLinesThreshold failed: %v", err)
	}
	expectedHash, _ := getCommitHash(t, repo, "HEAD~1")
	if hash != expectedHash {
		t.Errorf("GetCommitCrossingLinesThreshold(4) = %s, want %s", hash, expectedHash)
	}

	// Test threshold 10 (not reached)
	_, err = GetCommitCrossingLinesThreshold(repo, initialHash, 10)
	if err == nil {
		t.Errorf("GetCommitCrossingLinesThreshold(10) expected error, got nil")
	}
}

func TestUpdateAnchor(t *testing.T) {
	t.Run("time-based no-drift", func(t *testing.T) {
		anchor := "2026-01-01T00:00:00Z"
		interval := 7 // days
		// Mock time.Now is tricky in Go without libraries, but we can check the logic
		newAnchor, err := UpdateAnchor("", "", "days", anchor, interval)
		if err != nil {
			t.Fatalf("UpdateAnchor failed: %v", err)
		}
		// Since 2026-01-01 is far in the past, it should have skipped ahead to now + interval
		// We just check it's a valid date and later than anchor
		parsed, err := time.Parse(time.RFC3339, newAnchor)
		if err != nil {
			t.Fatalf("failed to parse new anchor: %v", err)
		}
		if !parsed.After(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)) {
			t.Errorf("new anchor %v is not after old anchor", parsed)
		}
	})

	t.Run("commits", func(t *testing.T) {
		repo, cleanup, err := setupGitRepo(t, true)
		if err != nil {
			t.Fatalf("failed to setup git repo: %v", err)
		}
		defer cleanup()

		for i := 1; i <= 3; i++ {
			addCommit(t, repo, fmt.Sprintf("f%d.txt", i), "c", "msg")
		}

		anchor, _ := getCommitHash(t, repo, "HEAD~3")
		newAnchor, err := UpdateAnchor(repo, "", "commits", anchor, 2)
		if err != nil {
			t.Fatalf("UpdateAnchor failed: %v", err)
		}
		expected, _ := getCommitHash(t, repo, "HEAD~1")
		if newAnchor != expected {
			t.Errorf("UpdateAnchor(commits) = %s, want %s", newAnchor, expected)
		}
	})

	t.Run("commits-missing-fallback", func(t *testing.T) {
		repo, cleanup, err := setupGitRepo(t, true)
		if err != nil {
			t.Fatalf("failed to setup git repo: %v", err)
		}
		defer cleanup()

		newAnchor, err := UpdateAnchor(repo, "", "commits", "0123456789abcdef0123456789abcdef01234567", 10)
		if err != nil {
			t.Fatalf("UpdateAnchor failed: %v", err)
		}
		head, _ := getCommitHash(t, repo, "HEAD")
		if newAnchor != head {
			t.Errorf("UpdateAnchor(missing) = %s, want HEAD %s", newAnchor, head)
		}
	})
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

	t.Run("Reuses provided log", func(t *testing.T) {
		log, err := activity.Open(tmpDir)
		if err != nil {
			t.Fatalf("failed to open activity log: %v", err)
		}
		// We don't defer log.Close() here because we want to verify it's still open

		_, err = EvaluateTasksCompletedMetric(tmpDir, "now", "T-test", log)
		if err != nil {
			t.Errorf("EvaluateTasksCompletedMetric failed: %v", err)
		}

		// Verify we can still write to the log
		err = log.WriteRecurrenceAnchorResolution("T-test", "now", "resolved")
		if err != nil {
			t.Errorf("log was prematurely closed: %v", err)
		}
		log.Close()
	})
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

func TestEvaluateTasksCompletedMetricWithTaskIDAnchors(t *testing.T) {
	// Create a temporary directory for the activity log
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup activity log with task completions
	log, err := activity.Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open activity log: %v", err)
	}

	now := time.Now().UTC()
	t1 := now.Add(-3 * time.Hour)
	t2 := now.Add(-2 * time.Hour)
	t3 := now.Add(-1 * time.Hour)

	// Write completion entries for specific tasks
	log.WriteEntry(activity.Entry{Timestamp: t1, TaskID: "T1111-task1", Type: activity.EventTaskCompleted, Report: "first"})
	log.WriteEntry(activity.Entry{Timestamp: t2, TaskID: "T2222-task2", Type: activity.EventTaskCompleted, Report: "second"})
	log.WriteEntry(activity.Entry{Timestamp: t3, TaskID: "T3333-task3", Type: activity.EventTaskCompleted, Report: "third"})
	log.WriteEntry(activity.Entry{Timestamp: now.Add(-30 * time.Minute), TaskID: "T1111-task1", Type: activity.EventTaskCompleted, Report: "task1 again"})

	log.Close()

	// Test 1: Resolve using task ID anchor (T1111-task1)
	// Should use the latest completion time of T1111-task1 (30 minutes ago)
	// and count completions since then (inclusive), which should be 1 (itself)
	count, err := EvaluateTasksCompletedMetric(tmpDir, "T1111-task1", "", nil)
	if err != nil {
		t.Errorf("EvaluateTasksCompletedMetric with task ID anchor failed: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 completion since T1111-task1 last completion, got %d", count)
	}

	// Test 2: Resolve using task ID anchor (T2222-task2)
	// Should use the completion time of T2222-task2 (2 hours ago)
	// and count completions since then (inclusive), which should be 3:
	// T2222-task2 (itself), T3333-task3, and T1111-task1 again
	count, err = EvaluateTasksCompletedMetric(tmpDir, "T2222-task2", "", nil)
	if err != nil {
		t.Errorf("EvaluateTasksCompletedMetric with task ID anchor T2222-task2 failed: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 completions since T2222-task2 completion, got %d", count)
	}

	// Test 3: Try with non-existent task ID (should fail)
	_, err = EvaluateTasksCompletedMetric(tmpDir, "T9999-nonexistent", "", nil)
	if err == nil {
		t.Errorf("expected error for non-existent task ID, got nil")
	}

	// Test 4: Task ID that looks like a date should still resolve as task ID if it exists in log
	// Create a task ID that starts with T but looks date-ish
	tmpDir2, err := os.MkdirTemp("", "activity-test-2-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir2)

	log2, err := activity.Open(tmpDir2)
	if err != nil {
		t.Fatalf("failed to open activity log: %v", err)
	}

	anchorTime := now.Add(-4 * time.Hour)
	log2.WriteEntry(activity.Entry{Timestamp: anchorTime, TaskID: "T4444-anchor", Type: activity.EventTaskCompleted})
	log2.WriteEntry(activity.Entry{Timestamp: now.Add(-2 * time.Hour), TaskID: "T5555-later", Type: activity.EventTaskCompleted})
	log2.Close()

	count, err = EvaluateTasksCompletedMetric(tmpDir2, "T4444-anchor", "", nil)
	if err != nil {
		t.Errorf("failed to resolve task ID anchor: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 completions since T4444-anchor (itself + T5555-later), got %d", count)
	}
}

func TestUpdateAnchorWithTaskIDAnchors(t *testing.T) {
	// Create a temporary directory for the activity log
	tmpDir, err := os.MkdirTemp("", "activity-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup activity log with task completions
	log, err := activity.Open(tmpDir)
	if err != nil {
		t.Fatalf("failed to open activity log: %v", err)
	}

	// Use a fixed base time for deterministic testing
	baseTime := time.Date(2026, 2, 1, 12, 0, 0, 0, time.UTC)
	// Write 5 completion entries at 10-minute intervals
	for i := 0; i < 5; i++ {
		ts := baseTime.Add(-time.Duration((5-i)*10) * time.Minute)
		log.WriteEntry(activity.Entry{
			Timestamp: ts,
			TaskID:    fmt.Sprintf("T%04d-task%d", 1000+i, i),
			Type:      activity.EventTaskCompleted,
		})
	}
	log.Close()

	// Test 1: UpdateAnchor with task ID anchor
	// T1000-task0 completed at baseTime - 50 minutes
	// We want to find the 2nd completion from there
	// Completions are at: -50, -40, -30, -20, -10 minutes
	// So the 2nd completion is at -40 minutes
	newAnchor, err := UpdateAnchor(tmpDir, tmpDir, "tasks_completed", "T1000-task0", 2)
	if err != nil {
		t.Errorf("UpdateAnchor failed: %v", err)
	}

	// Verify newAnchor is a valid RFC3339 timestamp
	newTime, err := time.Parse(time.RFC3339, newAnchor)
	if err != nil {
		t.Errorf("new anchor is not a valid RFC3339 timestamp: %s (%v)", newAnchor, err)
	}

	// The 2nd completion from T1000-task0 (at -50) is at -40 minutes
	expectedTime := baseTime.Add(-40 * time.Minute)
	timeDiff := newTime.Sub(expectedTime).Abs()
	if timeDiff > 1*time.Second {
		t.Errorf("new anchor time %v is not equal to expected %v (diff: %v)", newTime, expectedTime, timeDiff)
	}

	// Test 2: UpdateAnchor with date anchor (not task ID)
	// Use the timestamp of T1002-task2 as a date anchor (-30 minutes)
	anchor2Time := baseTime.Add(-30 * time.Minute)
	dateAnchor := anchor2Time.Format(time.RFC3339)
	newAnchor, err = UpdateAnchor(tmpDir, tmpDir, "tasks_completed", dateAnchor, 2)
	if err != nil {
		t.Errorf("UpdateAnchor with date anchor failed: %v", err)
	}

	newTime, err = time.Parse(time.RFC3339, newAnchor)
	if err != nil {
		t.Errorf("new anchor is not a valid RFC3339 timestamp: %s (%v)", newAnchor, err)
	}

	// The 2nd completion from -30 minutes is at -20 minutes
	// (1st is at -30 itself, 2nd is at -20)
	expectedTime = baseTime.Add(-20 * time.Minute)
	timeDiff = newTime.Sub(expectedTime).Abs()
	if timeDiff > 1*time.Second {
		t.Errorf("new anchor time %v is not equal to expected %v (diff: %v)", newTime, expectedTime, timeDiff)
	}
}
