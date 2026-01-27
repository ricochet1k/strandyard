/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ricochet1k/memmd/pkg/task"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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
		return runComplete(taskID)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}

func runComplete(taskID string) error {
	// Load all tasks to find the one we want
	parser := task.NewParser()
	tasks, err := parser.LoadTasks("tasks")
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	// Find the task by ID
	t, exists := tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	taskFile := t.FilePath

	// Check if already completed
	if t.Meta.Completed {
		fmt.Printf("Task %s is already marked as completed\n", taskID)
		return nil
	}

	// Read the file content
	content, readErr := os.ReadFile(taskFile)
	if readErr != nil {
		return fmt.Errorf("failed to read task file: %w", readErr)
	}

	// Update metadata
	t.Meta.Completed = true
	t.Meta.DateEdited = time.Now()

	// Split frontmatter and body
	contentStr := string(content)
	parts := strings.SplitN(contentStr, "---", 3)
	if len(parts) < 3 {
		return fmt.Errorf("invalid task file format: missing frontmatter delimiters")
	}

	// Serialize updated frontmatter
	frontmatterBytes, marshalErr := yaml.Marshal(&t.Meta)
	if marshalErr != nil {
		return fmt.Errorf("failed to serialize frontmatter: %w", marshalErr)
	}

	// Reconstruct file
	var newContent strings.Builder
	newContent.WriteString("---\n")
	newContent.Write(frontmatterBytes)
	newContent.WriteString("---")
	newContent.WriteString(parts[2])

	// Write back to file
	if writeErr := os.WriteFile(taskFile, []byte(newContent.String()), 0644); writeErr != nil {
		return fmt.Errorf("failed to write task file: %w", writeErr)
	}

	fmt.Printf("✓ Task %s marked as completed\n", taskID)
	fmt.Println("\nRun 'memmd validate' to update master lists")

	return nil
}
