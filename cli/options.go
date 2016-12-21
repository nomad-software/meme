package cli

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/nomad-software/meme/data"
)

type Options struct {
	Image  string
	Top    string
	Bottom string
	Help   bool
}

func ParseOptions() Options {
	var opt Options

	flag.StringVar(&opt.Image, "img", "", "One of the above default images, a URL or a local file.")
	flag.StringVar(&opt.Top, "top", "", "The text at the top of the meme.")
	flag.StringVar(&opt.Bottom, "btm", "", "The text at the bottom of the meme.")
	flag.BoolVar(&opt.Help, "help", false, "Show help.")
	flag.Parse()

	return opt
}

func (this *Options) Valid() bool {

	if this.Image == "" {
		fmt.Fprintln(os.Stderr, color.RedString("An image is required."))
		return false
	}

	if (this.Top + this.Bottom) == "" {
		fmt.Fprintln(os.Stderr, color.RedString("At least one piece of text is required."))
		return false
	}

	return true
}

func (this *Options) PrintUsage() {
	var banner string = ` _ __ ___   ___ _ __ ___   ___
| '_ ' _ \ / _ \ '_ ' _ \ / _ \
| | | | | |  __/ | | | | |  __/
|_| |_| |_|\___|_| |_| |_|\___|

`
	images := data.AssetNames()
	sort.Sort(sort.StringSlice(images))

	color.Cyan(banner)
	fmt.Println("  Default images:")
	fmt.Println("")

	for _, name := range images {
		if strings.HasPrefix(name, data.IMAGE_PATH) {
			name = path.Base(name)
			name = strings.TrimSuffix(name, data.IMAGE_EXTENSION)
			color.Green("    " + name)
		}
	}

	fmt.Println("")
	flag.Usage()
}
