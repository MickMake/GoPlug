# GoPlug - A working PoC for GoLang's plugin system.

This started out as a PoC for an itch I had. I needed a decent plugin system to support arbitrary modules that could be loaded/unloaded at run-time.

The goal was to determine how simple I could get a plugin system to be. Turns out GoLang's plugin system has flaws that make it unusable for ad-hoc loading/unloading of modules.
Specifically; when modules are built with one version against a different version of the "master".

Anyway, this is a complete working example of a simple Go based plugin system.

I might extend this to support HasiCorp's gRPC based "plugins" which is, of course, a lot more mature.

