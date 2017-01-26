package cli

import (
	"flag"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/output"
)

var (
	imageIds []string
)

// Initialise the package.
func init() {
	for _, asset := range data.AssetNames() {
		if strings.HasPrefix(asset, data.ImagePath) {
			id := strings.TrimSuffix(filepath.Base(asset), data.ImageExtension)
			imageIds = append(imageIds, id)
		}
	}

	sort.Sort(sort.StringSlice(imageIds))
}

// Options holds the options passed on the command line.
type Options struct {
	Anim      bool
	Bottom    string
	ClientID  string
	Help      bool
	Image     string
	ImageType string
	MaxAnim   bool
	OutName   string
	Top       string
}

// ParseOptions parses the command line options.
func ParseOptions() Options {
	var opt Options
	var text string

	flag.BoolVar(&opt.Help, "h", false, "Show help.")
	flag.BoolVar(&opt.Help, "help", false, "Show help.")
	flag.StringVar(&opt.ClientID, "cid", "", "The client id of an application registered with imgur.com. If specified, the new meme will be uploaded to imgur.com instead of being saved locally. (See README for full details.)")
	flag.StringVar(&opt.Image, "i", "", "One of the built-in templates, a URL or the path to a local file (gif, jpeg or png.) You can also use '-' to read an image from stdin.")
	flag.StringVar(&opt.OutName, "o", "", "The optional name of the output file (png). If omitted, a temporary file will be used.")
	flag.StringVar(&text, "t", "", "The meme text. Separate the top and bottom banners using a pipe character.")
	flag.BoolVar(&opt.Anim, "gif", false, "Assume the image is a gif. Animations will be preserved and the output will be a gif.")
	flag.BoolVar(&opt.MaxAnim, "max", false, "Used with -gif option to render potential reductions at max quality. This is at the expense of file size.")
	flag.Parse()

	parsed := strings.Split(text, "|")
	if len(parsed) == 1 {
		opt.Top = parsed[0]
	} else {
		opt.Top = parsed[0]
		opt.Bottom = parsed[1]
	}

	return opt
}

// Valid validates the command line options and returns true if they are valid,
// false if not.
func (opt *Options) Valid() bool {

	if opt.Image == "" {
		output.Error("An image is required")
	}

	if !opt.Anim && opt.OutName != "" {
		if !strings.HasSuffix(strings.ToLower(opt.OutName), ".png") {
			output.Error("The output file name must have the suffix of .png")
		}
	}

	if opt.Anim && opt.OutName != "" {
		if !strings.HasSuffix(strings.ToLower(opt.OutName), ".gif") {
			output.Error("The output file name must have the suffix of .gif")
		}
	}

	if (opt.Top + opt.Bottom) == "" {
		output.Error("At least one piece of text is required")
	}

	return true
}

// PrintUsage prints who to use this command.
func (opt *Options) PrintUsage() {
	var banner = ` _ __ ___   ___ _ __ ___   ___
| '_ ' _ \ / _ \ '_ ' _ \ / _ \
| | | | | |  __/ | | | | |  __/
|_| |_| |_|\___|_| |_| |_|\___|

`
	color.Green(banner)
	fmt.Println("")
	flag.Usage()
	fmt.Println("")

	fmt.Println("  Templates")
	fmt.Println("")
	for _, name := range imageIds {
		color.Cyan("    " + name)
	}
	fmt.Println("")

	fmt.Println("  Examples")
	fmt.Println("")
	color.Cyan("    meme -i kirk-khan -t \"|khaaaan\"")
	color.Cyan("    meme -i brace-yourselves -t \"Brace yourselves|The memes are coming!\"")
	color.Cyan("    meme -i http://i.imgur.com/FsWetC0.jpg -t \"|China\"")
	color.Cyan("    meme -i ~/Pictures/face.png -t \"Hello\"")
	fmt.Println("")
}
