package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Print an error and exit the program if the error is non nil.
func OnError(err error, text string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(text+": %s", err.Error()))
		os.Exit(1)
	}
}

// Print an error and exit the program.
func Error(text string) {
	fmt.Fprintln(os.Stderr, color.RedString(text))
	os.Exit(1)
}

// Print information.
func Info(format string, args ...interface{}) {
	fmt.Printf(color.GreenString(format)+"\n", args...)
}
