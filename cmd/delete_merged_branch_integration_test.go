package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const integrationEnv = "true"

func gitCmd(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%v: %s", err, out)
	}
	return nil
}

func setupRepo(t *testing.T) (func(), string) {
	t.Helper()
	dir, err := os.MkdirTemp("", "gitauto")
	if err != nil {
		t.Fatal(err)
	}
	cleanup := func() { _ = os.RemoveAll(dir) } //nolint:errcheck // cleanup

	if err := gitCmd(dir, "init"); err != nil {
		t.Fatal(err)
	}
	if err := gitCmd(dir, "config", "user.name", "tester"); err != nil {
		t.Fatal(err)
	}
	if err := gitCmd(dir, "config", "user.email", "tester@example.com"); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := gitCmd(dir, "add", "."); err != nil {
		t.Fatal(err)
	}
	if err := gitCmd(dir, "commit", "-m", "init"); err != nil {
		t.Fatal(err)
	}

	return cleanup, dir
}

func TestDeleteMergedBranchCommandIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != integrationEnv {
		t.Skip("integration test")
	}

	cleanup, dir := setupRepo(t)
	defer cleanup()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }() //nolint:errcheck // revert dir
	if chdirErr := os.Chdir(dir); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	defaultBranchBytes, err := exec.Command("git", "branch", "--show-current").Output()
	if err != nil {
		t.Fatal(err)
	}
	defaultBranch := strings.TrimSpace(string(defaultBranchBytes))

	if gitErr := gitCmd(dir, "checkout", "-b", "feature"); gitErr != nil {
		t.Fatal(gitErr)
	}
	if writeErr := os.WriteFile(filepath.Join(dir, "feature.txt"), []byte("a"), 0o644); writeErr != nil {
		t.Fatal(writeErr)
	}
	if gitErr := gitCmd(dir, "add", "feature.txt"); gitErr != nil {
		t.Fatal(gitErr)
	}
	if gitErr := gitCmd(dir, "commit", "-m", "feature"); gitErr != nil {
		t.Fatal(gitErr)
	}
	if gitErr := gitCmd(dir, "checkout", defaultBranch); gitErr != nil {
		t.Fatal(gitErr)
	}
	if gitErr := gitCmd(dir, "merge", "feature"); gitErr != nil {
		t.Fatal(gitErr)
	}

	c := NewDeleteMergedBranchCmd()
	if runErr := c.RunE(c, []string{}); runErr != nil {
		t.Fatal(runErr)
	}

	out, err := exec.Command("git", "branch", "--list", "feature").Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "" {
		t.Fatalf("branch not deleted: %s", out)
	}
}
