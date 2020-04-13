package main

import (
	"log"
	"os"
)

// Version is GDP's version.
const Version string = "v0.2.3"

// Usage is GDP's usage.
const Usage string = "usage: gdp deploy|publish [options]"

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
