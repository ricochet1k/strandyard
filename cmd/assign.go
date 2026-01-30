/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// assignCmd represents the assign command
var assignCmd = &cobra.Command{
	Use:   "assign",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runAssign(cmd.OutOrStdout())
	},
}

func init() {
	rootCmd.AddCommand(assignCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assignCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// assignCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runAssign(w io.Writer) error {
	_, err := fmt.Fprintln(w, "assign called")
	return err
}
