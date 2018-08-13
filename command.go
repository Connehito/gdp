package main

import (
	"errors"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

// Gdp is the interface which has methods deploying and publising.
type Gdp interface {
	IsMasterBranch() bool
	IsExistTagInLocal(tag string) bool
	IsExistTagInRemote(tag string) bool
	GetMergeCommitList(toTag string) (string, error)
	GetLatestTag() string
	Deploy(tag string) error
	Publish(tag string, commits string) error
}

// Command implements Git interface.
type Command struct {
}

// NewCommand is GdpCommand's constructor.
func NewCommand() (Gdp, error) {
	if err := availableCommand("git"); err != nil {
		return nil, err
	}
	if err := availableCommand("hub"); err != nil {
		return nil, err
	}

	if !isExistsCredential() {
		return nil, errors.New("please setup hub's credential file")
	}

	return &Command{}, nil
}

// IsMasterBranch checks current branch is master or not.
func (c *Command) IsMasterBranch() bool {
	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").CombinedOutput()
	if err != nil {
		return false
	}

	branch := strings.TrimRight(string(out), "\n")
	if branch != "master" {
		return false
	}

	return true
}

// IsExistTagInLocal checks the tag exist or not in local repository.
func (c *Command) IsExistTagInLocal(tag string) bool {
	out, err := exec.Command("git", "show", tag).CombinedOutput()
	if err != nil {
		return false
	}

	if string(out) == "" {
		return false
	}

	return true
}

// IsExistTagInRemote checks the tag exist or not in remote(origin) repository.
func (c *Command) IsExistTagInRemote(tag string) bool {
	out, err := exec.Command("git", "ls-remote", "--tags", "origin", tag).CombinedOutput()
	if err != nil {
		return false
	}

	if string(out) == "" {
		return false
	}

	return true
}

// GetMergeCommitList gets merge-commits list from previous tag to the tag
func (c *Command) GetMergeCommitList(toTag string) (string, error) {
	fromTag := getPreviousTag(toTag)
	if fromTag != "" {
		fromTag = fromTag + ".."
	}

	format := "--pretty=format:- %an: %b"
	out, err := exec.Command("git", "log", "--merges", "--first-parent", format, fromTag+toTag).CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}

	return strings.TrimRight(string(out), "\n"), nil
}

// GetLatestTag gets lastest tag name.
func (c *Command) GetLatestTag() string {
	out, err := exec.Command("git", "describe", "--abbrev=0", "--tags").CombinedOutput()
	if err != nil {
		return "" // No Tag
	}

	return strings.TrimRight(string(out), "\n")
}

// Deploy adds the tag and push the tag to remote(origin) repository.
func (c *Command) Deploy(tag string) error {
	out, err := exec.Command("git", "tag", tag).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	out, err = exec.Command("git", "push", "origin", tag).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	return nil
}

// Publish creates the release note in Github.
func (c *Command) Publish(tag string, message string) error {
	out, err := exec.Command("hub", "release", "create", "-m", message, tag).CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	return nil
}

func availableCommand(name string) error {
	out, err := exec.Command(name, "--version").CombinedOutput()
	if err != nil {
		if string(out) == "" {
			return err
		}

		return errors.New(string(out))
	}

	return nil
}

func isExistsCredential() bool {
	u, _ := user.Current()
	_, err := os.Stat(u.HomeDir + "/.config/hub")

	return err == nil
}

func getPreviousTag(tag string) string {
	out, err := exec.Command("git", "describe", "--abbrev=0", "--tags", tag+"^").CombinedOutput()
	if err != nil {
		return "" // No Tag
	}

	return strings.TrimRight(string(out), "\n")
}
