/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ricochet1k/streamyard/pkg/task"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	listScope          string
	listChildren       string
	listRole           string
	listPriority       string
	listCompleted      bool
	listBlocked        bool
	listBlocks         bool
	listOwnerApproval  bool
	listLabel          string
	listSort           string
	listOrder          string
	listFormat         string
	listColumns        string
	listGroup          string
	listMDTable        bool
	listUseMasterLists bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks with filtering and formatting options",
	RunE: func(cmd *cobra.Command, args []string) error {
		paths, err := resolveProjectPaths(projectName)
		if err != nil {
			return err
		}
		opts, err := listOptionsFromFlags(cmd)
		if err != nil {
			return err
		}
		return runList(cmd.OutOrStdout(), paths.TasksDir, opts)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&listScope, "scope", "all", "scope of tasks to list: all|root|free")
	listCmd.Flags().StringVar(&listChildren, "children", "", "list direct children of the given task ID")
	listCmd.Flags().StringVar(&listRole, "role", "", "filter by role name")
	listCmd.Flags().StringVar(&listPriority, "priority", "", "filter by priority: high|medium|low")
	listCmd.Flags().BoolVar(&listCompleted, "completed", false, "list only completed tasks (default: uncompleted)")
	listCmd.Flags().BoolVar(&listBlocked, "blocked", false, "filter by blocked status (has blockers)")
	listCmd.Flags().BoolVar(&listBlocks, "blocks", false, "filter by blocks status (has blocks)")
	listCmd.Flags().BoolVar(&listOwnerApproval, "owner-approval", false, "filter by owner approval")
	listCmd.Flags().StringVar(&listLabel, "label", "", "reserved for future labels support")
	listCmd.Flags().StringVar(&listSort, "sort", "", "sort by: id|priority|created|edited|role")
	listCmd.Flags().StringVar(&listOrder, "order", "asc", "sort order: asc|desc")
	listCmd.Flags().StringVar(&listFormat, "format", "table", "output format: table|md|json")
	listCmd.Flags().StringVar(&listColumns, "columns", "", "comma-separated list of columns to include")
	listCmd.Flags().StringVar(&listGroup, "group", "none", "group by: none|priority|parent|role")
	listCmd.Flags().BoolVar(&listMDTable, "md-table", false, "use markdown table output (with --format md)")
	listCmd.Flags().BoolVar(&listUseMasterLists, "use-master-lists", false, "use master lists for root/free scopes when no filters")
}

func listOptionsFromFlags(cmd *cobra.Command) (task.ListOptions, error) {
	opts := task.ListOptions{
		Scope:          strings.ToLower(strings.TrimSpace(listScope)),
		Parent:         strings.TrimSpace(listChildren),
		Role:           strings.TrimSpace(listRole),
		Priority:       strings.ToLower(strings.TrimSpace(listPriority)),
		Label:          strings.TrimSpace(listLabel),
		Sort:           strings.ToLower(strings.TrimSpace(listSort)),
		Order:          strings.ToLower(strings.TrimSpace(listOrder)),
		Format:         strings.ToLower(strings.TrimSpace(listFormat)),
		Group:          strings.ToLower(strings.TrimSpace(listGroup)),
		MdTable:        listMDTable,
		UseMasterLists: listUseMasterLists,
	}

	if !cmd.Flags().Changed("completed") {
		opts.Completed = boolPtr(false)
	}

	if cmd.Flags().Changed("completed") {
		opts.Completed = boolPtr(listCompleted)
	}
	if cmd.Flags().Changed("blocked") {
		opts.Blocked = boolPtr(listBlocked)
	}
	if cmd.Flags().Changed("blocks") {
		opts.Blocks = boolPtr(listBlocks)
	}
	if cmd.Flags().Changed("owner-approval") {
		opts.OwnerApproval = boolPtr(listOwnerApproval)
	}

	if strings.TrimSpace(listColumns) != "" {
		raw := strings.Split(listColumns, ",")
		opts.Columns = make([]string, 0, len(raw))
		for _, col := range raw {
			col = strings.ToLower(strings.TrimSpace(col))
			if col != "" {
				opts.Columns = append(opts.Columns, col)
			}
		}
	}

	if opts.Format != "json" {
		opts.Color = term.IsTerminal(int(os.Stdout.Fd()))
	}

	return opts, nil
}

func runList(w io.Writer, tasksRoot string, opts task.ListOptions) error {
	if opts.Label != "" {
		return fmt.Errorf("label filter is not supported yet")
	}
	switch opts.Scope {
	case "all", "root", "free":
	default:
		return fmt.Errorf("invalid scope %q (expected all, root, or free)", opts.Scope)
	}
	if opts.Priority != "" && !task.IsValidPriority(opts.Priority) {
		return fmt.Errorf("invalid priority %q (expected high, medium, or low)", opts.Priority)
	}
	switch opts.Sort {
	case "", "id", "priority", "created", "edited", "role":
	default:
		return fmt.Errorf("invalid sort %q (expected id, priority, created, edited, or role)", opts.Sort)
	}
	switch opts.Order {
	case "asc", "desc":
	default:
		return fmt.Errorf("invalid order %q (expected asc or desc)", opts.Order)
	}
	switch opts.Format {
	case "table", "md", "json":
	default:
		return fmt.Errorf("invalid format %q (expected table, md, or json)", opts.Format)
	}
	switch opts.Group {
	case "none", "priority", "parent", "role":
	default:
		return fmt.Errorf("invalid group %q (expected none, priority, parent, or role)", opts.Group)
	}
	if opts.Scope == "free" {
		if opts.Parent != "" {
			return fmt.Errorf("invalid flag combination: --scope free cannot be used with --children")
		}
		if opts.Group == "parent" {
			return fmt.Errorf("invalid flag combination: --scope free cannot be used with --group parent")
		}
	}
	if opts.Parent != "" && opts.Scope != "all" {
		return fmt.Errorf("invalid flag combination: --children cannot be used with --scope %s", opts.Scope)
	}

	tasks, err := task.ListTasks(tasksRoot, opts)
	if err != nil {
		return err
	}
	output, err := task.FormatList(tasks, opts)
	if err != nil {
		return err
	}
	if output != "" {
		fmt.Fprintln(w, output)
	}
	return nil
}

func runListWithProject(w io.Writer, projectName string, opts task.ListOptions) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}
	return runList(w, paths.TasksDir, opts)
}

func boolPtr(value bool) *bool {
	v := value
	return &v
}
