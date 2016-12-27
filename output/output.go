package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
)

var (
	stdout = colorable.NewColorableStdout()
	stderr = colorable.NewColorableStderr()
)

// Print an error and exit the program if the error is non nil.
func OnError(err error, text string) {
	if err != nil {
		fmt.Fprintln(stderr, color.RedString(text+": %s", err.Error()))
		os.Exit(1)
	}
}

// Print an error and exit the program.
func Error(text string) {
	fmt.Fprintln(stderr, color.RedString(text))
	os.Exit(1)
}

// Print information.
func Info(format string, args ...interface{}) {
	fmt.Fprintf(stdout, color.GreenString(format)+"\n", args...)
}
