/*
Package cmd
Copyright © 2023 matsubara
*/
package cmd

import (
	"git-auto/usecase"

	"github.com/spf13/cobra"
)

// mergeddCmd represents the mergedd command
var mergeddCmd = &cobra.Command{
	Use:   "mergedd",
	Short: "Delete merged branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		u := usecase.NewGitUsecase()
		err := u.DeleteMergedBranches()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergeddCmd)
}
