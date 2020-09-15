package manager

import (
	"github.com/NoizeMe/go-man/pkg/logging"
	"github.com/NoizeMe/go-man/pkg/releases"
	"github.com/hashicorp/go-version"
	"github.com/mholt/archiver/v3"
	"os"
	"path/filepath"
	"sort"
)

// The Install function installs new instances of the Go SDK.
// As installation parameters the version number, operating system and platform architecture are considered when choosing the
// correct installation artifacts. The releaseType parameter is used to limit the amount of accepted versions. Feedback is
// directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Install(versionNumber *version.Version, operatingSystem, arch string, releaseType releases.ReleaseType) {
	logging.Printf("Installing %s %s-%s:", versionNumber, operatingSystem, arch)

	release, releasePresent, err := releases.GetForVersion(releaseType, versionNumber)
	logging.IfTaskError(err)
	logging.IfTaskErrorf(!releasePresent, "release with versionName %s not present", versionNumber)

	files := release.FindFiles(operatingSystem, arch, releases.ArchiveFile)
	logging.IfTaskErrorf(len(files) != 1, "release %s with %s-%s not present", versionNumber, operatingSystem, arch)

	file := files[0]
	destinationFile := filepath.Join(m.RootDirectory, file.Filename)
	destinationDirectory := filepath.Join(m.RootDirectory, file.Version)

	if stat, err := os.Stat(destinationDirectory); err == nil && stat.IsDir() {
		logging.TaskPrintf("Version %s already installed, skipping.", file.Version)
		return
	}

	logging.TaskPrintf("Downloading: %s", file.GetUrl())
	if !m.DryRun {
		logging.IfTaskError(file.Download(destinationFile, false))
	}

	logging.TaskPrintf("Verifying integrity: %s", file.Sha256)
	if !m.DryRun {
		same, err := file.VerifySame(destinationFile)
		logging.IfTaskError(err)
		logging.IfTaskErrorf(
			!same,
			"Downloaded file %s could not be verified because the checksums did not match",
			destinationFile,
		)
	}

	logging.TaskPrintf("Extracting: %s", file.Filename)
	if !m.DryRun {
		logging.IfTaskError(archiver.Unarchive(destinationFile, destinationDirectory))
	}

	logging.TaskPrintf("Verifying installation: %s", destinationDirectory)
	if !m.DryRun {
		detectedVersion, err := detectGoVersion(destinationDirectory)
		logging.IfTaskError(err)
		logging.IfTaskErrorf(detectedVersion != versionNumber, "Could not verify installation: %s", detectedVersion)

		m.InstalledVersions = append(m.InstalledVersions, versionNumber)
		sort.Sort(m.InstalledVersions)
	}
}
