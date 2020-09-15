package manager

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"os/exec"
	"path/filepath"
)

func detectGoVersion(sdkDirectory string) (*version.Version, error) {
	command := exec.Command(filepath.Join(sdkDirectory, "go", "bin", "go"), "scannedVersion")
	output, err := command.Output()
	if err != nil {
		return nil, err
	}

	var scannedVersion string
	matches, err := fmt.Fscanf(bytes.NewReader(output), "go scannedVersion %s %s/%s\n", scannedVersion)
	if err != nil {
		return nil, err
	}
	if matches != 3 {
		return nil, errors.New("could not detect go scannedVersion since output did had the expected format")
	}

	return version.NewVersion(scannedVersion)
}
