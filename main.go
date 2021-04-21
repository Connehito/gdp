package main

import (
	"log"
	"os"
)

// Version is GDP's version.
const Version string = "v0.2.6"

// Usage is GDP's usage.
const Usage string = `gdp is a CLI tool for pushing the tag associated with deployment and publishing the release note in GitHub.

Usage:
  gdp <command> [-t | --tag <TAG>] [-d | --dry-run] [-f | --force]

Available Commands:
  deploy   Add the tag to local repository and push the tag to remote(origin) repository
  publish  Create the release note in GitHub which based on the merge commits of the tag

Flags:
  -d, --dry-run  dry-run gdp
  -t, --tag      specify tag at semantic(e.g. v1.2.3 or 1.2.3) or date(e.g. 20180525.1 or release_20180525) format
  -f, --force    run gdp without validation
  -h, --help     help for gdp
  -v, --version  confirm gdp version

Example Usage:
  gdp deploy -t TAG -d   specify tag and dry-run
  gdp publish -t TAG -f  force(skipped validation)
  gdp deploy/publish     set tag automatically

Further Help:
  https://github.com/Connehito/gdp`

func main() {
	gdp, err := NewCommand()
	if err != nil {
		log.Fatal(err)
	}

	cli := &CLI{
		outStream: os.Stdout,
		errStream: os.Stderr,
		gdp:       gdp,
	}

	os.Exit(cli.Run(os.Args))
}
