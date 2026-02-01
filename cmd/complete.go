/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete <task-id> [report]",
	Short: "Mark a task as completed",
	Long: `Mark a task as completed by setting completed: true in the frontmatter.
Also updates the date_edited field to the current time.
Use --todo to check off a specific todo item instead of completing the entire task.
Use --role to validate that the current role matches the task role.
A report can be provided as a second argument or via stdin (e.g. heredoc).`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		var report string
		if len(args) > 1 {
			report = args[1]
		} else {
			// Check if stdin is a pipe/redirect (not a terminal)
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				bytes, err := io.ReadAll(os.Stdin)
				if err == nil {
					report = string(bytes)
				}
			}
		}
		report = strings.TrimSpace(report)

		todoNum, _ := cmd.Flags().GetInt("todo")
		role, _ := cmd.Flags().GetString("role")
		return runComplete(cmd.OutOrStdout(), projectName, taskID, todoNum, role, report)
	},
}

func init() {
	completeCmd.Flags().Int("todo", 0, "Check off a specific todo item (1-based index)")
	completeCmd.Flags().String("role", "", "Validate role matches task role")
	rootCmd.AddCommand(completeCmd)
}

func runComplete(w io.Writer, projectName, taskID string, todoNum int, role string, report string) error {
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
		return fmt.Errorf("task not found: %s", task.ShortID(taskID))
	}

	// Check if there are incomplete todos
	var incompleteTodos []task.TaskItem
	for _, todo := range t.TodoItems {
		if !todo.Checked {
			incompleteTodos = append(incompleteTodos, todo)
		}
	}

	// Validate role if --role flag is provided
	if role != "" {
		var taskRole string

		if todoNum > 0 {
			// For --todo, validate against specific todo role
			if todoNum <= len(t.TodoItems) {
				taskRole = t.TodoItems[todoNum-1].Role
			} else {
				taskRole = "INVALID_TODO_NUM"
			}
		} else {
			// For task completion, validate against task role
			taskRole = t.Meta.Role
			if taskRole == "" {
				taskRole = t.GetEffectiveRole()
			}
		}

		if taskRole != role {
			// Include all todos in error message if role doesn't match
			var errorMsg string
			if len(incompleteTodos) > 0 {
				target := "task"
				if todoNum > 0 {
					target = "todo"
				}
				errorMsg = fmt.Sprintf("role validation failed: %s has role '%s' but --role flag specifies '%s'\n\nIncomplete todos:\n",
					target, taskRole, role)
				for i, todo := range incompleteTodos {
					errorMsg += fmt.Sprintf("%d. [ ] (role: %s) %s\n", i+1, todo.Role, todo.Text)
				}
			} else {
				target := "task"
				if todoNum > 0 {
					target = "todo"
				}
				errorMsg = fmt.Sprintf("role validation failed: %s has role '%s' but --role flag specifies '%s'",
					target, taskRole, role)
			}
			return errors.New(errorMsg)
		}
	}

	// Handle todo completion if --todo flag is used
	if todoNum > 0 {
		if todoNum <= 0 || todoNum > len(t.TodoItems) {
			return fmt.Errorf("invalid todo number %d, task has %d todo items", todoNum, len(t.TodoItems))
		}

		todoIndex := todoNum - 1
		if t.TodoItems[todoIndex].Checked {
			fmt.Fprintf(w, "Todo item %d is already checked off\n", todoNum)
			return nil
		}

		// Mark the todo item as checked
		t.TodoItems[todoIndex].Checked = true
		if report != "" {
			t.TodoItems[todoIndex].Report = report
		}
		t.MarkDirty()

		if err := t.Write(); err != nil {
			return fmt.Errorf("failed to write task file: %w", err)
		}

		// Check if this was the last incomplete todo
		remainingIncomplete := 0
		for _, todo := range t.TodoItems {
			if !todo.Checked {
				remainingIncomplete++
			}
		}

		if remainingIncomplete == 0 {
			// Last todo completed, mark task as completed
			t.Meta.Completed = true
			t.MarkDirty()

			if err := t.Write(); err != nil {
				return fmt.Errorf("failed to write task file: %w", err)
			}

			fmt.Fprintf(w, "✓ Todo item %d checked off in task %s (last todo - task marked complete)\n", todoNum, task.ShortID(taskID))
			if report == "" {
				fmt.Fprintf(w, "💡 Next time, consider adding a report: strand complete %s --todo %d \"summary of work\"\n", task.ShortID(taskID), todoNum)
			}
			fmt.Fprintf(w, "💡 Consider committing your changes: git add -A && git commit -m \"complete: %s\"\n", task.ShortID(taskID))
			return nil
		}

		fmt.Fprintf(w, "- [x] %v\n", t.TodoItems[todoIndex].Text)
		fmt.Fprintf(w, "✓ Todo item %d checked off in task %s\n", todoNum, task.ShortID(taskID))
		if report == "" {
			fmt.Fprintf(w, "💡 Next time, consider adding a report: strand complete %s --todo %d \"summary of work\"\n", task.ShortID(taskID), todoNum)
		}
		fmt.Fprintf(w, "💡 Consider committing your changes: git add -A && git commit -m \"%v (%v) check off %v\"\n", t.Title(), task.ShortID(taskID), t.TodoItems[todoIndex].Text)
		return nil
	}

	// If not using --todo flag, check for incomplete todos
	if todoNum == 0 {
		if len(incompleteTodos) > 0 {
			errorMsg := fmt.Sprintf("cannot complete task %s: incomplete todos remain\n\n", task.ShortID(taskID))
			for i, todo := range incompleteTodos {
				errorMsg += fmt.Sprintf("%d. [ ] (role: %s) %s\n", i+1, todo.Role, todo.Text)
			}
			return errors.New(errorMsg)
		}
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
	if report != "" {
		if t.OtherContent != "" {
			t.OtherContent += "\n\n"
		}
		t.OtherContent += "## Completion Report\n" + report
	}
	t.MarkDirty()

	if err := t.Write(); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	fmt.Fprintf(w, "✓ Task %s marked as completed\n", task.ShortID(taskID))
	if report == "" {
		fmt.Fprintf(w, "💡 Next time, consider adding a report: strand complete %s \"summary of work\"\n", task.ShortID(taskID))
	}

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
		errors := validator.ValidateAndRepair()
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
