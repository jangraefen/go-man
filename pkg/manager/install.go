package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/mholt/archiver/v3"

	"github.com/NoizeMe/go-man/pkg/releases"
	"github.com/NoizeMe/go-man/pkg/tasks"
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

	installTask.DieOnError(downloadRelease(installTask, file, destinationFile))
	installTask.DieOnError(verifyDownload(installTask, file, destinationFile))
	installTask.DieOnError(extractRelease(installTask, file, destinationFile, destinationDirectory))
	installTask.DieOnError(verifyRelease(installTask, versionNumber, destinationDirectory))

	m.InstalledVersions = append(m.InstalledVersions, versionNumber)
	sort.Sort(m.InstalledVersions)
}

func downloadRelease(installTask *tasks.Task, file releases.ReleaseFile, destinationFile string) error {
	if _, err := os.Stat(destinationFile); err == nil || !os.IsNotExist(err) {
		installTask.Printf("Downloading: Skipping, since %s is already present", destinationFile)
		return nil
	}

	installTask.Printf("Downloading: %s", file.GetURL())
	return file.Download(destinationFile, false)
}

func verifyDownload(installTask *tasks.Task, file releases.ReleaseFile, destinationFile string) error {
	installTask.Printf("Verifying integrity: %s", file.Sha256)

	same, err := file.VerifySame(destinationFile)
	if err != nil {
		return err
	}
	if !same {
		return fmt.Errorf("downloaded file %s could not be verified because the checksums did not match", destinationFile)
	}

	return nil
}

func extractRelease(installTask *tasks.Task, file releases.ReleaseFile, destinationFile string, destinationDirectory string) error {
	if _, err := os.Stat(destinationDirectory); err == nil || !os.IsNotExist(err) {
		installTask.Printf("Extracting: Skipping, since %s is already extracted", file.Version)
		return nil
	}

	installTask.Printf("Extracting: %s", file.Filename)
	return archiver.Unarchive(destinationFile, destinationDirectory)
}

func verifyRelease(installTask *tasks.Task, versionNumber *version.Version, destinationDirectory string) error {
	installTask.Printf("Verifying installation: %s", destinationDirectory)

	detectedVersion, err := detectGoVersion(filepath.Join(destinationDirectory, "go"))
	if err != nil {
		return err
	}
	if !detectedVersion.Equal(versionNumber) {
		return fmt.Errorf("could not verify installation: %s", detectedVersion)
	}

	return nil
}
