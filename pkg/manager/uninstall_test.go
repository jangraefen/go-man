package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jangraefen/go-man/internal/fileutil"
	"github.com/jangraefen/go-man/pkg/tasks"
)

func TestGoManager_UninstallAll(t *testing.T) {
	validVersion := version.Must(version.NewVersion("1.15.2"))
	anotherValidVersion := version.Must(version.NewVersion("1.14.0"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, true, validVersion.String())
	setupInstallation(t, tempDir, true, "1.14")

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

	assert.NoError(t, sut.UninstallAll())
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
	assert.NoDirExists(t, filepath.Join(tempDir, "go1.14"))

	assert.NoError(t, sut.UninstallAll())
}

func TestGoManager_UninstallWithTwoPartVersion(t *testing.T) {
	validVersion := version.Must(version.NewVersion("1.16"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, true, "1.16")

	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{validVersion},
		SelectedVersion:   nil,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}

	assert.NoError(t, sut.Uninstall(validVersion))
	assert.NoDirExists(t, filepath.Join(tempDir, "go1.16"))
}

func TestGoManager_UninstallAll_WithBrokenInstallation(t *testing.T) {
	validVersion := version.Must(version.NewVersion("1.15.2"))
	invalidVersion := version.Must(version.NewVersion("1.14.0"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, true, validVersion.String())

	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{invalidVersion, validVersion},
		SelectedVersion:   nil,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}

	assert.Error(t, sut.UninstallAll())
}

func TestGoManager_Uninstall(t *testing.T) {
	invalidVersion := version.Must(version.NewVersion("42.1337.3"))
	validVersion := version.Must(version.NewVersion("1.15.2"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, true, validVersion.String())

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
	require.NoError(t, link(
		filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)),
		filepath.Join(tempDir, selectedDirectoryName),
	))

	assert.Error(t, sut.Uninstall(invalidVersion))

	assert.NoError(t, sut.Uninstall(validVersion))
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
	assert.False(t, fileutil.PathExists(filepath.Join(tempDir, selectedDirectoryName)))

	sut.InstalledVersions = version.Collection{validVersion}
	setupInstallation(t, tempDir, true, validVersion.String())

	assert.NoError(t, sut.Uninstall(validVersion))
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))

	assert.Error(t, sut.Uninstall(validVersion))
}
