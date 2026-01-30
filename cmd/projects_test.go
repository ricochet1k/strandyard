package cmd

import (
	"io"
	"os"
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

	mapPath, err := projectMapPath()
	if err != nil {
		t.Fatalf("projectMapPath failed: %v", err)
	}
	if _, err := os.Stat(mapPath); !os.IsNotExist(err) {
		t.Fatalf("expected no project map for local init")
	}
}
