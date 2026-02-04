package cmd

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveProjectPaths_LocalMemmd(t *testing.T) {
	// Replaces manual setup with shared helper that uses real init
	paths := setupTestProject(t, initOptions{StorageMode: storageLocal})

	base := filepath.Join(paths.GitRoot, ".strand")

	if paths.BaseDir != base {
		t.Fatalf("expected base %s, got %s", base, paths.BaseDir)
	}
	if paths.TasksDir != filepath.Join(base, "tasks") {
		t.Fatalf("unexpected tasks dir: %s", paths.TasksDir)
	}
	if paths.RolesDir != filepath.Join(base, "roles") {
		t.Fatalf("unexpected roles dir: %s", paths.RolesDir)
	}
	if paths.Storage != storageLocal {
		t.Fatalf("expected storage %s, got %s", storageLocal, paths.Storage)
	}
}

func TestResolveProjectPaths_LocalByName(t *testing.T) {
	// Test that --project <name> works for local projects from anywhere
	repo, _ := setupTestEnv(t)
	projectName := filepath.Base(repo)

	if err := runInit(io.Discard, initOptions{StorageMode: storageLocal}); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	// Resolve by explicit project name (simulating --project flag from outside repo)
	paths, err := projectPathsForName(projectName)
	if err != nil {
		t.Fatalf("projectPathsForName(%q) failed: %v", projectName, err)
	}

	expectedBase := filepath.Join(repo, ".strand")
	if paths.BaseDir != expectedBase {
		t.Fatalf("expected base %s, got %s", expectedBase, paths.BaseDir)
	}
	if paths.Storage != storageLocal {
		t.Fatalf("expected storage %s, got %s", storageLocal, paths.Storage)
	}
	if paths.GitRoot != repo {
		t.Fatalf("expected git root %s, got %s", repo, paths.GitRoot)
	}
	if paths.ProjectName != projectName {
		t.Fatalf("expected project name %s, got %s", projectName, paths.ProjectName)
	}
}

func TestResolveProjectPaths_LocalByNameFromCwd(t *testing.T) {
	repo, _ := setupTestEnv(t)
	projectName := filepath.Base(repo)

	base := filepath.Join(repo, ".strand")
	if err := ensureProjectDirs(base); err != nil {
		t.Fatalf("ensureProjectDirs failed: %v", err)
	}

	chdir(t, filepath.Dir(repo))

	paths, err := projectPathsForName(projectName)
	if err != nil {
		t.Fatalf("projectPathsForName(%q) failed: %v", projectName, err)
	}

	if paths.BaseDir != base {
		t.Fatalf("expected base %s, got %s", base, paths.BaseDir)
	}
	if paths.Storage != storageLocal {
		t.Fatalf("expected storage %s, got %s", storageLocal, paths.Storage)
	}
	if paths.GitRoot != repo {
		t.Fatalf("expected git root %s, got %s", repo, paths.GitRoot)
	}
}

func TestResolveProjectPaths_GlobalMapping(t *testing.T) {
	// Replaces manual setup with shared helper that uses real init
	_ = setupTestProject(t, initOptions{
		ProjectName: "alpha",
		StorageMode: storageGlobal,
	})

	// Explicitly verify resolving from git root works (which looks up the project map)
	paths, err := resolveProjectPaths("")
	if err != nil {
		t.Fatalf("resolveProjectPaths(\"\") failed: %v", err)
	}

	if paths.ProjectName != "alpha" {
		t.Fatalf("expected project name alpha, got %s", paths.ProjectName)
	}
	if paths.Storage != storageGlobal {
		t.Fatalf("expected storage %s, got %s", storageGlobal, paths.Storage)
	}
	// BaseDir check is implied by successful setup and valid paths return
}

