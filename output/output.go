package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func OnError(err error, text string) {
	if err != nil {
		fmt.Fprintln(os.Stderr, color.RedString(text+": %s", err.Error()))
		os.Exit(1)
	}
}

func Error(text string) {
	fmt.Fprintln(os.Stderr, color.RedString(text))
	os.Exit(1)
}

func Infoln(format string, args ...interface{}) {
	fmt.Printf(color.GreenString(format)+"\n", args...)
}
