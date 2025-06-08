package usecase

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
	if err := gitCmd(dir, "tag", "v1.0.0"); err != nil {
		t.Fatal(err)
	}
	return cleanup, dir
}

func TestVersionUpIntegration(t *testing.T) {
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

	u := NewGitUsecase()
	version, err := u.VersionUp(VersionMinor, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if version != "v1.1.0" {
		t.Fatalf("expected v1.1.0, got %s", version)
	}
	out, err := exec.Command("git", "-C", dir, "tag", "--list", "v1.1.0").Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(out) == "" {
		t.Fatalf("tag not created")
	}
}

func TestDeleteMergedBranchesIntegration(t *testing.T) {
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

	// create branch and merge
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

	u := NewGitUsecase()
	if delErr := u.DeleteMergedBranches(); delErr != nil {
		t.Fatal(delErr)
	}
	out, err := exec.Command("git", "branch", "--list", "feature").Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "" {
		t.Fatalf("branch not deleted: %s", out)
	}
}

func TestVersionUpPushIntegration(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != integrationEnv {
		t.Skip("integration test")
	}

	cleanup, dir := setupRepo(t)
	defer cleanup()

	remoteDir, err := os.MkdirTemp("", "gitauto-remote")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(remoteDir) }() //nolint:errcheck // cleanup

	if gitErr := gitCmd(remoteDir, "init", "--bare"); gitErr != nil {
		t.Fatal(gitErr)
	}
	if gitErr := gitCmd(dir, "remote", "add", "origin", remoteDir); gitErr != nil {
		t.Fatal(gitErr)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(wd) }() //nolint:errcheck // revert dir
	if chdirErr := os.Chdir(dir); chdirErr != nil {
		t.Fatal(chdirErr)
	}

	u := NewGitUsecase()
	version, err := u.VersionUp(VersionPatch, "", true)
	if err != nil {
		t.Fatal(err)
	}
	if version != "v1.0.1" {
		t.Fatalf("expected v1.0.1, got %s", version)
	}
	out, err := exec.Command("git", "-C", remoteDir, "tag", "--list", "v1.0.1").Output()
	if err != nil {
		t.Fatal(err)
	}
	if string(out) == "" {
		t.Fatalf("tag not pushed")
	}
}
