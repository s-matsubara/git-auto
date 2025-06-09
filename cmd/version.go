/*
Package cmd
Copyright © 2025 matsubara
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	// Version is the semantic version of the CLI, set during build.
	Version = "dev"
	// Commit is the git commit hash, set during build.
	Commit = "none"
	// BuildDate is the build timestamp, set during build.
	BuildDate = "unknown"
)

// NewVersionCmd returns a new cobra.Command that prints CLI version information.
func NewVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version, commit, and build date information for the CLI`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Version: %s\n", Version); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Commit: %s\n", Commit); err != nil {
				return err
			}
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Build Date: %s\n", BuildDate); err != nil {
				return err
			}

			return nil
		},
	}

	return versionCmd
}
