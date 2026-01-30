package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var strandBinary string

func TestMain(m *testing.M) {
	if err := buildStrandBinary(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build strand binary: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func buildStrandBinary() error {
	tmpDir, err := os.MkdirTemp("", "strand-e2e-build")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	// We don't remove tmpDir here because we need the binary for the tests.
	// The OS will clean it up eventually, or we could track it.
	// Since it's /tmp, it's fine.

	strandBinary = filepath.Join(tmpDir, "strand")
	if runtime.GOOS == "windows" {
		strandBinary += ".exe"
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	// Assuming running from test/e2e
	repoRoot := filepath.Clean(filepath.Join(wd, "../.."))

	cmd := exec.Command("go", "build", "-o", strandBinary, "./cmd/strand")
	cmd.Dir = repoRoot
	// Inherit environment to ensure go build works (GOCACHE etc)
	cmd.Env = os.Environ()

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}
