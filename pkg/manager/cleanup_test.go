package manager

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jangraefen/go-man/internal/httputil"
	"github.com/jangraefen/go-man/pkg/releases"
	"github.com/jangraefen/go-man/pkg/tasks"
)

func TestGoManager_Cleanup(t *testing.T) {
	stableRelease, err := releases.GetLatest(releases.IncludeStable)
	require.NoError(t, err)

	stableVersion := stableRelease.GetVersionNumber()
	unstableVersion := version.Must(version.NewVersion("1.11.0"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, true, stableVersion)
	setupInstallation(t, tempDir, true, unstableVersion)

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

	assert.NoError(t, sut.Cleanup())
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", stableVersion)))
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", unstableVersion)))

	assert.NoError(t, sut.Cleanup())
}

func TestGoManager_Cleanup_WithInvalid(t *testing.T) {
	unstableVersion := version.Must(version.NewVersion("1.11.0"))

	tempDir := t.TempDir()
	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{unstableVersion},
		SelectedVersion:   nil,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}

	assert.Error(t, sut.Cleanup())
}

func TestGoManager_Cleanup_WithHTTPError(t *testing.T) {
	t.Cleanup(func() {
		httputil.Client = http.DefaultClient
	})

	unstableVersion := version.Must(version.NewVersion("1.11.0"))

	tempDir := t.TempDir()
	sut := &GoManager{
		RootDirectory:     tempDir,
		InstalledVersions: version.Collection{unstableVersion},
		SelectedVersion:   nil,
		task: &tasks.Task{
			ErrorExitCode: 1,
			Output:        os.Stdout,
			Error:         os.Stderr,
		},
	}

	httputil.Client = httputil.StaticResponseClient(404, []byte("not found"), nil)
	delete(releases.ReleaseListCache, releases.IncludeStable)

	assert.Error(t, sut.Cleanup())
}
