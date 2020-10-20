package manager

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/hashicorp/go-version"

	"github.com/NoizeMe/go-man/internal/archiveutil"
	"github.com/NoizeMe/go-man/internal/fileutil"
	"github.com/NoizeMe/go-man/internal/httputil"
	"github.com/NoizeMe/go-man/pkg/releases"
)

// Install is a function that installs new instances of the Go SDK.
// As installation parameters the version number, operating system and platform architecture are considered when choosing the
// correct installation artifacts. The releaseType parameter is used to limit the amount of accepted versions. Feedback is
// directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Install(versionNumber *version.Version, operatingSystem, arch string, releaseType releases.ReleaseType) error {
	m.task.Printf("Installing %s %s-%s:", versionNumber, operatingSystem, arch)
	installTask := m.task.Step()

	release, releasePresent, err := releases.GetForVersion(releaseType, versionNumber)
	if err != nil {
		return err
	}
	if !releasePresent {
		return fmt.Errorf("release with version %s not present", versionNumber)
	}

	files := release.FindFiles(operatingSystem, arch, releases.ArchiveFile)
	if len(files) != 1 {
		return fmt.Errorf("release %s with %s-%s not present", versionNumber, operatingSystem, arch)
	}

	file := files[0]
	downloadedArchive := filepath.Join(m.RootDirectory, file.Filename)
	extractionDirectory := filepath.Join(m.RootDirectory, fmt.Sprintf("extracting-%s", file.Version))
	sdkDirectory := filepath.Join(m.RootDirectory, file.Version)

	if fileutil.PathExists(sdkDirectory) {
		return fmt.Errorf("installation skipped, since %s is already present", sdkDirectory)
	}

	defer fileutil.TryRemove(downloadedArchive)
	defer fileutil.TryRemove(extractionDirectory)

	downloadDescription := "Downloading distribution"
	downloadFunction := func() error { return downloadRelease(file, downloadedArchive) }
	if err := installTask.Track(downloadDescription, downloadFunction); err != nil {
		return err
	}

	checksumDescription := "Verifying download integrity"
	checksumFunction := func() error { return verifyDownload(file, downloadedArchive) }
	if err := installTask.Track(checksumDescription, checksumFunction); err != nil {
		return err
	}

	extractDescription := "Extracting distribution"
	extractFunction := func() error { return extractRelease(downloadedArchive, extractionDirectory) }
	if err := installTask.Track(extractDescription, extractFunction); err != nil {
		return err
	}

	verifyDescription := "Verifying installation"
	verifyFunction := func() error { return verifyRelease(versionNumber, extractionDirectory) }
	if err := installTask.Track(verifyDescription, verifyFunction); err != nil {
		return err
	}

	moveDescription := "Moving installation to final location"
	moveFunction := func() error { return fileutil.MoveDirectory(filepath.Join(extractionDirectory, "go"), sdkDirectory) }
	if err := installTask.Track(moveDescription, moveFunction); err != nil {
		fileutil.TryRemove(sdkDirectory)
		return err
	}

	m.InstalledVersions = append(m.InstalledVersions, versionNumber)
	sort.Sort(m.InstalledVersions)

	return nil
}

func downloadRelease(file releases.ReleaseFile, destinationFile string) error {
	downloaded, err := httputil.GetFile(file.GetURL(), destinationFile, false)
	if err != nil {
		return err
	}
	if !downloaded {
		return fmt.Errorf("download skipping, since %s is already present", destinationFile)
	}

	return nil
}

func verifyDownload(file releases.ReleaseFile, destinationFile string) error {
	same, err := file.VerifySame(destinationFile)
	if err != nil {
		return err
	}
	if !same {
		return fmt.Errorf("downloaded file %s could not be verified because the checksums did not match", destinationFile)
	}

	return nil
}

func extractRelease(destinationFile string, destinationDirectory string) error {
	extracted, err := archiveutil.Extract(destinationFile, destinationDirectory, false)
	if err != nil {
		return err
	}
	if !extracted {
		return fmt.Errorf("extraction skipping, since %s is already present", destinationDirectory)
	}

	return nil
}

func verifyRelease(versionNumber *version.Version, destinationDirectory string) error {
	detectedVersion, err := detectGoVersion(filepath.Join(destinationDirectory, "go"))
	if err != nil {
		return err
	}
	if !detectedVersion.Equal(versionNumber) {
		return fmt.Errorf("could not verify installation: %s", detectedVersion)
	}

	return nil
}
