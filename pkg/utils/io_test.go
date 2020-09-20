package utils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathExists(t *testing.T) {
	directory := t.TempDir()

	assert.False(t, PathExists(filepath.Join("not", "existent", "directory")))
	assert.True(t, PathExists(directory))
}

func TestTryRemove(t *testing.T) {
	directory := t.TempDir()

	assert.False(t, TryRemove(filepath.Join("not", "existent", "directory")))
	assert.True(t, TryRemove(directory))
}
