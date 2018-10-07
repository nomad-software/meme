package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	// Stdout is a color friendly pipe.
	Stdout = colorable.NewColorableStdout()

	// Stderr is a color friendly pipe.
	Stderr = colorable.NewColorableStderr()
)

// OnError prints an error if err is not nil and exits the program.
func OnError(err error, text string) {
	if err != nil {
		fmt.Fprintln(Stderr, color.RedString(text+": %s", err.Error()))
		os.Exit(1)
	}
}

// Error prints an error and exits the program.
func Error(text string) {
	fmt.Fprintln(Stderr, color.RedString(text))
	os.Exit(1)
}

// Info prints information.
func Info(format string, args ...interface{}) {
	fmt.Fprintf(Stdout, color.GreenString(format)+"\n", args...)
}
