package manager

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
)

func detectGoVersion(sdkDirectory string) (*version.Version, error) {
	goBinaryPath := filepath.Join(sdkDirectory, "go", "bin", "go")

	command := exec.Command(goBinaryPath, "version")
	output, err := command.Output()
	if err != nil {
		return nil, err
	}

	var scannedVersion, osAndArch string
	matches, err := fmt.Fscanf(bytes.NewReader(output), "go version %s %s\n", &scannedVersion, &osAndArch)
	if err != nil {
		return nil, err
	}
	if matches != 2 {
		return nil, errors.New("could not detect go version since output did had the expected format")
	}

	return version.NewVersion(strings.TrimPrefix(scannedVersion, "go"))
}
