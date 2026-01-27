/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import "github.com/spf13/cobra"

// issueCmd groups issue-related subcommands.
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage issue-style tasks",
}

func init() {
	rootCmd.AddCommand(issueCmd)
}
