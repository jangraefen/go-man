package manager

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/mholt/archiver/v3"

	"github.com/NoizeMe/go-man/pkg/releases"
)

// Install is a function that installs new instances of the Go SDK.
// As installation parameters the version number, operating system and platform architecture are considered when choosing the
// correct installation artifacts. The releaseType parameter is used to limit the amount of accepted versions. Feedback is
// directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Install(versionNumber *version.Version, operatingSystem, arch string, releaseType releases.ReleaseType) {
	m.task.Printf("Installing %s %s-%s:", versionNumber, operatingSystem, arch)

	release, releasePresent, err := releases.GetForVersion(releaseType, versionNumber)
	m.task.SubDieOnError(err)
	m.task.SubDieIff(!releasePresent, "release with versionName %s not present", versionNumber)

	files := release.FindFiles(operatingSystem, arch, releases.ArchiveFile)
	m.task.SubDieIff(len(files) != 1, "release %s with %s-%s not present", versionNumber, operatingSystem, arch)

	file := files[0]
	destinationFile := filepath.Join(m.RootDirectory, file.Filename)
	destinationDirectory := filepath.Join(m.RootDirectory, file.Version)

	if _, err := os.Stat(destinationFile); err != nil && os.IsNotExist(err) {
		m.task.SubPrintf("Downloading: %s", file.GetURL())
		if !m.DryRun {
			m.task.SubDieOnError(file.Download(destinationFile, false))
		}
	} else {
		m.task.SubPrintf("Downloading: Skipping, since %s is already present", destinationFile)
	}

	m.task.SubPrintf("Verifying integrity: %s", file.Sha256)
	if !m.DryRun {
		same, err := file.VerifySame(destinationFile)
		m.task.SubDieOnError(err)
		m.task.SubDieIff(
			!same,
			"Downloaded file %s could not be verified because the checksums did not match",
			destinationFile,
		)
	}

	if _, err := os.Stat(destinationDirectory); err != nil && os.IsNotExist(err) {
		m.task.SubPrintf("Extracting: %s", file.Filename)
		if !m.DryRun {
			m.task.SubDieOnError(archiver.Unarchive(destinationFile, destinationDirectory))
		}
	} else {
		m.task.SubPrintf("Extracting: Skipping, since %s is already extracted", file.Version)
	}

	m.task.SubPrintf("Verifying installation: %s", destinationDirectory)
	if !m.DryRun {
		detectedVersion, err := detectGoVersion(destinationDirectory)
		m.task.SubDieOnError(err)
		m.task.SubDieIff(!detectedVersion.Equal(versionNumber), "Could not verify installation: %s", detectedVersion)

		m.InstalledVersions = append(m.InstalledVersions, versionNumber)
		sort.Sort(m.InstalledVersions)
	}
}
