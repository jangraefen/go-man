package utils

import (
	"github.com/mholt/archiver/v3"
)

// ExtractArchive is a function that extracts any given archive file into a given destination directory.
// If set to overwrite, any existing files in the destination directory are deleted before extracting the archive.
// The functions returns false if nothing was done because it should not overwrite and the destination directory already
// existed.
func ExtractArchive(archiveFile, destinationDirection string, overwrite bool) (bool, error) {
	if PathExists(destinationDirection) && !overwrite {
		return false, nil
	}

	TryRemove(destinationDirection)

	return true, archiver.Unarchive(archiveFile, destinationDirection)
}
