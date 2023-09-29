package GoPlug

import (
	"log"
	"sync"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

// ListPlugins implements the interface method
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) ListPlugins() {
	m.Loaders.StorePrint()
}

// GetPlugin implements the interface method
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) GetPlugin(pluginPath utils.FilePath) (*GoPlugLoader.PluginItem, Return.Error) {
	return m.Loaders.StoreGet(pluginPath.GetPath())
	// return m.store.Get(pluginPath.GetPath())
}

// GetPluginByName implements the interface method
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) GetPluginByName(name string) (*GoPlugLoader.PluginItem, Return.Error) {
	return m.Loaders.StoreGet(name)
}

// GetPlugins implements the interface method
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) GetPlugins() GoPlugLoader.PluginItems {
	return m.Loaders.StoreGetAll()
}

// RegisterPlugins - .
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) RegisterPlugins() Return.Error {
	for range Only.Once {
		log.Printf("[INFO]: Registering %d plugins.", m.Loaders.StoreSize())
		var items GoPlugLoader.PluginItems
		items, m.Error = m.Loaders.PluginRegister()
		if m.Error.IsError() {
			break
		}
		log.Printf("[INFO]: Registered %d/%d plugins.", len(items), m.Loaders.StoreSize())
	}

	return m.Error
}

// UnregisterPlugins - .
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) UnregisterPlugins() Return.Error {
	for range Only.Once {
		log.Printf("[INFO]: Unregistering %d plugins.", m.Loaders.StoreSize())
		m.Error = m.Loaders.PluginUnregister()
	}

	return m.Error
}

// LoadPlugin implements the interface method
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) LoadPlugin(pluginPath utils.FilePath) Return.Error {
	for range Only.Once {
		m.Error = pluginPath.FileExists()
		if m.Error.IsError() {
			break
		}

		pluginPath.ShortenPaths()
		base := pluginPath.SetAltPath(m.GetDir(), "[PluginDir]")
		log.Printf("[INFO]: Plugin(%s): Loading", base)

		// load
		var plug GoPlugLoader.PluginItem
		plug, m.Error = m.Loaders.PluginLoad(pluginPath)
		if m.Error.IsError() {
			log.Printf("[ERROR]: Plugin(%s): Load failed: %s", base, m.Error.String())
			break
		}
		log.Printf("[INFO]: Plugin(%s): Loaded OK - Native:%v RPC:%v\n",
			base, plug.IsNativePlugin(), plug.IsRpcPlugin())
	}

	return m.Error
}

// UnloadPlugin implements the interface method
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) UnloadPlugin(pluginPath utils.FilePath) Return.Error {
	for range Only.Once {
		pluginPath.ShortenPaths()
		base := pluginPath.SetAltPath(m.Loaders.GetDir(), "[PluginDir]")
		log.Printf("[INFO]: Plugin(%s): Unloading", base)

		m.Error = m.Loaders.PluginUnload(pluginPath)
		if m.Error.IsError() {
			break
		}

		log.Printf("[INFO]: Plugin(%s): Unloaded OK", base)
	}

	return m.Error
}

// GetInterface -
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) GetInterface(id string) (GoPlugLoader.PluginItemInterface, Return.Error) {
	var raw GoPlugLoader.PluginItemInterface
	var err Return.Error

	for range Only.Once {
		if _, ok := m.Plugins[id]; !ok {
			err.SetError("Plugin ID not found in registered plugins!")
			break
		}

		err = Return.Ok
	}

	return raw, err
}

// Scan -
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) Scan() Return.Error {

	for range Only.Once {
		m.Error = m.Loaders.PluginScan(m.FileGlob)
		if m.Error.IsError() {
			break
		}

		m.Initialized = true
		m.Error = Return.Ok
	}

	return m.Error
}

// Register -
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) Register() Return.Error {
	var err Return.Error

	for range Only.Once {
		items := m.Loaders.StoreGetAll()

		for index, info := range items {
			identity := info.Pluggable.Identify()
			log.Printf("[%d] Registering plugin client for name=%s, type=%s\n%s\n",
				index, identity.Name, info.Pluggable.GetPluginType(), identity)
			log.Printf("")
		}

		err = Return.Ok
	}

	return err
}

// Dispose -
// ---------------------------------------------------------------------------------------------------- //
func (m *PluginManager) Dispose() {
	var wg sync.WaitGroup
	for _, pinfo := range m.Plugins {
		m.UnloadPlugin(pinfo.Path)
	}

	wg.Wait()
}
