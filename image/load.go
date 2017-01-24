package image

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/image/stream"
	"github.com/nomad-software/meme/output"
)

var (
	imageMap = make(map[string]string)
)

// Initialise the package.
func init() {
	for _, asset := range data.AssetNames() {
		if strings.HasPrefix(asset, data.ImagePath) {
			id := strings.TrimSuffix(filepath.Base(asset), data.ImageExtension)
			imageMap[id] = asset
		}
	}
}

// Load an image from the passed string.
// The string can be a embedded asset id, an image URL or a local file.
func Load(opt cli.Options) stream.Stream {
	var s io.Reader

	if isAsset(opt.Image) {
		s = loadAsset(opt.Image)

	} else if isURL(opt.Image) {
		s = downloadURL(opt.Image)

	} else if isLocalFile(opt.Image) {
		s = readFile(opt.Image)

	} else if isStdin(opt.Image) {
		s = readStdin()

	} else {
		output.Error("Image not recognised")
	}

	return stream.NewStream(s)
}

// Return true if the passed string is an embedded asset id, false if not.
func isAsset(id string) bool {
	_, ok := imageMap[id]
	return ok
}

// Load and return an embedded asset (image) by id.
// The id is assumed to exist.
func loadAsset(id string) io.Reader {
	asset, _ := imageMap[id]
	st, _ := data.Asset(asset)

	return bytes.NewReader(st)
}

// Return true if the passed string is an image URL, false if not.
func isURL(url string) bool {
	return strings.HasPrefix(url, "http")
}

// Download the image located at the passed image URL, decode and return it.
func downloadURL(url string) io.Reader {
	res, err := http.Get(url)
	output.OnError(err, "Request error")
	defer res.Body.Close()

	if res.StatusCode != 200 {
		output.Error("Could not access URL")
	}

	st, err := ioutil.ReadAll(res.Body)
	output.OnError(err, "Could not read response body")

	return bytes.NewReader(st)
}

// Return true if the passed string is a file that exists on the local
// filesystem, false if not.
func isLocalFile(path string) bool {
	path, err := homedir.Expand(path)
	output.OnError(err, "Could not expand path")

	_, err = os.Stat(path)
	return err == nil
}

// Read and return a file on the local filesystem.
// The file is assumed to exist.
func readFile(path string) io.Reader {
	path, err := homedir.Expand(path)
	output.OnError(err, "Could not expand path")

	st, err := ioutil.ReadFile(path)
	output.OnError(err, "Could not read local file")

	return bytes.NewReader(st)
}

// return true if the passed string is '-' meaning we should read the image
// from stdin.
func isStdin(path string) bool {
	return path == "-"
}

// Read the image from stdin.
func readStdin() io.Reader {
	st, err := ioutil.ReadAll(os.Stdin)
	output.OnError(err, "Could not read stdin")

	return bytes.NewReader(st)
}
