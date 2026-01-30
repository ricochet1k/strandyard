/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ricochet1k/streamyard/pkg/task"
	"github.com/spf13/cobra"
)

var nextRole string

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Print the next free task",
	Long: `Print the next free task from the free-tasks list.
Also prints the full role (from metadata or first TODO) so that the output
contains all the information an agent needs to execute the task without
looking anything else up.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNext(cmd.OutOrStdout(), projectName, nextRole)
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
	nextCmd.Flags().StringVar(&nextRole, "role", "", "optional: filter tasks by role")
}

func runNext(w io.Writer, projectName, roleFilter string) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	// Ensure free-tasks exists; if not, run repair to generate it
	freePath := paths.FreeTasksFile
	if _, err := os.Stat(freePath); os.IsNotExist(err) {
		// Run repair to generate lists
		if err := runRepair(w, paths.TasksDir, paths.RootTasksFile, freePath, "text"); err != nil {
			return fmt.Errorf("unable to generate master lists: %w", err)
		}
	}

	// Read free tasks list
	data, err := os.ReadFile(freePath)
	if err != nil {
		return fmt.Errorf("unable to read %s: %w", freePath, err)
	}

	parser := task.NewParser()
	tasks, err := parser.LoadTasks(paths.TasksDir)
	if err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	parsed := task.ParseFreeList(string(data), tasks)
	if len(parsed.TaskIDs) == 0 {
		fmt.Fprintln(w, "No free tasks found")
		return nil
	}

	type candidate struct {
		task *task.Task
		path string
	}
	candidatesParsed := []candidate{}
	for _, taskID := range parsed.TaskIDs {
		t, exists := tasks[taskID]
		if !exists {
			continue
		}

		// If role filter specified, check if it matches
		if roleFilter != "" {
			taskRole := t.GetEffectiveRole()
			if taskRole != roleFilter {
				continue
			}
		}

		candidatesParsed = append(candidatesParsed, candidate{
			task: t,
			path: filepath.ToSlash(t.FilePath),
		})
	}

	if len(candidatesParsed) == 0 {
		if roleFilter != "" {
			fmt.Fprintf(w, "No free tasks found for role: %s\n", roleFilter)
		} else {
			fmt.Fprintln(w, "No free tasks found")
		}
		return nil
	}

	sort.Slice(candidatesParsed, func(i, j int) bool {
		pi := task.PriorityRank(candidatesParsed[i].task.Meta.Priority)
		pj := task.PriorityRank(candidatesParsed[j].task.Meta.Priority)
		if pi != pj {
			return pi < pj
		}
		return candidatesParsed[i].path < candidatesParsed[j].path
	})

	selectedTask := candidatesParsed[0].task

	role := selectedTask.GetEffectiveRole()
	if role != "" {
		rolePath := filepath.Join(paths.RolesDir, role+".md")
		roleData, err := os.ReadFile(rolePath)
		if err == nil {
			roleDoc := string(roleData)
			fmt.Fprintf(w, "Your role is %s. Here's the description of that role:\n\n", role)
			fmt.Fprint(w, roleDoc)
			if !strings.HasSuffix(roleDoc, "\n") {
				fmt.Fprint(w, "\n")
			}
			fmt.Fprint(w, "\n---\n")
		} else {
			fmt.Fprintf(w, "Your role is %s. This role appears to be missing, ask the user what to do.\n\n", role)
		}
	} else {
		fmt.Fprint(w, "This task has no role, ask the user what to do.\n\n")
	}

	fmt.Fprintf(w, "\nYour task is %s. Here's the description of that task:\n\n", selectedTask.ID)

	// Print task content
	fmt.Fprint(w, selectedTask.Content)

	return nil
}
