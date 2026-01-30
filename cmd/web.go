package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/ricochet1k/strandyard/pkg/web"
	"github.com/spf13/cobra"
)

var (
	webPort   int
	webNoOpen bool
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web dashboard server",
	Long:  "Start a web server that watches all strand projects and serves the dashboard UI.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runWeb()
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().IntVar(&webPort, "port", 8686, "port to listen on")
	webCmd.Flags().BoolVar(&webNoOpen, "no-open", false, "don't auto-open browser")
}

func runWeb() error {
	// Read environment variables
	envRoot := os.Getenv("STRAND_ROOT")
	envStorage := os.Getenv("STRAND_STORAGE")

	// Resolve current project
	var currentProject string
	if paths, err := resolveProjectPaths(""); err == nil {
		currentProject = paths.ProjectName
		if currentProject == "" {
			currentProject = "local"
		}
	}

	// Discover all projects
	projects, err := discoverAllProjects(envRoot, envStorage)
	if err != nil {
		return err
	}

	if len(projects) == 0 {
		return fmt.Errorf("no projects found (run 'strand init' first)")
	}

	cfg := web.ServerConfig{
		Port:           webPort,
		Projects:       projects,
		CurrentProject: currentProject,
		AutoOpen:       !webNoOpen,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	return web.Serve(ctx, cfg)
}

func discoverAllProjects(envRoot, envStorage string) ([]web.ProjectInfo, error) {
	var projects []web.ProjectInfo
	filter := strings.ToLower(strings.TrimSpace(envStorage))

	// Global projects
	if filter == "" || filter == "global" {
		if dir, err := projectsDir(); err == nil {
			entries, _ := os.ReadDir(dir)
			for _, e := range entries {
				if !e.IsDir() {
					continue
				}
				name := e.Name()
				baseDir := filepath.Join(dir, name)
				if hasProjectStructure(baseDir) {
					gitRoot := findGitRootForProject(name)
					projects = append(projects, web.ProjectInfo{
						Name:          name,
						StorageRoot:   baseDir,
						TasksRoot:     filepath.Join(baseDir, "tasks"),
						RolesRoot:     filepath.Join(baseDir, "roles"),
						TemplatesRoot: filepath.Join(baseDir, "templates"),
						GitRoot:       gitRoot,
						Storage:       "global",
					})
				}
			}
		}
	}

	// Local project
	if filter == "" || filter == "local" {
		gitRoot := envRoot
		if gitRoot == "" {
			gitRoot, _ = gitRootDir()
		}

		if gitRoot != "" {
			localDir := filepath.Join(gitRoot, ".strand")
			if info, err := os.Stat(localDir); err == nil && info.IsDir() {
				if hasProjectStructure(localDir) {
					projects = append(projects, web.ProjectInfo{
						Name:          "local",
						StorageRoot:   localDir,
						TasksRoot:     filepath.Join(localDir, "tasks"),
						RolesRoot:     filepath.Join(localDir, "roles"),
						TemplatesRoot: filepath.Join(localDir, "templates"),
						GitRoot:       gitRoot,
						Storage:       "local",
					})
				}
			}
		}
	}

	return projects, nil
}

func hasProjectStructure(dir string) bool {
	for _, sub := range []string{"tasks", "roles", "templates"} {
		path := filepath.Join(dir, sub)
		if info, err := os.Stat(path); err != nil || !info.IsDir() {
			return false
		}
	}
	return true
}

func findGitRootForProject(projectName string) string {
	cfg, _ := loadProjectMap()
	for gitRoot, proj := range cfg.Repos {
		if proj == projectName {
			return gitRoot
		}
	}
	return ""
}
