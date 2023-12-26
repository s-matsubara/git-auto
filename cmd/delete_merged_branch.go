/*
Package cmd
Copyright © 2023 matsubara
*/
package cmd

import (
	"git-auto/usecase"

	"github.com/spf13/cobra"
)

// deleteMergedBranchCmd represents the deleteMergedBranch command
var deleteMergedBranchCmd = &cobra.Command{
	Use:     "delete-merged-branch",
	Aliases: []string{"mergedd"},
	Short:   "Delete merged branch",
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
	rootCmd.AddCommand(deleteMergedBranchCmd)
}
