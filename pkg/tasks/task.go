package tasks

import (
	"io"
)

type Task struct {
	ErrorExitCode int
	Output        io.Writer
	Error         io.Writer
	indention     uint
}

// Step is a function that returns a sub-task for of the receiving Task.
func (t Task) Step() *Task {
	return &Task{
		ErrorExitCode: t.ErrorExitCode,
		Output:        t.Output,
		Error:         t.Error,
		indention:     t.indention + 1,
	}
}
