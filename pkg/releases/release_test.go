package releases

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestRelease_GetVersionNumber(t *testing.T) {
	sut := &Release{}
	assert.Nil(t, sut.GetVersionNumber())

	sut.Version = ""
	assert.Nil(t, sut.GetVersionNumber())

	sut.Version = "go1.15.2"
	assert.NotNil(t, sut.GetVersionNumber())
	assert.Equal(t, version.Must(version.NewVersion("1.15.2")), sut.GetVersionNumber())

	sut.Version = "1.15.2"
	assert.NotNil(t, sut.GetVersionNumber())
	assert.Equal(t, version.Must(version.NewVersion("1.15.2")), sut.GetVersionNumber())
}

func TestRelease_FindFiles(t *testing.T) {
	sut := &Release{}
	assert.Len(t, sut.FindFiles("windows", "amd64", ArchiveFile), 0)

	sut.Files = []ReleaseFile{}
	assert.Len(t, sut.FindFiles("windows", "amd64", ArchiveFile), 0)

	expected := ReleaseFile{
		OS:      "windows",
		Arch:    "amd64",
		Version: "go1.15.2",
		Kind:    ArchiveFile,
	}
	sut.Files = []ReleaseFile{expected}

	assert.Len(t, sut.FindFiles("windows", "amd64", ArchiveFile), 1)
	assert.Equal(t, expected, sut.FindFiles("windows", "amd64", ArchiveFile)[0])

	assert.Len(t, sut.FindFiles("windows", "amd64", InstallerFile), 0)
	assert.Len(t, sut.FindFiles("windows", "386", ArchiveFile), 0)
	assert.Len(t, sut.FindFiles("darwin", "amd64", ArchiveFile), 0)

	expected = ReleaseFile{
		OS:      "",
		Arch:    "",
		Version: "go1.15.2",
		Kind:    SourceFile,
	}
	sut.Files = []ReleaseFile{expected}

	assert.Len(t, sut.FindFiles("", "", SourceFile), 1)
	assert.Equal(t, expected, sut.FindFiles("", "", SourceFile)[0])
	assert.Len(t, sut.FindFiles("windows", "amd64", SourceFile), 1)
	assert.Equal(t, expected, sut.FindFiles("windows", "amd64", SourceFile)[0])
	assert.Len(t, sut.FindFiles("xxx", "yyy", SourceFile), 1)
	assert.Equal(t, expected, sut.FindFiles("xxx", "yyy", SourceFile)[0])
}

func TestReleaseFile_GetURL(t *testing.T) {
	sut := &ReleaseFile{}
	assert.Equal(t, "", sut.GetURL())

	sut.Filename = ""
	assert.Equal(t, "", sut.GetURL())

	sut.Filename = "go1.15.2.windows-amd64.zip"
	assert.Equal(t, "https://golang.org/dl/go1.15.2.windows-amd64.zip", sut.GetURL())
}

func TestReleaseFile_Download(t *testing.T) {
	tempDir := t.TempDir()
	targetFile := filepath.Join(tempDir, "download.file")

	sut := &ReleaseFile{}
	assert.Error(t, sut.Download(targetFile, false))

	sut.Filename = ""
	assert.Error(t, sut.Download(targetFile, false))

	sut.Filename = "go1.15.2.windows-amd64.zip"
	assert.NoError(t, sut.Download(targetFile, false))

	stat, err := os.Stat(targetFile)
	assert.Nil(t, err)
	assert.Greater(t, stat.Size(), int64(0))

	assert.NoError(t, sut.Download(targetFile, true))
	statNew, err := os.Stat(targetFile)
	assert.NoError(t, err)
	assert.Equal(t, stat.ModTime(), statNew.ModTime())
}

func TestReleaseFile_VerifySame(t *testing.T) {
	tempDir := t.TempDir()
	mockFile := filepath.Join(tempDir, "mock.file")
	targetFile := filepath.Join(tempDir, "verify.file")

	sut := &ReleaseFile{
		Filename: "go1.15.2.windows-amd64.zip",
	}

	assert.NoError(t, sut.Download(targetFile, false))
	assert.NoError(t, ioutil.WriteFile(mockFile, []byte("NOT_THE_EXPECTED_CONTENT"), 0600))

	same, err := sut.VerifySame("I_DO_NOT_EXIST.txt")
	assert.Error(t, err)
	assert.False(t, same)

	same, err = sut.VerifySame(targetFile)
	assert.NoError(t, err)
	assert.False(t, same)

	sut.Sha256 = "e72782cc6de233188c75b06849368826eaa1b8bd9e1cd766db9466a12b7138ca"

	same, err = sut.VerifySame(mockFile)
	assert.NoError(t, err)
	assert.False(t, same)

	same, err = sut.VerifySame(targetFile)
	assert.NoError(t, err)
	assert.True(t, same)
}

func TestCollection_Len(t *testing.T) {
	sut := Collection{}
	assert.Len(t, sut, 0)

	sut = append(sut, nil)
	assert.Len(t, sut, 1)

	sut = append(sut, &Release{})
	assert.Len(t, sut, 2)
}

func TestCollection_Less(t *testing.T) {
	sut := Collection{
		{Version: "go1.15.1"},
		{Version: "go1.15.1"},
		{Version: "go1.15.2"},
	}

	assert.False(t, sut.Less(0, 0))
	assert.False(t, sut.Less(0, 1))
	assert.True(t, sut.Less(0, 2))

	assert.False(t, sut.Less(1, 0))
	assert.False(t, sut.Less(1, 1))
	assert.True(t, sut.Less(1, 2))

	assert.False(t, sut.Less(2, 0))
	assert.False(t, sut.Less(2, 1))
	assert.False(t, sut.Less(2, 2))
}

func TestCollection_Swap(t *testing.T) {
	sut := Collection{
		{Version: "go1.15.1"},
		{Version: "go1.15.2"},
		{Version: "go1.15.3"},
	}

	sut.Swap(0, 2)
	assert.Equal(t, "go1.15.3", sut[0].Version)
	assert.Equal(t, "go1.15.2", sut[1].Version)
	assert.Equal(t, "go1.15.1", sut[2].Version)
}

func TestCollection_Sort(t *testing.T) {
	sut := Collection{
		{Version: "go1.15.3"},
		{Version: "go1.15.2"},
		{Version: "go1.15.1"},
	}

	sort.Sort(sut)
	assert.Equal(t, "go1.15.1", sut[0].Version)
	assert.Equal(t, "go1.15.2", sut[1].Version)
	assert.Equal(t, "go1.15.3", sut[2].Version)
}
