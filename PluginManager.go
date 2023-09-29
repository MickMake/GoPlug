package GoPlug

import (
	"log"
	"os"
	"strings"

	"github.com/MickMake/GoUnify/Only"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

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
	SetPrefix(prefix string) Return.Error

	// SetIdentity - Set the base dir where to load plugins.
	// If the dir does not exist, or it's not a dir an error will be returned.
	SetIdentity(config Plugin.Identity) Return.Error

	SetImplementor(impl goplugin.Plugin) Return.Error

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

	GetInterface(id string) (GoPlugLoader.PluginItemInterface, Return.Error)

	IsError() bool
	GetError() Return.Error
	PrintError()

	NameToPluginPath(id string) (*utils.FilePath, Return.Error)

	SetLogfile(s string) Return.Error
}

//
// PluginManager
// ---------------------------------------------------------------------------------------------------- //
type PluginManager struct {
	Config      *Plugin.Identity             `json:"config"`      //
	PluginDir   utils.FilePath               `json:"plugin_dir"`  //
	CmdFile     utils.FilePath               `json:"cmd_file"`    //
	FileGlob    string                       `json:"file_glob"`   // glob match for plugin filenames
	Prefix      string                       `json:"prefix"`      //
	Plugins     GoPlugLoader.PluginInfoMap   `json:"-"`           // Info for found plugins
	Initialized bool                         `json:"initialized"` // has been Initialized
	Loaders     GoPlugLoader.LoaderInterface `json:"-"`           //
	Validator   Plugin.Validator             `json:"-"`           //
	Logger      *utils.Logger                `json:"-"`           //
	Logfile     *utils.FilePath              `json:"logfile"`     //
	Error       Return.Error                 `json:"-"`           //
	pluginImpl  goplugin.Plugin              // Plugin implementation dummy interface
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
			Prefix:      "goplug-",
			Plugins:     make(GoPlugLoader.PluginInfoMap),
			Initialized: true,
			pluginImpl:  impl,
			Loaders:     GoPlugLoader.NewLoaders(&base, &file, config, &l),
			Validator:   Plugin.NewBaseValidatorChain(&Plugin.IdentityValidator{}),
			Logger:      &l,
			Error:       err,
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
	m.Prefix = prefix
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

	prefix := strings.TrimSuffix(glob, "*")
	m.SetPrefix(prefix)
	return Return.Ok
}

func (m *PluginManager) SetIdentity(identity Plugin.Identity) Return.Error {
	for range Only.Once {
		m.Error = identity.IsValid()
		if m.Error.IsError() {
			break
		}

		m.Config = &identity
	}

	return m.Error
}

func (m *PluginManager) SetImplementor(impl goplugin.Plugin) Return.Error {
	m.pluginImpl = impl
	return Return.Ok
}

func (m *PluginManager) NameToPluginPath(id string) (*utils.FilePath, Return.Error) {
	return m.Loaders.NameToPluginPath(id)
}
