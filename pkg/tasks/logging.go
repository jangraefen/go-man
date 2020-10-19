package tasks

import (
	"fmt"
	"os"
)

// Printf is a function that logs any string to system out.
// It provides the same formatting as the fmt package does.
func (t Task) Printf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.Output, t.logTemplate('+', format), args...)
}

// Fatalf is a function that logs any string to system err and causes the application to exit.
// It provides the same formatting as the fmt package does.
func (t Task) Fatalf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.Error, t.logTemplate('-', format), args...)
	os.Exit(t.ErrorExitCode)
}

// FatalOnError is a function that logs any error to system err and causes the application to exit, if a condition matches.
func (t Task) FatalOnError(err error) {
	t.FatalIff(err != nil, "%s", err)
}

// FatalIff is a function that logs any string to system err and causes the application to exit, if a condition matches.
// It provides the same formatting as the fmt package does.
func (t Task) FatalIff(condition bool, format string, args ...interface{}) {
	if condition {
		t.Fatalf(format, args...)
	}
}

func (t Task) logTemplate(prefixRune rune, format string) string {
	switch {
	case t.indention == 0:
		return format + "\n"
	case t.indention == 1:
		return fmt.Sprintf("[%c] %s\n", prefixRune, format)
	default:
		indentionString := " "
		for i := uint(1); i < t.indention; i++ {
			indentionString += "  "
		}
		return fmt.Sprintf("%s[%c] %s\n", indentionString, prefixRune, format)
	}
}
