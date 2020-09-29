package fileutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// PathExists is a function that checks if a given file or directory exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// TryRemove is a function that removes any given file or directory, without throwing errors.
// If an error did occur, this function returns false instead of throwing an error.
func TryRemove(path string) bool {
	return PathExists(path) && os.RemoveAll(path) == nil
}

// MoveDirectory moves a directory from one path to another recursively.
func MoveDirectory(fromDirectory, toDirectory string) error {
	if !PathExists(toDirectory) {
		if err := os.MkdirAll(toDirectory, 0755); err != nil {
			return err
		}
	}

	fileInfoList, err := ioutil.ReadDir(fromDirectory)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfoList {
		from := filepath.Join(fromDirectory, fileInfo.Name())
		to := filepath.Join(toDirectory, fileInfo.Name())

		if fileInfo.IsDir() {
			if err := MoveDirectory(from, to); err != nil {
				return err
			}
		} else {
			if err := os.Rename(from, to); err != nil {
				return err
			}
		}
	}

	return nil
}
