package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/fatih/color"
	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/font"
	"github.com/nomad-software/meme/image"
	"github.com/nomad-software/meme/output"
)

func main() {
	opt := cli.ParseOptions()

	// If a font was supplied on the command line, try to use it.
	if opt.Font != "" {
		if err := font.SetPath(opt.Font); err != nil {
			output.OnError(err, "Invalid font file")
		}
	}

	if opt.Help {
		opt.PrintUsage()

	} else if opt.ListTemplates {
		for _, id := range cli.ImageIds {
			fmt.Fprintln(output.Stdout, color.CyanString("%s", id))
		}

	} else if opt.Valid() {
		st := image.Load(opt)
		st = image.RenderImage(opt, st)

		if opt.ClientID != "" {
			url := image.Upload(opt, st)
			output.Info(url)
		} else {
			file := image.Save(opt, st)
			output.Info(file)
		}
	}
}
