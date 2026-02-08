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
		return runClaim(cmd.OutOrStdout(), args[0])
	},
}

var claimCmd = &cobra.Command{
	Use:   "claim <task-id>",
	Short: "Claim a task by marking it in progress",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runClaim(cmd.OutOrStdout(), args[0])
	},
}

func init() {
	rootCmd.AddCommand(cancelCmd)
	rootCmd.AddCommand(markDuplicateCmd)
	rootCmd.AddCommand(markInProgressCmd)
	rootCmd.AddCommand(claimCmd)
}

func runClaim(w io.Writer, inputID string) error {
	return runSetStatus(w, inputID, task.StatusInProgress, "")
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

	if status == task.StatusInProgress {
		if err := db.ClaimTask(taskID); err != nil {
			return err
		}
	} else {
		if err := db.SetStatusWithReport(taskID, status, report); err != nil {
			return err
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
