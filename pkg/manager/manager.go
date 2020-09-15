package manager

import (
	"github.com/hashicorp/go-version"
	"io/ioutil"
	"os"
	"path/filepath"
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
	// The dry run flag controls if functions of this struct actual perform any changes or just give feedback to the user
	// what would have happened, if they would be actually called.
	DryRun bool
	// The collection of all currently installed versions of the Go SDK.
	InstalledVersions version.Collection
	// The currently selected version. The selected version is the release that is synced to the "selected" directory. Might
	// be nil, if no version is currently selected.
	SelectedVersion *version.Version
}

func NewManager(rootDirectory string, dryRun bool) (*GoManager, error) {
	var selectedVersion *version.Version
	var installedVersions version.Collection

	fileInfos, err := ioutil.ReadDir(rootDirectory)
	if err != nil {
		return nil, err
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() || fileInfo.Mode()&os.ModeSymlink != 0 {
			detectedVersion, err := detectGoVersion(filepath.Join(rootDirectory, fileInfo.Name()))
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
		DryRun:            dryRun,
		InstalledVersions: installedVersions,
		SelectedVersion:   selectedVersion,
	}, nil
}
