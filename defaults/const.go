package defaults

import _ "embed"

// Need to execute `go generate -v -x defaults/const.go` OR `go generate -v -x ./...`
//go:generate cp ../README.md README.md
//go:generate cp ../EXAMPLES.md EXAMPLES.md
//go:generate cp ../CHANGELOG.md CHANGELOG.md

//go:embed README.md
var Readme string

//go:embed EXAMPLES.md
var Examples string

//go:embed CHANGELOG.md
var ChangeLog string

const (
	Description   = "GoPlug - An easy, flexible pluggable plugin package for GoLang."
	BinaryName    = "GoPlug"
	BinaryVersion = "1.0.1"
	SourceRepo    = "github.com/MickMake/" + BinaryName
	BinaryRepo    = "github.com/MickMake/" + BinaryName

	EnvPrefix = "GOPLUG"

	HelpSummary = `
GoPlug - An easy, flexible pluggable plugin package for GoLang.

Turns out GoLang's plugin system has flaws that make it unusable for ad-hoc loading/unloading of modules.
Specifically:
- Modules have to be built with the same Go version and build args as the "master".
- Windows is unsupported, (all my GoLang apps are cross-architecture).

This package supports two plugin variants:
- Native GoLang plugins.
- Hashicorp's gRPC based plugins.

Its goal is to be as simple as possible for quick onboarding, but extensible enough to support complex setups.

You can write a plugin system using the same config, but supporting either or both - at the same time!

It will also support pluggable plugin systems, so future plugin variants can be supported.

`
)
