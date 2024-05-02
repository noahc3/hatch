package log

import (
	"fmt"
	"os"
	"strings"
)

var Quiet bool = false
var DebugOutput bool = false

func Info(format string, args ...interface{}) {
	if !Quiet {
		fmt.Fprintf(os.Stdout, format, args...)
	}
}

func Error(format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)

	if !Quiet {
		fmt.Fprintf(os.Stderr, msg)
	}

	return fmt.Errorf(strings.Trim(msg, " \n"))
}

func Debug(format string, args ...interface{}) {
	if DebugOutput {
		fmt.Fprintf(os.Stdout, format, args...)
	}
}

func Nl() {
	if !Quiet {
		fmt.Println()
	}
}
