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

// Errorf is a function that logs any string to system err and causes the application to exit.
// It provides the same formatting as the fmt package does.
func (t Task) Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.Error, format+"\n", args...)
	os.Exit(t.ErrorExitCode)
}

// TaskPrintf is a function that logs any string to system out, but indenting it.
// It provides the same formatting as the fmt package does.
func (t Task) TaskPrintf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.Output, "[+] "+format+"\n", args...)
}

// TaskErrorf is a function that logs any string to system err, but indenting it, and causes the application to exit.
// It provides the same formatting as the fmt package does.
func (t Task) TaskErrorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.Error, "[-] "+format+"\n", args...)
	os.Exit(t.ErrorExitCode)
}

// IfError is a function that logs any error to system err and causes the application to exit, if a condition matches.
func (t Task) IfError(err error) {
	t.IfErrorf(err != nil, "%s", err)
}

// IfErrorf is a function that logs any string to system err and causes the application to exit, if a condition matches.
// It provides the same formatting as the fmt package does.
func (t Task) IfErrorf(condition bool, format string, args ...interface{}) {
	if condition {
		t.Errorf(format, args...)
	}
}

// IfTaskError is a function that logs any error to system err, but indenting it, and causes the application to exit, if a condition matches.
func (t Task) IfTaskError(err error) {
	t.IfTaskErrorf(err != nil, "%s", err)
}

// IfTaskErrorf is a function that logs any string to system err, but indenting it, and causes the application to exit, if a matches.
// It provides the same formatting as the fmt package does.
func (t Task) IfTaskErrorf(condition bool, format string, args ...interface{}) {
	if condition {
		t.TaskErrorf(format, args...)
	}
}
