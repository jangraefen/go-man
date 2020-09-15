package main

import (
	"github.com/NoizeMe/go-man/pkg/logging"
	"github.com/NoizeMe/go-man/pkg/manager"
	goreleases "github.com/NoizeMe/go-man/pkg/releases"
	"github.com/hashicorp/go-version"
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

	select_        = root.SubCommand("select", "Selects the active installation of the Go SDK.")
	selectVersions = select_.Args(
		"[version]",
		"The version that should be selected.",
	)

	unselect = root.SubCommand("unselect", "Unselects the active installation of the Go SDK.")
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
	case unselect.Parsed():
		handleUnselect(*dryRun)
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
		latest, err := goreleases.GetLatest(goreleases.SelectReleaseType(all))
		logging.IfError(err)

		versionNames = []string{latest.GetVersionNumber().String()}
	}

	for _, versionName := range versionNames {
		parsedVersion, err := version.NewVersion(versionName)
		if err != nil {
			logging.IfError(err)
		}

		goManager, err := manager.NewManager(gomanRoot(), dryRun)
		logging.IfError(err)
		goManager.Install(parsedVersion, operatingSystem, arch, goreleases.SelectReleaseType(all))
	}
}

func handleRemove(dryRun bool, all bool, versionNames []string) {
	root := gomanRoot()

	logging.IfErrorf(!all && len(versionNames) == 0, "No versionNames to remove, skipping.")
	logging.IfErrorf(all && len(versionNames) > 0, "Both all flag and versionNames given, skipping.")

	goManager, err := manager.NewManager(root, dryRun)
	logging.IfError(err)

	if all {
		goManager.RemoveAll()
	} else {
		for _, versionName := range versionNames {
			versionNumber, err := version.NewVersion(versionName)
			logging.IfError(err)

			goManager.Remove(versionNumber)
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

	goManager, err := manager.NewManager(gomanRoot(), dryRun)
	logging.IfError(err)
	goManager.Select(parsedVersion)
}

func handleUnselect(dryRun bool) {
	goManager, err := manager.NewManager(gomanRoot(), dryRun)
	logging.IfError(err)
	goManager.Unselect()
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
