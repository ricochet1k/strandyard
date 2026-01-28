/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
Shows the task's role (from metadata or first TODO) and the task content.
Use --role to filter tasks by specific role.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runNext(nextRole)
	},
}

func init() {
	rootCmd.AddCommand(nextCmd)
	nextCmd.Flags().StringVar(&nextRole, "role", "", "optional: filter tasks by role")
}

func runNext(roleFilter string) error {
	// Ensure free-tasks exists; if not, run repair to generate it
	freePath := "tasks/free-tasks.md"
	if _, err := os.Stat(freePath); os.IsNotExist(err) {
		// Run repair to generate lists
		if err := runRepair("tasks", "tasks/root-tasks.md", freePath, "text"); err != nil {
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
	baseDir := filepath.Dir(freePath)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "- ") {
			if path := parseListPath(strings.TrimSpace(strings.TrimPrefix(l, "- "))); path != "" {
				if !filepath.IsAbs(path) {
					path = filepath.Join(baseDir, path)
				}
				candidates = append(candidates, path)
			}
		}
	}

	if len(candidates) == 0 {
		fmt.Println("No free tasks found")
		return nil
	}

	// Parse tasks using the task library
	parser := task.NewParser()
	type candidate struct {
		task *task.Task
		path string
	}
	candidatesParsed := []candidate{}

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

		candidatesParsed = append(candidatesParsed, candidate{
			task: t,
			path: filepath.ToSlash(t.FilePath),
		})
	}

	if len(candidatesParsed) == 0 {
		if roleFilter != "" {
			fmt.Printf("No free tasks found for role: %s\n", roleFilter)
		} else {
			fmt.Println("No free tasks found")
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

	// Print role info (compact)
	role := selectedTask.GetEffectiveRole()
	if role != "" {
		fmt.Printf("Role: %s\n\n", role)
	} else {
		fmt.Println("Role: (none)")
	}

	// Print task content
	fmt.Print(selectedTask.Content)

	return nil
}

func parseListPath(entry string) string {
	if !strings.HasPrefix(entry, "[") {
		return entry
	}
	open := strings.Index(entry, "](")
	close := strings.LastIndex(entry, ")")
	if open == -1 || close == -1 || close <= open+2 {
		return ""
	}
	return strings.TrimSpace(entry[open+2 : close])
}
