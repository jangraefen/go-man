package releases

import (
	"fmt"
	"sort"

	"github.com/hashicorp/go-version"

	"github.com/NoizeMe/go-man/internal/httputil"
)

// The ReleaseType type is a string that describes what kind of release types should be returned by a release list.
type ReleaseType string

const (
	releaseListURLTemplate = "https://golang.org/dl/?mode=json&include=%s"

	// IncludeAll is the release type that will include each and every release of Go that was ever distributed publicly.
	IncludeAll = ReleaseType("all")
	// IncludeStable is the release type that will include each release that is currently considered stable.
	IncludeStable = ReleaseType("stable")
)

var (
	// ReleaseListCache is a map that caches the last fetched release list. Visible mostly for testing.
	ReleaseListCache = map[ReleaseType]Collection{}
)

// SelectReleaseType is a function that returns the release type that matches the input parameters best.
// For convenience can be used to get the correct release type by describing what kind of releases are desired and the
// correct release type is then selected by this function. By default IncludeStable is returned.
func SelectReleaseType(unstable bool) ReleaseType {
	if unstable {
		return IncludeAll
	}

	return IncludeStable
}

// ListAll is a function that retrieves a list of all Golang releases from the official website.
// This list is retrieved by querying a JSON endpoint that is provided by the official Golang website. If the endpoint
// responds with any other status code than 200, an error is returned.
func ListAll(releaseType ReleaseType) (Collection, error) {
	if _, ok := ReleaseListCache[releaseType]; !ok {
		newReleaseList := Collection{}
		if err := httputil.GetJSON(fmt.Sprintf(releaseListURLTemplate, releaseType), &newReleaseList); err != nil {
			return nil, err
		}

		ReleaseListCache[releaseType] = newReleaseList
	}

	return ReleaseListCache[releaseType], nil
}

// GetLatest is a function that retrieves the latest release of the Golang SDK.
func GetLatest(releaseType ReleaseType) (*Release, error) {
	releases, err := ListAll(releaseType)
	if err != nil {
		return nil, err
	}

	sort.Sort(releases)
	return releases[releases.Len()-1], nil
}

// GetForVersion is a function that returns the Golang release with a given version, if such a release exists.
// A list of releases is retrieved, honoring the given release type as a filter, and then scanned for a release that has the
// same version number as the version variable. If no such release can be found, an empty release object is returned and the
// boolean return value will be set to false.
func GetForVersion(releaseType ReleaseType, version *version.Version) (*Release, bool, error) {
	releases, err := ListAll(releaseType)
	if err != nil {
		return nil, false, err
	}

	for _, release := range releases {
		if version.Equal(release.GetVersionNumber()) {
			return release, true, nil
		}
	}

	return nil, false, nil
}
