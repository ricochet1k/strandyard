/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ricochet1k/memmd/pkg/task"
	"github.com/spf13/cobra"
)

var nextRole string

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Print the next free task",
	Long: `Print the next free task from the free-tasks list.
Shows the task's role (from metadata or first TODO) and the task content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNext(nextRole)
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
	nextCmd.Flags().StringVar(&nextRole, "role", "", "optional: filter tasks by role")
}

func runNext(roleFilter string) error {
	// Ensure free-tasks exists; if not, run validate to generate it
	freePath := "tasks/free-tasks.md"
	if _, err := os.Stat(freePath); os.IsNotExist(err) {
		// Run validate to generate lists
		if err := runValidate("tasks", "tasks/root-tasks.md", freePath, "text"); err != nil {
			return fmt.Errorf("unable to generate master lists: %w", err)
		}
	}

	// Read free tasks list
	data, err := os.ReadFile(freePath)
	if err != nil {
		return fmt.Errorf("unable to read %s: %w", freePath, err)
	}

	// Parse task paths from free-tasks.md
	lines := strings.Split(string(data), "\n")
	candidates := []string{}
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "- ") {
			candidates = append(candidates, strings.TrimSpace(strings.TrimPrefix(l, "- ")))
		}
	}

	if len(candidates) == 0 {
		fmt.Println("No free tasks found")
		return nil
	}

	// Parse tasks using the task library
	parser := task.NewParser()
	var selectedTask *task.Task

	for _, candidatePath := range candidates {
		t, err := parser.ParseFile(candidatePath)
		if err != nil {
			// Skip tasks that fail to parse
			continue
		}

		// If role filter specified, check if it matches
		if roleFilter != "" {
			taskRole := t.GetEffectiveRole()
			if taskRole != roleFilter {
				continue
			}
		}

		// Select first matching task
		selectedTask = t
		break
	}

	if selectedTask == nil {
		if roleFilter != "" {
			fmt.Printf("No free tasks found for role: %s\n", roleFilter)
		} else {
			fmt.Println("No free tasks found")
		}
		return nil
	}

	// Print role
	role := selectedTask.GetEffectiveRole()
	if role != "" {
		fmt.Printf("Role: %s\n\n", role)
		rolePath := fmt.Sprintf("roles/%s.md", role)
		roleData, err := os.ReadFile(rolePath)
		if err != nil {
			fmt.Printf("Role file not found: %s\n\n", rolePath)
		} else {
			fmt.Print(string(roleData))
			fmt.Print("\n\n")
		}
	} else {
		fmt.Println("Role: (none)\n")
	}

	// Print task content
	fmt.Print(selectedTask.Content)

	return nil
}
