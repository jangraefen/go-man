package manager

import (
	"fmt"
	"github.com/NoizeMe/go-man/pkg/logging"
	"github.com/hashicorp/go-version"
	"os"
	"path/filepath"
)

// The RemoveAll function removes all current installations of the Go SDK.
func (m *GoManager) RemoveAll() {
	for _, versionNumber := range m.InstalledVersions {
		m.Remove(versionNumber)
	}
}

// The Remove function removes an existing installation of the Go SDK.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Remove(versionNumber *version.Version) {
	versionDirectory := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s", versionNumber))
	versionArchive := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s*", versionNumber))

	logging.Printf("Deleting %s", versionNumber)

	logging.TaskPrintf("Removing SDK: %s", versionDirectory)
	if !m.DryRun {
		logging.IfTaskError(os.RemoveAll(versionDirectory))
	}

	matches, err := filepath.Glob(versionArchive)
	logging.IfTaskError(err)

	for _, match := range matches {
		logging.TaskPrintf("Removing SDK archive: %s", match)
		if !m.DryRun {
			logging.IfTaskError(os.Remove(match))
		}
	}

	if !m.DryRun {
		for index, installedVersion := range m.InstalledVersions {
			if installedVersion.Equal(versionNumber) {
				m.InstalledVersions = append(m.InstalledVersions[:index], m.InstalledVersions[index+1:]...)

				break
			}
		}

		if versionNumber.Equal(m.SelectedVersion) {
			m.Unselect()
		}
	}
}
