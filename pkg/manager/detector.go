package manager

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
)

func detectGoVersion(sdkDirectory string) (*version.Version, error) {
	versionPath := filepath.Join(sdkDirectory, "VERSION")

	versionContent, err := ioutil.ReadFile(versionPath)
	if err != nil {
		return nil, err
	}

	return version.NewVersion(strings.TrimPrefix(string(versionContent), "go"))
}
