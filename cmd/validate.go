/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ricochet1k/memmd/pkg/task"
	"github.com/spf13/cobra"
)

var (
	repairPath  string
	repairRoots string
	repairFree  string
	repairFmt   string
)

// repairCmd represents the repair command
var repairCmd = &cobra.Command{
	Use:   "repair",
	Short: "Repair task tree and regenerate master lists",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := resolveProjectPaths(projectName)
		if err != nil {
			return err
		}
		tasksRoot := repairPath
		rootsFile := repairRoots
		freeFile := repairFree
		if !cmd.Flags().Changed("path") {
			tasksRoot = paths.TasksDir
		}
		if !cmd.Flags().Changed("roots") {
			rootsFile = paths.RootTasksFile
		}
		if !cmd.Flags().Changed("free") {
			freeFile = paths.FreeTasksFile
		}
		return runRepair(tasksRoot, rootsFile, freeFile, repairFmt)
	},
}

func init() {
	rootCmd.AddCommand(repairCmd)
	repairCmd.Flags().StringVar(&repairPath, "path", "tasks", "path to tasks directory")
	repairCmd.Flags().StringVar(&repairRoots, "roots", "tasks/root-tasks.md", "path to write root tasks list")
	repairCmd.Flags().StringVar(&repairFree, "free", "tasks/free-tasks.md", "path to write free tasks list")
	repairCmd.Flags().StringVar(&repairFmt, "format", "text", "output format for repair errors: text|json")
}

func runRepair(tasksRoot, rootsFile, freeFile, outFormat string) error {
	// Parse all tasks
	parser := task.NewParser()
	tasks, err := parser.LoadTasks(tasksRoot)
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	// Update parent blockers from incomplete subtasks
	if _, err := task.UpdateBlockersFromChildren(tasks); err != nil {
		return fmt.Errorf("failed to update blockers from subtasks: %w", err)
	}

	// Update parent TODO entries from subtasks
	if _, err := task.UpdateAllParentTodoEntries(tasks); err != nil {
		return fmt.Errorf("failed to update parent TODO entries: %w", err)
	}

	// Fix missing references, then validate tasks
	rolesDir := filepath.Join(filepath.Dir(tasksRoot), "roles")
	validator := task.NewValidatorWithRoles(tasks, rolesDir)
	fixed := validator.FixMissingReferences()
	errors := validator.Validate()

	// Persist repaired tasks at the end of the run
	if _, err := task.WriteDirtyTasks(tasks); err != nil {
		return fmt.Errorf("failed to write repaired tasks: %w", err)
	}

	// Generate master lists
	if err := task.GenerateMasterLists(tasks, tasksRoot, rootsFile, freeFile); err != nil {
		return fmt.Errorf("failed to generate master lists: %w", err)
	}

	// Report repairs to missing references
	if len(fixed) > 0 && outFormat != "json" {
		for _, e := range fixed {
			fmt.Println("ERROR:", e.Error())
		}
	}

	// Report errors
	if len(errors) > 0 {
		if outFormat == "json" {
			errMsgs := make([]string, len(errors))
			for i, e := range errors {
				errMsgs[i] = e.Error()
			}
			payload := map[string]interface{}{"errors": errMsgs}
			if len(fixed) > 0 {
				fixedMsgs := make([]string, len(fixed))
				for i, e := range fixed {
					fixedMsgs[i] = e.Error()
				}
				payload["fixed"] = fixedMsgs
			}
			b, _ := json.MarshalIndent(payload, "", "  ")
			fmt.Println(string(b))
		} else {
			for _, e := range errors {
				fmt.Println("ERROR:", e.Error())
			}
		}
		return fmt.Errorf("repair failed: %d error(s)", len(errors))
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
		payload := map[string]interface{}{"roots": roots, "free": free}
		if len(fixed) > 0 {
			fixedMsgs := make([]string, len(fixed))
			for i, e := range fixed {
				fixedMsgs[i] = e.Error()
			}
			payload["fixed"] = fixedMsgs
		}
		b, _ := json.MarshalIndent(payload, "", "  ")
		fmt.Println(string(b))
	} else {
		fmt.Println("repair: ok")
		fmt.Printf("Repaired %d tasks\n", len(tasks))
		fmt.Printf("Master lists updated: %s, %s\n", rootsFile, freeFile)
	}

	return nil
}
