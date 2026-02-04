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

	"github.com/ricochet1k/strandyard/pkg/task"
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

	freePath := paths.FreeTasksFile
	if _, err := os.Stat(freePath); os.IsNotExist(err) {
		if err := runRepair(w, paths.TasksDir, paths.RootTasksFile, freePath, "text"); err != nil {
			return fmt.Errorf("unable to generate master lists: %w", err)
		}
	}

	data, err := os.ReadFile(freePath)
	if err != nil {
		return fmt.Errorf("unable to read %s: %w", freePath, err)
	}

	db := task.NewTaskDB(paths.TasksDir)
	if err := db.LoadAllIfEmpty(); err != nil {
		return fmt.Errorf("failed to load tasks: %w", err)
	}

	parsed := task.ParseFreeList(string(data), db.GetAll())
	if len(parsed.TaskIDs) == 0 {
		fmt.Fprintln(w, "No free tasks found")
		return nil
	}

	type candidate struct {
		task *task.Task
		path string
	}
	var candidatesParsed []candidate
	var hasOwnerTasks bool

	for _, taskID := range parsed.TaskIDs {
		t, err := db.Get(taskID)
		if err != nil {
			continue
		}

		taskRole := t.GetEffectiveRole()
		if taskRole == "owner" {
			hasOwnerTasks = true
		}

		if roleFilter != "" {
			if taskRole != roleFilter {
				continue
			}
		} else if taskRole == "owner" {
			continue
		}

		candidatesParsed = append(candidatesParsed, candidate{
			task: t,
			path: filepath.ToSlash(t.FilePath),
		})
	}

	if len(candidatesParsed) == 0 {
		if roleFilter != "" {
			fmt.Fprintf(w, "No free tasks found for role: %s\n", roleFilter)
		} else if hasOwnerTasks {
			fmt.Fprintln(w, "No free tasks found. There are owner tasks remaining; try `strand next --role owner`.")
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
	fmt.Fprint(w, selectedTask.Content())

	for i, todo := range selectedTask.TodoItems {
		if !todo.Checked {
			fmt.Fprintf(w, "\n\nYou should focus on TODO #%v which is: %v\n", i+1, todo.Text)
			fmt.Fprintf(w, "\nMark the TODO completed with `strand complete %v --role %v --todo %v \"report\"`\n", selectedTask.ID, role, i+1)
			break
		}
	}

	return nil
}
