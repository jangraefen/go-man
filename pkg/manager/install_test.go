package manager

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"

	"github.com/NoizeMe/go-man/pkg/releases"
	"github.com/NoizeMe/go-man/pkg/tasks"
	"github.com/NoizeMe/go-man/pkg/utils"
)

func TestGoManager_Install(t *testing.T) {
	t.Cleanup(func() {
		utils.Client = http.DefaultClient
	})

	validVersion := version.Must(version.NewVersion("1.15.2"))
	invalidVersion := version.Must(version.NewVersion("42.1337.3"))

	tempDir := t.TempDir()

	sut, err := NewManager(&tasks.Task{
		ErrorExitCode: 1,
		Output:        os.Stdout,
		Error:         os.Stderr,
	}, tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, sut)

	assert.Error(t, sut.Install(invalidVersion, runtime.GOOS, runtime.GOARCH, releases.IncludeAll))

	assert.NoError(t, sut.Install(validVersion, runtime.GOOS, runtime.GOARCH, releases.IncludeAll))
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))

	assert.Error(t, sut.Install(validVersion, runtime.GOOS, runtime.GOARCH, releases.IncludeAll))
	assert.DirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
}

func TestGoManager_Install_WithInvalidTarget(t *testing.T) {
	validVersion := version.Must(version.NewVersion("1.15.2"))
	tempDir := t.TempDir()
	sut, err := NewManager(&tasks.Task{
		ErrorExitCode: 1,
		Output:        os.Stdout,
		Error:         os.Stderr,
	}, tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, sut)

	assert.Error(t, sut.Install(validVersion, "foobar", runtime.GOARCH, releases.IncludeAll))
	assert.Error(t, sut.Install(validVersion, runtime.GOOS, "foobar", releases.IncludeAll))
}

func TestGoManager_Install_WithHTTPError(t *testing.T) {
	t.Cleanup(func() {
		utils.Client = http.DefaultClient
	})

	validVersion := version.Must(version.NewVersion("1.15.2"))
	tempDir := t.TempDir()
	sut, err := NewManager(&tasks.Task{
		ErrorExitCode: 1,
		Output:        os.Stdout,
		Error:         os.Stderr,
	}, tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, sut)

	utils.Client = utils.StaticResponseClient(404, []byte("not found"), nil)
	assert.Error(t, sut.Install(validVersion, runtime.GOOS, runtime.GOARCH, releases.IncludeAll))

	utils.Client = utils.StaticResponseClient(0, nil, errors.New("failure"))
	assert.Error(t, sut.Install(validVersion, runtime.GOOS, runtime.GOARCH, releases.IncludeAll))
}
