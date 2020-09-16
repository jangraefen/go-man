package tasks

import (
	"fmt"
	"io"
	"os"
)

type Task struct {
	ErrorExitCode int
	Output        io.Writer
	Error         io.Writer
}

// Printf is a function that logs any string to system out.
// It provides the same formatting as the fmt package does.
func (t Task) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.Output, format+"\n", args...)
}

// Dief is a function that logs any string to system err and causes the application to exit.
// It provides the same formatting as the fmt package does.
func (t Task) Dief(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.Error, format+"\n", args...)
	os.Exit(t.ErrorExitCode)
}

// SubPrintf is a function that logs any string to system out, but indenting it.
// It provides the same formatting as the fmt package does.
func (t Task) SubPrintf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.Output, "[+] "+format+"\n", args...)
}

// SubDief is a function that logs any string to system err, but indenting it, and causes the application to exit.
// It provides the same formatting as the fmt package does.
func (t Task) SubDief(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.Error, "[-] "+format+"\n", args...)
	os.Exit(t.ErrorExitCode)
}

// DieOnError is a function that logs any error to system err and causes the application to exit, if a condition matches.
func (t Task) DieOnError(err error) {
	t.DieIff(err != nil, "%s", err)
}

// DieIff is a function that logs any string to system err and causes the application to exit, if a condition matches.
// It provides the same formatting as the fmt package does.
func (t Task) DieIff(condition bool, format string, args ...interface{}) {
	if condition {
		t.Dief(format, args...)
	}
}

// SubDieOnError is a function that logs any error to system err, but indenting it, and causes the application to exit, if a condition matches.
func (t Task) SubDieOnError(err error) {
	t.SubDieIff(err != nil, "%s", err)
}

// SubDieIff is a function that logs any string to system err, but indenting it, and causes the application to exit, if a matches.
// It provides the same formatting as the fmt package does.
func (t Task) SubDieIff(condition bool, format string, args ...interface{}) {
	if condition {
		t.SubDief(format, args...)
	}
}
