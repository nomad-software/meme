package image

import (
	"image"
	"image/png"
	"os"
	"path"

	"github.com/nomad-software/meme/output"
)

// Save the passed image to disk.
func Save(img image.Image) string {
	name := fileName()

	file, err := os.Create(name)
	output.OnError(err, "Could not create image file")
	defer file.Close()

	png.Encode(file, img)

	return name
}

// Generate a temporary file name.
func fileName() string {
	dir := os.TempDir()
	return path.Join(dir, "meme.png")
}
