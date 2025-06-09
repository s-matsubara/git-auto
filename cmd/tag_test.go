package cmd

import "testing"

func TestTagCommandNoArgs(t *testing.T) {
	cmd := NewTagCmd()
	if err := cmd.RunE(cmd, []string{}); err == nil {
		t.Fatal("expected error")
	}
}
