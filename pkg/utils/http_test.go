package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJSON(t *testing.T) {
	var collection interface{} = nil
	assert.Error(t, GetJSON("http://example.org/404.txt", &collection))
	assert.Nil(t, collection)

	assert.NoError(t, GetJSON("https://golang.org/dl/?mode=json", &collection))
	assert.NotNil(t, collection)
}

func TestGetFile(t *testing.T) {
	tempDir := t.TempDir()
	destinationFile := filepath.Join(tempDir, "destination.txt")

	downloaded, err := GetFile("http://example.org/404.txt", destinationFile, false)
	assert.Error(t, err)
	assert.True(t, downloaded)

	downloaded, err = GetFile("https://golang.org/dl/?mode=json", destinationFile, false)
	assert.NoError(t, err)
	assert.True(t, downloaded)

	downloaded, err = GetFile("https://golang.org/dl/?mode=json", destinationFile, false)
	assert.NoError(t, err)
	assert.False(t, downloaded)

	downloaded, err = GetFile("https://golang.org/dl/?mode=json", destinationFile, true)
	assert.NoError(t, err)
	assert.True(t, downloaded)
}
