package releases

import (
	"crypto/sha256"
	"fmt"
	"github.com/hashicorp/go-version"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// The FileKind type is a string that describes what nature a ReleaseFile has.
type FileKind string

const (
	fileUrlTemplate = "https://golang.org/dl/%s"

	// The SourceFile file kind describes source archives of the Golang SDK release.
	SourceFile = FileKind("source")
	// The ArchiveFile file kind describes binary distribution archives of the Golang SDK release.
	ArchiveFile = FileKind("archive")
	// The InstallerFile file kind describes installer executable of the Golang SDK release.
	InstallerFile = FileKind("installer")
)

// The Release struct holds all relevant information about a released Go version.
// The official golang website offers a JSON-based endpoint that serves a list of all Golang releases that are structured
// just like this struct, so it can be used to read that endpoint.
type Release struct {
	// The version name of the release. It has the pattern of `go1.X.Y`.
	Version string `json:"version"`
	// Flag that marks a release as stable. A release is considered stable if it is one of the latest two releases and has
	// the latest patch version.
	Stable bool `json:"stable"`
	// A slice of all files that are associated with the release. This should never be empty.
	Files []ReleaseFile `json:"files"`
}

// The GetVersionNumber function returns the version number for a Golang release.
// Since the Version field of a release is prefixed by the string "go", this method returns a substring of this fields that
// is stripped of that exact prefix, to allow easier processing.
func (r Release) GetVersionNumber() *version.Version {
	if r.Version == "" {
		return nil
	}

	return version.Must(version.NewVersion(strings.TrimPrefix(r.Version, "go")))
}

// The FindFiles function returns a sub-slice of all files that match the given operating system and architecture.
// If a file is not specific for an operating system or architecture, it will also be included.
func (r Release) FindFiles(os, arch string, kind FileKind) []ReleaseFile {
	var filteredFiles []ReleaseFile

	for _, file := range r.Files {
		if (file.OS == "" || file.OS == os) && (file.Arch == "" || file.Arch == arch) && file.Kind == kind {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

// The ReleaseFile struct holds all information about a file that is part of a released Go version.
// With each release, a couple of files are distributed, like an installer, a source archive or a binary distribution for
// different operating systems and architectures. This struct holds information exactly about these files, as part of a
// Release struct that is obtained by querying the endpoint by the official Golang website.
type ReleaseFile struct {
	// The filename of the file as it can be found on the download mirror, including the extension.
	Filename string `json:"filename"`
	// The operating system that the file is target at. May be empty if not applicable.
	OS string `json:"os"`
	// The processor architecture that the file is target at. May be empty if not applicable.
	Arch string `json:"arch"`
	// The version of the release the file belongs to.
	Version string `json:"version"`
	// The sha256 checksum that can be used to verify the integrity of the file.
	Sha256 string `json:"sha256"`
	// The size in bytes of the file.
	Size int32 `json:"size"`
	// The kind that this file belongs to.
	Kind FileKind `json:"kind"`
}

// The GetUrl function returns the URL where the file can be downloaded from.
func (f ReleaseFile) GetUrl() string {
	if len(f.Filename) == 0 {
		return ""
	}

	return fmt.Sprintf(fileUrlTemplate, f.Filename)
}

// The Download function loads the receiving release file to the given destination file.
// If the destination file already exists it is overwritten, except if the skipIfPresent flag is set to true. In this case
// nothing will be downloaded.
func (f ReleaseFile) Download(destinationFile string, skipIfPresent bool) error {
	if stat, err := os.Stat(destinationFile); err == nil && !stat.IsDir() {
		if skipIfPresent {
			return nil
		} else if err = os.Remove(destinationFile); err != nil {
			return err
		}
	}

	directory, _ := filepath.Split(destinationFile)
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	response, err := http.Get(f.GetUrl())
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status while retrieving release file: %s", response.Status)
	}

	defer response.Body.Close()

	file, err := os.Create(destinationFile)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := io.Copy(file, response.Body); err != nil {
		return err
	}

	return nil
}

// The VerifySame function checks if a given file has the correct checksum.
// It first builds the sha256 of the given file and then compares that value against the Sha256 attribute.
func (f ReleaseFile) VerifySame(fileName string) (bool, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return false, err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	return f.Sha256 == checksum, nil
}

// The Collection type is a type-alias for a slice of releases, which also implements sort.Interface.
type Collection []*Release

// See sort.Interface for more details.
func (c Collection) Len() int {
	return len(c)
}

// See sort.Interface for more details.
func (c Collection) Less(i, j int) bool {
	return c[i].GetVersionNumber().LessThan(c[j].GetVersionNumber())
}

// See sort.Interface for more details.
func (c Collection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
