package font

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nomad-software/meme/cli"
	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/output"
)

var (
	// Path is the location of the font file.
	Path string
)

// Override the font path at runtime.
func SetPath(opt cli.Options) {
	if opt.Font == "" {
		return
	}

	// direct file
	if _, err := os.Stat(opt.Font); err == nil {
		Path = opt.Font
		return
	}

	// try fc-match (Linux/macOS with fontconfig)
	if fc, err := exec.LookPath("fc-match"); err == nil {
		out, err := exec.Command(fc, "-f", "%{file}\\n", opt.Font).Output()
		if err == nil {
			f := strings.TrimSpace(string(out))
			if f != "" {
				if _, err := os.Stat(f); err == nil {
					Path = f
					return
				}
			}
		}
	}

	output.Error(fmt.Sprintf("Invalid font: %s", opt.Font))
}

// Write the embedded font to the temporary directory.
func init() {
	if Path != "" {
		return
	}

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
