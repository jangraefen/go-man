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

func TestGoManager_Select(t *testing.T) {
	validVersion := version.Must(version.NewVersion("1.15.2"))
	anotherValidVersion := version.Must(version.NewVersion("1.14.9"))
	invalidVersion := version.Must(version.NewVersion("42.1337.3"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, true, validVersion.String())
	setupInstallation(t, tempDir, true, anotherValidVersion.String())

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

	assert.NoError(t, sut.Select(validVersion))
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", anotherValidVersion)))
	assert.True(t, fileutil.PathExists(filepath.Join(tempDir, selectedDirectoryName)))

	assert.NoError(t, sut.Select(anotherValidVersion))
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", anotherValidVersion)))
	assert.True(t, fileutil.PathExists(filepath.Join(tempDir, selectedDirectoryName)))

	assert.Error(t, sut.Select(invalidVersion))
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", anotherValidVersion)))
	assert.True(t, fileutil.PathExists(filepath.Join(tempDir, selectedDirectoryName)))
}

func TestGoManager_Select_WithLinkFailure(t *testing.T) {
	invalidVersion := version.Must(version.NewVersion("1.14.9"))

	tempDir := t.TempDir()

	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{invalidVersion},
		SelectedVersion:   nil,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}

	require.NoError(t, os.MkdirAll(filepath.Join(tempDir, selectedDirectoryName), 0700))
	assert.Error(t, sut.Select(invalidVersion))

	setupInstallation(t, tempDir, true, invalidVersion.String())
	assert.Error(t, sut.Select(invalidVersion))
}

func TestGoManager_Select_WithFailingUnselect(t *testing.T) {
	validVersion := version.Must(version.NewVersion("1.15.2"))
	invalidVersion := version.Must(version.NewVersion("1.14.9"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, true, validVersion.String())

	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{validVersion, invalidVersion},
		SelectedVersion:   invalidVersion,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}

	assert.Error(t, sut.Select(validVersion))
}

func TestGoManager_Select_WithTwoPartVersion(t *testing.T) {
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

	assert.NoError(t, sut.Select(validVersion))
	assert.DirExists(t, filepath.Join(tempDir, "go1.16"))
	assert.True(t, fileutil.PathExists(filepath.Join(tempDir, selectedDirectoryName)))
}

func TestGoManager_Unselect(t *testing.T) {
	validVersion := version.Must(version.NewVersion("1.15.2"))

	tempDir := t.TempDir()
	sdkPath := filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion))
	selectedPath := filepath.Join(tempDir, selectedDirectoryName)

	setupInstallation(t, tempDir, true, validVersion.String())

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

	assert.Error(t, sut.Unselect())
	assert.False(t, fileutil.PathExists(selectedPath))

	require.NoError(t, link(sdkPath, selectedPath))
	sut.SelectedVersion = validVersion

	assert.NoError(t, sut.Unselect())
	assert.False(t, fileutil.PathExists(selectedPath))
	assert.DirExists(t, sdkPath)
}

func TestGoManager_Unselect_WithoutExistingDirectory(t *testing.T) {
	invalidVersion := version.Must(version.NewVersion("1.14.9"))

	tempDir := t.TempDir()

	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{invalidVersion},
		SelectedVersion:   invalidVersion,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}

	assert.Error(t, sut.Unselect())
}
