package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/image"
	"github.com/nomad-software/meme/image/renderer"
	"github.com/nomad-software/meme/output"
)

func main() {

	options := cli.ParseOptions()

	if options.Help {
		options.PrintUsage()

	} else if options.Valid() {
		img := image.Load(options.Image)
		img = renderer.Render(options, img)

		if options.ClientID != "" {
			url := image.Upload(options, img)
			output.Info(url)
		} else {
			file := image.Save(options, img)
			output.Info(file)
		}
	}
}
