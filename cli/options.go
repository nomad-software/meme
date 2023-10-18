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
	ImageIds []string
)

// Initialise the package.
func init() {
	images, err := data.Files.ReadDir(data.ImagePath)
	output.OnError(err, "Could not read embedded images")

	for _, image := range images {
		id := strings.TrimSuffix(filepath.Base(image.Name()), data.ImageExtension)
		ImageIds = append(ImageIds, id)
	}

	sort.Strings(ImageIds)
}

// Options holds the options passed on the command line.
type Options struct {
	Gif           bool
	Bottom        string
	ClientID      string
	Help          bool
	Image         string
	ImageType     string
	OutName       string
	Shake         bool
	Top           string
	Trigger       bool
	ListTemplates bool
}

// ParseOptions parses the command line options.
func ParseOptions() Options {
	var opt Options
	var text string

	flag.BoolVar(&opt.Help, "h", false, "Show help.\n")
	flag.BoolVar(&opt.Help, "help", false, "Show help.\n")
	flag.StringVar(&opt.ClientID, "cid", "", "The client id of an application registered with imgur.com.\nIf specified, the new meme will be uploaded to imgur.com.\n(See README for full details.)\n")
	flag.StringVar(&opt.Image, "i", "", "A built-in template, a URL or the path to a local file.\nYou can also use '-' to read an image from stdin.\n")
	flag.StringVar(&opt.OutName, "o", "", "The optional name of the output file.\nIf omitted, a temporary file will be created.\n")
	flag.StringVar(&text, "t", "", "The meme text. Separate the top and bottom banners using a pipe '|'.\n")
	flag.BoolVar(&opt.Gif, "gif", false, "Gif animations will be preserved and the output will be a gif.\nDoes nothing for other image types.\n")
	flag.BoolVar(&opt.Shake, "shake", false, "Shake the image to intensify it. Always outputs a gif.\n")
	flag.BoolVar(&opt.Trigger, "trigger", false, "Shake the image and add a triggered banner. Always outputs a gif.\n")
	flag.BoolVar(&opt.ListTemplates, "list-templates", false, "List all of the built in templates.\n")
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

	if !(opt.Gif || opt.Trigger || opt.Shake) && opt.OutName != "" {
		if !strings.HasSuffix(strings.ToLower(opt.OutName), ".png") {
			output.Error("The output file name must have the suffix of .png")
		}
	}

	if (opt.Gif || opt.Trigger || opt.Shake) && opt.OutName != "" {
		if !strings.HasSuffix(strings.ToLower(opt.OutName), ".gif") {
			output.Error("The output file name must have the suffix of .gif")
		}
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
	for x, name := range ImageIds {
		if (x % 2) == 0 {
			fmt.Fprint(output.Stdout, color.CyanString("    %-30s", name))
		} else {
			fmt.Fprintln(output.Stdout, color.CyanString("%s", name))
		}
	}

	if len(ImageIds)%3 != 0 {
		fmt.Println("")
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
