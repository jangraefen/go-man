package manager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectGoVersion(t *testing.T) {
	version, err := detectGoVersion(filepath.Join("not", "existent", "directory"))
	assert.Error(t, err)
	assert.Nil(t, version)

	version, err = detectGoVersion(os.Getenv("GOROOT"))
	assert.NoError(t, err)
	assert.NotNil(t, version)
}
