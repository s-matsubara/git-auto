/*
Package cmd
Copyright © 2025 matsubara
*/
package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCommand(t *testing.T) {
	oldVersion := Version
	oldCommit := Commit
	oldBuildDate := BuildDate

	Version = "1.0.0"
	Commit = "abc123"
	BuildDate = "2023-01-01"

	defer func() {
		Version = oldVersion
		Commit = oldCommit
		BuildDate = oldBuildDate
	}()

	buf := new(bytes.Buffer)

	cmd := NewVersionCmd()
	cmd.SetOut(buf)

	err := cmd.Execute()
	assert.NoError(t, err, "Expected no error when executing version command")

	output := buf.String()
	assert.Contains(t, output, "Version: 1.0.0", "Expected output to contain version")
	assert.Contains(t, output, "Commit: abc123", "Expected output to contain commit")
	assert.Contains(t, output, "Build Date: 2023-01-01", "Expected output to contain build date")
}
