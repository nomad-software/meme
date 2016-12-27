package image

import (
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/output"
)

// Save the passed image to disk.
func Save(opt cli.Options, img image.Image) string {
	var name string

	if opt.Name != "" {
		name = opt.Name
	} else {
		name = tempName()
	}

	file, err := os.Create(name)
	output.OnError(err, "Could not create image file")
	defer file.Close()

	png.Encode(file, img)

	return name
}

// Generate a temporary file name.
func tempName() string {
	dir := os.TempDir()
	return filepath.Join(dir, "meme.png")
}
