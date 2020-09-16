package manager

import (
	"github.com/hashicorp/go-version"

	"github.com/NoizeMe/go-man/pkg/releases"
)

// Cleanup is a function that removes all Go SDK installations that are currently not considered stable.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m GoManager) Cleanup() {
	m.task.Printf("Scanning for non-stable versions")
	var versionsToRemove version.Collection

	for _, installedVersion := range m.InstalledVersions {
		_, exists, err := releases.GetForVersion(releases.IncludeStable, installedVersion)
		m.task.IfTaskError(err)

		if !exists {
			m.task.TaskPrintf("Marked %s for removal", installedVersion)
			versionsToRemove = append(versionsToRemove, installedVersion)
		}
	}

	for _, versionToRemove := range versionsToRemove {
		m.Uninstall(versionToRemove)
	}
}
