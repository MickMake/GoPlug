package GoPlug

import (
	"log"

	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoPlug/pluggable"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoUnify/Only"
)

// LoadPlugins implements the interface method
func (bm *PluginManager) LoadPlugins() Return.Error {
	for range Only.Once {
		// scan plugin base dir
		var paths utils.PluginPaths
		paths, bm.Error = bm.loader.Scan("so")
		if bm.Error.IsError() {
			break
		}

		log.Printf("[INFO]: Found %d possible plugins", paths.Length())

		if !paths.AnyPaths() {
			// No plugin files found
			break
		}

		for _, pDir := range paths {
			log.Printf("[INFO]: %d plugin files found in %s", pDir.Length(), pDir.Dir.String())
			for _, pName := range pDir.Get() {
				if !pName.HasExtension("so") {
					continue
				}
				bm.Error = bm.loadPlugin(pName)
				if bm.Error.IsError() {
					continue
				}
				bm.Error = bm.initPlugin(pName)
				if bm.Error.IsError() {
					continue
				}
			}
		}

		log.Printf("[INFO]: %d plugins loaded", bm.store.Size())
	}

	return bm.Error
}

// LoadPlugin implements the interface method
func (bm *PluginManager) LoadPlugin(p string) Return.Error {
	for range Only.Once {
		name := bm.stringToPluginPath(p)
		if bm.Error.IsError() {
			break
		}

		bm.Error = bm.loadPlugin(name)
		if bm.Error.IsError() {
			break
		}

		bm.Error = bm.initPlugin(name)
		if bm.Error.IsError() {
			break
		}
	}

	return bm.Error
}

// UnloadPlugin implements the interface method
func (bm *PluginManager) UnloadPlugin(p string) Return.Error {
	for range Only.Once {
		name := bm.stringToPluginPath(p)
		if bm.Error.IsError() {
			break
		}

		bm.Error = bm.unloadPlugin(name)
	}

	return bm.Error
}

// InitPlugin implements the interface method
func (bm *PluginManager) InitPlugin(p string) Return.Error {
	for range Only.Once {
		name := bm.stringToPluginPath(p)
		if bm.Error.IsError() {
			break
		}

		bm.Error = bm.initPlugin(name)
	}

	return bm.Error
}

func (bm *PluginManager) loadPlugin(pluginPath utils.PluginPath) Return.Error {
	for range Only.Once {
		bm.Error = pluginPath.FileExists()
		if bm.Error.IsError() {
			break
		}

		base := pluginPath.SetAltPath(bm.loader.GetDir(), "[PluginDir]")
		log.Printf("[INFO]: Plugin(%s): Loading", base)

		// load
		var plug *pluggable.PluginItem
		plug, bm.Error = bm.loader.Load(pluginPath)
		if bm.Error.IsError() {
			log.Printf("[ERROR]: Plugin(%s): Load failed: %s", base, bm.Error.String())
			break
		}
		log.Printf("[INFO]: Plugin(%s): Loaded OK - Name: %s Version: %s Description: '%s'",
			base, plug.Config.Name, plug.Config.Version, plug.Config.Description)

		// validate
		var result interface{}
		result, bm.Error = bm.validator.Validate(plug.Config)
		if bm.Error.IsError() {
			log.Printf("[ERROR]: Plugin(%s): Validation failed: %s", base, bm.Error.String())
			break
		}
		log.Printf("[INFO]: Plugin(%s): Validation OK", base)

		// Convert validate result to plugin identity object
		identity, ok := result.(*pluggable.PluginIdentity)
		if !ok {
			log.Printf("[ERROR]: Plugin(%s): Config load failed", base)
			bm.Error.SetError("Failed to convert validation result to plugin identity")
			break
		}
		log.Printf("[INFO]: Plugin(%s): Config loaded", base)

		plug.File = pluginPath
		plug.Config = identity

		// Save
		bm.Error = bm.store.Put(plug, true)
		if bm.Error.IsError() {
			log.Printf("[ERROR]: Plugin(%s): Storage failed: %s", base, bm.Error.String())
			break
		}
		log.Printf("[INFO]: Plugin(%s): Stored OK", base)
	}

	return bm.Error
}

func (bm *PluginManager) initPlugin(pluginPath utils.PluginPath) Return.Error {
	for range Only.Once {
		var plug *pluggable.PluginItem
		plug, bm.Error = bm.store.Get(pluginPath.GetPath())
		if bm.Error.IsError() {
			break
		}

		base := pluginPath.SetAltPath(bm.loader.GetDir(), "[PluginDir]")
		bm.Error = bm.loader.Init(plug)
		if bm.Error.IsError() {
			bm.Error.SetError("Plugin(%s): Initialisation failed: %s", base, bm.Error)
			break
		}
	}

	return bm.Error
}

// UnloadPlugin implements the interface method
func (bm *PluginManager) unloadPlugin(pluginPath utils.PluginPath) Return.Error {
	for range Only.Once {
		base := pluginPath.SetAltPath(bm.loader.GetDir(), "[PluginDir]")
		log.Printf("[INFO]: Plugin(%s): Unloading", base)

		var plug *pluggable.PluginItem
		plug, bm.Error = bm.store.Get(pluginPath.GetPath())
		if bm.Error.IsError() {
			break
		}

		_, bm.Error = bm.store.Remove(plug.File.GetPath())
		if bm.Error.IsError() {
			bm.Error.SetError("[INFO]: Plugin(%s): Unload FAILED", base)
			log.Printf("[INFO]: Plugin(%s): Unload FAILED", base)
			break
		}
		log.Printf("[INFO]: Plugin(%s): Unloaded OK", base)
	}

	return bm.Error
}
