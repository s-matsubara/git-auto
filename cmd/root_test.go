package cmd

import "testing"

func TestExecuteHelp(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	Execute()
}
