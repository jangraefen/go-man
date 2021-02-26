package manager

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func Test_toVersionName(t *testing.T) {
	assert.Equal(t, "1.16", toVersionName(version.Must(version.NewVersion("1.16"))))
	assert.Equal(t, "1.16", toVersionName(version.Must(version.NewVersion("1.16.0.0.0"))))
	assert.Equal(t, "1.16.0.1", toVersionName(version.Must(version.NewVersion("1.16.0.1.0"))))
}
