/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"strconv"

	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
)

// todoCmd represents the todo command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage task TODO items",
	Long:  `Add, remove, edit, and check/uncheck TODO items within a task.`,
}

var todoAddCmd = &cobra.Command{
	Use:   "add <task-id> <text>",
	Short: "Add a new TODO item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTodoAdd(cmd.OutOrStdout(), projectName, args[0], args[1])
	},
}

var todoRemoveCmd = &cobra.Command{
	Use:   "remove <task-id> <index>",
	Short: "Remove a TODO item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		idx, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid index: %w", err)
		}
		return runTodoRemove(cmd.OutOrStdout(), projectName, args[0], idx)
	},
}

var todoEditCmd = &cobra.Command{
	Use:   "edit <task-id> <index> <text>",
	Short: "Edit a TODO item's text",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		idx, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid index: %w", err)
		}
		return runTodoEdit(cmd.OutOrStdout(), projectName, args[0], idx, args[2])
	},
}

var todoCheckCmd = &cobra.Command{
	Use:   "check <task-id> <index> [report]",
	Short: "Check off a TODO item",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		idx, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid index: %w", err)
		}
		report := ""
		if len(args) > 2 {
			report = args[2]
		}
		return runTodoCheck(cmd.OutOrStdout(), projectName, args[0], idx, report)
	},
}

var todoUncheckCmd = &cobra.Command{
	Use:   "uncheck <task-id> <index>",
	Short: "Uncheck a TODO item",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		idx, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid index: %w", err)
		}
		return runTodoUncheck(cmd.OutOrStdout(), projectName, args[0], idx)
	},
}

var todoReorderCmd = &cobra.Command{
	Use:   "reorder <task-id> <old-index> <new-index>",
	Short: "Reorder a TODO item",
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
		return runTodoReorder(cmd.OutOrStdout(), projectName, args[0], oldIdx, newIdx)
	},
}

var todoListCmd = &cobra.Command{
	Use:   "list <task-id>",
	Short: "List TODO items for a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTodoList(cmd.OutOrStdout(), projectName, args[0])
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)
	todoCmd.AddCommand(todoAddCmd)
	todoCmd.AddCommand(todoRemoveCmd)
	todoCmd.AddCommand(todoEditCmd)
	todoCmd.AddCommand(todoCheckCmd)
	todoCmd.AddCommand(todoUncheckCmd)
	todoCmd.AddCommand(todoReorderCmd)
	todoCmd.AddCommand(todoListCmd)
}

func runTodoAdd(w io.Writer, projectName, inputID, text string) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	_, taskID, err := db.GetResolved(inputID)
	if err != nil {
		return err
	}

	if err := db.AddTodo(taskID, text); err != nil {
		return err
	}

	if _, err := db.SaveDirty(); err != nil {
		return err
	}

	fmt.Fprintf(w, "✓ Added todo to task %s\n", task.ShortID(taskID))
	return nil
}

func runTodoRemove(w io.Writer, projectName, inputID string, idx int) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	_, taskID, err := db.GetResolved(inputID)
	if err != nil {
		return err
	}

	if err := db.RemoveTodo(taskID, idx); err != nil {
		return err
	}

	if _, err := db.SaveDirty(); err != nil {
		return err
	}

	fmt.Fprintf(w, "✓ Removed todo %d from task %s\n", idx, task.ShortID(taskID))
	return nil
}

func runTodoEdit(w io.Writer, projectName, inputID string, idx int, text string) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	_, taskID, err := db.GetResolved(inputID)
	if err != nil {
		return err
	}

	if err := db.EditTodo(taskID, idx, text); err != nil {
		return err
	}

	if _, err := db.SaveDirty(); err != nil {
		return err
	}

	fmt.Fprintf(w, "✓ Edited todo %d in task %s\n", idx, task.ShortID(taskID))
	return nil
}

func runTodoCheck(w io.Writer, projectName, inputID string, idx int, report string) error {
	// Re-use runComplete logic
	return runComplete(w, projectName, inputID, idx, "", report)
}

func runTodoUncheck(w io.Writer, projectName, inputID string, idx int) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	_, taskID, err := db.GetResolved(inputID)
	if err != nil {
		return err
	}

	if err := db.UncheckTodo(taskID, idx); err != nil {
		return err
	}

	if _, err := db.SaveDirty(); err != nil {
		return err
	}

	fmt.Fprintf(w, "✓ Unchecked todo %d in task %s\n", idx, task.ShortID(taskID))
	return nil
}

func runTodoReorder(w io.Writer, projectName, inputID string, oldIdx, newIdx int) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	_, taskID, err := db.GetResolved(inputID)
	if err != nil {
		return err
	}

	if err := db.ReorderTodo(taskID, oldIdx, newIdx); err != nil {
		return err
	}

	if _, err := db.SaveDirty(); err != nil {
		return err
	}

	fmt.Fprintf(w, "✓ Reordered todo %d to %d in task %s\n", oldIdx, newIdx, task.ShortID(taskID))
	return nil
}

func runTodoList(w io.Writer, projectName, inputID string) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	t, taskID, err := db.GetResolved(inputID)
	if err != nil {
		return err
	}

	if len(t.TodoItems) == 0 {
		fmt.Fprintf(w, "Task %s has no todo items\n", task.ShortID(taskID))
		return nil
	}

	fmt.Fprintf(w, "TODOs for task %s:\n", task.ShortID(taskID))
	for i, item := range t.TodoItems {
		status := "[ ]"
		if item.Checked {
			status = "[x]"
		}
		roleStr := ""
		if item.Role != "" {
			roleStr = fmt.Sprintf(" (role: %s)", item.Role)
		}
		fmt.Fprintf(w, "%d. %s%s %s\n", i+1, status, roleStr, item.Text)
	}

	return nil
}
