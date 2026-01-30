package cmd

import (
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func testToken(parts ...string) string {
	h := fnv.New32a()
	for _, part := range parts {
		_, _ = h.Write([]byte(part))
	}
	return fmt.Sprintf("%08x", h.Sum32())[:6]
}

func testRoleName(t *testing.T, suffix string) string {
	name := strings.TrimSpace(t.Name())
	return "role-" + testToken("role", name, suffix)
}

func initGitRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	// resolve symlinks for consistency
	repo, err := filepath.EvalSymlinks(repo)
	if err != nil {
		t.Fatalf("failed to eval symlinks for temp dir: %v", err)
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = repo
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %s", string(output))
	}
	return repo
}

func chdir(t *testing.T, dir string) {
	t.Helper()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to %s: %v", dir, err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(prev)
	})
}

// setupTestEnv prepares a test environment with a git repo and a config dir.
// It sets the STRAND_CONFIG_DIR environment variable.
// It returns the path to the git repo root and the config dir.
func setupTestEnv(t *testing.T) (string, string) {
	t.Helper()
	repo := initGitRepo(t)
	chdir(t, repo)

	configDir := t.TempDir()
	t.Setenv("STRAND_CONFIG_DIR", configDir)

	return repo, configDir
}

// setupTestProject initializes a project using the real runInit.
func setupTestProject(t *testing.T, opts initOptions) projectPaths {
	t.Helper()
	repo, _ := setupTestEnv(t)

	if err := runInit(io.Discard, opts); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// If local storage, verify by resolving from current dir (empty project name)
	if opts.StorageMode == storageLocal {
		paths, err := resolveProjectPaths("")
		if err != nil {
			t.Fatalf("resolveProjectPaths failed for local: %v", err)
		}
		return paths
	}

	projectName := opts.ProjectName
	if projectName == "" {
		projectName = filepath.Base(repo)
	}

	// resolveProjectPaths might depend on CWD being the git root, which setupTestEnv ensures.
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		t.Fatalf("resolveProjectPaths failed: %v", err)
	}
	return paths
}
