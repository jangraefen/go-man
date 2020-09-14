package main

import (
	"fmt"
	"github.com/NoizeMe/go-man/pkg/logging"
	goreleases "github.com/NoizeMe/go-man/pkg/releases"
	"github.com/NoizeMe/go-man/pkg/selectors"
	"github.com/hashicorp/go-version"
	"github.com/mholt/archiver/v3"
	"github.com/posener/cmd"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	root   = cmd.New()
	dryRun = root.Bool(
		"dryrun",
		false,
		"If set, all actions are done as if they were successful, but no persistent changes will be performed.",
	)

	releases    = root.SubCommand("releases", "Lists all available releases of the Go SDK.")
	releasesAll = releases.Bool(
		"all",
		false,
		"If set, not only stable but all releases are listed.",
	)

	install    = root.SubCommand("install", "Installs one or more new versions of the Go SDK.")
	installAll = install.Bool(
		"all",
		false,
		"If set, not only stable but all releases are installable.",
	)
	installOS = install.String(
		"os",
		runtime.GOOS,
		"Defines for which operating system the Go SDK should be downloaded. By default, the current OS is chosen.",
	)
	installArch = install.String(
		"arch",
		runtime.GOARCH,
		"Defines for which architecture the Go SDK should be downloaded. By default, the current architecture is chosen.",
	)
	installVersions = install.Args(
		"[versions]",
		"The versions that should be installed. May be 'latest' or any version number.",
	)

	remove    = root.SubCommand("remove", "Remove an existing installation of the Go SDK.")
	removeAll = remove.Bool(
		"all",
		false,
		"If set, all installed versions will be deleted.",
	)
	removeVersions = remove.Args(
		"[versions]",
		"The versions that should be removed.",
	)

	select_        = root.SubCommand("select", "Selects the active installation of the Go SDK.")
	selectVersions = select_.Args(
		"[version]",
		"The version that should be selected.",
	)
)

func main() {
	// Parse the command line arguments. Any errors will get caught be the library and will cause the usage to be printed.
	// The program will exit afterwards.
	_ = root.Parse()

	switch {
	case releases.Parsed():
		handleReleases(*releasesAll)
	case install.Parsed():
		handleInstall(*dryRun, *installAll, *installOS, *installArch, *installVersions)
	case remove.Parsed():
		handleRemove(*dryRun, *removeAll, *removeVersions)
	case select_.Parsed():
		handleSelect(*dryRun, *selectVersions)
	}
}

func handleReleases(all bool) {
	logging.Printf("List of available releases:")

	releaseList, err := goreleases.ListAll(goreleases.SelectReleaseType(all))
	logging.IfTaskError(err)

	for _, r := range releaseList {
		logging.TaskPrintf("%s", r.GetVersionNumber())
	}
}

