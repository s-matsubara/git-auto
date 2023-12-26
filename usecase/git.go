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
		break
	default:
		version = target
		break
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

func (u *gitUsecase) tagVersionUp(tag string, target string) (string, error) {
	check := u.checkVersionPrefix(tag)
	version := u.getIgnoreVersionPrefix(tag)
	version = u.incrementVersion(version, target)

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

	cmd := exec.Command("git", "push", tag)
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

func (u *gitUsecase) checkVersionPrefix(tag string) bool {
	re := regexp.MustCompile("^v")
	return re.MatchString(tag)
}

func (u *gitUsecase) getIgnoreVersionPrefix(tag string) string {
	return strings.Replace(tag, "v", "", -1)
}

func (u *gitUsecase) incrementVersion(version string, target string) string {
	versions := strings.Split(version, ".")

	switch target {
	case "major":
		num, _ := strconv.Atoi(versions[0])
		versions[0] = strconv.Itoa(num + 1)
		versions[1] = "0"
		versions[2] = "0"
		break
	case "minor":
		num, _ := strconv.Atoi(versions[1])
		versions[1] = strconv.Itoa(num + 1)
		versions[2] = "0"
	case "patch":
		num, _ := strconv.Atoi(versions[2])
		versions[2] = strconv.Itoa(num + 1)
	}

	return strings.Join(versions, ".")
}
