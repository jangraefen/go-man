# go-man

A manager for Go SDK installations.

[![Build Status](https://img.shields.io/github/workflow/status/jangraefen/go-man/Build?logo=GitHub)](https://github.com/jangraefen/go-man/actions?query=workflow:Build)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/jangraefen/go-man)](https://pkg.go.dev/mod/github.com/jangraefen/go-man)
[![Coverage](https://img.shields.io/codecov/c/github/jangraefen/go-man?logo=codecov)](https://codecov.io/gh/jangraefen/go-man)
[![Go Report Card](https://goreportcard.com/badge/github.com/jangraefen/go-man)](https://goreportcard.com/report/github.com/jangraefen/go-man)

Many other language have small helper tools that help to manage one or many installations of the SDK required to develop in
that language. Popular examples are [sdkman](https://sdkman.io/) for Java, [rvm](https://rvm.io/) for Ruby or
[nvm](https://github.com/nvm-sh/nvm) for Node.js. go-man aims to provide a similar feature for Go developers.

## Usage

To get an overview on how to use gmn, run `gmn -help` or `gmn <sub-command> -help`. Currently, the following subcommands are
implemented:

- `gmn cleanup` Removes all Go installations, that are not considered stable.
- `gmn install [flags] [versions...]` Installs one or more new Go releases
	- `-arch value` Processor architecture for that Go will be installed (defaults to your current arch)
	- `-os value` Operating system for that Go will be installed (defaults to your current OS)
	- `-unstable` Unlocks the installation of unstable Go versions
- `gmn list [flags]` Lists of all available Go releases
	- `-unstable` Unlocks the listing of unstable Go versions
- `gmn select [version]` Selects the default Go installation
- `gmn uninstall [flags] [versions...]` Uninstall an existing Go installation
	- `-all` If set, all installations of Go will be uninstalled
- `gmn unselect` Unselects the default Go installation
