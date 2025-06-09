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

// NewTagCmd represents the tag command.
func NewTagCmd() *cobra.Command {
	cmd := &cobra.Command{
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

	cmd.Flags().BoolP(flagPush, flagPushShort, false, flagPush)
	cmd.Flags().StringP(flagMessage, flagMessageShort, "", flagMessage)

	return cmd
}
