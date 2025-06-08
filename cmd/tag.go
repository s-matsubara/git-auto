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

const (
	flagPush         = "push"
	flagPushShort    = "p"
	flagMessage      = "message"
	flagMessageShort = "m"
)

// tagCmd represents the tag command.
var tagCmd = &cobra.Command{
	Use:   "tag [<version>] [major|minor|patch]",
	Short: "Auto increment tag version",
	RunE: func(cmd *cobra.Command, args []string) error {
		isPush, err := cmd.Flags().GetBool(flagPush)
		if err != nil {
			return err
		}

		msg, err := cmd.Flags().GetString(flagMessage)
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
	tagCmd.Flags().BoolP(flagPush, flagPushShort, false, flagPush)
	tagCmd.Flags().StringP(flagMessage, flagMessageShort, "", flagMessage)
}
