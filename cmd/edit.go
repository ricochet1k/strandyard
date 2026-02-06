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
	editBlocks   []string
	editEvery    []string
	editStatus   string
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
	editCmd.Flags().StringSliceVarP(&editBlockers, "blocker", "b", nil, "blocker task ID(s); replaces existing blockers")
	editCmd.Flags().StringSliceVar(&editBlocks, "blocks", nil, "task ID(s) this task blocks; replaces existing blocks")
	editCmd.Flags().StringSliceVar(&editEvery, "every", nil, `recurrence rule: "<amount> <metric> [from <anchor>]" (repeatable)
metrics: days, weeks, months, commits, lines_changed, tasks_completed
examples: "10 days", "50 commits from HEAD", "20 tasks_completed from T1a1a"`)
	editCmd.Flags().StringVarP(&editStatus, "status", "s", "", fmt.Sprintf("task status: %s", task.FormatStatusListForUser()))
}

func runEdit(cmd *cobra.Command, inputID, newBody string) error {
	w := cmd.OutOrStdout()
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	db := task.NewTaskDB(paths.TasksDir)
	t, taskID, err := db.GetResolved(inputID)
	if err != nil {
		return err
	}

	if cmd.Flags().Changed("role") {
		rName := strings.TrimSpace(editRole)
		if err := role.ValidateRole(paths.RolesDir, rName); err != nil {
			return err
		}
		if err := db.SetRole(taskID, rName); err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("priority") {
		if err := db.SetPriority(taskID, editPriority); err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("parent") {
		parent := strings.TrimSpace(editParent)
		if parent != "" {
			resolvedParent, err := db.ResolveID(parent)
			if err != nil {
				return fmt.Errorf("parent task %s does not exist: %w", parent, err)
			}
			if err := db.SetParent(taskID, resolvedParent); err != nil {
				return err
			}
		} else {
			if err := db.ClearParent(taskID); err != nil {
				return err
			}
		}
	}

	if cmd.Flags().Changed("blocker") {
		newBlockers, err := db.ResolveIDs(normalizeTaskIDs(editBlockers))
		if err != nil {
			return err
		}

		currentBlockers := make(map[string]bool)
		for _, b := range t.Meta.Blockers {
			currentBlockers[b] = true
		}
		newBlockerSet := make(map[string]bool)
		for _, b := range newBlockers {
			newBlockerSet[b] = true
		}

		for _, b := range t.Meta.Blockers {
			if !newBlockerSet[b] {
				if err := db.RemoveBlocker(taskID, b); err != nil {
					return fmt.Errorf("failed to remove blocker %s: %w", b, err)
				}
			}
		}
		for _, b := range newBlockers {
			if !currentBlockers[b] {
				if err := db.AddBlocker(taskID, b); err != nil {
					return fmt.Errorf("failed to add blocker %s: %w", b, err)
				}
			}
		}
	}

	if cmd.Flags().Changed("blocks") {
		newBlocks, err := db.ResolveIDs(normalizeTaskIDs(editBlocks))
		if err != nil {
			return err
		}

		currentBlocks := make(map[string]bool)
		for _, b := range t.Meta.Blocks {
			currentBlocks[b] = true
		}
		newBlocksSet := make(map[string]bool)
		for _, b := range newBlocks {
			newBlocksSet[b] = true
		}

		for _, b := range t.Meta.Blocks {
			if !newBlocksSet[b] {
				if err := db.RemoveBlocked(taskID, b); err != nil {
					return fmt.Errorf("failed to remove blocked %s: %w", b, err)
				}
			}
		}
		for _, b := range newBlocks {
			if !currentBlocks[b] {
				if err := db.AddBlocked(taskID, b); err != nil {
					return fmt.Errorf("failed to add blocked %s: %w", b, err)
				}
			}
		}
	}

	if cmd.Flags().Changed("every") {
		resolvedEvery, err := validateEvery(editEvery, paths.BaseDir, db.GetAll())
		if err != nil {
			os.Exit(2)
		}
		t.Meta.Every = resolvedEvery
		t.MarkDirty()
	}

	if isStdinRedirected() {
		if err := db.SetBody(taskID, newBody); err != nil {
			return err
		}
		if !cmd.Flags().Changed("title") {
			if title := task.ExtractTitle(newBody); title != "" {
				if err := db.SetTitle(taskID, title); err != nil {
					return err
				}
			}
		}
	}

	if cmd.Flags().Changed("status") {
		if err := db.SetStatus(taskID, editStatus); err != nil {
			return err
		}
	}

	if cmd.Flags().Changed("title") {
		if err := db.SetTitle(taskID, editTitle); err != nil {
			return err
		}
	}

	if t.Dirty {
		if _, err := db.SaveDirty(); err != nil {
			return fmt.Errorf("failed to write task file: %w", err)
		}
		fmt.Fprintf(w, "✓ Task %s updated\n", task.ShortID(taskID))

		// Refetch task to get updated parent
		t, err = db.Get(taskID)
		if err != nil {
			return err
		}

		if cmd.Flags().Changed("parent") || t.Meta.Parent != "" {
			if _, err := db.UpdateParentTodosForChild(taskID); err != nil {
				return fmt.Errorf("failed to update parent task TODO entries: %w", err)
			}
			if _, err := db.SaveDirty(); err != nil {
				return fmt.Errorf("failed to write dirty tasks: %w", err)
			}
		}

		// TODO: This should not be necessary
		if err := runRepair(w, paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, "text"); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(w, "No changes to task %s\n", task.ShortID(taskID))
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
