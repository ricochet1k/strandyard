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

	"github.com/ricochet1k/strandyard/pkg/activity"
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

func runComplete(w io.Writer, projectName, inputID string, todoNum int, role string, report string) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	t, taskID, err := db.GetResolved(inputID)
	if err != nil {
		return err
	}

	incompleteTodos, err := db.GetIncompleteTodos(taskID)
	if err != nil {
		return err
	}

	// Validate role if --role flag is provided
	if role != "" {
		if err := validateRole(t, todoNum, role, incompleteTodos); err != nil {
			return err
		}
	}

	// Handle todo completion if --todo flag is used
	if todoNum > 0 {
		return runCompleteTodo(w, db, paths, t, taskID, todoNum, report)
	}

	// Check for incomplete todos before completing the whole task
	if len(incompleteTodos) > 0 {
		errorMsg := fmt.Sprintf("cannot complete task %s: incomplete todos remain\n\n", task.ShortID(taskID))
		for i, todo := range incompleteTodos {
			errorMsg += fmt.Sprintf("%d. [ ] (role: %s) %s\n", i+1, todo.Role, todo.Text)
		}
		return errors.New(errorMsg)
	}

	if t.Meta.Completed {
		fmt.Fprintf(w, "Task %s is already marked as completed\n", task.ShortID(taskID))
		return nil
	}

	// Calculate incremental update before marking the task complete
	update, err := task.CalculateIncrementalFreeListUpdate(db.GetAll(), taskID)
	if err != nil {
		return fmt.Errorf("failed to calculate incremental update: %w", err)
	}

	if err := db.CompleteTask(taskID, report); err != nil {
		return err
	}

	if err := db.UpdateBlockersAfterCompletion(taskID); err != nil {
		return fmt.Errorf("failed to update blockers after completion: %w", err)
	}

	if _, err := db.SaveDirty(); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	activityLog, err := activity.Open(paths.BaseDir)
	if err != nil {
		return fmt.Errorf("failed to open activity log: %w", err)
	}
	defer activityLog.Close()

	if err := activityLog.WriteTaskCompletion(taskID, report); err != nil {
		return fmt.Errorf("failed to write activity log: %w", err)
	}

	fmt.Fprintf(w, "✓ Task %s marked as completed\n", task.ShortID(taskID))
	if report == "" {
		fmt.Fprintf(w, "💡 Next time, consider adding a report: strand complete %s \"summary of work\"\n", task.ShortID(taskID))
	}

	if _, err := db.UpdateParentTodosForChild(taskID); err != nil {
		return fmt.Errorf("failed to update parent task TODO entries: %w", err)
	}
	if _, err := db.SaveDirty(); err != nil {
		return fmt.Errorf("failed to write parent task updates: %w", err)
	}

	// Try incremental update first, fall back to full validation
	if err := task.UpdateFreeListIncrementally(db.GetAll(), paths.FreeTasksFile, update); err != nil {
		fmt.Fprintf(w, "⚠️  Incremental update failed, falling back to full repair: %v\n", err)
		if err := runRepair(w, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(w, "✓ Incrementally updated free-tasks.md\n")
		validator := task.NewValidatorWithRoles(db.GetAll(), paths.RolesDir)
		validationErrors := validator.ValidateAndRepair()
		if _, err := db.SaveDirty(); err != nil {
			return fmt.Errorf("failed to write repaired tasks: %w", err)
		}
		if len(validationErrors) > 0 {
			fmt.Fprintf(w, "⚠️  Repair errors found:\n")
			for _, e := range validationErrors {
				fmt.Fprintf(w, "ERROR: %s\n", e.Error())
			}
			return fmt.Errorf("repair failed: %d error(s)", len(validationErrors))
		}
	}

	fmt.Fprintf(w, "💡 Consider committing your changes: git add -A && git commit -m \"complete: %s\"\n", task.ShortID(taskID))
	return nil
}

