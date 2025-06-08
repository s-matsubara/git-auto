/*
Package cmd
Copyright © 2023 matsubara
*/
package cmd

import (
	"errors"
	"git-auto/usecase"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command.
var tagCmd = &cobra.Command{
	Use:   "tag [<version>] [major|minor|patch]",
	Short: "Auto increment tag version",
	RunE: func(cmd *cobra.Command, args []string) error {
		isPush, err := cmd.Flags().GetBool("push")
		if err != nil {
			return err
		}

		msg, err := cmd.Flags().GetString("message")
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return errors.New("not found argument")
		}

		target := args[0]

		u := usecase.NewGitUsecase()
		_, err = u.VersionUp(target, msg, isPush)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.Flags().BoolP("push", "p", false, "push")
	tagCmd.Flags().StringP("message", "m", "", "message")
}
