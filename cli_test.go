package main

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

// Tests for flag
func TestRun_Version(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &Command{},
	}

	args := strings.Split("gdp -v", " ")
	code := cli.Run(args)
	if code != ExitSuccess {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitSuccess)
	}

	expected := fmt.Sprintf("gdp version %s\n", Version)
	if out.String() != expected {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

func TestRun_Help(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &Command{},
	}

	args := strings.Split("gdp -h", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := Usage
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

func TestRun_FewArg(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &Command{},
	}

	args := strings.Split("gdp", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitSuccess)
	}

	expected := "Too few argument."
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

func TestRun_ManyArg(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &Command{},
	}

	args := strings.Split("gdp deploy publish test", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "Too many argument."
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

func TestRun_NotSupportedArg(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &Command{},
	}

	args := strings.Split("gdp deploy -s", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "flag provided but not defined: -s"
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

func TestRun_NoSubCommand(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &Command{},
	}

	args := strings.Split("gdp -t", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "flag needs an argument"
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

func TestRun_InvalidSubCommand(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &Command{},
	}

	args := strings.Split("gdp pub", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "Invalid sub command."
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

// Tests for deploy
type FakeGdpDeploy struct {
	Gdp
}

func (f *FakeGdpDeploy) IsMasterBranch() bool {
	return true
}

func (f *FakeGdpDeploy) IsExistTagInLocal(tag string) bool {
	return false
}

func (f *FakeGdpDeploy) GetLatestTag() string {
	return "v1.2.3"
}

func (f *FakeGdpDeploy) GetMergeCommitList(toTag string) (string, error) {
	list := "- itosho: initial commit\n"
	list = list + "- itosho: fix bug"
	return list, nil
}

func (f *FakeGdpDeploy) Deploy(tag string) error {
	return nil
}

func TestRun_Deploy(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpDeploy{},
	}
	Set(time.Date(2020, 4, 1, 17, 00, 00, 0, time.Local))

	args := strings.Split("gdp deploy -t v1.2.4", " ")
	code := cli.Run(args)
	if code != ExitSuccess {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "gdp deploy done."
	if !strings.Contains(out.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

func TestRun_DeployDryRun(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpDeploy{},
	}

	args := strings.Split("gdp deploy -t v1.2.4 -dry-run", " ")
	code := cli.Run(args)
	if code != ExitSuccess {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "gdp deploy done(dry-run mode)."
	if !strings.Contains(out.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

func TestRun_DeployNoSpecifiedTag(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpDeploy{},
	}
	Set(time.Date(2020, 4, 1, 17, 00, 00, 0, time.Local))

	args := strings.Split("gdp deploy", " ")

	code := cli.Run(args)
	if code != ExitSuccess {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "v1.2.4"
	if !strings.Contains(out.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

type FakeGdpDeployForce struct {
	Gdp
}

func (f *FakeGdpDeployForce) GetMergeCommitList(toTag string) (string, error) {
	list := "- itosho: initial commit\n"
	list = list + "- itosho: fix bug"
	return list, nil
}

func (f *FakeGdpDeployForce) Deploy(tag string) error {
	return nil
}

func TestRun_DeployForce(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpDeployForce{},
	}
	Set(time.Date(2020, 4, 1, 17, 00, 00, 0, time.Local))

	args := strings.Split("gdp deploy -t v1.2.4 -f", " ")
	code := cli.Run(args)
	if code != ExitSuccess {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "gdp deploy done."
	if !strings.Contains(out.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

type FakeGdpDeployNotMasterBranch struct {
	Gdp
}

func (f *FakeGdpDeployNotMasterBranch) IsMasterBranch() bool {
	return false
}

func TestRun_DeployNotMasterBranch(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpDeployNotMasterBranch{},
	}

	args := strings.Split("gdp deploy -t v1.2.4", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "Branch is not master."
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

type FakeGdpDeployExistTagInLocal struct {
	Gdp
}

func (f *FakeGdpDeployExistTagInLocal) IsMasterBranch() bool {
	return true
}

func (f *FakeGdpDeployExistTagInLocal) IsExistTagInLocal(tag string) bool {
	return true
}

func TestRun_DeployExistTagInLocal(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpDeployExistTagInLocal{},
	}

	args := strings.Split("gdp deploy -t v1.2.4", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "Tag is already exist in local."
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

type FakeGdpDeployErrorInGetNextTag struct {
	Gdp
}

func (f *FakeGdpDeployErrorInGetNextTag) IsMasterBranch() bool {
	return true
}

func (f *FakeGdpDeployErrorInGetNextTag) IsExistTagInLocal(tag string) bool {
	return false
}

func (f *FakeGdpDeployErrorInGetNextTag) GetLatestTag() string {
	return "v1.2.semantic"
}

func TestRun_DeployErrorInGetNextTag(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpDeployErrorInGetNextTag{},
	}

	args := strings.Split("gdp deploy", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "Getting release tag error"
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

type FakeGdpDeployErrorInGetMergeCommitList struct {
	Gdp
}

func (f *FakeGdpDeployErrorInGetMergeCommitList) IsMasterBranch() bool {
	return true
}

func (f *FakeGdpDeployErrorInGetMergeCommitList) IsExistTagInLocal(tag string) bool {
	return false
}

func (f *FakeGdpDeployErrorInGetMergeCommitList) GetMergeCommitList(toTag string) (string, error) {
	return "", errors.New("error occurred")
}

func TestRun_ErrorInGetMergeCommitList(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpDeployErrorInGetMergeCommitList{},
	}

	args := strings.Split("gdp deploy -t v1.2.4", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "Getting merge commit error"
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

type FakeGdpDeployErrorInDeploy struct {
	Gdp
}

func (f *FakeGdpDeployErrorInDeploy) IsMasterBranch() bool {
	return true
}

func (f *FakeGdpDeployErrorInDeploy) IsExistTagInLocal(tag string) bool {
	return false
}

func (f *FakeGdpDeployErrorInDeploy) GetMergeCommitList(toTag string) (string, error) {
	list := "- itosho: initial commit\n"
	list = list + "- itosho: fix bug"
	return list, nil
}

func (f *FakeGdpDeployErrorInDeploy) Deploy(tag string) error {
	return errors.New("error occurred")
}

func TestRun_DeployErrorInDeploy(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpDeployErrorInDeploy{},
	}
	Set(time.Date(2020, 4, 1, 17, 00, 00, 0, time.Local))

	args := strings.Split("gdp deploy -t v1.2.4", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "Deploy execution error"
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

// Tests for publish
type FakeGdpPublish struct {
	Gdp
}

func (f *FakeGdpPublish) IsExistTagInRemote(tag string) bool {
	return true
}

func (f *FakeGdpPublish) GetLatestTag() string {
	return "v1.2.3"
}

func (f *FakeGdpPublish) GetMergeCommitList(toTag string) (string, error) {
	list := "- itosho: initial commit\n"
	list = list + "- itosho: fix bug"
	return list, nil
}

func (f *FakeGdpPublish) Publish(tag string, commits string) error {
	return nil
}

func TestRun_Publish(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpPublish{},
	}

	args := strings.Split("gdp publish -t v1.2.3", " ")
	code := cli.Run(args)
	if code != ExitSuccess {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "gdp publish done."
	if !strings.Contains(out.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

func TestRun_PublishDryRun(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpPublish{},
	}

	args := strings.Split("gdp publish -t v1.2.3 -dry-run", " ")
	code := cli.Run(args)
	if code != ExitSuccess {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "gdp publish done(dry-run mode)."
	if !strings.Contains(out.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

func TestRun_PublishNoSpecifiedTag(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpPublish{},
	}

	args := strings.Split("gdp publish", " ")
	code := cli.Run(args)
	if code != ExitSuccess {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "v1.2.3"
	if !strings.Contains(out.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

type FakeGdpPublishForce struct {
	Gdp
}

func (f *FakeGdpPublishForce) GetMergeCommitList(toTag string) (string, error) {
	list := "- itosho: initial commit\n"
	list = list + "- itosho: fix bug"
	return list, nil
}

func (f *FakeGdpPublishForce) Publish(tag string, commits string) error {
	return nil
}

func TestRun_PublishForce(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpPublishForce{},
	}

	args := strings.Split("gdp publish -t v1.2.3 -f", " ")
	code := cli.Run(args)
	if code != ExitSuccess {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "gdp publish done."
	if !strings.Contains(out.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}

type FakeGdpPublishNotExistTagInRemote struct {
	Gdp
}

func (f *FakeGdpPublishNotExistTagInRemote) IsExistTagInRemote(tag string) bool {
	return false
}

func TestRun_PublishNotExistTagInRemote(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpPublishNotExistTagInRemote{},
	}

	args := strings.Split("gdp publish -t v1.2.3", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "Tag is not exist in remote."
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", err.String(), expected)
	}
}

type FakeGdpPublishErrorInPublish struct {
	Gdp
}

func (f *FakeGdpPublishErrorInPublish) IsExistTagInRemote(tag string) bool {
	return true
}

func (f *FakeGdpPublishErrorInPublish) GetMergeCommitList(toTag string) (string, error) {
	list := "- itosho: initial commit\n"
	list = list + "- itosho: fix bug"
	return list, nil
}

func (f *FakeGdpPublishErrorInPublish) Publish(tag string, commits string) error {
	return errors.New("error occurred")
}

func TestRun_PublishErrorInPublish(t *testing.T) {
	out, err := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{
		outStream: out,
		errStream: err,
		gdp:       &FakeGdpPublishErrorInPublish{},
	}

	args := strings.Split("gdp publish -t v1.2.3", " ")
	code := cli.Run(args)
	if code != ExitError {
		t.Errorf("ExitCode=%d, Expected=%d", code, ExitError)
	}

	expected := "Publish execution error"
	if !strings.Contains(err.String(), expected) {
		t.Errorf("Output=%q, Expected=%q", out.String(), expected)
	}
}
