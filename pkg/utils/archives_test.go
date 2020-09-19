package utils

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractArchive(t *testing.T) {
	archiveFile := getTestFile(t, "test.zip")
	missingDirectory := filepath.Join(t.TempDir(), "missing")
	destinationDirectory := filepath.Join(t.TempDir(), "extracted")

	extracted, err := ExtractArchive("MISSING_FILE.zip", missingDirectory, false)
	assert.Error(t, err)
	assert.False(t, extracted)

	extracted, err = ExtractArchive(archiveFile, destinationDirectory, false)
	assert.NoError(t, err)
	assert.True(t, extracted)

	extracted, err = ExtractArchive(archiveFile, destinationDirectory, false)
	assert.NoError(t, err)
	assert.False(t, extracted)

	extracted, err = ExtractArchive(archiveFile, destinationDirectory, true)
	assert.NoError(t, err)
	assert.True(t, extracted)
}

func getTestFile(t *testing.T, fileName string) string {
	t.Helper()

	_, file, _, _ := runtime.Caller(0)
	dir, _ := filepath.Split(file)

	return filepath.Clean(filepath.Join(
		dir,
		"..",
		"..",
		"test",
		fileName,
	))
}
