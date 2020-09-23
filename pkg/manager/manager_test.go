package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"text/template"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"

	"github.com/NoizeMe/go-man/pkg/tasks"
)

func TestNewManager(t *testing.T) {
	task := &tasks.Task{ErrorExitCode: 1, Output: os.Stdout, Error: os.Stderr}
	rootDirectory := t.TempDir()

	validVersion := version.Must(version.NewVersion("1.15.2"))
	anotherValidVersion := version.Must(version.NewVersion("1.14.9"))

	manager, err := NewManager(task, rootDirectory)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Empty(t, manager.InstalledVersions)
	assert.Nil(t, manager.SelectedVersion)

	setupInstallation(t, rootDirectory, true, validVersion)
	setupInstallation(t, rootDirectory, true, anotherValidVersion)
	assert.NoError(t, link(
		filepath.Join(rootDirectory, "go1.15.2"),
		filepath.Join(rootDirectory, selectedDirectoryName),
	))

	manager, err = NewManager(task, rootDirectory)
	assert.NoError(t, err)
	assert.NotNil(t, manager)
	assert.Len(t, manager.InstalledVersions, 2)
	assert.Contains(t, manager.InstalledVersions, validVersion)
	assert.Contains(t, manager.InstalledVersions, anotherValidVersion)
	assert.True(t, manager.SelectedVersion.Equal(validVersion))

	rootDirectory = filepath.Join("not", "existent", "directory")

	manager, err = NewManager(task, rootDirectory)
	assert.Error(t, err)
	assert.Nil(t, manager)
}

func setupInstallation(t *testing.T, rootDirectory string, valid bool, goVersion fmt.Stringer) {
	t.Helper()

	folderPath := filepath.Join(rootDirectory, fmt.Sprintf("go%s", goVersion))
	binPath := filepath.Join(folderPath, "go", "bin")

	if err := os.MkdirAll(binPath, 0700); err != nil {
		assert.FailNow(t, "Could not create installation directory", folderPath)
		return
	}

	tmpl, err := template.New(getExecutableTemplateName(t)).ParseFiles(getExecutableTemplatePath(t))
	if err != nil {
		assert.FailNow(t, "Could not render go binary template", err)
		return
	}

	executablePath := filepath.Join(binPath, getExecutableTemplateTarget("go"))
	file, err := os.OpenFile(executablePath, os.O_WRONLY|os.O_CREATE, 0744)
	if err != nil {
		assert.FailNow(t, "Could not create go binary file", err)
		return
	}
	defer func() {
		_ = file.Close()
	}()

	parameters := struct {
		GOVersion string
		GOOS      string
		GOArch    string
		Valid     bool
	}{goVersion.String(), runtime.GOOS, runtime.GOARCH, valid}

	if err := tmpl.Execute(file, parameters); err != nil {
		assert.FailNow(t, "Could not render go binary template", err)
		return
	}
}

func getExecutableTemplateTarget(name string) string {
	if runtime.GOOS == "windows" { //nolint:gosec
		return name + ".bat"
	}

	return name
}

func getExecutableTemplatePath(t *testing.T) string {
	t.Helper()

	_, file, _, _ := runtime.Caller(0)
	dir, _ := filepath.Split(file)

	fileName := getExecutableTemplateName(t)

	return filepath.Clean(filepath.Join(
		dir,
		"..",
		"..",
		"test",
		fileName,
	))
}

func getExecutableTemplateName(t *testing.T) string {
	t.Helper()

	if runtime.GOOS == "windows" {
		return "go.bat"
	}

	return "go.sh"
}
