package manager

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"

	"github.com/NoizeMe/go-man/internal/fileutil"
)

func detectGoVersion(sdkDirectory string) (*version.Version, error) {
	versionPath := filepath.Join(sdkDirectory, "VERSION")

	if !fileutil.PathExists(versionPath) {
		return nil, fmt.Errorf("could not locate version file: %s", versionPath)
	}

	versionContent, err := ioutil.ReadFile(versionPath)
	if err != nil {
		return nil, err
	}

	return version.NewVersion(strings.TrimPrefix(string(versionContent), "go"))
}
