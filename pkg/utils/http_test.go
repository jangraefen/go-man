package utils

import (
	"errors"
	"io/ioutil"
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

func TestStaticResponseClient(t *testing.T) {
	sut := StaticResponseClient(404, []byte("not found"), nil)

	response, err := sut.Get("http://example.org")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 404, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	assert.NoError(t, err)
	assert.Equal(t, []byte("not found"), body)

	sut = StaticResponseClient(0, nil, errors.New("failure"))

	response, err = sut.Get("http://example.org")
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, `Get "http://example.org": failure`, err.Error())
}
