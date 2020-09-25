// +build !windows

package manager

import (
	"fmt"
	"os"

	"github.com/NoizeMe/go-man/internal/fileutil"
)

func link(sourceDirectory, targetDirectory string) error {
	if fileutil.PathExists(targetDirectory) {
		return fmt.Errorf("%s: file or directory already exists", sourceDirectory)
	}

	return os.Symlink(sourceDirectory, targetDirectory)
}

func unlink(directory string) error {
	return os.Remove(directory)
}
