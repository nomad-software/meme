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

	if _, err := os.Stat(Path); os.IsNotExist(err) {
		file, err := os.Create(Path)
		output.OnError(err, "Could not create font file")
		defer file.Close()

		stream, err := data.Files.ReadFile(data.Font)
		output.OnError(err, "Could not read embedded font")

		_, err = file.Write(stream)
		output.OnError(err, "Could not write font file")
	}
}
