package task

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
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

		cmd = exec.Command("git", "commit", "-m", "initial commit")
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

	commits, err := evaluateGitMetric(repoUnborn, "commits", "HEAD~0")
	if err != nil {
		t.Errorf("evaluateGitMetric for unborn HEAD commits returned an error: %v", err)
	}
	if commits != 0 {
		t.Errorf("evaluateGitMetric for unborn HEAD commits returned %d, expected 0", commits)
	}

	lines, err := evaluateGitMetric(repoUnborn, "lines_changed", "HEAD~0")
	if err != nil {
		t.Errorf("evaluateGitMetric for unborn HEAD lines_changed returned an error: %v", err)
	}
	if lines != 0 {
		t.Errorf("evaluateGitMetric for unborn HEAD lines_changed returned %d, expected 0", lines)
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
	commits, err = evaluateGitMetric(repoValid, "commits", firstCommitHash)
	if err != nil {
		t.Fatalf("evaluateGitMetric for commits returned an error: %v", err)
	}
	if diff := cmp.Diff(3, commits); diff != "" {
		t.Errorf("evaluateGitMetric for commits mismatch (-want +got):\n%s", diff)
	}

	// Test lines_changed metric
	// "HEAD~3..HEAD" should include 3 commits:
	// - second commit: file1.txt (1 addition)
	// - third commit: file2.txt (2 additions)
	// - fourth commit: file3.txt (2 additions)
	// Total additions: 1 + 2 + 2 = 5 lines
	lines, err = evaluateGitMetric(repoValid, "lines_changed", firstCommitHash)
	if err != nil {
		t.Fatalf("evaluateGitMetric for lines_changed returned an error: %v", err)
	}
	if diff := cmp.Diff(5, lines); diff != "" {
		t.Errorf("evaluateGitMetric for lines_changed mismatch (-want +got):\n%s", diff)
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
	detachedCommits, err := evaluateGitMetric(detachedRepo, "commits", initialCommitHash)
	if err != nil {
		t.Errorf("evaluateGitMetric for detached HEAD commits returned an error: %v", err)
	}
	if detachedCommits != 0 {
		t.Errorf("evaluateGitMetric for detached HEAD commits returned %d, expected 0", detachedCommits)
	}
	detachedLines, err := evaluateGitMetric(detachedRepo, "lines_changed", initialCommitHash)
	if err != nil {
		t.Errorf("evaluateGitMetric for detached HEAD lines_changed returned an error: %v", err)
	}
	if detachedLines != 0 {
		t.Errorf("evaluateGitMetric for detached HEAD lines_changed returned %d, expected 0", detachedLines)
	}

	// Test with a non-existent anchor (should now return 0, nil)
	commits, err = evaluateGitMetric(repoValid, "commits", "nonexistenthash")
	if err != nil {
		t.Errorf("evaluateGitMetric for nonexistent anchor returned an error: %v", err)
	}
	if commits != 0 {
		t.Errorf("evaluateGitMetric for nonexistent anchor returned %d, expected 0", commits)
	}

	// Test unsupported metric type
	_, err = evaluateGitMetric(repoValid, "unsupported", "HEAD~1")
	if err == nil {
		t.Errorf("evaluateGitMetric for unsupported metric type expected an error, got nil")
	} else if !strings.Contains(err.Error(), "unsupported git metric type: unsupported") {
		t.Errorf("evaluateGitMetric for unsupported metric type returned unexpected error: %v", err)
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
