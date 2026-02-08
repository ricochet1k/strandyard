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
	"time"

	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
)

var nextRole string
var nextClaim bool
var nextClaimTimeout time.Duration

type nextOptions struct {
	Claim        bool
	ClaimTimeout time.Duration
	Now          func() time.Time
}

// nextCmd represents the next command
var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Print the next free task",
	Long: `Print the next free task from the free-tasks list.
Also prints the full role (from metadata or first TODO) so that the output
contains all the information an agent needs to execute the task without
looking anything else up.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNextWithOptions(cmd.OutOrStdout(), projectName, nextRole, nextOptions{
			Claim:        nextClaim,
			ClaimTimeout: nextClaimTimeout,
		})
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
	nextCmd.Flags().StringVar(&nextRole, "role", "", "optional: filter tasks by role")
	nextCmd.Flags().BoolVar(&nextClaim, "claim", false, "claim the selected task by marking it in_progress")
	nextCmd.Flags().DurationVar(&nextClaimTimeout, "claim-timeout", time.Hour, "timeout before an in-progress claim is treated as open again")
}

func runNext(w io.Writer, projectName, roleFilter string) error {
	return runNextWithOptions(w, projectName, roleFilter, nextOptions{ClaimTimeout: time.Hour})
}

func runNextWithOptions(w io.Writer, projectName, roleFilter string, opts nextOptions) error {
	if opts.ClaimTimeout <= 0 {
		return fmt.Errorf("--claim-timeout must be greater than 0")
	}
	if opts.Now == nil {
		opts.Now = time.Now
	}
	now := opts.Now().UTC()

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

	claimStateChanged := false

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

		if t.Meta.IsInProgress() {
			if now.Sub(t.Meta.DateEdited) >= opts.ClaimTimeout {
				if err := db.SetStatus(taskID, task.StatusOpen); err != nil {
					return fmt.Errorf("failed to reopen expired claim for %s: %w", taskID, err)
				}
				claimStateChanged = true
			} else {
				continue
			}
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

	if opts.Claim {
		if err := db.MarkInProgress(selectedTask.ID); err != nil {
			return fmt.Errorf("failed to claim task %s: %w", selectedTask.ID, err)
		}
		claimStateChanged = true
	}

	if claimStateChanged {
		if _, err := db.SaveDirty(); err != nil {
			return fmt.Errorf("failed to persist task claim state: %w", err)
		}
		if err := task.GenerateMasterLists(db.GetAll(), paths.TasksDir, paths.RootTasksFile, paths.FreeTasksFile); err != nil {
			return fmt.Errorf("failed to update master lists: %w", err)
		}
	}

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

	// Print ancestors if this task has parents
	ancestors := db.GetAncestors(selectedTask.ID)
	if len(ancestors) > 0 {
		fmt.Fprint(w, "\nAncestors:\n")
		for _, ancestor := range ancestors {
			fmt.Fprintf(w, "  %s: %s\n", ancestor[0], ancestor[1])
		}
		fmt.Fprint(w, "\n")
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
