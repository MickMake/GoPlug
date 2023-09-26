package GoPlug

import (
	"log"
	"sync"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

// ---------------------------------------------------------------------------------------------------- //

// ListPlugins implements the interface method
func (m *PluginManager) ListPlugins() {
	m.Loaders.StorePrint()
}

// GetPlugin implements the interface method
func (m *PluginManager) GetPlugin(pluginPath utils.FilePath) (*GoPlugLoader.PluginItem, Return.Error) {
	return m.Loaders.StoreGet(pluginPath.GetPath())
	// return m.store.Get(pluginPath.GetPath())
}

// GetPluginByName implements the interface method
func (m *PluginManager) GetPluginByName(name string) (*GoPlugLoader.PluginItem, Return.Error) {
	return m.Loaders.StoreGet(name)
}

// GetPlugins implements the interface method
func (m *PluginManager) GetPlugins() GoPlugLoader.PluginItems {
	return m.Loaders.StoreGetAll()
}

// ---------------------------------------------------------------------------------------------------- //

func (m *PluginManager) RegisterPlugins() Return.Error {
	for range Only.Once {
		log.Printf("[INFO]: Registering %d plugins.", m.Loaders.StoreSize())
		var items GoPlugLoader.PluginItems
		items, m.Error = m.Loaders.PluginRegisterAll()
		if m.Error.IsError() {
			break
		}
		log.Printf("[INFO]: Registered %d/%d plugins.", len(items), m.Loaders.StoreSize())
	}

	return m.Error
}

func (m *PluginManager) UnregisterPlugins() Return.Error {
	for range Only.Once {
		log.Printf("[INFO]: Unregistering %d plugins.", m.Loaders.StoreSize())
		m.Error = m.Loaders.PluginUnregisterAll()
	}

	return m.Error
}

// ---------------------------------------------------------------------------------------------------- //

// LoadPlugin implements the interface method
func (m *PluginManager) LoadPlugin(pluginPath utils.FilePath) Return.Error {
	for range Only.Once {
		m.Error = pluginPath.FileExists()
		if m.Error.IsError() {
			break
		}

		base := pluginPath.SetAltPath(m.GetDir(), "[PluginDir]")
		log.Printf("[INFO]: Plugin(%s): Loading", base)

		// load
		var plug GoPlugLoader.PluginItem
		plug, m.Error = m.Loaders.PluginLoad(pluginPath)
		if m.Error.IsError() {
			log.Printf("[ERROR]: Plugin(%s): Load failed: %s", base, m.Error.String())
			break
		}
		log.Printf("[INFO]: Plugin(%s): Loaded OK\n", base)
		// log.Printf("[INFO]: Plugin(%s): Loaded OK - Name: %s Version: %s Description: '%s'",
		// 	base, plug.Common.identity.Name, plug.Common.identity.Version, plug.Common.identity.Description)
		//
		// // validate
		// var result any
		// result, m.Error = m.validator.Validate(plug.Common.identity)
		// if m.Error.IsError() {
		// 	log.Printf("[ERROR]: Plugin(%s): Validation failed: %s", base, m.Error.String())
		// 	break
		// }
		// log.Printf("[INFO]: Plugin(%s): Validation OK", base)
		//
		// // Convert validate result to plugin identity object
		// identity, ok := result.(*PluginIdentity)
		// if !ok {
		// 	log.Printf("[ERROR]: Plugin(%s): Config load failed", base)
		// 	m.Error.SetError("Failed to convert validation result to plugin identity")
		// 	break
		// }
		// log.Printf("[INFO]: Plugin(%s): Config loaded", base)
		//
		// plug.Common.filename = pluginPath
		// plug.Common.identity = identity
		//
		// // Save
		// m.Error = m.store.Put(&plug, true)
		// if m.Error.IsError() {
		// 	log.Printf("[ERROR]: Plugin(%s): Storage failed: %s", base, m.Error.String())
		// 	break
		// }
		// log.Printf("[INFO]: Plugin(%s): Stored OK", base)
		//
		// m.Error = m.Loaders.PluginInit(plug)
		// if m.Error.IsError() {
		// 	log.Printf("[ERROR]: Plugin(%s): Initialisation failed: %s", base, m.Error.String())
		// 	break
		// }

		log.Printf("[INFO]: Plugin(%s): Done\n%s\n", base, plug.Data.Common.String())
	}

	return m.Error
}

// UnloadPlugin implements the interface method
func (m *PluginManager) UnloadPlugin(pluginPath utils.FilePath) Return.Error {
	for range Only.Once {
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

// ---------------------------------------------------------------------------------------------------- //

func (m *PluginManager) GetInterface(id string) (any, Return.Error) {
	var raw any
	var err Return.Error

	for range Only.Once {
		if _, ok := m.Plugins[id]; !ok {
			err.SetError("Plugin ID not found in registered plugins!")
			break
		}

		// // Grab registerd plugin.Client
		// item := m.Plugins[id].Item
		// raw, err = item.Rpc.GetInterface()
		// if err.IsError() {
		// 	break
		// }

		// // Grab registerd plugin.Client
		// item := m.Plugins[id].Item
		// client := item.Rpc.client.GpClient
		//
		// // Connect via RPC
		// rpcClient, e := client.Client()
		// if e != nil {
		// 	err.SetError(e)
		// 	break
		// }
		//
		// // Request the plugin
		// raw, e = rpcClient.Dispense(id)
		// if e != nil {
		// 	err.SetError(e)
		// 	break
		// }

		err = Return.Ok
	}

	return raw, err
}
func (m *PluginManager) GetRpcInterface(id string) (GoPlugLoader.RpcPluginInterface, Return.Error) {
	var plug GoPlugLoader.RpcPluginInterface
	var err Return.Error

	for range Only.Once {
		// var p any
		// p, err = m.GetInterface(id)
		// if err.IsError() {
		// 	break
		// }
		//
		// plug = p.(GoPlugLoader.RpcPluginInterface)
		//
		// m.Logger.Info("Plugin struct name: %s", plug.Serve().GetStructName())
	}

	return plug, err
}
func (m *PluginManager) GetNativeInterface(id string) (GoPlugLoader.NativePluginInterface, Return.Error) {
	var plug GoPlugLoader.NativePluginInterface
	var err Return.Error

	for range Only.Once {
		var p any
		p, err = m.GetInterface(id)
		if err.IsError() {
			break
		}

		// @TODO - Not currently getting the natve interface.

		plug = p.(GoPlugLoader.NativePluginInterface)

		m.Logger.Info("Plugin struct name: %s", plug.GetStructName())
	}

	return plug, err
}

// ---------------------------------------------------------------------------------------------------- //

func (m *PluginManager) Scan() Return.Error {

	for range Only.Once {
		// plugins, err = m.Base.Scan(m.fileGlob)
		// plugins, e := plugin.Discover(m.fileGlob, m.Base.GetPath())
		m.Error = m.Loaders.PluginScan(m.FileGlob)
		if m.Error.IsError() {
			break
		}

		// // Native plugin files
		// loader := m.Loaders.GetLoader(NativeLoaderName)
		// if loader != nil {
		// 	for _, foo := range loader.GetFiles() {
		// 		for _, plug := range foo.Paths {
		// 			id := FileToId(plug.GetBase(), m.fileGlob)
		// 			m.Plugins[id] = &PluginInfo{
		// 				Types: NativePluginType,
		// 				// PluginTypes{
		// 				// 	Rpc:    false,
		// 				// 	Native: true,
		// 				// },
		// 				ID:   id,
		// 				Path: plug,
		// 				Item: PluginItem{
		// 					// Common: PluginCommon{},
		// 					// Native: NativePlugin{},
		// 					// Rpc:    RpcPlugin{},
		// 				},
		// 			}
		// 		}
		// 	}
		// }
		//
		// // RPC plugin files
		// loader = m.Loaders.GetLoader(RpcLoaderName)
		// // if loader.IsLoaderType(RpcLoaderName) {
		// if loader != nil {
		// 	for _, foo := range loader.GetFiles() {
		// 		for _, plug := range foo.Paths {
		// 			id := FileToId(plug.GetBase(), m.fileGlob)
		// 			m.Plugins[id] = &PluginInfo{
		// 				Types: RpcPluginType,
		// 				// 	PluginTypes{
		// 				// 	Rpc:    true,
		// 				// 	Native: false,
		// 				// },
		// 				ID:   id,
		// 				Path: plug,
		// 				Item: PluginItem{
		// 					// Common: PluginCommon{},
		// 					// Native: NativePlugin{},
		// 					// Rpc:    RpcPlugin{},
		// 				},
		// 			}
		// 		}
		// 	}
		// }

		m.Initialized = true
		m.Error = Return.Ok
	}

	return m.Error
}

// ---------------------------------------------------------------------------------------------------- //

func (m *PluginManager) Register() Return.Error {
	var err Return.Error

	for range Only.Once {
		items := m.Loaders.StoreGetAll()

		for index, info := range items {
			identity := info.Data.Identify()
			log.Printf("[%d] Registering plugin client for name=%s, type=%s\n%s\n",
				index, identity.Name, info.Data.Common.GetPluginType(), identity)
			log.Printf("")
		}

		// for id, info := range m.Plugins {
		// 	log.Printf("Registering plugin client for type=%s, id=%s, impl=%s",
		// 		m.pluginType, id, info.Path)
		//
		// 	// var l Logger
		// 	// l, err = NewLogger(fmt.Sprintf("Manager(%s)", id), "")
		// 	// if err.IsError() {
		// 	// 	break
		// 	// }
		// 	//
		// 	// l := hclog.New(&hclog.LoggerOptions{
		// 	// 	Name:                     "PLUGIN:" + id,
		// 	// 	Level:                    m.logger.GetLevel(),
		// 	// 	Output:                   nil,
		// 	// 	Mutex:                    nil,
		// 	// 	JSONFormat:               false,
		// 	// 	IncludeLocation:          true,
		// 	// 	AdditionalLocationOffset: 0,
		// 	// 	TimeFormat:               "2006-01-02 15:04:05",
		// 	// 	TimeFn:                   nil,
		// 	// 	DisableTime:              false,
		// 	// 	Color:                    hclog.ForceColor,
		// 	// 	ColorHeaderOnly:          true,
		// 	// 	ColorHeaderAndFields:     true,
		// 	// 	Exclude:                  nil,
		// 	// 	IndependentLevels:        false,
		// 	// 	SubloggerHook:            nil,
		// 	// })
		// 	//
		// 	// cfg := plugin.ClientConfig{
		// 	// 	HandshakeConfig:  config.HandshakeConfig,
		// 	// 	Plugins:          m.pluginMap(id),
		// 	// 	VersionedPlugins: nil,
		// 	// 	Cmd:              exec.Command(info.Path),
		// 	// 	Reattach:         nil,
		// 	// 	RunnerFunc:       nil,
		// 	// 	SecureConfig:     nil,
		// 	// 	TLSConfig:        nil,
		// 	// 	Managed:          true,
		// 	// 	MinPort:          0,
		// 	// 	MaxPort:          0,
		// 	// 	StartTimeout:     time.Second * 30,
		// 	// 	Stderr:           os.Stderr,
		// 	// 	SyncStdout:       nil,
		// 	// 	SyncStderr:       os.Stderr,
		// 	// 	AllowedProtocols: nil,
		// 	// 	Logger:           l,
		// 	// 	AutoMTLS:         false,
		// 	// 	GRPCDialOptions:  nil,
		// 	// 	SkipHostEnv:      false,
		// 	// 	UnixSocketConfig: nil,
		// 	// }
		// 	// // create new client
		// 	// client := plugin.NewClient(&cfg)
		// 	//
		// 	// if _, ok := m.Plugins[id]; !ok {
		// 	// 	// if not found, ignore?
		// 	// 	continue
		// 	// }
		// 	// pinfo := m.Plugins[id]
		// 	// pinfo.RpcClient = client
		//
		// 	if m.pluginType.Native {
		// 		loader := m.Loaders.GetLoader(NativeLoaderName)
		// 		if loader == nil {
		// 			continue
		// 		}
		// 		// if !loader.IsLoaderType(NativeLoaderName) {
		// 		// 	break
		// 		// }
		// 		loader.Load(m.pluginType, info.Path)
		// 	}
		//
		// 	if m.pluginType.Native {
		// 		loader := m.Loaders.GetLoader(RpcLoaderName)
		// 		if loader == nil {
		// 			continue
		// 		}
		// 		// if !loader.IsLoaderType(RpcLoaderName) {
		// 		// 	break
		// 		// }
		// 		loader.Load(m.pluginType, info.Path)
		// 	}
		//
		// 	// if info.Path.HasExtension(NativePluginExtensions...) {
		// 	// 	loader := m.Loaders.GetLoader(NativeLoaderName)
		// 	// 	if loader == nil {
		// 	// 		continue
		// 	// 	}
		// 	// 	loader.Load(pluginTypes, info.Path)
		// 	// 	continue
		// 	// }
		// 	//
		// 	// loader := m.Loaders.GetLoader(RpcLoaderName)
		// 	// if loader == nil {
		// 	// 	continue
		// 	// }
		// 	// loader.Load(pluginTypes, info.Path)
		// }

		err = Return.Ok
	}

	return err
}

// ---------------------------------------------------------------------------------------------------- //

func (m *PluginManager) Dispose() {
	var wg sync.WaitGroup
	for _, pinfo := range m.Plugins {
		// if pinfo.Item.Rpc.IsConfigured() {
		// 	wg.Add(1)
		// 	log.Printf("Unloading RPC plugin: ")
		// 	go func(client *plugin.Client) {
		// 		client.Kill()
		// 		wg.Done()
		// 	}(pinfo.Item.Rpc.RpcClient.GpClient)
		// }

		// if pinfo.Item.Native.IsConfigured() {
		// }
		m.UnloadPlugin(pinfo.Path)
	}

	wg.Wait()
}
