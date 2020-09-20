package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"

	"github.com/NoizeMe/go-man/pkg/tasks"
)

func TestGoManager_Uninstall(t *testing.T) {
	invalidVersion := version.Must(version.NewVersion("42.1337.3"))
	validVersion := version.Must(version.NewVersion("1.15.2"))

	tempDir := t.TempDir()

	setupInstallation(t, tempDir, validVersion)

	sut, err := NewManager(&tasks.Task{
		ErrorExitCode: 1,
		Output:        os.Stdout,
		Error:         os.Stderr,
	}, tempDir)

	assert.NoError(t, err)
	assert.NotNil(t, sut)

	assert.Error(t, sut.Uninstall(invalidVersion))

	assert.NoError(t, sut.Uninstall(validVersion))
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
	assert.NoFileExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s-amd64-windows.zip", validVersion)))

	assert.Error(t, sut.Uninstall(validVersion))
	assert.NoDirExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s", validVersion)))
	assert.NoFileExists(t, filepath.Join(tempDir, fmt.Sprintf("go%s-amd64-windows.zip", validVersion)))
}

func setupInstallation(t *testing.T, rootDirectory string, goVersion *version.Version) {
	t.Helper()

	folderPath := filepath.Join(rootDirectory, fmt.Sprintf("go%s", goVersion))
	archivePath := filepath.Join(rootDirectory, fmt.Sprintf("go%s-amd64-windows.zip", goVersion))

	if err := os.MkdirAll(folderPath, 0700); err != nil {
		assert.FailNowf(t, "Could not create installation directory %s", archivePath)
		return
	}

	file, err := os.Create(archivePath)
	if err != nil {
		assert.FailNowf(t, "Could not touch file %s", archivePath)
		return
	}
	file.Close()
}
