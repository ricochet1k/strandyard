/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ricochet1k/strandyard/pkg/task"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	searchSort    string
	searchOrder   string
	searchFormat  string
	searchColumns string
	searchGroup   string
	searchMDTable bool
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search tasks by title, description, and todos",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.TrimSpace(strings.Join(args, " "))
		if query == "" {
			return fmt.Errorf("search query cannot be empty")
		}

		paths, err := resolveProjectPaths(projectName)
		if err != nil {
			return err
		}

		opts, err := searchOptionsFromFlags(query)
		if err != nil {
			return err
		}

		tasks, err := task.SearchTasks(paths.TasksDir, opts)
		if err != nil {
			return err
		}
		output, err := task.FormatList(tasks, opts.ListOptions)
		if err != nil {
			return err
		}
		if output != "" {
			fmt.Fprintln(cmd.OutOrStdout(), output)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringVar(&searchSort, "sort", "", "sort by: id|priority|created|edited|role")
	searchCmd.Flags().StringVar(&searchOrder, "order", "asc", "sort order: asc|desc")
	searchCmd.Flags().StringVar(&searchFormat, "format", "table", "output format: table|md|json")
	searchCmd.Flags().StringVar(&searchColumns, "columns", "", "comma-separated list of columns to include")
	searchCmd.Flags().StringVar(&searchGroup, "group", "none", "group by: none|priority|parent|role")
	searchCmd.Flags().BoolVar(&searchMDTable, "md-table", false, "use markdown table output (with --format md)")
}

func searchOptionsFromFlags(query string) (task.SearchOptions, error) {
	opts := task.SearchOptions{
		Query: query,
		ListOptions: task.ListOptions{
			Sort:    strings.ToLower(strings.TrimSpace(searchSort)),
			Order:   strings.ToLower(strings.TrimSpace(searchOrder)),
			Format:  strings.ToLower(strings.TrimSpace(searchFormat)),
			Group:   strings.ToLower(strings.TrimSpace(searchGroup)),
			MdTable: searchMDTable,
		},
	}

	if strings.TrimSpace(searchColumns) != "" {
		raw := strings.Split(searchColumns, ",")
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

	switch opts.Sort {
	case "", "id", "priority", "created", "edited", "role":
	default:
		return task.SearchOptions{}, fmt.Errorf("invalid sort %q (expected id, priority, created, edited, or role)", opts.Sort)
	}
	switch opts.Order {
	case "asc", "desc":
	default:
		return task.SearchOptions{}, fmt.Errorf("invalid order %q (expected asc or desc)", opts.Order)
	}
	switch opts.Format {
	case "table", "md", "json":
	default:
		return task.SearchOptions{}, fmt.Errorf("invalid format %q (expected table, md, or json)", opts.Format)
	}
	switch opts.Group {
	case "none", "priority", "parent", "role":
	default:
		return task.SearchOptions{}, fmt.Errorf("invalid group %q (expected none, priority, parent, or role)", opts.Group)
	}

	return opts, nil
}

func runSearchWithProject(w io.Writer, projectName string, opts task.SearchOptions) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}
	tasks, err := task.SearchTasks(paths.TasksDir, opts)
	if err != nil {
		return err
	}
	output, err := task.FormatList(tasks, opts.ListOptions)
	if err != nil {
		return err
	}
	if output != "" {
		fmt.Fprintln(w, output)
	}
	return nil
}
