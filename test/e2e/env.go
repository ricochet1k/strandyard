package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestEnv provides an isolated test environment for E2E tests
type TestEnv struct {
	t        testing.TB
	rootDir  string
	baseDir  string
	tasksDir string
}

// NewTestEnv creates a new isolated test environment
func NewTestEnv(t testing.TB) *TestEnv {
	t.Helper()

	rootDir, err := os.MkdirTemp("", "strand-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	baseDir := filepath.Join(rootDir, ".strand")
	tasksDir := filepath.Join(baseDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatalf("Failed to create tasks dir: %v", err)
	}

	rolesDir := filepath.Join(baseDir, "roles")
	if err := os.MkdirAll(rolesDir, 0755); err != nil {
		t.Fatalf("Failed to create roles dir: %v", err)
	}

	templatesDir := filepath.Join(baseDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		t.Fatalf("Failed to create templates dir: %v", err)
	}

	gitInit := exec.Command("git", "init")
	gitInit.Dir = rootDir
	if output, err := gitInit.CombinedOutput(); err != nil {
		t.Fatalf("Failed to init git repo: %v\nOutput: %s", err, string(output))
	}

	env := &TestEnv{
		t:        t,
		rootDir:  rootDir,
		baseDir:  baseDir,
		tasksDir: tasksDir,
	}

	// Register cleanup
	t.Cleanup(env.Cleanup)

	return env
}

// Cleanup removes the test environment
func (e *TestEnv) Cleanup() {
	if e.rootDir != "" {
		os.RemoveAll(e.rootDir)
	}
}

// Root returns the root directory of the test environment
func (e *TestEnv) Root() string {
	return e.rootDir
}

// TasksDir returns the tasks directory
func (e *TestEnv) TasksDir() string {
	return e.tasksDir
}

// Path returns a path relative to the test environment root
func (e *TestEnv) Path(relPath string) string {
	return filepath.Join(e.rootDir, relPath)
}
