package cmd

import "testing"

func TestTagCommandNoArgs(t *testing.T) {
	if err := tagCmd.RunE(tagCmd, []string{}); err == nil {
		t.Fatal("expected error")
	}
}
