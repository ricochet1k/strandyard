package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	storageGlobal = "global"
	storageLocal  = "local"
)

type projectPaths struct {
	BaseDir       string
	TasksDir      string
	RolesDir      string
	TemplatesDir  string
	RootTasksFile string
	FreeTasksFile string
	ProjectName   string
	GitRoot       string
	Storage       string
}

type projectMap struct {
	Repos       map[string]string `json:"repos"`
	LocalPaths  map[string]string `json:"local_paths,omitempty"`
}

func resolveProjectPaths(projectName string) (projectPaths, error) {
	if strings.TrimSpace(projectName) != "" {
		return projectPathsForName(strings.TrimSpace(projectName))
	}

	gitRoot, err := gitRootDir()
	if err != nil {
		return projectPaths{}, err
	}

	localDir := filepath.Join(gitRoot, ".strand")
	if info, err := os.Stat(localDir); err == nil && info.IsDir() {
		return projectPathsFromBase(localDir, "", gitRoot, storageLocal)
	}

	cfg, err := loadProjectMap()
	if err != nil {
		return projectPaths{}, err
	}
	project, ok := cfg.Repos[gitRoot]
	if !ok || strings.TrimSpace(project) == "" {
		return projectPaths{}, fmt.Errorf("no strand project found for %s (run strand init or pass --project)", gitRoot)
	}
	return projectPathsForName(project)
}

func projectPathsForName(projectName string) (projectPaths, error) {
	// Check if this is a local project registered in projects.json
	cfg, err := loadProjectMap()
	if err != nil {
		return projectPaths{}, err
	}
	if gitRoot, ok := cfg.LocalPaths[projectName]; ok {
		localDir := filepath.Join(gitRoot, ".strand")
		if info, err := os.Stat(localDir); err == nil && info.IsDir() {
			return projectPathsFromBase(localDir, projectName, gitRoot, storageLocal)
		}
	}

	if paths, ok, err := projectPathsForLocalDir(projectName); ok || err != nil {
		return paths, err
	}

	// Fall back to global projects directory
	dir, err := projectsDir()
	if err != nil {
		return projectPaths{}, err
	}
	base := filepath.Join(dir, projectName)
	info, err := os.Stat(base)
	if err != nil {
		if os.IsNotExist(err) {
			return projectPaths{}, fmt.Errorf("project %q not found in %s", projectName, dir)
		}
		return projectPaths{}, err
	}
	if !info.IsDir() {
		return projectPaths{}, fmt.Errorf("project path %s is not a directory", base)
	}
	return projectPathsFromBase(base, projectName, "", storageGlobal)
}

func projectPathsForLocalDir(projectName string) (projectPaths, bool, error) {
	if projectName == "" {
		return projectPaths{}, false, nil
	}
	if filepath.IsAbs(projectName) {
		return projectPathsForLocalCandidate(projectName, projectName)
	}

	searchRoots, err := projectSearchRoots()
	if err != nil {
		return projectPaths{}, false, err
	}

	useOnlyCwd := strings.ContainsRune(projectName, filepath.Separator)
	for _, root := range searchRoots {
		if useOnlyCwd && root.kind != searchRootCwd {
			continue
		}
		candidate := filepath.Join(root.path, projectName)
		paths, ok, err := projectPathsForLocalCandidate(candidate, projectName)
		if ok || err != nil {
			return paths, ok, err
		}
	}

	return projectPaths{}, false, nil
}

type searchRootKind int

const (
	searchRootCwd searchRootKind = iota
	searchRootGitRoot
	searchRootGitParent
)

type projectSearchRoot struct {
	path string
	kind searchRootKind
}

func projectSearchRoots() ([]projectSearchRoot, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	roots := []projectSearchRoot{{path: cwd, kind: searchRootCwd}}

	gitRoot, err := gitRootDir()
	if err != nil {
		return roots, nil
	}
	if gitRoot != cwd {
		roots = append(roots, projectSearchRoot{path: gitRoot, kind: searchRootGitRoot})
	}
	parent := filepath.Dir(gitRoot)
	if parent != "" && parent != gitRoot {
		roots = append(roots, projectSearchRoot{path: parent, kind: searchRootGitParent})
	}

	return roots, nil
}

