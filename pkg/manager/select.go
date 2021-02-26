package manager

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/go-version"

	"github.com/jangraefen/go-man/internal/fileutil"
	"github.com/jangraefen/go-man/pkg/tasks"
)

// Select is a function that selects an existing installation of the Go SDK as the active one.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Select(versionNumber *version.Version) error {
	versionName := toVersionName(versionNumber)
	m.task.Printf("Selecting version as active: %s", versionName)

	versionDirectory := filepath.Join(m.RootDirectory, fmt.Sprintf("go%s", versionName))
	if !fileutil.PathExists(versionDirectory) {
		return fmt.Errorf("version %v was not found", versionName)
	}

	selectTask := m.task.Step()
	if m.SelectedVersion != nil {
		if err := m.unselect(selectTask); err != nil {
			return err
		}
	}

	linkDescription := "Linking selection directory"
	linkFunction := func() error { return link(versionDirectory, filepath.Join(m.RootDirectory, selectedDirectoryName)) }
	if err := selectTask.Track(linkDescription, linkFunction); err != nil {
		return err
	}

	m.SelectedVersion = versionNumber
	return nil
}

// Unselect is a function that unselects an existing installation of the Go SDK as the active one.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Unselect() error {
	m.task.Printf("Unselect current selected version")
	if m.SelectedVersion == nil {
		return errors.New("could not unselect because no version is selected")
	}

	return m.unselect(m.task.Step())
}

func (m *GoManager) unselect(task *tasks.Task) error {
	unlinkDescription := "Unlinking selection directory"
	unlinkFunction := func() error { return unlink(filepath.Join(m.RootDirectory, selectedDirectoryName)) }
	if err := task.Track(unlinkDescription, unlinkFunction); err != nil {
		return err
	}

	m.SelectedVersion = nil
	return nil
}