func handleInstall(dryRun, all bool, operatingSystem, arch string, versionNames []string) {
	if len(versionNames) == 0 {
		latest, err := goreleases.GetLatest()
		logging.IfError(err)

		versionNames = []string{latest.GetVersionNumber().String()}
	}

	for _, versionName := range versionNames {
		parsedVersion, err := version.NewVersion(versionName)
		if err != nil {
			logging.IfTaskError(err)
		}

		logging.Printf("Installing %s %s-%s:", parsedVersion, operatingSystem, arch)

		release, releasePresent, err := goreleases.GetForVersion(goreleases.SelectReleaseType(all), parsedVersion)
		logging.IfTaskError(err)
		logging.IfTaskErrorf(!releasePresent, "release with versionName %s not present", parsedVersion)

		files := release.FindFiles(operatingSystem, arch, goreleases.ArchiveFile)
		logging.IfTaskErrorf(len(files) == 0, "release %s with %s-%s not present", parsedVersion, operatingSystem, arch)

		for _, file := range files {
			root := gomanRoot()
			destinationFile := filepath.Join(root, file.Filename)
			destinationDirectory := filepath.Join(root, file.Version)

			if stat, err := os.Stat(destinationDirectory); err == nil && stat.IsDir() {
				logging.TaskPrintf("Version %s already installed, skipping.", file.Version)
				continue
			}

			logging.TaskPrintf("Downloading: %s", file.GetUrl())
			if !dryRun {
				logging.IfTaskError(file.Download(destinationFile, false))
			}

			logging.TaskPrintf("Verifying integrity: %s", file.Sha256)
			if !dryRun {
				same, err := file.VerifySame(destinationFile)
				logging.IfTaskError(err)
				logging.IfTaskErrorf(
					!same,
					"Downloaded file %s could not be verified because the checksums did not match",
					destinationFile,
				)
			}

			logging.TaskPrintf("Extracting: %s", file.Filename)
			if !dryRun {
				logging.IfTaskError(archiver.Unarchive(destinationFile, destinationDirectory))
			}

			goBinaryPath := filepath.Join(destinationDirectory, "go", "bin", "go")
			logging.TaskPrintf("Verifying installation: %s", goBinaryPath)

			command := exec.Command(goBinaryPath, "version")
			output, err := command.Output()
			logging.IfTaskError(err)

			actualOutput := strings.TrimSpace(string(output))
			expectedOutput := fmt.Sprintf("go version %s %s/%s", file.Version, file.OS, file.Arch)
			logging.IfTaskErrorf(expectedOutput != actualOutput, "Could not verify installation: %s", actualOutput)
		}
	}
}

func handleRemove(dryRun bool, all bool, versionNames []string) {
	root := gomanRoot()

	logging.IfErrorf(!all && len(versionNames) == 0, "No versionNames to remove, skipping.")
	logging.IfErrorf(all && len(versionNames) > 0, "Both all flag and versionNames given, skipping.")

	if all {
		fileInfos, err := ioutil.ReadDir(root)
		logging.IfError(err)

		for _, fileInfo := range fileInfos {
			if fileInfo.IsDir() && strings.HasPrefix(fileInfo.Name(), "go") {
				versionNames = append(versionNames, strings.TrimPrefix(fileInfo.Name(), "go"))
			}
		}
	}

	for _, versionName := range versionNames {
		parsedVersion, err := version.NewVersion(versionName)
		if err != nil {
			logging.IfTaskError(err)
		}

		versionDirectory := filepath.Join(root, fmt.Sprintf("go%s", parsedVersion))
		versionArchive := filepath.Join(root, fmt.Sprintf("go%s*", parsedVersion))

		logging.Printf("Deleting %s", parsedVersion)

		logging.TaskPrintf("Removing SDK: %s", versionDirectory)
		if !dryRun {
			logging.IfTaskError(os.RemoveAll(versionDirectory))
		}

		matches, err := filepath.Glob(versionArchive)
		logging.IfTaskError(err)

		for _, match := range matches {
			logging.TaskPrintf("Removing SDK archive: %s", match)
			if !dryRun {
				logging.IfTaskError(os.Remove(match))
			}
		}
	}
}

func handleSelect(dryRun bool, versionNames []string) {
	logging.IfErrorf(len(versionNames) == 0, "No version to select, skipping.")
	logging.IfErrorf(len(versionNames) > 1, "More then one version to select, skipping.")

	parsedVersion, err := version.NewVersion(versionNames[0])
	if err != nil {
		logging.IfError(err)
	}

	logging.Printf("Selecting %s as the active Go version", parsedVersion)

	root := gomanRoot()
	versionDirectory := filepath.Join(root, fmt.Sprintf("go%s", parsedVersion))

	stat, err := os.Stat(versionDirectory)
	logging.IfTaskError(err)
	logging.IfTaskErrorf(!stat.IsDir(), "%s is not a directory", versionDirectory)

	if !dryRun {
		logging.IfTaskError(selectors.SyncToCurrent(filepath.Split(versionDirectory)))
	}
}

func gomanRoot() string {
	root := os.Getenv("GOMANROOT")
	if len(root) > 0 {
		return root
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/", ".goman")
	}

	return filepath.Join(homeDir, ".goman")
}
