package manager

import (
	"fmt"
	"github.com/NoizeMe/go-man/pkg/logging"
	"github.com/hashicorp/go-version"
	"path/filepath"
)

// The Select function selects an existing installation of the Go SDK as the active one.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Select(versionNumber *version.Version) {
	logging.Printf("Selecting version as active: %s", versionNumber)

	if !m.DryRun && m.SelectedVersion != nil {
		m.Unselect()
	}

	versionDirectory := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s", versionNumber))
	selectedDirectory := filepath.Join(m.RootDirectory, selectedDirectoryName)

	logging.Printf("Linking %s to %s", fmt.Sprintf("go%s", versionNumber), selectedDirectoryName)
	if !m.DryRun {
		err := link(versionDirectory, selectedDirectory)
		logging.IfTaskError(err)

		m.SelectedVersion = versionNumber
	}
}

// The Unselect function unselects an existing installation of the Go SDK as the active one.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Unselect() {
	logging.Printf("Unselect current selected version")
	logging.IfTaskErrorf(m.SelectedVersion == nil, "could not unselect because no version is selected")

	if m.DryRun {
		return
	}

	selectedDirectory := filepath.Join(m.RootDirectory, selectedDirectoryName)

	logging.TaskPrintf("Unlinking directory: %s", selectedDirectory)
	err := unlink(selectedDirectory)
	logging.IfTaskError(err)

	m.SelectedVersion = nil
}
