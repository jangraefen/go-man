package manager

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectGoVersion(t *testing.T) {
	goVersion, err := detectGoVersion(filepath.Join("not", "existent", "directory"))
	assert.Error(t, err)
	assert.Nil(t, goVersion)

	goVersion, err = detectGoVersion(runtime.GOROOT())
	assert.NoError(t, err)
	assert.NotNil(t, goVersion)

	rootDirectory := t.TempDir()
	setupInstallation(t, rootDirectory, false, "1.15.2")

	goVersion, err = detectGoVersion(filepath.Join(rootDirectory, "go1.15.2", "go"))
	assert.Error(t, err)
	assert.Nil(t, goVersion)
}
