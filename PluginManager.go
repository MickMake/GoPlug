package GoPlug

import (
	"log"
	"os"
	"strings"

	"github.com/MickMake/GoUnify/Only"
	plugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

// DefaultManager is the default plugin manager
//goland:noinspection GoUnusedGlobalVariable
// var DefaultManager, _ = NewPluginManager(".", &PluginManagerIdentity{})

//
// Manager - Used to load, organize and maintain the plugins.
// ---------------------------------------------------------------------------------------------------- //
type Manager interface {
	SetPluginTypes(pluginTypes Plugin.Types) Return.Error

	// SetDir - Set the base dir where to load plugins.
	// If the dir does not exist, or it's not a dir an error will be returned.
	SetDir(dir string) Return.Error

	GetDir() string

	SetFileGlob(glob string) Return.Error

	// SetIdentity - Set the base dir where to load plugins.
	// If the dir does not exist, or it's not a dir an error will be returned.
	SetIdentity(config Plugin.Identity) Return.Error

	SetImplementor(impl plugin.Plugin) Return.Error

	// ListPlugins - Print out all the plugins found.
	ListPlugins()

	// RegisterPlugins - Load all the plugins from the base plugin dir.
	RegisterPlugins() Return.Error

	// UnregisterPlugins - Unload all the plugins from the base plugin dir.
	UnregisterPlugins() Return.Error

	// LoadPlugin - Load plugin with the specified name.
	LoadPlugin(pluginPath utils.FilePath) Return.Error

	// UnloadPlugin - Unload plugin with the specified name.
	UnloadPlugin(pluginPath utils.FilePath) Return.Error

	// GetPlugin - Get the plugin with the specified name.
	GetPlugin(pluginPath utils.FilePath) (*GoPlugLoader.PluginItem, Return.Error)

	// GetPluginByName - Get the plugin with the specified name.
	GetPluginByName(name string) (*GoPlugLoader.PluginItem, Return.Error)

	GetPlugins() GoPlugLoader.PluginItems

	// CheckPlugin - Get the plugin with the specified name.
	CheckPlugin(pluginPath utils.FilePath) (*GoPlugLoader.PluginItem, Return.Error)

	// BuildPlugins - Build plugins.
	BuildPlugins() Return.Error

	Scan() Return.Error

	Register() Return.Error

	Dispose()

	GetInterface(id string) (any, Return.Error)

	GetNativeInterface(id string) (GoPlugLoader.NativePluginInterface, Return.Error)

	GetRpcInterface(id string) (GoPlugLoader.RpcPluginInterface, Return.Error)

	IsError() bool
	GetError() Return.Error
	PrintError()

	pluginMap(id string) map[string]plugin.Plugin
	NameToPluginPath(id string) (*utils.FilePath, Return.Error)

	SetLogfile(s string) Return.Error
}

//
// PluginManager
// ---------------------------------------------------------------------------------------------------- //
type PluginManager struct {
	Config      *Plugin.Identity
	PluginDir   utils.FilePath
	CmdFile     utils.FilePath
	FileGlob    string                     // glob match for plugin filenames
	Plugins     GoPlugLoader.PluginInfoMap // Info for found plugins
	Initialized bool                       // has been Initialized
	PluginImpl  plugin.Plugin              // Plugin implementation dummy interface
	Loaders     GoPlugLoader.LoaderInterface
	Validator   Plugin.Validator
	Logger      *utils.Logger
	Logfile     *utils.FilePath
	Error       Return.Error
}

// NewPluginManager is constructor of PluginManager
func NewPluginManager(config *Plugin.Identity) (Manager, Return.Error) {
	var manager Manager
	var err Return.Error

	for range Only.Once {
		err.SetPrefix("GoPlugManager: ")

		if config == nil {
			err.SetError("PluginIdentity is nil")
			break
		}

		var l utils.Logger
		l, err = utils.NewLogger("GoPlugManager", "GoPlugManager.log")
		if err.IsError() {
			break
		}

		var file utils.FilePath
		file, err = utils.NewFile(os.Args[0])
		if err.IsError() {
			break
		}

		var base utils.FilePath
		base, err = utils.NewDir(file.GetDir())
		if err.IsError() {
			break
		}

		var impl GoPlugLoader.RpcDefaultStruct

		manager = &PluginManager{
			Config:      config,
			PluginDir:   base,
			CmdFile:     file,
			FileGlob:    "goplug-*",
			Plugins:     make(GoPlugLoader.PluginInfoMap),
			Initialized: true,
			// PluginImpl:  &GoPlugLoader.RpcPlugin{},
			PluginImpl: impl,
			Loaders:    GoPlugLoader.NewLoaders(&base, &file, config, &l),
			Validator:  Plugin.NewBaseValidatorChain(&Plugin.IdentityValidator{}),
			Logger:     &l,
			Error:      err,
			// validator: Plugin.NewBaseValidatorChain(&Plugin.JSONFileValidator{}, &Plugin.IdentityValidator{}, &Plugin.LocalSourceValidator{}),
		}

		err = manager.SetPluginTypes(config.PluginTypes)
		if err.IsError() {
			break
		}
	}

	return manager, err
}

func (m *PluginManager) IsError() bool {
	return m.Error.IsError()
}

func (m *PluginManager) GetError() Return.Error {
	return m.Error
}

func (m *PluginManager) PrintError() {
	m.Error.Print()
}

func (m *PluginManager) SetPluginTypes(pluginTypes Plugin.Types) Return.Error {
	return m.Loaders.SetPluginTypes(pluginTypes)
}

func (m *PluginManager) SetPrefix(prefix string) Return.Error {
	return m.Loaders.SetPrefix(prefix)
}

func (m *PluginManager) SetDir(dir string) Return.Error {
	for range Only.Once {
		m.PluginDir, m.Error = utils.NewDir(dir)
		if m.Error.IsError() {
			break
		}

		m.Error = m.Loaders.SetDir(dir)
		if m.Error.IsError() {
			break
		}

		log.Printf("[INFO]: Plugin BaseDir is '%s'", m.PluginDir.GetDir())
	}

	return m.Error
}

// GetDir implements the interface method
func (m *PluginManager) GetDir() string {
	return m.PluginDir.GetDir()
}

func (m *PluginManager) SetLogfile(dir string) Return.Error {
	for range Only.Once {
		var fp utils.FilePath
		fp, m.Error = utils.NewFile(dir)
		if m.Error.IsError() {
			break
		}
		m.Logfile = &fp

		m.Error = m.Loaders.SetLogfile(*m.Logfile)
		if m.Error.IsError() {
			break
		}

		m.Error = m.Logger.SetLogFile(m.Logfile.GetPath())
		if m.Error.IsError() {
			break
		}

		log.Printf("[INFO]: Plugin BaseDir is '%s'", m.PluginDir.GetDir())
	}

	return m.Error
}

func (m *PluginManager) SetFileGlob(glob string) Return.Error {
	m.FileGlob = glob
	return Return.Ok
}

func (m *PluginManager) SetIdentity(config Plugin.Identity) Return.Error {
	for range Only.Once {
		m.Error = config.IsValid()
		if m.Error.IsError() {
			break
		}

		m.Config = &config
		// bm.Plugin = NewRpcPlugin(nil)
		// bm.Error = bm.rpc.SetDir(bm.GetBaseDir())
		// if bm.Error.IsError() {
		// 	break
		// }
		//
		// m.Config = &config
		// bm.Plugin, bm.Error = NewRpcPlugin(&config, nil, bm.logger)
		// bm.Plugin.SetConfigDir(bm.GetBaseDir())
	}

	return m.Error
}

func (m *PluginManager) SetImplementor(impl plugin.Plugin) Return.Error {
	m.PluginImpl = impl
	return Return.Ok
}

// pluginMap should be used by clients for the map of plugins.
func (m *PluginManager) pluginMap(id string) map[string]plugin.Plugin {
	pmap := map[string]plugin.Plugin{}

	// for _, pinfo := range m.Plugins {
	// 	pmap[pinfo.ID] = m.pluginImpl
	// }

	pmap[id] = m.PluginImpl

	return pmap
}

func (m *PluginManager) NameToPluginPath(id string) (*utils.FilePath, Return.Error) {
	return m.Loaders.NameToPluginPath(id)
}

// ---------------------------------------------------------------------------------------------------- //

func FileToId(file string, glob string) string {
	globAsterix := strings.LastIndex(glob, "*")
	trim := glob[0:globAsterix]
	return strings.TrimPrefix(file, trim)
}
