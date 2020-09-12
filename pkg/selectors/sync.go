package selectors

import (
	"path/filepath"
)

func SyncToCurrent(directory, versionFolderName string) error {
	versionDirectory := filepath.Join(directory, versionFolderName)
	currentDirectory := filepath.Join(directory, "go-current")

	return symlink(versionDirectory, currentDirectory)
}
