package cmd

import (
	"fmt"
	"io"
	"sort"

	"github.com/ricochet1k/strandyard/pkg/template"
	"github.com/spf13/cobra"
)

// templatesCmd represents the templates command
var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available task templates with descriptions",
	Long: `The 'templates' command lists all available task templates found in the '.strand/templates/' directory,
along with their short descriptions. This helps users choose the appropriate template
when creating new tasks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTemplates(cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}

func runTemplates(w io.Writer) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	templates, err := template.LoadTemplates(paths.TemplatesDir)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Available Task Templates:")
	fmt.Fprintln(w, "--------------------------")

	var names []string
	for name := range templates {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		t := templates[name]
		desc := t.Meta.Description
		if desc == "" {
			desc = "(no description found)"
		}
		fmt.Fprintf(w, "%-20s %s\n", name, desc)
	}
	return nil
}
