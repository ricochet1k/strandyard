/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//go:embed assets/agents.md
var agentsDoc string

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Print portable agent instructions",
	Long:  "Print a backend-agnostic subset of memmd agent instructions for reuse.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprint(cmd.OutOrStdout(), agentsDoc)
	},
}

func init() {
	rootCmd.AddCommand(agentsCmd)
}
