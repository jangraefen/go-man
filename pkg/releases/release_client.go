package releases

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// The ReleaseType type is a string that describes what kind of release types should be returned by a release list.
type ReleaseType string

const (
	releaseListUrlTemplate = "https://golang.org/dl/?mode=json&include=%s"

	// The IncludeAll release type will include each and every release of Go that was ever distributed publicly.
	IncludeAll = ReleaseType("all")
	// The IncludeStable release type will include each release that is currently considered stable.
	IncludeStable = ReleaseType("stable")
)

// The SelectReleaseType function returns the release type that matches the input parameters best.
// For convenience can be used to get the correct release type by describing what kind of releases are desired and the
// correct release type is then selected by this function. By default IncludeStable is returned.
func SelectReleaseType(all bool) ReleaseType {
	if all {
		return IncludeAll
	}

	return IncludeStable
}

// The ListAll function retrieves a list of all Golang releases from the official website.
// This list is retrieved by querying a JSON endpoint that is provided by the official Golang website. If the endpoint
// responds with any other status code than 200, an error is returned.
func ListAll(releaseType ReleaseType) ([]Release, error) {
	response, err := http.Get(fmt.Sprintf(releaseListUrlTemplate, releaseType))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status while retrieving releases: %s", response.Status)
	}

	defer response.Body.Close()

	versions := make([]Release, 0)
	err = json.NewDecoder(response.Body).Decode(&versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}

// The GetLatest function retrieves the latest stable release of the Golang SDK.
func GetLatest() (Release, error) {
	releases, err := ListAll(IncludeStable)
	if err != nil {
		return Release{}, err
	}

	// TODO This is working by accident, not by design. It would be much better to actually parse the version numbers here.
	return releases[0], nil
}

// The GetForVersion function returns the Golang release with a given version, if such a release exists.
// A list of releases is retrieved, honoring the given release type as a filter, and then scanned for a release that has the
// same version number as the version variable. If no such release can be found, an empty release object is returned and the
// boolean return value will be set to false.
func GetForVersion(releaseType ReleaseType, version string) (Release, bool, error) {
	releases, err := ListAll(releaseType)
	if err != nil {
		return Release{}, false, err
	}

	for _, release := range releases {
		if version == release.GetVersionNumber() {
			return release, true, nil
		}
	}

	return Release{}, false, nil
}
