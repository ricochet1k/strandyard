package e2e

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEnv provides an isolated test environment for E2E tests
type TestEnv struct {
	t        testing.TB
	rootDir  string
	tasksDir string
}

// NewTestEnv creates a new isolated test environment
func NewTestEnv(t testing.TB) *TestEnv {
	t.Helper()

	rootDir, err := os.MkdirTemp("", "memmd-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	tasksDir := filepath.Join(rootDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatalf("Failed to create tasks dir: %v", err)
	}

	rolesDir := filepath.Join(rootDir, "roles")
	if err := os.MkdirAll(rolesDir, 0755); err != nil {
		t.Fatalf("Failed to create roles dir: %v", err)
	}

	env := &TestEnv{
		t:        t,
		rootDir:  rootDir,
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
