package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-version"
	"github.com/posener/cmd"

	"github.com/NoizeMe/go-man/pkg/manager"
	"github.com/NoizeMe/go-man/pkg/releases"
	"github.com/NoizeMe/go-man/pkg/tasks"
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

	selectz        = root.SubCommand("select", "Selects the active installation of the Go SDK.")
	selectVersions = selectz.Args(
		"[version]",
		"The version that should be selected.",
	)

	unselect = root.SubCommand("unselect", "Unselects the active installation of the Go SDK.")

	cleanup = root.SubCommand("cleanup", "Removes all installations of the Go SDK that are not considered stable.")
)

func main() {
	task := &tasks.Task{
		ErrorExitCode: 1,
		Output:        os.Stdout,
		Error:         os.Stderr,
	}

	// Parse the command line arguments. Any errors will get caught be the library and will cause the usage to be printed.
	// The program will exit afterwards.
	_ = root.Parse()

	if stat, err := os.Stat(gomanRoot()); err != nil && os.IsNotExist(err) || !stat.IsDir() {
		task.IfError(os.MkdirAll(gomanRoot(), 0755))
	}

	switch {
	case list.Parsed():
		handleList(task, *listAll)
	case install.Parsed():
		handleInstall(task, *dryRun, *installAll, *installOS, *installArch, *installVersions)
	case uninstall.Parsed():
		handleUninstall(task, *dryRun, *uninstallAll, *uninstallVersions)
	case selectz.Parsed():
		handleSelect(task, *dryRun, *selectVersions)
	case unselect.Parsed():
		handleUnselect(task, *dryRun)
	case cleanup.Parsed():
		handleCleanup(task, *dryRun)
	}
}

func handleList(task *tasks.Task, all bool) {
	task.Printf("List of available releases:")

	releaseList, err := releases.ListAll(releases.SelectReleaseType(all))
	task.IfTaskError(err)

	for _, r := range releaseList {
		task.Printf("%s", r.GetVersionNumber())
	}
}

func handleInstall(task *tasks.Task, dryRun, all bool, operatingSystem, arch string, versionNames []string) {
	task.IfErrorf(len(versionNames) == 0, "No versions given to install, skipping")

	if len(versionNames) == 1 && versionNames[0] == "latest" {
		latest, err := releases.GetLatest(releases.SelectReleaseType(all))
		task.IfError(err)

		versionNames = []string{latest.GetVersionNumber().String()}
	}

	for _, versionName := range versionNames {
		parsedVersion, err := version.NewVersion(versionName)
		if err != nil {
			task.IfError(err)
		}

		goManager, err := manager.NewManager(task, gomanRoot(), dryRun)
		task.IfError(err)
		goManager.Install(parsedVersion, operatingSystem, arch, releases.SelectReleaseType(all))
	}
}

func handleUninstall(task *tasks.Task, dryRun bool, all bool, versionNames []string) {
	root := gomanRoot()

	task.IfErrorf(!all && len(versionNames) == 0, "No versions to uninstall, skipping.")
	task.IfErrorf(all && len(versionNames) > 0, "Both all flag and versions given, skipping.")

	goManager, err := manager.NewManager(task, root, dryRun)
	task.IfError(err)

	if all {
		goManager.UninstallAll()
	} else {
		for _, versionName := range versionNames {
			versionNumber, err := version.NewVersion(versionName)
			task.IfError(err)

			goManager.Uninstall(versionNumber)
		}
	}
}

func handleSelect(task *tasks.Task, dryRun bool, versionNames []string) {
	task.IfErrorf(len(versionNames) == 0, "No version to select, skipping.")
	task.IfErrorf(len(versionNames) > 1, "More then one version to select, skipping.")

	parsedVersion, err := version.NewVersion(versionNames[0])
	if err != nil {
		task.IfError(err)
	}

	goManager, err := manager.NewManager(task, gomanRoot(), dryRun)
	task.IfError(err)
	goManager.Select(parsedVersion)
}

func handleUnselect(task *tasks.Task, dryRun bool) {
	goManager, err := manager.NewManager(task, gomanRoot(), dryRun)
	task.IfError(err)
	goManager.Unselect()
}

func handleCleanup(task *tasks.Task, dryRun bool) {
	goManager, err := manager.NewManager(task, gomanRoot(), dryRun)
	task.IfError(err)
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