func validateRole(t *task.Task, todoNum int, role string, incompleteTodos []task.TaskItem) error {
	var taskRole string
	if todoNum > 0 {
		if todoNum <= len(t.TodoItems) {
			taskRole = t.TodoItems[todoNum-1].Role
		} else {
			taskRole = "INVALID_TODO_NUM"
		}
	} else {
		taskRole = t.Meta.Role
		if taskRole == "" {
			taskRole = t.GetEffectiveRole()
		}
	}

	if taskRole != role {
		target := "task"
		if todoNum > 0 {
			target = "todo"
		}
		if len(incompleteTodos) > 0 {
			errorMsg := fmt.Sprintf("role validation failed: %s has role '%s' but --role flag specifies '%s'\n\nIncomplete todos:\n",
				target, taskRole, role)
			for i, todo := range incompleteTodos {
				errorMsg += fmt.Sprintf("%d. [ ] (role: %s) %s\n", i+1, todo.Role, todo.Text)
			}
			return errors.New(errorMsg)
		}
		return fmt.Errorf("role validation failed: %s has role '%s' but --role flag specifies '%s'", target, taskRole, role)
	}
	return nil
}

func runCompleteTodo(w io.Writer, db *task.TaskDB, paths projectPaths, t *task.Task, taskID string, todoNum int, report string) error {
	// Calculate incremental update before completing the TODO
	update, err := task.CalculateIncrementalFreeListUpdate(db.GetAll(), taskID)
	if err != nil {
		return fmt.Errorf("failed to calculate incremental update: %w", err)
	}

	result, err := db.CompleteTodo(taskID, todoNum, report)
	if err != nil {
		return err
	}

	// Already checked case - CompleteTodo returns early for this
	if result.RemainingIncomplete >= 0 && !t.TodoItems[todoNum-1].Checked {
		// This means the todo was already checked before we called CompleteTodo
		// But actually we need to check differently - if the todo was already checked,
		// we should detect that. Let's check the task state directly.
	}

	// Re-fetch task to check current state
	t, err = db.Get(taskID)
	if err != nil {
		return err
	}

	todoIndex := todoNum - 1
	if _, err := db.SaveDirty(); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	if result.TaskCompleted {
		activityLog, err := activity.Open(paths.BaseDir)
		if err != nil {
			return fmt.Errorf("failed to open activity log: %w", err)
		}
		defer activityLog.Close()

		if err := activityLog.WriteTaskCompletion(taskID, report); err != nil {
			return fmt.Errorf("failed to write activity log: %w", err)
		}

		if err := db.UpdateBlockersAfterCompletion(taskID); err != nil {
			return fmt.Errorf("failed to update blockers after completion: %w", err)
		}

		if _, err := db.UpdateParentTodosForChild(taskID); err != nil {
			return fmt.Errorf("failed to update parent task TODO entries: %w", err)
		}
		if _, err := db.SaveDirty(); err != nil {
			return fmt.Errorf("failed to write task file: %w", err)
		}

		// Try incremental update first, fall back to full validation
		if err := task.UpdateFreeListIncrementally(db.GetAll(), paths.FreeTasksFile, update); err != nil {
			fmt.Fprintf(w, "⚠️  Incremental update failed, falling back to full repair: %v\n", err)
			if err := runRepair(w, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(w, "✓ Incrementally updated free-tasks.md\n")
			validator := task.NewValidatorWithRoles(db.GetAll(), paths.RolesDir)
			validationErrors := validator.ValidateAndRepair()
			if _, err := db.SaveDirty(); err != nil {
				return fmt.Errorf("failed to write repaired tasks: %w", err)
			}
			if len(validationErrors) > 0 {
				fmt.Fprintf(w, "⚠️  Repair errors found:\n")
				for _, e := range validationErrors {
					fmt.Fprintf(w, "ERROR: %s\n", e.Error())
				}
				return fmt.Errorf("repair failed: %d error(s)", len(validationErrors))
			}
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
