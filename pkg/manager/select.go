package manager

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-version"
)

// Select is a function that selects an existing installation of the Go SDK as the active one.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Select(versionNumber *version.Version) {
	m.task.Printf("Selecting version as active: %s", versionNumber)

	if !m.DryRun && m.SelectedVersion != nil {
		m.Unselect()
	}

	versionDirectory := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s", versionNumber))
	selectedDirectory := filepath.Join(m.RootDirectory, selectedDirectoryName)

	m.task.SubPrintf("Linking %s to %s", fmt.Sprintf("go%s", versionNumber), selectedDirectoryName)
	if !m.DryRun {
		err := link(versionDirectory, selectedDirectory)
		m.task.SubDieOnError(err)

		m.SelectedVersion = versionNumber
	}
}

// Unselect is a function that unselects an existing installation of the Go SDK as the active one.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Unselect() {
	m.task.Printf("Unselect current selected version")
	m.task.SubDieIff(m.SelectedVersion == nil, "could not unselect because no version is selected")

	if m.DryRun {
		return
	}

	selectedDirectory := filepath.Join(m.RootDirectory, selectedDirectoryName)

	m.task.SubPrintf("Unlinking directory: %s", selectedDirectory)
	err := unlink(selectedDirectory)
	m.task.SubDieOnError(err)

	m.SelectedVersion = nil
}
