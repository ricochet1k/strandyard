/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ricochet1k/strandyard/pkg/role"
	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
)

var (
	editTitle    string
	editRole     string
	editPriority string
	editParent   string
	editBlockers []string
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit <task-id>",
	Short: "Edit a task's metadata and description",
	Long: `Edit a task's metadata (title, role, priority, parent, blockers) and description.
The description is read from standard input. It is recommended to use heredocs to provide the description.
Example:
  strand edit T3k7x --priority high <<EOF
  # New Title
  New description here.
  EOF

If only metadata flags are provided and stdin is a terminal, the description remains unchanged.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		body, err := readStdin()
		if err != nil {
			return err
		}
		return runEdit(cmd, taskID, body)
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringVarP(&editTitle, "title", "t", "", "task title")
	editCmd.Flags().StringVarP(&editRole, "role", "r", "", "role responsible for the task")
	editCmd.Flags().StringVarP(&editParent, "parent", "p", "", "parent task ID")
	editCmd.Flags().StringVar(&editPriority, "priority", "", "priority: high, medium, or low")
}

func runEdit(cmd *cobra.Command, taskID, newBody string) error {
	w := cmd.OutOrStdout()
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	// Load all tasks to find the one we want and resolve IDs
	parser := task.NewParser()
	tasks, err := parser.LoadTasks(paths.TasksDir)
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	resolvedID, err := task.ResolveTaskID(tasks, taskID)
	if err != nil {
		return err
	}
	t, ok := tasks[resolvedID]
	if !ok {
		return fmt.Errorf("task not found: %s", resolvedID)
	}

	// Track if any changes were made
	changes := false

	// Update flags if they were changed
	if cmd.Flags().Changed("role") {
		rName := strings.TrimSpace(editRole)
		if err := role.ValidateRole(paths.RolesDir, rName); err != nil {
			return err
		}
		t.Meta.Role = rName
		t.MarkDirty()
		changes = true
	}

	if cmd.Flags().Changed("priority") {
		priority := task.NormalizePriority(editPriority)
		if !task.IsValidPriority(priority) {
			return fmt.Errorf("invalid priority: %s", editPriority)
		}
		t.Meta.Priority = priority
		t.MarkDirty()
		changes = true
	}

	if cmd.Flags().Changed("parent") {
		parent := strings.TrimSpace(editParent)
		if parent != "" {
			resolvedParent, err := task.ResolveTaskID(tasks, parent)
			if err != nil {
				return fmt.Errorf("parent task %s does not exist: %w", parent, err)
			}
			t.Meta.Parent = resolvedParent
		} else {
			t.Meta.Parent = ""
		}
		t.MarkDirty()
		changes = true
	}

	if cmd.Flags().Changed("blocker") {
		panic("Unimplemented: TODO: Must update blockers AND child blocks list")
		// blockers, err := resolveTaskIDs(tasks, normalizeTaskIDs(editBlockers))
		// if err != nil {
		// 	return err
		// }
		// t.Meta.Blockers = blockers
		// t.MarkDirty()
		// changes = true
	}

	// Body from stdin
	if isStdinRedirected() {
		t.SetBody(newBody)
		// If title flag not provided, try extracting from stdin
		if !cmd.Flags().Changed("title") {
			if title := task.ExtractTitle(newBody); title != "" {
				t.SetTitle(title)
			}
		}
		changes = true
	}

	// Title update (must be done AFTER SetBody if both are used, so title flag wins)
	if cmd.Flags().Changed("title") {
		t.SetTitle(editTitle)
		changes = true
	}

	if changes {
		if err := t.Write(); err != nil {
			return fmt.Errorf("failed to write task file: %w", err)
		}
		fmt.Fprintf(w, "✓ Task %s updated\n", task.ShortID(resolvedID))

		// Update parent TODOs if parent changed or if it's a child
		if cmd.Flags().Changed("parent") || t.Meta.Parent != "" {
			if t.Meta.Parent != "" {
				if _, err := task.UpdateParentTodoEntries(tasks, t.Meta.Parent); err != nil {
					return fmt.Errorf("failed to update parent task TODO entries: %w", err)
				}
			}
			if _, err := task.WriteDirtyTasks(tasks); err != nil {
				return fmt.Errorf("failed to write dirty tasks: %w", err)
			}
		}

		// TODO: This should not be necessary
		if err := runRepair(w, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(w, "No changes to task %s\n", task.ShortID(resolvedID))
	}

	return nil
}

func isStdinRedirected() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) == 0
}
