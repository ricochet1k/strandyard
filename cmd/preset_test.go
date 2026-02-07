package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPresetRefresh(t *testing.T) {
	_, _ = setupTestEnv(t)

	// Create a dummy preset repo
	presetDir := t.TempDir()
	for _, d := range []string{"tasks", "roles", "templates"} {
		if err := os.Mkdir(filepath.Join(presetDir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	// Add some files to preset
	os.WriteFile(filepath.Join(presetDir, "roles", "dev.md"), []byte("dev role"), 0o644)
	os.WriteFile(filepath.Join(presetDir, "templates", "task.md"), []byte("task template"), 0o644)

	// Initialize project with this preset
	opts := initOptions{
		StorageMode: storageLocal,
		Preset:      presetDir,
	}
	if err := runInit(io.Discard, opts); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Verify initial files
	paths, err := resolveProjectPaths("")
	if err != nil {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(filepath.Join(paths.RolesDir, "dev.md"))
	if !strings.Contains(string(content), "dev") {
		t.Errorf("expected it to contain 'dev', got %q", string(content))
	}

	// Update preset files
	os.WriteFile(filepath.Join(presetDir, "roles", "dev.md"), []byte("updated dev role"), 0o644)
	os.WriteFile(filepath.Join(presetDir, "templates", "task.md"), []byte("updated task template"), 0o644)
	os.WriteFile(filepath.Join(presetDir, "tasks", "new-task.md"), []byte("new task"), 0o644)

	// Refresh from preset
	var buf bytes.Buffer
	if err := runPresetRefresh(&buf, presetDir); err != nil {
		t.Fatalf("runPresetRefresh failed: %v", err)
	}

	// Verify roles and templates updated
	content, _ = os.ReadFile(filepath.Join(paths.RolesDir, "dev.md"))
	if string(content) != "updated dev role" {
		t.Errorf("expected 'updated dev role', got %q", string(content))
	}
	content, _ = os.ReadFile(filepath.Join(paths.TemplatesDir, "task.md"))
	if string(content) != "updated task template" {
		t.Errorf("expected 'updated task template', got %q", string(content))
	}

	// Verify tasks NOT updated (new-task.md should not exist in project)
	if _, err := os.Stat(filepath.Join(paths.TasksDir, "new-task.md")); err == nil {
		t.Error("new-task.md from preset should NOT have been copied to project tasks dir")
	}
}

func TestPresetRefreshGit(t *testing.T) {
	_, _ = setupTestEnv(t)

	// Create a dummy preset git repo
	presetDir := t.TempDir()
	cmd := exec.Command("git", "init")
	cmd.Dir = presetDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %s", string(out))
	}

	for _, d := range []string{"tasks", "roles", "templates"} {
		if err := os.Mkdir(filepath.Join(presetDir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	os.WriteFile(filepath.Join(presetDir, "roles", "dev.md"), []byte("---\ndescription: dev\n---\n# Dev"), 0o644)
	os.WriteFile(filepath.Join(presetDir, "templates", "task.md"), []byte("task template"), 0o644)
	os.WriteFile(filepath.Join(presetDir, "tasks", "T1234-dummy.md"), []byte("---\nrole: dev\npriority: medium\n---\n# Dummy"), 0o644)

	cmd = exec.Command("git", "add", ".")
	cmd.Dir = presetDir
	cmd.Run()
	cmd = exec.Command("git", "commit", "-m", "initial")
	cmd.Dir = presetDir
	cmd.Run()

	// Initialize project with this preset
	opts := initOptions{
		StorageMode: storageLocal,
		Preset:      presetDir,
	}
	if err := runInit(io.Discard, opts); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Verify initial files
	paths, err := resolveProjectPaths("")
	if err != nil {
		t.Fatal(err)
	}
	content, _ := os.ReadFile(filepath.Join(paths.RolesDir, "dev.md"))
	if !strings.Contains(string(content), "dev") {
		t.Errorf("expected it to contain 'dev', got %q", string(content))
	}

	// Update preset files and commit
	os.WriteFile(filepath.Join(presetDir, "roles", "dev.md"), []byte("---\ndescription: updated dev\n---\n# Updated Dev"), 0o644)
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = presetDir
	cmd.Run()
	cmd = exec.Command("git", "commit", "-m", "update")
	cmd.Dir = presetDir
	cmd.Run()

	// Refresh from preset (should clone as it's not detected as a directory easily if we use the path directly and it has .git?
	// Actually os.Stat will say it's a directory. But if we pass it as a file:// URL or just the path, applyPreset handles it.)

	var buf bytes.Buffer
	if err := runPresetRefresh(&buf, presetDir); err != nil {
		t.Fatalf("runPresetRefresh failed: %v", err)
	}

	// Verify roles updated
	content, _ = os.ReadFile(filepath.Join(paths.RolesDir, "dev.md"))
	if !strings.Contains(string(content), "description: updated dev") {
		t.Errorf("expected it to contain 'description: updated dev', got %q", string(content))
	}
}

func TestPresetRefreshNotInitialized(t *testing.T) {
	_, _ = setupTestEnv(t)
	// Do NOT runInit

	if err := runPresetRefresh(io.Discard, "some-preset"); err == nil {
		t.Error("expected error when refreshing in uninitialized project, got nil")
	}
}

func TestPresetRefreshMissingDirectory(t *testing.T) {
	_, _ = setupTestEnv(t)

	// Create a preset with missing templates directory
	presetDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(presetDir, "roles"), 0o755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(presetDir, "roles", "dev.md"), []byte("dev role"), 0o644)

	// Initialize project
	opts := initOptions{
		StorageMode: storageLocal,
	}
	if err := runInit(io.Discard, opts); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Try to refresh from incomplete preset
	var buf bytes.Buffer
	err := runPresetRefresh(&buf, presetDir)
	if err == nil {
		t.Fatal("expected error when preset is missing templates directory, got nil")
	}
	if !strings.Contains(err.Error(), "missing required directories") {
		t.Errorf("expected error about missing directories, got: %v", err)
	}
	if !strings.Contains(err.Error(), "templates") {
		t.Errorf("expected error to mention 'templates', got: %v", err)
	}
}

func TestPresetRefreshInvalidGitURL(t *testing.T) {
	_, _ = setupTestEnv(t)

	// Initialize project
	opts := initOptions{
		StorageMode: storageLocal,
	}
	if err := runInit(io.Discard, opts); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Try to refresh from invalid git URL
	var buf bytes.Buffer
	err := runPresetRefresh(&buf, "https://github.com/nonexistent/repo-that-does-not-exist-12345.git")
	if err == nil {
		t.Fatal("expected error when cloning invalid git URL, got nil")
	}
	if !strings.Contains(err.Error(), "failed to clone") {
		t.Errorf("expected error about clone failure, got: %v", err)
	}
}

func TestPresetRefreshEmptyPath(t *testing.T) {
	_, _ = setupTestEnv(t)

	// Initialize project
	opts := initOptions{
		StorageMode: storageLocal,
	}
	if err := runInit(io.Discard, opts); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Try to refresh with empty preset path
	var buf bytes.Buffer
	err := runPresetRefresh(&buf, "")
	if err == nil {
		t.Fatal("expected error when preset path is empty, got nil")
	}
	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("expected error about empty path, got: %v", err)
	}
}

func TestPresetRefreshVerboseOutput(t *testing.T) {
	_, _ = setupTestEnv(t)

	// Create a proper preset
	presetDir := t.TempDir()
	for _, d := range []string{"tasks", "roles", "templates"} {
		if err := os.Mkdir(filepath.Join(presetDir, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	os.WriteFile(filepath.Join(presetDir, "roles", "dev.md"), []byte("dev role"), 0o644)
	os.WriteFile(filepath.Join(presetDir, "templates", "task.md"), []byte("task template"), 0o644)

	// Initialize project
	opts := initOptions{
		StorageMode: storageLocal,
		Preset:      presetDir,
	}
	if err := runInit(io.Discard, opts); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Refresh and capture output
	var buf bytes.Buffer
	if err := runPresetRefresh(&buf, presetDir); err != nil {
		t.Fatalf("runPresetRefresh failed: %v", err)
	}

	output := buf.String()
	expectedPhrases := []string{
		"Using local preset directory",
		"Validating preset structure",
		"✓ Preset structure validated",
		"Refreshing roles/",
		"Refreshing templates/",
		"✓ Refresh complete",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("expected output to contain %q, got:\n%s", phrase, output)
		}
	}
}
