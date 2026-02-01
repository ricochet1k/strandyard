package task

import (
	"testing"
)

func TestEvaluateGitMetricFlagInjection(t *testing.T) {
	repo, cleanup, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup git repo: %v", err)
	}
	defer cleanup()

	// A malicious anchor that looks like a flag.
	// We use --version because it's safe and predictable.
	maliciousAnchor := "--version"

	// If not hardened, this might execute 'git rev-list --count --version..HEAD'
	// which might be interpreted differently depending on git version,
	// but the goal is to ensure it doesn't return success or a version string.
	// We expect it to return 0 and no error (because of the error handling in EvaluateGitMetric
	// that treats unknown revisions as no-op).
	// However, if it WAS treated as a flag, it might return a version string which would fail parsing.

	count, err := EvaluateGitMetric(repo, "commits", maliciousAnchor, "", nil)
	if err != nil {
		t.Errorf("expected no error (no-op), got: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0 for malicious anchor, got %d", count)
	}

	// Test with an anchor that would definitely cause a usage error if interpreted as a flag.
	maliciousAnchor2 := "--invalid-flag"
	count, err = EvaluateGitMetric(repo, "commits", maliciousAnchor2, "", nil)
	if err != nil {
		t.Errorf("expected no error (no-op), got: %v", err)
	}
	if count != 0 {
		t.Errorf("expected count 0 for malicious anchor, got %d", count)
	}
}

func TestResolveGitHashFlagInjection(t *testing.T) {
	repo, cleanup, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup git repo: %v", err)
	}
	defer cleanup()

	maliciousAnchor := "--version"
	_, err = ResolveGitHash(repo, maliciousAnchor)
	if err == nil {
		t.Fatal("expected error for flag-like anchor in ResolveGitHash, got nil")
	}

	// If it was treated as a flag, it might have returned version info.
	// If it was treated as a revision, it should say it's not a valid revision.
}

func TestGetCommitAtOffsetFlagInjection(t *testing.T) {
	repo, cleanup, err := setupGitRepo(t, true)
	if err != nil {
		t.Fatalf("failed to setup git repo: %v", err)
	}
	defer cleanup()

	maliciousAnchor := "--version"
	_, err = GetCommitAtOffset(repo, maliciousAnchor, 1)
	if err == nil {
		t.Fatal("expected error for flag-like anchor in GetCommitAtOffset, got nil")
	}
}
