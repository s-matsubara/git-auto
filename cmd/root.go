/*
Package cmd
Copyright © 2023 matsubara
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd represents the base command when called without any subcommands.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git-auto",
		Short: "Auto git commands",
	}

	cmd.AddCommand(NewTagCmd())
	cmd.AddCommand(NewDeleteMergedBranchCmd())
	cmd.AddCommand(NewVersionCmd())

	return cmd
}

var rootCmd = NewRootCmd()

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
