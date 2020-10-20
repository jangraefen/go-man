package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"

	"github.com/NoizeMe/go-man/internal/fileutil"
)

// UninstallAll is a function that removes all current installations of the Go SDK.
func (m *GoManager) UninstallAll() error {
	installedVersions := make(version.Collection, len(m.InstalledVersions))
	copy(installedVersions, m.InstalledVersions)

	for _, versionNumber := range installedVersions {
		if err := m.Uninstall(versionNumber); err != nil {
			return err
		}
	}

	return nil
}

// Uninstall is a function that removes an existing installation of the Go SDK.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Uninstall(versionNumber *version.Version) error {
	m.task.Printf("Uninstalling %s", versionNumber)
	uninstallTask := m.task.Step()

	if versionNumber.Equal(m.SelectedVersion) {
		if err := m.unselect(uninstallTask); err != nil {
			return err
		}
	}

	removeDescription := "Deleting installation directory"
	removeFunction := func() error {
		versionDirectory := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s", versionNumber))

		if !fileutil.PathExists(versionDirectory) {
			return fmt.Errorf("no directory %s to uninstall from", versionDirectory)
		}

		return os.RemoveAll(versionDirectory)
	}

	if err := uninstallTask.Track(removeDescription, removeFunction); err != nil {
		return err
	}

	for index, installedVersion := range m.InstalledVersions {
		if installedVersion.Equal(versionNumber) {
			m.InstalledVersions = append(m.InstalledVersions[:index], m.InstalledVersions[index+1:]...)
			break
		}
	}

	return nil
}
