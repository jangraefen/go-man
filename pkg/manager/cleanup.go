package manager

import (
	"github.com/hashicorp/go-version"

	"github.com/NoizeMe/go-man/pkg/releases"
)

// Cleanup is a function that removes all Go SDK installations that are currently not considered stable.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m GoManager) Cleanup() {
	m.task.Printf("Removing all non-stable versions")
	cleanupTask := m.task.Step()

	versionsToRemove, err := filterNonStableVersions(m.InstalledVersions)
	cleanupTask.DieOnError(err)

	for _, versionToRemove := range versionsToRemove {
		m.Uninstall(versionToRemove)
	}
}

func filterNonStableVersions(versions version.Collection) (version.Collection, error) {
	filtered := version.Collection{}

	for _, v := range versions {
		_, exists, err := releases.GetForVersion(releases.IncludeStable, v)
		if err != nil {
			return nil, err
		}
		if !exists {
			filtered = append(filtered, v)
		}
	}

	return filtered, nil
}
