package usecase

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type GitUsecase interface {
	VersionUp(target, atgMsg string, isPush bool) (string, error)
	DeleteMergedBranches() error
}

type gitUsecase struct{}

func NewGitUsecase() GitUsecase {
	return &gitUsecase{}
}

func (u *gitUsecase) VersionUp(target, tagMsg string, isPush bool) (string, error) {
	var version string
	var err error

	switch target {
	case "major", "minor", "patch":
		var tag string
		tag, err = u.getCurrentTag()
		if err != nil {
			return "", err
		}

		version, err = u.tagVersionUp(tag, target)
		if err != nil {
			return "", err
		}
	default:
		version = target
	}

	err = u.setTag(version, tagMsg)
	if err != nil {
		return "", err
	}

	if isPush {
		err = u.pushTag(version)
		if err != nil {
			return "", err
		}
	}

	return version, nil
}

func (u *gitUsecase) DeleteMergedBranches() error {
	branches, err := u.getMergedBranches()
	if err != nil {
		return err
	}

	for _, branch := range branches {
		err = u.deleteBranch(branch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *gitUsecase) tagVersionUp(tag, target string) (string, error) {
	check := u.checkVersionPrefix(tag)
	version := u.getIgnoreVersionPrefix(tag)
	version, err := u.incrementVersion(version, target)
	if err != nil {
		return "", err
	}

	if check {
		version = "v" + version
	}

	return version, nil
}

func (u *gitUsecase) setTag(tag, msg string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("git", "tag", tag)
	if msg != "" {
		cmd = exec.Command("git", "tag", "-am", msg, tag)
	}

	p, err := os.Getwd()
	if err != nil {
		return err
	}
	cmd.Dir = p
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (u *gitUsecase) pushTag(tag string) error {
	p, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "push", "origin", tag)
	cmd.Dir = p
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (u *gitUsecase) getCurrentTag() (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}

	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = p
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		return "", errors.New(stderr.String())
	}

	return strings.ReplaceAll(stdout.String(), "\n", " "), nil
}

func (u *gitUsecase) getMergedBranches() ([]string, error) {
	p, err := os.Getwd()
	if err != nil {
		return []string{}, err
	}

	cmd := exec.Command("git", "branch", "--merged")
	cmd.Dir = p
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		return []string{}, errors.New(stderr.String())
	}

	re := regexp.MustCompile(`^\*|main|master|development|staging|production`)

	branches := strings.Split(stdout.String(), "\n")
	targetBranches := make([]string, 0, len(branches))
	for _, branch := range branches {
		branch = strings.ReplaceAll(branch, " ", "")
		if re.MatchString(branch) {
			continue
		}

		if branch == "" {
			continue
		}

		targetBranches = append(targetBranches, branch)
	}

	return targetBranches, nil
}

func (u *gitUsecase) deleteBranch(branch string) error {
	p, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "branch", "-D", branch)
	cmd.Dir = p
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (u *gitUsecase) checkVersionPrefix(tag string) bool {
	re := regexp.MustCompile("^v")
	return re.MatchString(tag)
}

func (u *gitUsecase) getIgnoreVersionPrefix(tag string) string {
	return strings.ReplaceAll(tag, "v", "")
}

func (u *gitUsecase) incrementVersion(version, target string) (string, error) {
	version = strings.ReplaceAll(version, " ", "")
	versions := strings.Split(version, ".")

	re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !re.MatchString(version) {
		return "", errors.New("does not match version")
	}

	switch target {
	case "major":
		num, err := strconv.Atoi(versions[0])
		if err != nil {
			return "", err
		}
		versions[0] = strconv.Itoa(num + 1)
		versions[1] = "0"
		versions[2] = "0"
	case "minor":
		num, err := strconv.Atoi(versions[1])
		if err != nil {
			return "", err
		}
		versions[1] = strconv.Itoa(num + 1)
		versions[2] = "0"
	case "patch":
		num, err := strconv.Atoi(versions[2])
		if err != nil {
			return "", err
		}
		versions[2] = strconv.Itoa(num + 1)
	}

	return strings.Join(versions, "."), nil
}
