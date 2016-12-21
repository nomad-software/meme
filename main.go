package main

import "github.com/nomad-software/meme/cli"

func main() {

	options := cli.ParseOptions()

	if options.Help {
		options.PrintUsage()

	} else if options.Valid() {

	}
}
