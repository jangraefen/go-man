package main

import (
	"fmt"
	"github.com/NoizeMe/go-man/pkg/logging"
	goreleases "github.com/NoizeMe/go-man/pkg/releases"
	"github.com/mholt/archiver/v3"
	"github.com/posener/cmd"
	"os"
	"path/filepath"
	"runtime"
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

func handleInstall(dryRun, all bool, operatingSystem, arch string, versions []string) {
	if len(versions) == 0 {
		latest, err := goreleases.GetLatest()
		logging.IfError(err)

		versions = []string{latest.GetVersionNumber()}
	}

	for _, version := range versions {
		logging.Printf("Installing %s %s-%s:", version, operatingSystem, arch)

		release, releasePresent, err := goreleases.GetForVersion(goreleases.SelectReleaseType(all), version)
		logging.IfTaskError(err)
		logging.IfTaskErrorf(!releasePresent, "release with version %s not present", version)

		files := release.FindFiles(operatingSystem, arch, goreleases.ArchiveFile)
		logging.IfTaskErrorf(len(files) == 0, "release %s with %s-%s not present", version, operatingSystem, arch)

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
		}
	}
}

func handleRemove(dryRun bool, all bool, versions []string) {
	logging.IfErrorf(!all && len(versions) == 0, "No versions to remove, skipping.")
	logging.IfErrorf(all && len(versions) > 0, "Both all flag and versions given, skipping.")

	for _, version := range versions {
		root := gomanRoot()
		versionDirectory := filepath.Join(root, fmt.Sprintf("go%s", version))
		versionArchive := filepath.Join(root, fmt.Sprintf("go%s*", version))

		logging.Printf("Deleting %s", version)

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

func gomanRoot() string {
	root := os.Getenv("GOMANROOT")
	if len(root) > 0 {
		return root
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/", ".gomanroot")
	}

	return filepath.Join(homeDir, ".gomanroot")
}
