package logging

import (
	"fmt"
	"os"
)

const (
	errorExitCode = 2
)

// The Printf function logs any string to system out.
// It provides the same formatting as the fmt package does.
func Printf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// The Errorf function logs any string to system err and causes the application to exit.
// It provides the same formatting as the fmt package does.
func Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(errorExitCode)
}

// The TaskPrintf function logs any string to system out, but indenting it.
// It provides the same formatting as the fmt package does.
func TaskPrintf(format string, args ...interface{}) {
	fmt.Printf("[+] "+format+"\n", args...)
}

// The TaskErrorf function logs any string to system err, but indenting it, and causes the application to exit.
// It provides the same formatting as the fmt package does.
func TaskErrorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "[-] "+format+"\n", args...)
	os.Exit(errorExitCode)
}

// The IfError function logs any error to system err and causes the application to exit, if the given condition matches.
func IfError(err error) {
	IfErrorf(err != nil, "%s", err)
}

// The Errorf function logs any string to system err and causes the application to exit, if the given condition matches.
// It provides the same formatting as the fmt package does.
func IfErrorf(condition bool, format string, args ...interface{}) {
	if condition {
		Errorf(format, args...)
	}
}

// The IfError function logs any error to system err, but indenting it, and causes the application to exit, if the given
// condition matches.
func IfTaskError(err error) {
	IfTaskErrorf(err != nil, "%s", err)
}

// The TaskErrorf function logs any string to system err, but indenting it, and causes the application to exit, if the given
// condition matches. It provides the same formatting as the fmt package does.
func IfTaskErrorf(condition bool, format string, args ...interface{}) {
	if condition {
		TaskErrorf(format, args...)
	}
}
