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
	selectTask := m.task.Step()

	if !m.DryRun && m.SelectedVersion != nil {
		m.Unselect()
	}

	versionDirectory := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s", versionNumber))
	selectedDirectory := filepath.Join(m.RootDirectory, selectedDirectoryName)

	selectTask.Printf("Linking %s to %s", fmt.Sprintf("go%s", versionNumber), selectedDirectoryName)
	if !m.DryRun {
		err := link(versionDirectory, selectedDirectory)
		selectTask.DieOnError(err)

		m.SelectedVersion = versionNumber
	}
}

// Unselect is a function that unselects an existing installation of the Go SDK as the active one.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Unselect() {
	m.task.Printf("Unselect current selected version")
	unselectTask := m.task.Step()

	unselectTask.DieIff(m.SelectedVersion == nil, "could not unselect because no version is selected")

	if m.DryRun {
		return
	}

	selectedDirectory := filepath.Join(m.RootDirectory, selectedDirectoryName)

	unselectTask.Printf("Unlinking directory: %s", selectedDirectory)
	err := unlink(selectedDirectory)
	unselectTask.DieOnError(err)

	m.SelectedVersion = nil
}
