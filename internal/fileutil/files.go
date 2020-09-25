package fileutil

import (
	"os"
)

// PathExists is a function that checks if a given file or directory exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// TyRemove is a function that removes any given file or directory, without throwing errors.
// If an error did occur, this function returns false instead of throwing an error.
func TryRemove(path string) bool {
	return PathExists(path) && os.RemoveAll(path) == nil
}
