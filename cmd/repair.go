/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
)

var (
	repairFmt string
	repairAll bool
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
		return runRepair(cmd.OutOrStdout(), paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile, repairFmt)
	},
}

func init() {
	rootCmd.AddCommand(repairCmd)
	repairCmd.Flags().StringVar(&repairFmt, "format", "text", "output format for repair errors: text|json")
	repairCmd.Flags().BoolVar(&repairAll, "all", false, "output format for repair errors: text|json")
}

func runRepair(w io.Writer, tasksRoot, rootsFile, freeFile, outFormat string) error {
	db := task.NewTaskDB(tasksRoot)
	if err := db.LoadAllIfEmpty(); err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	if _, err := db.SyncBlockersFromChildren(); err != nil {
		return fmt.Errorf("failed to update blockers from subtasks: %w", err)
	}

	if _, err := task.UpdateAllParentTodoEntries(db.GetAll()); err != nil {
		return fmt.Errorf("failed to update parent TODO entries: %w", err)
	}

	rolesDir := filepath.Join(filepath.Dir(tasksRoot), "roles")
	validator := task.NewValidatorWithRoles(db.GetAll(), rolesDir)
	fixed := validator.FixMissingReferences()
	validationErrors := validator.ValidateAndRepair()

	var repairedCount int
	var err error
	if repairAll {
		fmt.Printf("Writing all tasks...")
		repairedCount, err = db.SaveAll()
	} else {
		repairedCount, err = db.SaveDirty()
	}
	if err != nil {
		return fmt.Errorf("failed to write repaired tasks: %w", err)
	}

	if err := task.GenerateMasterLists(db.GetAll(), tasksRoot, rootsFile, freeFile); err != nil {
		return fmt.Errorf("failed to generate master lists: %w", err)
	}

	if len(fixed) > 0 && outFormat != "json" {
		for _, e := range fixed {
			fmt.Fprintln(w, "ERROR:", e.Error())
		}
	}

	if len(validationErrors) > 0 {
		if outFormat == "json" {
			errMsgs := make([]string, len(validationErrors))
			for i, e := range validationErrors {
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
			fmt.Fprintln(w, string(b))
		} else {
			for _, e := range validationErrors {
				fmt.Fprintln(w, "ERROR:", e.Error())
			}
		}
		return fmt.Errorf("repair failed: %d error(s)", len(validationErrors))
	}

	if outFormat == "json" {
		var roots, free []string
		for id, t := range db.GetAll() {
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
		fmt.Fprintln(w, string(b))
	} else {
		fmt.Fprintln(w, "repair: ok")
		fmt.Fprintf(w, "Repaired %d tasks\n", repairedCount)
	}

	return nil
}
