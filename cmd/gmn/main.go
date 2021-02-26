package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-version"
	"github.com/posener/cmd"
	"github.com/posener/complete/v2/predict"

	"github.com/jangraefen/go-man/internal/fileutil"
	"github.com/jangraefen/go-man/pkg/manager"
	"github.com/jangraefen/go-man/pkg/releases"
	"github.com/jangraefen/go-man/pkg/tasks"
)

var (
	root = cmd.New(
		cmd.OptName("gmn"),
		cmd.OptDetails("A manager for Go installations"),
	)

	list         = root.SubCommand("list", "Lists of all available Go releases")
	listUnstable = list.Bool(
		"unstable",
		false,
		"Unlocks the listing of unstable Go versions",
	)

	install         = root.SubCommand("install", "Installs one or more new Go releases")
	installUnstable = install.Bool(
		"unstable",
		false,
		"Unlocks the installation of unstable Go versions",
	)
	installOS = install.String(
		"os",
		runtime.GOOS,
		"Operating system for that Go will be installed",
		predict.OptValues("freebsd", "darwin", "linux", "windows"),
		predict.OptCheck(),
	)
	installArch = install.String(
		"arch",
		runtime.GOARCH,
		"Processor architecture for that Go will be installed",
		predict.OptValues("386", "amd64", "armv61", "ppc64le", "s390x"),
		predict.OptCheck(),
	)
	installVersions = install.Args(
		"[versions...]",
		"Versions of Go that will be installed. 'latest' or any version number",
	)

	uninstall    = root.SubCommand("uninstall", "Uninstall an existing Go installation")
	uninstallAll = uninstall.Bool(
		"all",
		false,
		"If set, all installations of Go will be uninstalled",
	)
	uninstallVersions = uninstall.Args(
		"[versions...]",
		"The versions that should be uninstalled",
	)

	selectz        = root.SubCommand("select", "Selects the default Go installation")
	selectVersions = selectz.Args(
		"[version]",
		"The version that should be selected",
	)

	unselect = root.SubCommand("unselect", "Unselects the default Go installation")

	cleanup = root.SubCommand("cleanup", "Removes all Go installations, that are not considered stable")
)

func main() {
	task := &tasks.Task{
		ErrorExitCode: 1,
		Output:        os.Stdout,
		Error:         os.Stderr,
	}

	if !fileutil.PathExists(gomanRoot()) {
		task.FatalOnError(os.MkdirAll(gomanRoot(), 0755))
	}

	// Parse the command line arguments. Any errors will get caught be the library and will cause the usage to be printed.
	// The program will exit afterwards.
	_ = root.Parse()

	switch {
	case list.Parsed():
		handleList(task, *listUnstable)
	case install.Parsed():
		handleInstall(task, *installUnstable, *installOS, *installArch, *installVersions)
	case uninstall.Parsed():
		handleUninstall(task, *uninstallAll, *uninstallVersions)
	case selectz.Parsed():
		handleSelect(task, *selectVersions)
	case unselect.Parsed():
		handleUnselect(task)
	case cleanup.Parsed():
		handleCleanup(task)
	}
}

func handleList(task *tasks.Task, all bool) {
	task.Printf("List of available releases:")
	listTask := task.Step()

	releaseList, err := releases.ListAll(releases.SelectReleaseType(all))
	listTask.FatalOnError(err)

	for _, r := range releaseList {
		listTask.Printf("%s", r.GetVersionName())
	}
}

func handleInstall(task *tasks.Task, unstable bool, operatingSystem, arch string, versionNames []string) {
	task.FatalIff(len(versionNames) == 0, "No versions given to install, skipping")

	if len(versionNames) == 1 && versionNames[0] == "latest" {
		latest, err := releases.GetLatest(releases.SelectReleaseType(unstable))
		task.FatalOnError(err)

		versionNames = []string{latest.GetVersionNumber().String()}
	}

	for _, versionName := range versionNames {
		parsedVersion, err := version.NewVersion(versionName)
		if err != nil {
			task.FatalOnError(err)
		}

		goManager, err := manager.NewManager(task, gomanRoot())
		task.FatalOnError(err)
		task.FatalOnError(goManager.Install(parsedVersion, operatingSystem, arch, releases.SelectReleaseType(unstable)))
	}
}

func handleUninstall(task *tasks.Task, all bool, versionNames []string) {
	root := gomanRoot()

	task.FatalIff(!all && len(versionNames) == 0, "No versions to uninstall, skipping.")
	task.FatalIff(all && len(versionNames) > 0, "Both all flag and versions given, skipping.")

	goManager, err := manager.NewManager(task, root)
	task.FatalOnError(err)

	if all {
		task.FatalOnError(goManager.UninstallAll())
	} else {
		for _, versionName := range versionNames {
			versionNumber, err := version.NewVersion(versionName)
			task.FatalOnError(err)
			task.FatalOnError(goManager.Uninstall(versionNumber))
		}
	}
}

func handleSelect(task *tasks.Task, versionNames []string) {
	task.FatalIff(len(versionNames) == 0, "No version to select, skipping.")
	task.FatalIff(len(versionNames) > 1, "More then one version to select, skipping.")

	parsedVersion, err := version.NewVersion(versionNames[0])
	if err != nil {
		task.FatalOnError(err)
	}

	goManager, err := manager.NewManager(task, gomanRoot())
	task.FatalOnError(err)
	task.FatalOnError(goManager.Select(parsedVersion))
}

func handleUnselect(task *tasks.Task) {
	goManager, err := manager.NewManager(task, gomanRoot())
	task.FatalOnError(err)
	task.FatalOnError(goManager.Unselect())
}

func handleCleanup(task *tasks.Task) {
	goManager, err := manager.NewManager(task, gomanRoot())
	task.FatalOnError(err)
	task.FatalOnError(goManager.Cleanup())
}

func gomanRoot() string {
	root := os.Getenv("GMNROOT")
	if len(root) > 0 {
		return root
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("/", ".gmn")
	}

	return filepath.Join(homeDir, ".gmn")
}
