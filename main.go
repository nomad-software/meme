package main

import (
	"fmt"

	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/image"
)

func main() {

	options := cli.ParseOptions()

	if options.Help {
		options.PrintUsage()

	} else if options.Valid() {

		image := image.Load(options.Image)
		fmt.Printf("%v\n", image.Bounds())
	}
}
