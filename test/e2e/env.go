package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var strandBinary string

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

	// Use the built binary to initialize the project
	// This ensures we test the actual init logic and directory structure
	initCmd := exec.Command(strandBinary, "init", "--storage=local")
	initCmd.Dir = rootDir
	initCmd.Env = append(os.Environ(), "STRAND_CONFIG_DIR="+rootDir) // Isolate global config just in case, though --storage=local shouldn't use it
	if output, err := initCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to run strand init: %v\nOutput: %s", err, string(output))
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
