package font

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nomad-software/meme/data"
	"github.com/nomad-software/meme/output"
)

var (
	// Path is the location of the font file.
	Path string
)

// SetPath overrides the font path at runtime. If the file exists it will be used directly.
func SetPath(p string) error {
	if p == "" {
		return nil
	}

	// direct file
	if _, err := os.Stat(p); err == nil {
		Path = p
		return nil
	}

	// try fc-match (Linux/macOS with fontconfig)
	if fc, err := exec.LookPath("fc-match"); err == nil {
		out, err := exec.Command(fc, "-f", "%{file}\\n", p).Output()
		if err == nil {
			f := strings.TrimSpace(string(out))
			if f != "" {
				if _, err := os.Stat(f); err == nil {
					Path = f
					return nil
				}
			}
		}
	}

	return fmt.Errorf("font not found: %s", p)
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
