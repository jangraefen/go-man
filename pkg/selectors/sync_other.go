// +build !windows

package selectors

import (
	"os"
)

func symlink(sourceDirectory, targetDirectory string) error {
	return os.Symlink(sourceDirectory, targetDirectory)
}
