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
	VersionUp(target string, atgMsg string, isPush bool) (string, error)
	DeleteMergedBranches() error
}

type gitUsecase struct{}

func NewGitUsecase() GitUsecase {
	return &gitUsecase{}
}

func (u *gitUsecase) VersionUp(target string, tagMsg string, isPush bool) (string, error) {
	var version string
	var err error

	switch target {
	case "major", "minor", "patch":
		tag, err := u.getCurrentTag()
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

func (u *gitUsecase) tagVersionUp(tag string, target string) (string, error) {
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

func (u *gitUsecase) setTag(tag string, msg string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("git", "tag", tag)
	if msg != "" {
		cmd = exec.Command("git", "tag", "-am", msg, tag)
	}

	p, _ := os.Getwd()
	cmd.Dir = p
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (u *gitUsecase) pushTag(tag string) error {
	p, _ := os.Getwd()

	cmd := exec.Command("git", "push", "origin", tag)
	cmd.Dir = p
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return errors.New(stderr.String())
	}

	return nil
}

func (u *gitUsecase) getCurrentTag() (string, error) {
	p, _ := os.Getwd()

	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = p
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return "", errors.New(stderr.String())
	}

	return strings.Replace(stdout.String(), "\n", " ", -1), nil
}

func (u *gitUsecase) getMergedBranches() ([]string, error) {
	p, _ := os.Getwd()

	cmd := exec.Command("git", "branch", "--merged")
	cmd.Dir = p
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return []string{}, errors.New(stderr.String())
	}

	re := regexp.MustCompile(`^\*|main|master|development|staging|production`)

	var targetBranches []string
	branches := strings.Split(stdout.String(), "\n")
	for _, branch := range branches {
		branch = strings.Replace(branch, " ", "", -1)
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
	p, _ := os.Getwd()

	cmd := exec.Command("git", "branch", "-D", branch)
	cmd.Dir = p
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

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
	return strings.Replace(tag, "v", "", -1)
}

func (u *gitUsecase) incrementVersion(version string, target string) (string, error) {
	version = strings.Replace(version, " ", "", -1)
	versions := strings.Split(version, ".")

	re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	if !re.MatchString(version) {
		return "", errors.New("does not match version")
	}

	switch target {
	case "major":
		num, _ := strconv.Atoi(versions[0])
		versions[0] = strconv.Itoa(num + 1)
		versions[1] = "0"
		versions[2] = "0"
	case "minor":
		num, _ := strconv.Atoi(versions[1])
		versions[1] = strconv.Itoa(num + 1)
		versions[2] = "0"
	case "patch":
		num, _ := strconv.Atoi(versions[2])
		versions[2] = strconv.Itoa(num + 1)
	}

	return strings.Join(versions, "."), nil
}
