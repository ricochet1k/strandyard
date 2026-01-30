/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

//go:embed assets/agents.md
var agentsDoc string

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Print portable agent instructions",
	Long:  "Print a backend-agnostic subset of strand agent instructions for reuse.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAgents(cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(agentsCmd)
}

func runAgents(w io.Writer) error {
	_, err := fmt.Fprint(w, agentsDoc)
	return err
}
