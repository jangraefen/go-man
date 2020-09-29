package manager

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NoizeMe/go-man/pkg/tasks"
)

func TestNewManager(t *testing.T) {
	task := &tasks.Task{ErrorExitCode: 1, Output: os.Stdout, Error: os.Stderr}
	rootDirectory := t.TempDir()

	validVersion := version.Must(version.NewVersion("1.15.2"))
	anotherValidVersion := version.Must(version.NewVersion("1.14.9"))
	invalidVersion := version.Must(version.NewVersion("1.11.0"))

	manager, err := NewManager(task, rootDirectory)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Empty(t, manager.InstalledVersions)
	assert.Nil(t, manager.SelectedVersion)

	setupInstallation(t, rootDirectory, true, validVersion)
	setupInstallation(t, rootDirectory, true, anotherValidVersion)
	setupInstallation(t, rootDirectory, false, invalidVersion)
	require.NoError(t, link(
		filepath.Join(rootDirectory, "go1.15.2"),
		filepath.Join(rootDirectory, selectedDirectoryName),
	))

	manager, err = NewManager(task, rootDirectory)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Len(t, manager.InstalledVersions, 2)
	assert.Contains(t, manager.InstalledVersions, validVersion)
	assert.Contains(t, manager.InstalledVersions, anotherValidVersion)
	assert.NotContains(t, manager.InstalledVersions, invalidVersion)
	assert.True(t, manager.SelectedVersion.Equal(validVersion))

	rootDirectory = filepath.Join("not", "existent", "directory")

	manager, err = NewManager(task, rootDirectory)
	assert.Error(t, err)
	assert.Nil(t, manager)
}

func setupInstallation(t *testing.T, rootDirectory string, valid bool, goVersion fmt.Stringer) {
	t.Helper()

	goVersionString := fmt.Sprintf("go%s", goVersion)

	sdkPath := filepath.Join(rootDirectory, goVersionString)
	versionPath := filepath.Join(sdkPath, "VERSION")

	versionContent := goVersionString
	if !valid {
		versionContent = "invalidContent"
	}

	require.NoError(t, os.MkdirAll(sdkPath, 0700), "Could not create installation directory", sdkPath)
	require.NoError(t, ioutil.WriteFile(versionPath, []byte(versionContent), 0600))
}