func TestRunInit_GlobalStorage(t *testing.T) {
	root, _ := setupTestEnv(t)

	// We test runInit directly here, effectively what setupTestProject does but explicit
	if err := runInit(io.Discard, initOptions{ProjectName: "beta", StorageMode: storageGlobal}); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	projectsRoot, err := projectsDir()
	if err != nil {
		t.Fatalf("projectsDir failed: %v", err)
	}
	base := filepath.Join(projectsRoot, "beta")

	for _, dir := range []string{"tasks", "roles", "templates"} {
		if _, err := os.Stat(filepath.Join(base, dir)); err != nil {
			t.Fatalf("expected %s dir to exist: %v", dir, err)
		}
	}

	cfg, err := loadProjectMap()
	if err != nil {
		t.Fatalf("loadProjectMap failed: %v", err)
	}
	if cfg.Repos[root] != "beta" {
		t.Fatalf("expected mapping for %s to beta, got %s", root, cfg.Repos[root])
	}
}

func TestRunInit_LocalStoragePreset(t *testing.T) {
	repo, _ := setupTestEnv(t)

	preset := t.TempDir()
	for _, dir := range []string{"tasks", "roles", "templates"} {
		if err := os.MkdirAll(filepath.Join(preset, dir), 0o755); err != nil {
			t.Fatalf("failed to create preset %s: %v", dir, err)
		}
	}
	roleName := testRoleName(t, "preset")
	roleFile := filepath.Join(preset, "roles", roleName+".md")
	if err := os.WriteFile(roleFile, []byte("# "+strings.Title(roleName)+"\n"), 0o644); err != nil {
		t.Fatalf("failed to write preset role file: %v", err)
	}

	if err := runInit(io.Discard, initOptions{Preset: preset, StorageMode: storageLocal}); err != nil {
		t.Fatalf("runInit failed: %v", err)
	}

	base := filepath.Join(repo, ".strand")
	if _, err := os.Stat(filepath.Join(base, "roles", roleName+".md")); err != nil {
		t.Fatalf("expected preset role file to be copied: %v", err)
	}

	cfg, err := loadProjectMap()
	if err != nil {
		t.Fatalf("loadProjectMap failed: %v", err)
	}
	projectName := filepath.Base(repo)
	if cfg.LocalPaths[projectName] != repo {
		t.Fatalf("expected local project %q to be registered at %q, got %q", projectName, repo, cfg.LocalPaths[projectName])
	}
}

func TestResolveProjectPaths_LocalSiblingByName(t *testing.T) {
	base := t.TempDir()
	configDir := filepath.Join(base, "config")
	t.Setenv("STRAND_CONFIG_DIR", configDir)

	repo := filepath.Join(base, "strandyard")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatalf("failed to create repo: %v", err)
	}
	initGitRepoAt(t, repo)
	if resolved, err := filepath.EvalSymlinks(repo); err == nil {
		repo = resolved
	}
	if err := ensureProjectDirs(filepath.Join(repo, ".strand")); err != nil {
		t.Fatalf("ensureProjectDirs failed: %v", err)
	}

	otherRepo := filepath.Join(base, "other")
	if err := os.MkdirAll(otherRepo, 0o755); err != nil {
		t.Fatalf("failed to create other repo: %v", err)
	}
	initGitRepoAt(t, otherRepo)
	chdir(t, otherRepo)

	paths, err := projectPathsForName("strandyard")
	if err != nil {
		t.Fatalf("projectPathsForName failed: %v", err)
	}

	expectedBase := filepath.Join(repo, ".strand")
	if paths.BaseDir != expectedBase {
		t.Fatalf("expected base %s, got %s", expectedBase, paths.BaseDir)
	}
	if paths.GitRoot != repo {
		t.Fatalf("expected git root %s, got %s", repo, paths.GitRoot)
	}
}

func initGitRepoAt(t *testing.T, repo string) {
	t.Helper()
	cmd := exec.Command("git", "init")
	cmd.Dir = repo
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git init failed: %s", string(output))
	}
}
