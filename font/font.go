package font

import (
	"os"
	"path/filepath"

	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/output"
)

var (
	// Path is the location of the font file.
	Path string
)

// Write the embedded font to the temporary directory.
func init() {
	Path = filepath.Join(os.TempDir(), filepath.Base(data.Font))

	_, err := os.Stat(Path)
	if err != nil {
		file, err := os.Create(Path)
		output.OnError(err, "Could not create font file")
		defer file.Close()

		stream, err := data.Asset(data.Font)
		output.OnError(err, "Could not extract font")

		_, err = file.Write(stream)
		output.OnError(err, "Could not write font file")
	}
}
