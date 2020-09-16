package logging

import (
	"fmt"
	"os"
)

const (
	errorExitCode = 2
)

// Printf is a function that logs any string to system out.
// It provides the same formatting as the fmt package does.
func Printf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Errorf is a function that logs any string to system err and causes the application to exit.
// It provides the same formatting as the fmt package does.
func Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(errorExitCode)
}

// TaskPrintf is a function that logs any string to system out, but indenting it.
// It provides the same formatting as the fmt package does.
func TaskPrintf(format string, args ...interface{}) {
	fmt.Printf("[+] "+format+"\n", args...)
}

// TaskErrorf is a function that logs any string to system err, but indenting it, and causes the application to exit.
// It provides the same formatting as the fmt package does.
func TaskErrorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "[-] "+format+"\n", args...)
	os.Exit(errorExitCode)
}

// IfError is a function that logs any error to system err and causes the application to exit, if a condition matches.
func IfError(err error) {
	IfErrorf(err != nil, "%s", err)
}

// IfErrorf is a function that logs any string to system err and causes the application to exit, if a condition matches.
// It provides the same formatting as the fmt package does.
func IfErrorf(condition bool, format string, args ...interface{}) {
	if condition {
		Errorf(format, args...)
	}
}

// IfTaskError is a function that logs any error to system err, but indenting it, and causes the application to exit, if a condition matches.
func IfTaskError(err error) {
	IfTaskErrorf(err != nil, "%s", err)
}

// IfTaskErrorf is a function that logs any string to system err, but indenting it, and causes the application to exit, if a matches.
// It provides the same formatting as the fmt package does.
func IfTaskErrorf(condition bool, format string, args ...interface{}) {
	if condition {
		TaskErrorf(format, args...)
	}
}
