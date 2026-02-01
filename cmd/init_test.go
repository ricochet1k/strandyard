package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitWithMaliciousPreset(t *testing.T) {
	_, _ = setupTestEnv(t)

	// A malicious preset starting with a hyphen.
	// If the -- separator is missing, git clone might treat this as a flag.
	// We use an invalid flag that git clone doesn't recognize to see how it's handled,
	// or one that would have side effects.
	// Since we want to verify it's treated as a repository path, we expect it to fail
	// because that "path" doesn't exist, but it shouldn't be parsed as a flag.

	maliciousPreset := "--invalid-git-flag"

	var buf bytes.Buffer
	opts := initOptions{
		StorageMode: storageLocal,
		Preset:      maliciousPreset,
	}

	err := runInit(&buf, opts)
	if err == nil {
		t.Fatal("expected error when cloning malicious preset, but got nil")
	}

	// If it was treated as a flag, git might report "unknown option: --invalid-git-flag"
	// If it was treated as a repo, git might report "repository '--invalid-git-flag' does not exist"
	errStr := err.Error()
	if strings.Contains(errStr, "unknown option") || strings.Contains(errStr, "usage:") {
		t.Errorf("git clone seems to have treated preset as a flag: %s", errStr)
	}

	if !strings.Contains(errStr, "repository '--invalid-git-flag' does not exist") &&
		!strings.Contains(errStr, "does not exist") {
		t.Errorf("expected error about missing repository, got: %s", errStr)
	}
}

func TestInitNormal(t *testing.T) {
	repo, _ := setupTestEnv(t)

	// Create a dummy preset repo
	presetDir := t.TempDir()
	for _, d := range []string{"tasks", "roles", "templates"} {
		if err := os.Mkdir(filepath.Join(presetDir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	var buf bytes.Buffer
	opts := initOptions{
		StorageMode: storageLocal,
		Preset:      presetDir,
	}

	if err := runInit(&buf, opts); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(repo, ".strand", "tasks")); err != nil {
		t.Errorf("tasks dir not created: %v", err)
	}
}
