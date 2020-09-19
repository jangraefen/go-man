package manager

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/hashicorp/go-version"

	"github.com/NoizeMe/go-man/pkg/releases"
	"github.com/NoizeMe/go-man/pkg/utils"
)

// Install is a function that installs new instances of the Go SDK.
// As installation parameters the version number, operating system and platform architecture are considered when choosing the
// correct installation artifacts. The releaseType parameter is used to limit the amount of accepted versions. Feedback is
// directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Install(versionNumber *version.Version, operatingSystem, arch string, releaseType releases.ReleaseType) {
	m.task.Printf("Installing %s %s-%s:", versionNumber, operatingSystem, arch)
	installTask := m.task.Step()

	release, releasePresent, err := releases.GetForVersion(releaseType, versionNumber)
	installTask.DieOnError(err)
	installTask.DieIff(!releasePresent, "release with versionName %s not present", versionNumber)

	files := release.FindFiles(operatingSystem, arch, releases.ArchiveFile)
	installTask.DieIff(len(files) != 1, "release %s with %s-%s not present", versionNumber, operatingSystem, arch)

	file := files[0]
	destinationFile := filepath.Join(m.RootDirectory, file.Filename)
	destinationDirectory := filepath.Join(m.RootDirectory, file.Version)

	installTask.Printf("Downloading: %s", file.GetURL())
	downloaded, err := downloadRelease(file, destinationFile)
	installTask.DieOnError(err)
	if !downloaded {
		installTask.Printf("Downloading: Skipping, since %s is already present", destinationFile)
	}

	installTask.Printf("Verifying integrity: %s", file.Sha256)
	installTask.DieOnError(verifyDownload(file, destinationFile))

	installTask.Printf("Extracting: %s", file.Filename)
	extracted, err := extractRelease(destinationFile, destinationDirectory)
	installTask.DieOnError(err)
	if !extracted {
		installTask.Printf("Extracting: Skipping, since %s is already extracted", file.Version)
	}

	installTask.Printf("Verifying installation: %s", destinationDirectory)
	installTask.DieOnError(verifyRelease(versionNumber, destinationDirectory))

	m.InstalledVersions = append(m.InstalledVersions, versionNumber)
	sort.Sort(m.InstalledVersions)
}

func downloadRelease(file releases.ReleaseFile, destinationFile string) (bool, error) {
	return utils.DownloadFile(file.GetURL(), destinationFile, false)
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

func extractRelease(destinationFile string, destinationDirectory string) (bool, error) {
	return utils.ExtractArchive(destinationFile, destinationDirectory, false)
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
