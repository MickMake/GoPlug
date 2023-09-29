# GoPlug - An easy, flexible pluggable plugin package for GoLang.

This started out as a PoC for an itch I had.
I needed a decent plugin system to support arbitrary modules that could be loaded/unloaded at run-time.

Turns out GoLang's plugin system has flaws that make it unusable for ad-hoc loading/unloading of modules.
Specifically:
- When modules are built with one version against a different version of the "master".
- Windows is unsupported, (all my GoLang apps are cross-architecture).

This package now supports two variants of a "plugin" system:
- Native GoLang.
- Hashicorp's gRPC based plugin system.

Its goal is still to be as simple as possible for quick onboarding, but extensible enough to support complex setups.

You can write a plugin system using the same config, but supporting either or both - at the same time!

Have a look at the examples to see what can be done.
