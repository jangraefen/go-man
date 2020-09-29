package fileutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathExists(t *testing.T) {
	directory := t.TempDir()

	assert.False(t, PathExists(filepath.Join("not", "existent", "directory")))
	assert.True(t, PathExists(directory))
}

func TestTryRemove(t *testing.T) {
	directory := t.TempDir()

	assert.False(t, TryRemove(filepath.Join("not", "existent", "directory")))
	assert.True(t, TryRemove(directory))
}

func TestMoveDirectory(t *testing.T) {
	directory := t.TempDir()
	fromDirectory := filepath.Join(directory, "from")
	toDirectory := filepath.Join(directory, "to")
	toDirectoryWithoutPermission := getNoPermissionDirectory("to")

	assert.Error(t, MoveDirectory(fromDirectory, toDirectory))
	TryRemove(toDirectory)

	assert.Error(t, MoveDirectory(fromDirectory, toDirectoryWithoutPermission))
	TryRemove(toDirectoryWithoutPermission)

	require.NoError(t, os.MkdirAll(filepath.Join(fromDirectory, "sub"), 0700))
	require.NoError(t, ioutil.WriteFile(filepath.Join(fromDirectory, "sub", "file.txt"), []byte("content"), 0600))

	assert.NoError(t, MoveDirectory(fromDirectory, toDirectory))
	assert.FileExists(t, filepath.Join(toDirectory, "sub", "file.txt"))
	TryRemove(toDirectory)
}

func getNoPermissionDirectory(fileName string) string {
	if runtime.GOOS == "windows" { //nolint:goconst
		return filepath.Join(filepath.VolumeName("C:"), fileName)
	}

	return filepath.Join("/", fileName)
}
