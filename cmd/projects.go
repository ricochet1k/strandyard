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
	Repos map[string]string `json:"repos"`
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
			return projectMap{Repos: map[string]string{}}, nil
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
