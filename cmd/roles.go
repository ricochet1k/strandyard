package cmd

import (
	"fmt"
	"io"
	"sort"

	"github.com/ricochet1k/strandyard/pkg/role"
	"github.com/spf13/cobra"
)

// rolesCmd represents the roles command
var rolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "List available roles with descriptions",
	Long: `The 'roles' command lists all available roles found in the '.strand/roles/' directory,
along with their short descriptions. This helps users understand the responsibilities
of each role.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoles(cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(rolesCmd)
}

func runRoles(w io.Writer) error {
	paths, err := resolveProjectPaths(projectName)
	if err != nil {
		return err
	}

	roles, err := role.LoadRoles(paths.RolesDir)
	if err != nil {
		return err
	}

	fmt.Fprintln(w, "Available Roles:")
	fmt.Fprintln(w, "--------------------------")

	var names []string
	for name := range roles {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		r := roles[name]
		desc := r.Meta.Description
		if desc == "" {
			desc = "(no description found)"
		}
		fmt.Fprintf(w, "%-20s %s\n", name, desc)
	}
	return nil
}
