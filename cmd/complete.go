/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ricochet1k/streamyard/pkg/task"
	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete <task-id>",
	Short: "Mark a task as completed",
	Long: `Mark a task as completed by setting completed: true in the frontmatter.
Also updates the date_edited field to the current time.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		return runComplete(cmd.OutOrStdout(), projectName, taskID)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}

func runComplete(w io.Writer, projectName, taskID string) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	// Load all tasks to find the one we want
	parser := task.NewParser()
	tasks, err := parser.LoadTasks(paths.TasksDir)
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	resolvedID, err := task.ResolveTaskID(tasks, taskID)
	if err != nil {
		return err
	}
	taskID = resolvedID

	// Find the task by ID
	t, exists := tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	// Check if already completed
	if t.Meta.Completed {
		fmt.Fprintf(w, "Task %s is already marked as completed\n", task.ShortID(taskID))
		return nil
	}

	// Calculate incremental update before marking the task complete
	update, err := task.CalculateIncrementalFreeListUpdate(tasks, taskID)
	if err != nil {
		return fmt.Errorf("failed to calculate incremental update: %w", err)
	}

	// Update metadata
	t.Meta.Completed = true
	t.Meta.DateEdited = time.Now().UTC()
	t.MarkDirty()

	if err := t.Write(); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	fmt.Fprintf(w, "✓ Task %s marked as completed\n", task.ShortID(taskID))

	if strings.TrimSpace(t.Meta.Parent) != "" {
		if _, err := task.UpdateParentTodoEntries(tasks, t.Meta.Parent); err != nil {
			return fmt.Errorf("failed to update parent task TODO entries: %w", err)
		}
		if _, err := task.WriteDirtyTasks(tasks); err != nil {
			return fmt.Errorf("failed to write parent task updates: %w", err)
		}
	}

	// Try incremental update first, fall back to full validation
	if err := task.UpdateFreeListIncrementally(tasks, paths.FreeTasksFile, update); err != nil {
		fmt.Fprintf(w, "⚠️  Incremental update failed, falling back to full repair: %v\n", err)
		if err := runRepair(w, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(w, "✓ Incrementally updated free-tasks.md\n")
		// Still need to run repair for error checking, but skip master list generation
		validator := task.NewValidatorWithRoles(tasks, paths.RolesDir)
		errors := validator.Validate()
		if _, err := task.WriteDirtyTasks(tasks); err != nil {
			return fmt.Errorf("failed to write repaired tasks: %w", err)
		}
		if len(errors) > 0 {
			fmt.Fprintf(w, "⚠️  Repair errors found:\n")
			for _, e := range errors {
				fmt.Fprintf(w, "ERROR: %s\n", e.Error())
			}
			return fmt.Errorf("repair failed: %d error(s)", len(errors))
		}
	}

	fmt.Fprintf(w, "💡 Consider committing your changes: git add -A && git commit -m \"complete: %s\"\n", task.ShortID(taskID))

	return nil
}
