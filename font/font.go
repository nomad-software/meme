package font

import (
	"os"
	"path"

	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/output"
)

var (
	Path string
)

// Write the embedded font to the temporary directory.
func init() {
	Path = path.Join(os.TempDir(), path.Base(data.FONT))

	_, err := os.Stat(Path)
	if err != nil {
		file, err := os.Create(Path)
		output.OnError(err, "Unable to create file")
		defer file.Close()

		stream, err := data.Asset(data.FONT)
		output.OnError(err, "Can not extract font")

		_, err = file.Write(stream)
		output.OnError(err, "Can not write to file")
	}
}
