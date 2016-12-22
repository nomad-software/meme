package cli

import (
	"flag"
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/output"
)

var (
	ImageIds []string
)

// Initialise the package.
func init() {
	for _, asset := range data.AssetNames() {
		if strings.HasPrefix(asset, data.IMAGE_PATH) {
			id := strings.TrimSuffix(path.Base(asset), data.IMAGE_EXTENSION)
			ImageIds = append(ImageIds, id)
		}
	}

	sort.Sort(sort.StringSlice(ImageIds))
}

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
		output.Error("An image is required")
	}

	if (this.Top + this.Bottom) == "" {
		output.Error("At least one piece of text is required")
	}

	return true
}

func (this *Options) PrintUsage() {
	var banner string = ` _ __ ___   ___ _ __ ___   ___
| '_ ' _ \ / _ \ '_ ' _ \ / _ \
| | | | | |  __/ | | | | |  __/
|_| |_| |_|\___|_| |_| |_|\___|

`
	color.Cyan(banner)
	fmt.Println("  Default images:")
	fmt.Println("")

	for _, name := range ImageIds {
		color.Green("    " + name)
	}

	fmt.Println("")
	flag.Usage()
}
