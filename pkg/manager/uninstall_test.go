package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"

	"github.com/NoizeMe/go-man/pkg/tasks"
	"github.com/NoizeMe/go-man/pkg/utils"
)

func TestGoManager_UninstallAll(t *testing.T) {
	validVersion := version.Must(version.NewVersion("1.15.2"))
	anotherValidVersion := version.Must(version.NewVersion("1.14.0"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, validVersion)
	setupInstallation(t, tempDir, anotherValidVersion)

	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{validVersion, anotherValidVersion},
		SelectedVersion:   nil,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}
	assert.NotNil(t, sut)

	assert.NoError(t, sut.UninstallAll())
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", anotherValidVersion)))

	assert.NoError(t, sut.UninstallAll())
}

func TestGoManager_Uninstall(t *testing.T) {
	invalidVersion := version.Must(version.NewVersion("42.1337.3"))
	validVersion := version.Must(version.NewVersion("1.15.2"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, validVersion)

	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{validVersion},
		SelectedVersion:   validVersion,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}
	assert.NotNil(t, sut)
	assert.NoError(t, link(
		filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)),
		filepath.Join(tempDir, selectedDirectoryName),
	))

	assert.Error(t, sut.Uninstall(invalidVersion))

	assert.NoError(t, sut.Uninstall(validVersion))
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
	assert.False(t, utils.PathExists(filepath.Join(tempDir, selectedDirectoryName)))

	sut.InstalledVersions = version.Collection{validVersion}
	setupInstallation(t, tempDir, validVersion)

	assert.NoError(t, sut.Uninstall(validVersion))
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))

	assert.Error(t, sut.Uninstall(validVersion))
}

func setupInstallation(t *testing.T, rootDirectory string, goVersion *version.Version) {
	t.Helper()

	folderPath := filepath.Join(rootDirectory, fmt.Sprintf("go%s", goVersion))

	if err := os.MkdirAll(folderPath, 0700); err != nil {
		assert.FailNowf(t, "Could not create installation directory %s", folderPath)
		return
	}
}
