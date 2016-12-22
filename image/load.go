package image

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/output"
)

var (
	imageMap map[string]string = make(map[string]string)
	ImageIds []string
)

func init() {
	for _, asset := range data.AssetNames() {
		if strings.HasPrefix(asset, data.IMAGE_PATH) {
			id := strings.TrimSuffix(path.Base(asset), data.IMAGE_EXTENSION)
			imageMap[id] = asset
			ImageIds = append(ImageIds, id)
		}
	}

	sort.Sort(sort.StringSlice(ImageIds))
}

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

	output.Error("Image not recognised")
	panic("Not reached")
}

func isAsset(id string) bool {
	_, ok := imageMap[id]
	return ok
}

func load(id string) image.Image {
	output.Infoln("Loading: %s", id)

	asset, _ := imageMap[id]
	stream, _ := data.Asset(asset)
	img, _, err := image.Decode(bytes.NewReader(stream))
	output.OnError(err, "Could not decode image")

	return img
}

func isUrl(url string) bool {
	return strings.HasPrefix(url, "http")
}

func download(url string) image.Image {
	output.Infoln("Downloading: %s", url)

	res, err := http.Get(url)
	output.OnError(err, "Request error")
	defer res.Body.Close()

	if res.StatusCode != 200 {
		output.Error("Could not access URL")
	}

	img, _, err := image.Decode(res.Body)
	output.OnError(err, "Could not decode image")

	return img
}

func isLocal(path string) bool {
	path, err := homedir.Expand(path)
	output.OnError(err, "Could not expand path")

	_, err = os.Stat(path)
	return err == nil
}

func read(path string) image.Image {
	output.Infoln("Reading: %s", path)

	path, err := homedir.Expand(path)
	output.OnError(err, "Could not expand path")

	stream, err := ioutil.ReadFile(path)
	output.OnError(err, "Could not read local file")

	img, _, err := image.Decode(bytes.NewReader(stream))
	output.OnError(err, "Could not decode image")

	return img
}
