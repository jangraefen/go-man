package manager

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectGoVersion(t *testing.T) {
	version, err := detectGoVersion("/dev/null")
	assert.Error(t, err)
	assert.Nil(t, version)

	version, err = detectGoVersion(os.Getenv("GOROOT"))
	assert.NoError(t, err)
	assert.NotNil(t, version)
}
