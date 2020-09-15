package manager

import (
	"github.com/NoizeMe/go-man/pkg/logging"
	"github.com/NoizeMe/go-man/pkg/releases"
	"github.com/hashicorp/go-version"
)

// The Cleanup function removes all Go SDK installations that are currently not considered stable.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m GoManager) Cleanup() {
	logging.Printf("Scanning for non-stable versions")
	var versionsToRemove version.Collection

	for _, installedVersion := range m.InstalledVersions {
		_, exists, err := releases.GetForVersion(releases.IncludeStable, installedVersion)
		logging.IfTaskError(err)

		if !exists {
			logging.TaskPrintf("Marked %s for removal", installedVersion)
			versionsToRemove = append(versionsToRemove, installedVersion)
		}
	}

	for _, versionToRemove := range versionsToRemove {
		m.Remove(versionToRemove)
	}
}
