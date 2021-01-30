package manager

import (
	"github.com/hashicorp/go-version"

	"github.com/jangraefen/go-man/pkg/releases"
)

// Cleanup is a function that removes all Go SDK installations that are currently not considered stable.
// Feedback is directly printed to the stdout or stderr, so nothing is returned here.
func (m *GoManager) Cleanup() error {
	m.task.Printf("Removing all non-stable versions")

	versionsToRemove, err := filterNonStableVersions(m.InstalledVersions)
	if err != nil {
		return err
	}

	for _, versionToRemove := range versionsToRemove {
		if err := m.Uninstall(versionToRemove); err != nil {
			return err
		}
	}

	return nil
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
