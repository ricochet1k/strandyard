/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
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
		if err := applyPreset(w, baseDir, strings.TrimSpace(opts.Preset), []string{"tasks", "roles", "templates"}); err != nil {
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
