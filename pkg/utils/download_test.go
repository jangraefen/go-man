package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadFile(t *testing.T) {
	tempDir := t.TempDir()
	destinationFile := filepath.Join(tempDir, "destination.txt")

	downloaded, err := DownloadFile("http://example.org/404.txt", destinationFile, false)
	assert.Error(t, err)
	assert.True(t, downloaded)

	downloaded, err = DownloadFile("https://golang.org/dl/?mode=json", destinationFile, false)
	assert.NoError(t, err)
	assert.True(t, downloaded)

	downloaded, err = DownloadFile("https://golang.org/dl/?mode=json", destinationFile, false)
	assert.NoError(t, err)
	assert.False(t, downloaded)

	downloaded, err = DownloadFile("https://golang.org/dl/?mode=json", destinationFile, true)
	assert.NoError(t, err)
	assert.True(t, downloaded)
}
