package main

import (
	"flag"
	"fmt"
	"io"

	"github.com/mitchellh/colorstring"
)

// CLI has stdout/stderr's writer and Gdp's interface.
type CLI struct {
	outStream io.Writer
	errStream io.Writer
	gdp       Gdp
}

// Exit code.
const (
	ExitSuccess = iota
	ExitError
)

// Sub command name.
const (
	CommandDeploy  = "deploy"
	CommandPublish = "publish"
)

// Run invokes depoloy and publish's process.
func (cli *CLI) Run(args []string) int {
	var version bool
	var dryRun bool
	var force bool
	var tag string

	flags := flag.NewFlagSet("gdp", flag.ContinueOnError)
	flags.SetOutput(cli.errStream)
	flags.Usage = func() {
		fmt.Fprintln(cli.errStream, Usage)
	}

	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")
	flags.BoolVar(&dryRun, "dry-run", false, "")
	flags.BoolVar(&dryRun, "d", false, "")
	flags.BoolVar(&force, "force", false, "")
	flags.BoolVar(&force, "f", false, "")
	flags.StringVar(&tag, "tag", "", "")
	flags.StringVar(&tag, "t", "", "")

	if len(args) < 2 {
		printError(cli.errStream, "Too few argument.")
		printError(cli.errStream, Usage)
		return ExitError
	}

	parseIndex := 1
	if args[1] == CommandDeploy || args[1] == CommandPublish {
		parseIndex++
	}
	if err := flags.Parse(args[parseIndex:]); err != nil {
		return ExitError
	}

	if version {
		fmt.Fprintln(cli.outStream, fmt.Sprintf("gdp version %s", Version))
		return ExitSuccess
	}

	parsedArgs := flags.Args()
	len := len(parsedArgs)
	if len > 1 {
		printError(cli.errStream, "Too many argument.")
		printError(cli.errStream, Usage)
		return ExitError
	}

	subCommand := args[1]
	if subCommand != CommandDeploy && subCommand != CommandPublish {
		printError(cli.errStream, "Invalid sub command.")
		printError(cli.errStream, Usage)
		return ExitError
	}

	if tag == "" {
		latestTag := cli.gdp.GetLatestTag()
		if subCommand == CommandDeploy {
			next, err := GetNextVersion(latestTag)
			if err != nil {
				printError(cli.errStream, fmt.Sprintf("Getting release tag error: %s.", err.Error()))
				return ExitError
			}
			latestTag = next
		}
		tag = latestTag
	}

	// validation
	if !force {
		if subCommand == CommandDeploy {
			if !cli.gdp.IsMasterBranch() {
				printError(cli.errStream, fmt.Sprintf("Branch is not master."))
				return ExitError
			}
			if cli.gdp.IsExistTagInLocal(tag) {
				printError(cli.errStream, fmt.Sprintf("Tag is already exist in local."))
				return ExitError
			}
		} else {
			if !cli.gdp.IsExistTagInRemote(tag) {
				printError(cli.errStream, fmt.Sprintf("Tag is not exist in remote."))
				return ExitError
			}
		}
	}

	toTag := "HEAD"
	if subCommand == CommandPublish {
		toTag = tag
	}

	// show release note
	list, err := cli.gdp.GetMergeCommitList(toTag)
	if err != nil {
		printError(cli.errStream, fmt.Sprintf("Getting merge commit error: %s.", err.Error()))
		return ExitError
	}

	note := GetReleaseNote(tag, list)
	fmt.Fprintln(cli.outStream, "The release note is as follows.")
	fmt.Fprintln(cli.outStream, "====================================")
	fmt.Fprintln(cli.outStream, note)
	fmt.Fprintln(cli.outStream, "====================================")

	if dryRun {
		printSuccess(cli.outStream, fmt.Sprintf("gdp %s done(dry-run mode).", subCommand))
		return ExitSuccess
	}

	// execution
	if subCommand == CommandDeploy {
		if err := cli.gdp.Deploy(tag); err != nil {
			printError(cli.errStream, fmt.Sprintf("Deploy execution error: %s.", err.Error()))
			return ExitError
		}
	} else {
		if err := cli.gdp.Publish(tag, note); err != nil {
			printError(cli.errStream, fmt.Sprintf("Publish execution error: %s.", err.Error()))
			return ExitError
		}
	}

	printSuccess(cli.outStream, fmt.Sprintf("gdp %s done.", subCommand))
	message := "Do not be satisfied with 'released', let's face user's feedback in sincerity!"
	printSuccess(cli.outStream, message)

	return ExitSuccess
}

func printSuccess(w io.Writer, message string, args ...interface{}) {
	message = fmt.Sprintf("[green]%s[reset]", message)
	fmt.Fprintln(w, colorstring.Color(fmt.Sprintf(message, args...)))
}

func printError(w io.Writer, message string, args ...interface{}) {
	message = fmt.Sprintf("[red]%s[reset]", message)
	fmt.Fprintln(w, colorstring.Color(fmt.Sprintf(message, args...)))
}
