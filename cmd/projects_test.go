package cmd

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func initGitRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
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

func setHome(t *testing.T, dir string) {
	t.Helper()
	prev := os.Getenv("HOME")
	if err := os.Setenv("HOME", dir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Setenv("HOME", prev)
	})
}

func TestResolveProjectPaths_LocalMemmd(t *testing.T) {
	repo := initGitRepo(t)
	chdir(t, repo)
	root, err := gitRootDir()
	if err != nil {
		t.Fatalf("gitRootDir failed: %v", err)
	}

	base := filepath.Join(root, ".strand")
	if err := ensureProjectDirs(base); err != nil {
		t.Fatalf("failed to create .strand dirs: %v", err)
	}

	paths, err := resolveProjectPaths("")
	if err != nil {
		t.Fatalf("resolveProjectPaths failed: %v", err)
	}

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
	repo := initGitRepo(t)
	chdir(t, repo)
	setHome(t, t.TempDir())
	root, err := gitRootDir()
	if err != nil {
		t.Fatalf("gitRootDir failed: %v", err)
	}

	projectsRoot, err := projectsDir()
	if err != nil {
		t.Fatalf("projectsDir failed: %v", err)
	}
	base := filepath.Join(projectsRoot, "alpha")
	if err := ensureProjectDirs(base); err != nil {
		t.Fatalf("failed to create project dirs: %v", err)
	}

	if err := saveProjectMap(projectMap{Repos: map[string]string{root: "alpha"}}); err != nil {
		t.Fatalf("saveProjectMap failed: %v", err)
	}

	paths, err := resolveProjectPaths("")
	if err != nil {
		t.Fatalf("resolveProjectPaths failed: %v", err)
	}

	if paths.BaseDir != base {
		t.Fatalf("expected base %s, got %s", base, paths.BaseDir)
	}
	if paths.ProjectName != "alpha" {
		t.Fatalf("expected project name alpha, got %s", paths.ProjectName)
	}
	if paths.Storage != storageGlobal {
		t.Fatalf("expected storage %s, got %s", storageGlobal, paths.Storage)
	}
}

func TestRunInit_GlobalStorage(t *testing.T) {
	repo := initGitRepo(t)
	chdir(t, repo)
	setHome(t, t.TempDir())
	root, err := gitRootDir()
	if err != nil {
		t.Fatalf("gitRootDir failed: %v", err)
	}

	prevStorage := initStorageMode
	prevPreset := initPreset
	t.Cleanup(func() {
		initStorageMode = prevStorage
		initPreset = prevPreset
	})

	initStorageMode = storageGlobal
	initPreset = ""

	if err := runInit(io.Discard, initOptionsFromFlags("beta")); err != nil {
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
	repo := initGitRepo(t)
	chdir(t, repo)
	setHome(t, t.TempDir())

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

	prevStorage := initStorageMode
	prevPreset := initPreset
	t.Cleanup(func() {
		initStorageMode = prevStorage
		initPreset = prevPreset
	})

	initStorageMode = storageLocal
	initPreset = preset

	if err := runInit(io.Discard, initOptionsFromFlags("")); err != nil {
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
