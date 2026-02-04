/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [project_name]",
	Short: "Initialize strand storage",
	Long:  "Initialize the strand project storage.\n\nBy default this creates a global project under ~/.config/strand/projects/<project_name> and records a mapping from the current git root to the project name. Use --storage=local to place tasks, roles, and templates inside .strand/ at the git root instead.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		project := ""
		if len(args) > 0 {
			project = strings.TrimSpace(args[0])
		}
		opts := initOptionsFromFlags(project)
		return runInit(cmd.OutOrStdout(), opts)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	initCmd.Flags().StringVar(&initStorageMode, "storage", storageGlobal, "storage mode: global or local")
	initCmd.Flags().StringVar(&initPreset, "preset", "", "preset directory or git repo to seed tasks/roles/templates")
}

var (
	initStorageMode string
	initPreset      string
)

type initOptions struct {
	ProjectName string
	StorageMode string
	Preset      string
}

func initOptionsFromFlags(projectArg string) initOptions {
	return initOptions{
		ProjectName: strings.TrimSpace(projectArg),
		StorageMode: strings.TrimSpace(initStorageMode),
		Preset:      strings.TrimSpace(initPreset),
	}
}

func runInit(w io.Writer, opts initOptions) error {
	storage := strings.ToLower(strings.TrimSpace(opts.StorageMode))
	if storage == "" {
		storage = storageGlobal
	}
	switch storage {
	case storageGlobal, storageLocal:
	default:
		return fmt.Errorf("invalid storage mode %q (expected global or local)", storage)
	}

	gitRoot, err := gitRootDir()
	if err != nil {
		return err
	}

	projectName := strings.TrimSpace(opts.ProjectName)
	if projectName == "" {
		projectName = filepath.Base(gitRoot)
	}

	var baseDir string
	if storage == storageLocal {
		baseDir = filepath.Join(gitRoot, ".strand")
	} else {
		projectsRoot, err := projectsDir()
		if err != nil {
			return err
		}
		baseDir = filepath.Join(projectsRoot, projectName)
	}

	if info, err := os.Stat(baseDir); err == nil {
		if info.IsDir() {
			return fmt.Errorf("strand already initialized at %s", baseDir)
		}
		return fmt.Errorf("strand path exists and is not a directory: %s", baseDir)
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := ensureProjectDirs(baseDir); err != nil {
		return err
	}

	if strings.TrimSpace(opts.Preset) != "" {
		if err := applyPreset(baseDir, strings.TrimSpace(opts.Preset)); err != nil {
			return err
		}
	}

	cfg, err := loadProjectMap()
	if err != nil {
		return err
	}
	if storage == storageGlobal {
		cfg.Repos[gitRoot] = projectName
	} else {
		cfg.LocalPaths[projectName] = gitRoot
	}
	if err := saveProjectMap(cfg); err != nil {
		return err
	}

	fmt.Fprintf(w, "✓ Initialized strand at %s\n", baseDir)
	fmt.Fprintf(w, "✓ Linked %s to project %s\n", gitRoot, projectName)
	return nil
}

func applyPreset(baseDir, preset string) error {
	sourceDir := preset
	cleanup := func() {}

	if info, err := os.Stat(preset); err != nil || !info.IsDir() {
		tempDir, err := os.MkdirTemp("", "strand-preset-")
		if err != nil {
			return fmt.Errorf("failed to create temp dir: %w", err)
		}
		cleanup = func() {
			_ = os.RemoveAll(tempDir)
		}
		cmd := exec.Command("git", "clone", "--depth", "1", "--", preset, tempDir)
		if output, err := cmd.CombinedOutput(); err != nil {
			cleanup()
			return fmt.Errorf("failed to clone preset: %s", strings.TrimSpace(string(output)))
		}
		sourceDir = tempDir
	}
	defer cleanup()

	for _, name := range []string{"tasks", "roles", "templates"} {
		src := filepath.Join(sourceDir, name)
		dst := filepath.Join(baseDir, name)
		info, err := os.Stat(src)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("preset is missing %s directory", name)
			}
			return err
		}
		if !info.IsDir() {
			return fmt.Errorf("preset %s is not a directory", src)
		}
		if err := copyDir(src, dst); err != nil {
			return fmt.Errorf("failed to copy %s: %w", name, err)
		}
	}

	return nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			if rel == "." {
				return nil
			}
			return os.MkdirAll(target, 0o755)
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		return os.WriteFile(target, data, info.Mode())
	})
}
