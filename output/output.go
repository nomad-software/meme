package output

import (
	"fmt"
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
		Error(fmt.Sprintf(text+": %s", err.Error()))
	}
}

// Error prints an error and exits the program.
func Error(text string) {
	panic(text)
}

// Info prints information.
func Info(format string, args ...interface{}) {
	fmt.Fprintf(Stdout, color.GreenString(format)+"\n", args...)
}
