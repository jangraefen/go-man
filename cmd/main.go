package main

import (
	"fmt"
	goreleases "github.com/NoizeMe/go-man/pkg/releases"
	"github.com/posener/cmd"
	"os"
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

	install         = root.SubCommand("install", "This sub command is used to install new version of the Go SDK.")
	installVersions = install.Args(
		"[versions]",
		"The version that should be installed. May be 'latest' or any version number.",
	)
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
	}
}

func handleReleases(all bool) {
	fmt.Println("List of available releases:")

	releaseList, err := goreleases.ListAll(goreleases.SelectReleaseType(all))
	if err != nil {
		handleError("[-] %s", err)
	}

	for _, r := range releaseList {
		fmt.Printf("[+] %s\n", r.GetVersionNumber())
	}
}

func handleInstall(dryRun, all bool, operatingSystem, arch string, versions []string) {
	if len(versions) == 0 {
		latest, err := goreleases.GetLatest()
		if err != nil {
			handleError("[-] %s", err)
		}

		versions = []string{latest.GetVersionNumber()}
	}

	for _, version := range versions {
		fmt.Printf("Installing %s %s-%s:\n", version, operatingSystem, arch)

		release, releasePresent, err := goreleases.GetForVersion(goreleases.SelectReleaseType(all), version)
		if err != nil {
			handleError("[-] %s", err)
		}
		if !releasePresent {
			handleError("[-] %s", "release with version "+version+" not present")
		}

		files := release.FindFiles(operatingSystem, arch, goreleases.ArchiveFile)
		if len(files) == 0 {
			handleError("[-] %s", "release "+version+" with "+operatingSystem+"/"+arch+" not present")
		}

		for _, file := range files {
			fmt.Printf("[+] Considering file: %s\n", file.GetUrl())
		}
	}
}

func handleError(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(2)
}
