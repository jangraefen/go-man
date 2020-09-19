package utils

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathExists(t *testing.T) {
	directory, _ := ioutil.TempDir("", "goman-utils-io-TestPathExists")

	assert.False(t, PathExists(filepath.Join("not", "existent", "directory")))
	assert.True(t, PathExists(directory))
}

func TestTryRemove(t *testing.T) {
	directory, _ := ioutil.TempDir("", "goman-utils-io-TestTryRemove")

	assert.False(t, TryRemove(filepath.Join("not", "existent", "directory")))
	assert.True(t, TryRemove(directory))
}