func projectPathsForLocalCandidate(candidate, projectName string) (projectPaths, bool, error) {
	info, err := os.Stat(candidate)
	if err != nil || !info.IsDir() {
		return projectPaths{}, false, nil
	}

	baseDir := candidate
	gitRoot := candidate
	if filepath.Base(candidate) != ".strand" {
		strandDir := filepath.Join(candidate, ".strand")
		if info, err := os.Stat(strandDir); err == nil && info.IsDir() {
			baseDir = strandDir
		} else if !hasStrandLayout(candidate) {
			return projectPaths{}, false, nil
		}
	} else {
		gitRoot = filepath.Dir(candidate)
	}

	if !hasStrandLayout(baseDir) {
		return projectPaths{}, false, nil
	}

	paths, err := projectPathsFromBase(baseDir, projectName, gitRoot, storageLocal)
	if err != nil {
		return projectPaths{}, false, err
	}
	return paths, true, nil
}

func hasStrandLayout(baseDir string) bool {
	for _, name := range []string{"tasks", "roles", "templates"} {
		info, err := os.Stat(filepath.Join(baseDir, name))
		if err != nil || !info.IsDir() {
			return false
		}
	}
	return true
}

func projectPathsFromBase(base, projectName, gitRoot, storage string) (projectPaths, error) {
	tasksDir := filepath.Join(base, "tasks")
	rolesDir := filepath.Join(base, "roles")
	templatesDir := filepath.Join(base, "templates")
	return projectPaths{
		BaseDir:       base,
		TasksDir:      tasksDir,
		RolesDir:      rolesDir,
		TemplatesDir:  templatesDir,
		RootTasksFile: filepath.Join(tasksDir, "root-tasks.md"),
		FreeTasksFile: filepath.Join(tasksDir, "free-tasks.md"),
		ProjectName:   projectName,
		GitRoot:       gitRoot,
		Storage:       storage,
	}, nil
}

func gitRootDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", errors.New("unable to locate git root (run inside a git repo or use --project)")
	}
	root := strings.TrimSpace(string(output))
	if root == "" {
		return "", errors.New("unable to locate git root (run inside a git repo or use --project)")
	}
	if resolved, err := filepath.EvalSymlinks(root); err == nil {
		root = resolved
	}
	return root, nil
}

func configDir() (string, error) {
	if dir := os.Getenv("STRAND_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to resolve home directory: %w", err)
	}
	return filepath.Join(home, ".config", "strand"), nil
}

func projectsDir() (string, error) {
	base, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "projects"), nil
}

func projectMapPath() (string, error) {
	base, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "projects.json"), nil
}

func loadProjectMap() (projectMap, error) {
	path, err := projectMapPath()
	if err != nil {
		return projectMap{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return projectMap{Repos: map[string]string{}, LocalPaths: map[string]string{}}, nil
		}
		return projectMap{}, fmt.Errorf("failed to read project map: %w", err)
	}
	var cfg projectMap
	if err := json.Unmarshal(data, &cfg); err != nil {
		return projectMap{}, fmt.Errorf("failed to parse project map: %w", err)
	}
	if cfg.Repos == nil {
		cfg.Repos = map[string]string{}
	}
	if cfg.LocalPaths == nil {
		cfg.LocalPaths = map[string]string{}
	}
	return cfg, nil
}

func saveProjectMap(cfg projectMap) error {
	path, err := projectMapPath()
	if err != nil {
		return err
	}
	if cfg.Repos == nil {
		cfg.Repos = map[string]string{}
	}
	if cfg.LocalPaths == nil {
		cfg.LocalPaths = map[string]string{}
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize project map: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

func ensureProjectDirs(base string) error {
	for _, name := range []string{"tasks", "roles", "templates"} {
		if err := os.MkdirAll(filepath.Join(base, name), 0o755); err != nil {
			return fmt.Errorf("failed to create %s directory: %w", name, err)
		}
	}
	return nil
}
