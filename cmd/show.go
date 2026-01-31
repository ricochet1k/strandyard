package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show <task-id>",
	Short: "Print the full contents of a task",
	Long: `Print the full contents of a task by ID, short ID, or any valid prefix.
The output includes the complete markdown file content including frontmatter.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		return runShow(cmd.OutOrStdout(), projectName, taskID)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(w io.Writer, projectName, taskID string) error {
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

	// Read and print the full file contents
	content, err := os.ReadFile(t.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read task file: %w", err)
	}

	fmt.Fprint(w, string(content))
	return nil
}
