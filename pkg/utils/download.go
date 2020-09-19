package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url, destinationFile string, overwrite bool) (bool, error) {
	if PathExists(destinationFile) && !overwrite {
		return false, nil
	}

	TryRemove(destinationFile)

	response, err := http.Get(url) //nolint:gosec
	if err != nil {
		return true, err
	}
	if response.StatusCode != 200 {
		return true, fmt.Errorf("unexpected status while retrieving release file: %s", response.Status)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	directory, _ := filepath.Split(destinationFile)
	if err := os.MkdirAll(directory, 0755); err != nil {
		return true, err
	}

	file, err := os.Create(destinationFile)
	if err != nil {
		return true, err
	}

	defer func() {
		_ = file.Close()
	}()

	if _, err := io.Copy(file, response.Body); err != nil {
		return true, err
	}

	return true, nil
}
