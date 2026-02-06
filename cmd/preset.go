/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// presetCmd represents the preset command
var presetCmd = &cobra.Command{
	Use:   "preset",
	Short: "Manage project presets",
	Long:  `Manage project presets, including refreshing roles and templates from a preset source.`,
}

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh <preset>",
	Short: "Refresh roles and templates from a preset",
	Long: `Refresh roles and templates from a preset source (local directory or git URL).
This command will overwrite existing role and template files but will not touch the tasks directory.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		preset := args[0]
		return runPresetRefresh(cmd.OutOrStdout(), preset)
	},
}

func init() {
	rootCmd.AddCommand(presetCmd)
	presetCmd.AddCommand(refreshCmd)
}

func runPresetRefresh(w io.Writer, preset string) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return fmt.Errorf("project not initialized: %w", err)
	}

	fmt.Fprintf(w, "Refreshing roles and templates from preset %q...\n", preset)

	if err := applyPreset(w, paths.BaseDir, preset, []string{"roles", "templates"}); err != nil {
		return err
	}

	fmt.Fprintln(w, "✓ Refresh complete. Running repair...")

	return runRepair(w, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text")
}
