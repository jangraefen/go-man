package main

import (
	"github.com/NoizeMe/go-man/pkg/logging"
	"github.com/NoizeMe/go-man/pkg/manager"
	"github.com/NoizeMe/go-man/pkg/releases"
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
		"By passing this flag, all actions are done as if they were successful, but no changes will be performed.",
	)

	list    = root.SubCommand("list", "Lists all available releases of the Go SDK.")
	listAll = list.Bool(
		"all",
		false,
		"By passing this flag, non-stable releases are listed as well.",
	)

	install    = root.SubCommand("install", "Installs one or more new versions of the Go SDK.")
	installAll = install.Bool(
		"all",
		false,
		"By passing this flag, non-stable releases are installable as well.",
	)
	installOS = install.String(
		"os",
		runtime.GOOS,
		"Defines the operating system for that the Go SDK will be downloaded.",
	)
	installArch = install.String(
		"arch",
		runtime.GOARCH,
		"Defines the processor architecture for that Go SDK will be downloaded.",
	)
	installVersions = install.Args(
		"[versions...]",
		"Versions that should be installed. May be 'latest' or any version number.",
	)

	uninstall    = root.SubCommand("uninstall", "Uninstall an existing installation of the Go SDK.")
	uninstallAll = uninstall.Bool(
		"all",
		false,
		"If set, all installed versions will be removed.",
	)
	uninstallVersions = uninstall.Args(
		"[versions...]",
		"The versions that should be removed.",
	)

	select_        = root.SubCommand("select", "Selects the active installation of the Go SDK.")
	selectVersions = select_.Args(
		"[version]",
		"The version that should be selected.",
	)

	unselect = root.SubCommand("unselect", "Unselects the active installation of the Go SDK.")

	cleanup = root.SubCommand("cleanup", "Removes all installations of the Go SDK that are not considered stable.")
)

func main() {
	// Parse the command line arguments. Any errors will get caught be the library and will cause the usage to be printed.
	// The program will exit afterwards.
	_ = root.Parse()

	if stat, err := os.Stat(gomanRoot()); err != nil && os.IsNotExist(err) || !stat.IsDir() {
		logging.IfError(os.MkdirAll(gomanRoot(), 0755))
	}

	switch {
	case list.Parsed():
		handleList(*listAll)
	case install.Parsed():
		handleInstall(*dryRun, *installAll, *installOS, *installArch, *installVersions)
	case uninstall.Parsed():
		handleUninstall(*dryRun, *uninstallAll, *uninstallVersions)
	case select_.Parsed():
		handleSelect(*dryRun, *selectVersions)
	case unselect.Parsed():
		handleUnselect(*dryRun)
	case cleanup.Parsed():
		handleCleanup(*dryRun)
	}
}

func handleList(all bool) {
	logging.Printf("List of available releases:")

	releaseList, err := releases.ListAll(releases.SelectReleaseType(all))
	logging.IfTaskError(err)

	for _, r := range releaseList {
		logging.TaskPrintf("%s", r.GetVersionNumber())
	}
}

func handleInstall(dryRun, all bool, operatingSystem, arch string, versionNames []string) {
	logging.IfErrorf(len(versionNames) == 0, "No versions given to install, skipping")

	if len(versionNames) == 1 && versionNames[0] == "latest" {
		latest, err := releases.GetLatest(releases.SelectReleaseType(all))
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
		goManager.Install(parsedVersion, operatingSystem, arch, releases.SelectReleaseType(all))
	}
}

func handleUninstall(dryRun bool, all bool, versionNames []string) {
	root := gomanRoot()

	logging.IfErrorf(!all && len(versionNames) == 0, "No versionNames to uninstall, skipping.")
	logging.IfErrorf(all && len(versionNames) > 0, "Both all flag and versionNames given, skipping.")

	goManager, err := manager.NewManager(root, dryRun)
	logging.IfError(err)

	if all {
		goManager.UninstallAll()
	} else {
		for _, versionName := range versionNames {
			versionNumber, err := version.NewVersion(versionName)
			logging.IfError(err)

			goManager.Uninstall(versionNumber)
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

func handleCleanup(dryRun bool) {
	goManager, err := manager.NewManager(gomanRoot(), dryRun)
	logging.IfError(err)
	goManager.Cleanup()
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
