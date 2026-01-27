/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ricochet1k/memmd/pkg/task"
	"github.com/spf13/cobra"
)

var (
	validatePath  string
	validateRoots string
	validateFree  string
	validateFmt   string
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate task tree and regenerate master lists",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runValidate(validatePath, validateRoots, validateFree, validateFmt)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringVar(&validatePath, "path", "tasks", "path to tasks directory")
	validateCmd.Flags().StringVar(&validateRoots, "roots", "tasks/root-tasks.md", "path to write root tasks list")
	validateCmd.Flags().StringVar(&validateFree, "free", "tasks/free-tasks.md", "path to write free tasks list")
	validateCmd.Flags().StringVar(&validateFmt, "format", "text", "output format for validation errors: text|json")
}

func runValidate(tasksRoot, rootsFile, freeFile, outFormat string) error {
	// Parse all tasks
	parser := task.NewParser()
	tasks, err := parser.LoadTasks(tasksRoot)
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	// Validate tasks
	validator := task.NewValidator(tasks)
	errors := validator.Validate()

	// Generate master lists
	if err := task.GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		return fmt.Errorf("failed to generate master lists: %w", err)
	}

	// Report errors
	if len(errors) > 0 {
		if outFormat == "json" {
			errMsgs := make([]string, len(errors))
			for i, e := range errors {
				errMsgs[i] = e.Error()
			}
			b, _ := json.MarshalIndent(map[string]interface{}{"errors": errMsgs}, "", "  ")
			fmt.Println(string(b))
		} else {
			for _, e := range errors {
				fmt.Println("ERROR:", e.Error())
			}
		}
		return fmt.Errorf("validation failed: %d error(s)", len(errors))
	}

	// Success output
	if outFormat == "json" {
		// Get lists for output
		roots := []string{}
		free := []string{}
		for id, t := range tasks {
			if t.Meta.Parent == "" {
				roots = append(roots, id)
			}
			if len(t.Meta.Blockers) == 0 {
				free = append(free, id)
			}
		}
		b, _ := json.MarshalIndent(map[string]interface{}{"roots": roots, "free": free}, "", "  ")
		fmt.Println(string(b))
	} else {
		fmt.Println("validate: ok")
		fmt.Printf("Validated %d tasks\n", len(tasks))
		fmt.Printf("Master lists updated: %s, %s\n", rootsFile, freeFile)
	}

	return nil
}
