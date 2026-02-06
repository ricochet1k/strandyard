package cmd

import (
	"fmt"
	"io"

	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
)

var cancelCmd = &cobra.Command{
	Use:   "cancel <task-id> [reason]",
	Short: "Mark a task as cancelled",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		reason := ""
		if len(args) > 1 {
			reason = args[1]
		}
		return runSetStatus(cmd.OutOrStdout(), taskID, task.StatusCancelled, reason)
	},
}

var markDuplicateCmd = &cobra.Command{
	Use:   "mark-duplicate <task-id> <duplicate-of>",
	Short: "Mark a task as a duplicate of another task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		duplicateOf := args[1]
		return runSetStatus(cmd.OutOrStdout(), taskID, task.StatusDuplicate, "Duplicate of "+duplicateOf)
	},
}

var markInProgressCmd = &cobra.Command{
	Use:   "mark-in-progress <task-id>",
	Short: "Mark a task as in progress",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetStatus(cmd.OutOrStdout(), args[0], task.StatusInProgress, "")
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)
	rootCmd.AddCommand(markDuplicateCmd)
	rootCmd.AddCommand(markInProgressCmd)
}

func runSetStatus(w io.Writer, inputID, status, report string) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	_, taskID, err := db.GetResolved(inputID)
	if err != nil {
		return err
	}

	if err := db.SetStatus(taskID, status); err != nil {
		return err
	}

	if report != "" {
		if err := db.AppendCompletionReport(taskID, report); err != nil {
			return err
		}
	}

	// For statuses that count as "not active", update blockers
	if status == task.StatusCancelled || status == task.StatusDuplicate || status == task.StatusDone {
		if err := db.UpdateBlockersAfterCompletion(taskID); err != nil {
			return fmt.Errorf("failed to update blockers: %w", err)
		}
	}

	if _, err := db.SaveDirty(); err != nil {
		return fmt.Errorf("failed to save changes: %w", err)
	}

	// Trigger a repair to update master lists
	if err := runRepair(w, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
		return err
	}

	fmt.Fprintf(w, "✓ Task %s status set to %s\n", task.ShortID(taskID), status)
	return nil
}
