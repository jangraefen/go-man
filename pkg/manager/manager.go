package manager

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"

	"github.com/NoizeMe/go-man/pkg/tasks"
)

const (
	selectedDirectoryName = "go-selected"
)

// The GoManager is responsible for managing Go SDK installations.
// It know about the currently installed SDKs, the current active SDK and has means to install additional SDKs, change which
// SDK is currently selected and remove existing installations.
type GoManager struct {
	// The root directory stores the installed SDKs as well as any configuration files.
	RootDirectory string
	// The collection of all currently installed versions of the Go SDK.
	InstalledVersions version.Collection
	// The currently selected version. The selected version is the release that is synced to the "selected" directory. Might
	// be nil, if no version is currently selected.
	SelectedVersion *version.Version

	task *tasks.Task
}

// NewManager is a constructor for the GoManager struct.
// It reads through the given root directory and detects the current state and initializes the GoManager instance
// accordingly.
func NewManager(task *tasks.Task, rootDirectory string) (*GoManager, error) {
	var selectedVersion *version.Version
	var installedVersions version.Collection

	fileInfos, err := ioutil.ReadDir(rootDirectory)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() || fileInfo.Mode()&os.ModeSymlink != 0 {
			detectedVersion, err := detectGoVersion(filepath.Join(rootDirectory, fileInfo.Name(), "go"))
			if err != nil {
				continue
			}

			if fileInfo.Name() == selectedDirectoryName {
				selectedVersion = detectedVersion
			} else {
				installedVersions = append(installedVersions, detectedVersion)
			}
		}
	}

	return &GoManager{
		RootDirectory:     rootDirectory,
		InstalledVersions: installedVersions,
		SelectedVersion:   selectedVersion,
		task:              task,
	}, nil
}
