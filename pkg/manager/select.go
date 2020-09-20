package manager

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-version"

	"github.com/NoizeMe/go-man/pkg/utils"
)

// Select is a function that selects an existing installation of the Go SDK as the active one.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Select(versionNumber *version.Version) error {
	versionDirectory := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s", versionNumber))
	selectedDirectory := filepath.Join(m.RootDirectory, selectedDirectoryName)

	m.task.Printf("Selecting version as active: %s", versionNumber)
	selectTask := m.task.Step()

	if !utils.PathExists(versionDirectory) {
		return fmt.Errorf("version %v was not found", versionNumber)
	}

	if m.SelectedVersion != nil {
		if err := m.Unselect(); err != nil {
			return err
		}
	}

	selectTask.Printf("Linking %s to %s", fmt.Sprintf("go%s", versionNumber), selectedDirectoryName)

	if err := link(versionDirectory, selectedDirectory); err != nil {
		return err
	}

	m.SelectedVersion = versionNumber
	return nil
}

// Unselect is a function that unselects an existing installation of the Go SDK as the active one.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Unselect() error {
	m.task.Printf("Unselect current selected version")
	unselectTask := m.task.Step()

	if m.SelectedVersion == nil {
		return errors.New("could not unselect because no version is selected")
	}

	selectedDirectory := filepath.Join(m.RootDirectory, selectedDirectoryName)

	unselectTask.Printf("Unlinking directory: %s", selectedDirectory)
	if err := unlink(selectedDirectory); err != nil {
		return err
	}

	m.SelectedVersion = nil
	return nil
}
