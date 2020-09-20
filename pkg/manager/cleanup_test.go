package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"

	"github.com/NoizeMe/go-man/pkg/releases"
	"github.com/NoizeMe/go-man/pkg/tasks"
)

func TestGoManager_Cleanup(t *testing.T) {
	stableRelease, err := releases.GetLatest(releases.IncludeStable)
	assert.NoError(t, err)

	stableVersion := stableRelease.GetVersionNumber()
	unstableVersion := version.Must(version.NewVersion("1.11.0"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, stableVersion)
	setupInstallation(t, tempDir, unstableVersion)

	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{stableVersion, unstableVersion},
		SelectedVersion:   nil,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}
	assert.NotNil(t, sut)

	assert.NoError(t, sut.Cleanup())
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", stableVersion)))
	assert.FileExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s-windows-amd64.zip", stableVersion)))
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", unstableVersion)))
	assert.NoFileExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s-windows-amd64.zip", unstableVersion)))

	assert.NoError(t, sut.Cleanup())
}