package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
)

// UninstallAll is a function that removes all current installations of the Go SDK.
func (m *GoManager) UninstallAll() {
	installedVersions := make(version.Collection, len(m.InstalledVersions))
	copy(installedVersions, m.InstalledVersions)

	for _, versionNumber := range installedVersions {
		m.Uninstall(versionNumber)
	}
}

// Uninstall is a function that removes an existing installation of the Go SDK.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Uninstall(versionNumber *version.Version) {
	versionDirectory := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s", versionNumber))
	versionArchive := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s*", versionNumber))

	if versionNumber.Equal(m.SelectedVersion) {
		m.Unselect()
	}

	m.task.Printf("Removing %s", versionNumber)
	uninstallTask := m.task.Step()

	uninstallTask.Printf("Deleting SDK: %s", versionDirectory)
	uninstallTask.DieOnError(deleteVersionDirectory(versionDirectory))
	uninstallTask.DieOnError(deleteVersionArchive(versionArchive))

	for index, installedVersion := range m.InstalledVersions {
		if installedVersion.Equal(versionNumber) {
			m.InstalledVersions = append(m.InstalledVersions[:index], m.InstalledVersions[index+1:]...)
			break
		}
	}
}

func deleteVersionArchive(archivePattern string) error {
	matches, err := filepath.Glob(archivePattern)
	if err != nil {
		return err
	}

	for _, match := range matches {
		if err := os.Remove(match); err != nil {
			return err
		}
	}

	return nil
}

func deleteVersionDirectory(directory string) error {
	return os.RemoveAll(directory)
}
