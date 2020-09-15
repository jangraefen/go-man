// +build !windows

package manager

import (
	"os"
)

func link(sourceDirectory, targetDirectory string) error {
	return os.Symlink(sourceDirectory, targetDirectory)
}

func unlink(directory string) error {
	return os.Remove(directory)
}
