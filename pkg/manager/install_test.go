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
	copy2 "github.com/otiai10/copy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/NoizeMe/go-man/internal/fileutil"
	"github.com/NoizeMe/go-man/internal/httputil"
	"github.com/NoizeMe/go-man/pkg/releases"
	"github.com/NoizeMe/go-man/pkg/tasks"
)

func TestGoManager_Install(t *testing.T) {
	validVersion := version.Must(version.NewVersion("1.15.2"))
	invalidVersion := version.Must(version.NewVersion("42.1337.3"))

	tempDir := t.TempDir()

	sut, err := NewManager(&tasks.Task{
		ErrorExitCode: 1,
		Output:        os.Stdout,
		Error:         os.Stderr,
	}, tempDir)
	require.NoError(t, err)

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
	require.NoError(t, err)

	assert.Error(t, sut.Install(validVersion, "foobar", runtime.GOARCH, releases.IncludeAll))
	assert.Error(t, sut.Install(validVersion, runtime.GOOS, "foobar", releases.IncludeAll))
}

func TestGoManager_Install_WithHTTPError(t *testing.T) {
	t.Cleanup(func() {
		httputil.Client = http.DefaultClient
	})

	validVersion := version.Must(version.NewVersion("1.15.2"))
	tempDir := t.TempDir()
	sut, err := NewManager(&tasks.Task{
		ErrorExitCode: 1,
		Output:        os.Stdout,
		Error:         os.Stderr,
	}, tempDir)
	require.NoError(t, err)

	httputil.Client = httputil.StaticResponseClient(404, []byte("not found"), nil)
	delete(releases.ReleaseListCache, releases.IncludeAll)
	assert.Error(t, sut.Install(validVersion, runtime.GOOS, runtime.GOARCH, releases.IncludeAll))

	httputil.Client = httputil.StaticResponseClient(0, nil, errors.New("failure"))
	delete(releases.ReleaseListCache, releases.IncludeAll)
	assert.Error(t, sut.Install(validVersion, runtime.GOOS, runtime.GOARCH, releases.IncludeAll))
}

func TestDownloadRelease(t *testing.T) {
	t.Cleanup(func() {
		httputil.Client = http.DefaultClient
	})

	file := releases.ReleaseFile{Filename: "go1.15.2.src.tar.gz"}
	destinationFile := filepath.Join(t.TempDir(), "download.rel")

	assert.NoError(t, downloadRelease(file, destinationFile))
	assert.Error(t, downloadRelease(file, destinationFile))

	httputil.Client = httputil.StaticResponseClient(404, []byte("not found"), nil)
	fileutil.TryRemove(destinationFile)
	assert.Error(t, downloadRelease(file, destinationFile))

	httputil.Client = httputil.StaticResponseClient(0, nil, errors.New("failure"))
	fileutil.TryRemove(destinationFile)
	assert.Error(t, downloadRelease(file, destinationFile))
}

func TestVerifyDownload(t *testing.T) {
	file := releases.ReleaseFile{Filename: "go1.15.2.src.tar.gz", Sha256: "28bf9d0bcde251011caae230a4a05d917b172ea203f2a62f2c2f9533589d4b4d"}
	destinationFile := filepath.Join(t.TempDir(), "download.rel")

	require.NoError(t, downloadRelease(file, destinationFile))
	assert.NoError(t, verifyDownload(file, destinationFile))

	fileutil.TryRemove(destinationFile)

	assert.Error(t, verifyDownload(file, destinationFile))

	f, err := os.Create(destinationFile)
	require.NoError(t, err)
	_ = f.Close()

	assert.Error(t, verifyDownload(file, destinationFile))
}

func TestExtractRelease(t *testing.T) {
	file := releases.ReleaseFile{Filename: "go1.15.2.src.tar.gz"}
	destinationFile := filepath.Join(t.TempDir(), "download.tar.gz")
	destinationDirectory := filepath.Join(t.TempDir(), "extracted")

	require.NoError(t, downloadRelease(file, destinationFile))

	assert.NoError(t, extractRelease(destinationFile, destinationDirectory))
	assert.Error(t, extractRelease(destinationFile, destinationDirectory))

	fileutil.TryRemove(destinationFile)
	assert.Error(t, extractRelease(destinationFile, destinationDirectory))

	fileutil.TryRemove(destinationDirectory)
	assert.Error(t, extractRelease(getTestFile(t, "invalid.zip"), destinationDirectory))
}

func TestVerifyRelease(t *testing.T) {
	destinationDirectory := filepath.Join(t.TempDir(), "go-installation")
	require.NoError(t, copy2.Copy(
		filepath.Join(runtime.GOROOT(), "VERSION"),
		filepath.Join(destinationDirectory, "go", "VERSION"),
	))

	validVersion := version.Must(version.NewVersion("1.15.2"))
	invalidVersion := version.Must(version.NewVersion("42.1337.3"))

	assert.NoError(t, verifyRelease(validVersion, destinationDirectory))
	assert.Error(t, verifyRelease(invalidVersion, destinationDirectory))

	fileutil.TryRemove(destinationDirectory)
	assert.Error(t, verifyRelease(validVersion, destinationDirectory))
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
