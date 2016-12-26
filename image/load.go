package image

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/output"
)

var (
	imageMap map[string]string = make(map[string]string)
)

// Initialise the package.
func init() {
	for _, asset := range data.AssetNames() {
		if strings.HasPrefix(asset, data.IMAGE_PATH) {
			id := strings.TrimSuffix(path.Base(asset), data.IMAGE_EXTENSION)
			imageMap[id] = asset
		}
	}
}

// Load an image from the passed string.
// The string can be a embedded asset id, an image URL or a local file.
func Load(path string) image.Image {
	if isAsset(path) {
		return load(path)
	}

	if isUrl(path) {
		return download(path)
	}

	if isLocal(path) {
		return read(path)
	}

	if isStdin(path) {
		return loadStdin()
	}

	output.Error("Image not recognised")
	panic("Never reached")
}

// Return true if the passed string is an embedded asset id, false if not.
func isAsset(id string) bool {
	_, ok := imageMap[id]
	return ok
}

// Load and return an embedded asset (image) by id.
// The id is assumed to exist.
func load(id string) image.Image {
	asset, _ := imageMap[id]
	stream, _ := data.Asset(asset)

	return decode(bytes.NewReader(stream))
}

// return true if the passed string is '-'
func isStdin(path string) bool {
	return path == "-"
}

func loadStdin() image.Image {
	var data []byte
	data, err := ioutil.ReadAll(os.Stdin)
	output.OnError(err, "Stdin read error")
	return decode(bytes.NewReader(data))
}

// Return true if the passed string is an image URL, false if not.
func isUrl(url string) bool {
	return strings.HasPrefix(url, "http")
}

// Download the image located at the passed image URL, decode and return it.
func download(url string) image.Image {
	res, err := http.Get(url)
	output.OnError(err, "Request error")
	defer res.Body.Close()

	if res.StatusCode != 200 {
		output.Error("Could not access URL")
	}

	return decode(res.Body)
}

// Return true if the passed string is a file that exists on the local
// filesystem, false if not.
func isLocal(path string) bool {
	path, err := homedir.Expand(path)
	output.OnError(err, "Could not expand path")

	_, err = os.Stat(path)
	return err == nil
}

// Read and return a file on the local filesystem.
// The file is assumed to exist.
func read(path string) image.Image {
	path, err := homedir.Expand(path)
	output.OnError(err, "Could not expand path")

	stream, err := ioutil.ReadFile(path)
	output.OnError(err, "Could not read local file")

	return decode(bytes.NewReader(stream))
}

// Decode the passed byte stream and return an image.
func decode(r io.Reader) image.Image {
	img, _, err := image.Decode(r)
	output.OnError(err, "Could not decode image")
	return img
}
