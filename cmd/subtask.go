package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
)

var subtaskCmd = &cobra.Command{
	Use:   "subtask",
	Short: "Manage task subtasks",
}

var subtaskReorderCmd = &cobra.Command{
	Use:   "reorder <parent-task-id> <old-index> <new-index>",
	Short: "Reorder a parent's subtasks",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldIdx, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid old index: %w", err)
		}
		newIdx, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid new index: %w", err)
		}
		return runSubtaskReorder(cmd.OutOrStdout(), projectName, args[0], oldIdx, newIdx)
	},
}

func init() {
	rootCmd.AddCommand(subtaskCmd)
	subtaskCmd.AddCommand(subtaskReorderCmd)
}

func runSubtaskReorder(w io.Writer, projectName, inputParentID string, oldIdx, newIdx int) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	_, parentID, err := db.GetResolved(inputParentID)
	if err != nil {
		return err
	}

	if err := db.ReorderSubtask(parentID, oldIdx, newIdx); err != nil {
		return err
	}

	if _, err := db.SaveDirty(); err != nil {
		return fmt.Errorf("failed to write reordered subtasks: %w", err)
	}

	if err := task.GenerateMasterLists(db.GetAll(), paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile); err != nil {
		return fmt.Errorf("failed to update master lists: %w", err)
	}

	fmt.Fprintf(w, "✓ Reordered subtask %d to %d in task %s\n", oldIdx, newIdx, task.ShortID(parentID))
	return nil
}
