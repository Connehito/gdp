package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mitchellh/colorstring"
	"io"
	"os"
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

// Safety Hour.
const (
	SafetyHourStart  = 9
	SafetyHourEnd = 19
)

// Run invokes deploy and publish's process.
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
	if !force && !validate(cli, subCommand, tag) {
		return ExitError
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
		if !isSafetyHour() {
			fmt.Fprintln(cli.outStream, "It's past the regular time. Is this a hot-fix release?")
			fmt.Fprint(cli.outStream, "> ")
			if !yesOrNo(cli) {
				return ExitError
			}
		}

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

func validate(cli *CLI, subCommand string, tag string) bool {
	if subCommand == CommandDeploy {
		if !cli.gdp.IsMasterBranch() {
			printError(cli.errStream, fmt.Sprintf("Branch is not master."))
			return false
		}
		if cli.gdp.IsExistTagInLocal(tag) {
			printError(cli.errStream, fmt.Sprintf("Tag is already exist in local."))
			return false
		}
	} else {
		if !cli.gdp.IsExistTagInRemote(tag) {
			printError(cli.errStream, fmt.Sprintf("Tag is not exist in remote."))
			return false
		}
	}

	return true
}

func isSafetyHour() bool {
	t := Now()
	if t.Hour() >= SafetyHourStart && t.Hour() < SafetyHourEnd {
		return true
	}

	return false
}

func yesOrNo(cli *CLI) bool {
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadByte()
	if err != nil {
		return false
	}

	if s == []byte("Y")[0] || s == []byte("y")[0] {
		fmt.Fprintln(cli.outStream, "OK. Take time.")
		return true
	} else if s == []byte("N")[0] || s == []byte("n")[0] {
		printError(cli.errStream, fmt.Sprintf("Good choice."))
		return false
	}

	printError(cli.errStream, fmt.Sprintf("Please enter y or n."))
	return false
}
