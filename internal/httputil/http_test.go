package httputil

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetJSON(t *testing.T) {
	t.Cleanup(func() {
		Client = http.DefaultClient
	})

	var collection interface{} = nil

	assert.Error(t, GetJSON("http://example.org/404.txt", &collection))
	assert.Nil(t, collection)

	assert.NoError(t, GetJSON("https://golang.org/dl/?mode=json", &collection))
	assert.NotNil(t, collection)

	Client = StaticResponseClient(404, []byte("not found"), nil)
	collection = nil

	assert.Error(t, GetJSON("http://example.org/404.txt", &collection))
	assert.Nil(t, collection)

	Client = StaticResponseClient(200, []byte("not json"), nil)

	assert.Error(t, GetJSON("http://example.org/404.txt", &collection))
	assert.Nil(t, collection)

	Client = StaticResponseClient(0, nil, errors.New("failure"))

	assert.Error(t, GetJSON("http://example.org/404.txt", &collection))
	assert.Nil(t, collection)
}

func TestGetFile(t *testing.T) {
	tempDir := t.TempDir()
	destinationFile := filepath.Join(tempDir, "destination.txt")
	rootFile := getNoPermissionDirectory("destination.txt")

	t.Cleanup(func() {
		Client = http.DefaultClient
	})
	t.Cleanup(func() {
		_ = os.Remove(rootFile)
	})

	downloaded, err := GetFile("http://example.org/404.txt", destinationFile, false)
	assert.Error(t, err)
	assert.True(t, downloaded)

	downloaded, err = GetFile("https://golang.org/dl/?mode=json", rootFile, false)
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

	Client = StaticResponseClient(0, nil, errors.New("failure"))

	downloaded, err = GetFile("https://golang.org/dl/?mode=json", destinationFile, true)
	assert.Error(t, err)
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

	response, err = sut.Get("http://example.org") //nolint:bodyclose
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, `Get "http://example.org": failure`, err.Error())
}

func getNoPermissionDirectory(fileName string) string {
	if runtime.GOOS == "windows" { //nolint:goconst
		return filepath.Join(filepath.VolumeName("C:"), fileName)
	}

	return filepath.Join("/", fileName)
}
