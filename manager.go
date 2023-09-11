package GoPlug

import (
	"log"

	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoPlug/pluggable"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoUnify/Only"
)

// DefaultManager is the default plugin manager
//goland:noinspection GoUnusedGlobalVariable
var DefaultManager = NewPluginManager()

// ---------------------------------------------------------------------------------------------------- //

// Manager - Used to load, organize and maintain the plugins.
type Manager interface {
	// SetConfig - Set the base dir where to load plugins.
	// If the dir does not exist, or it's not a dir an error will be returned.
	SetConfig(config pluggable.PluginManagerIdentity) Return.Error

	// SetBaseDir - Set the base dir where to load plugins.
	// If the dir does not exist, or it's not a dir an error will be returned.
	SetBaseDir(dir string) Return.Error

	// LoadPlugins - Load all the plugins from the base plugin dir.
	// Any issues happened, an error will be returned.
	LoadPlugins() Return.Error

	// ListPlugins - .
	ListPlugins()

	// LoadPlugin - Load plugin with the specified name.
	// If failed to load, an error will be returned.
	LoadPlugin(name string) Return.Error

	// UnloadPlugin - Unload plugin with the specified name.
	// If failed to unload, an error will be returned.
	UnloadPlugin(name string) Return.Error

	// GetPlugin - Get the plugin with the specified name.
	// If plugin is not existing, an error will be returned.
	GetPlugin(name string) (*pluggable.PluginItem, Return.Error)

	// CheckPlugin - Get the plugin with the specified name.
	// If plugin is not existing, an error will be returned.
	CheckPlugin(name string) (*pluggable.PluginItem, Return.Error)

	// BuildPlugins - Build plugins.
	BuildPlugins() Return.Error
}

// ---------------------------------------------------------------------------------------------------- //

// PluginManager is implemented as default plugin manager
type PluginManager struct {
	// PluginConfig - config
	Config *pluggable.PluginManagerIdentity

	// Pluggable -
	Plugin pluggable.Plugin

	// For keeping values
	// valueStore    map[string]interface{}

	// The plugin loader
	loader Loader

	// The plugin validator
	validator Validator

	// The list to keep the loaded one
	store Store

	Error Return.Error
}

// NewPluginManager is constructor of PluginManager
func NewPluginManager() Manager {
	var err Return.Error
	err.SetPrefix("PluginManager: ")
	return &PluginManager{
		// Config: nil,
		// Plugin: pluggable.NewPlugin(nil),
		// valueStore: make(map[string]interface{}),
		loader: NewPluginLoader(),
		// validator: NewBaseValidatorChain(&JSONFileValidator{}, &IdentityValidator{}, &LocalSourceValidator{}),
		validator: NewBaseValidatorChain(&IdentityValidator{}),
		store:     NewBaseStore(),
		Error:     err,
	}
}

func (bm *PluginManager) SetConfig(config pluggable.PluginManagerIdentity) Return.Error {
	for range Only.Once {
		bm.Error = config.IsValid()
		if bm.Error.IsError() {
			break
		}

		bm.Config = &config
		bm.Plugin = pluggable.NewPlugin(nil)
		bm.Plugin.SetConfigDir(bm.GetBaseDir())
	}

	return bm.Error
}

// SetBaseDir implements the interface method
func (bm *PluginManager) SetBaseDir(dir string) Return.Error {
	for range Only.Once {
		bm.Error = bm.loader.SetDir(dir)
		if bm.Error.IsError() {
			break
		}

		bm.Error = bm.Plugin.SetConfigDir(bm.loader.GetDir())
		if bm.Error.IsError() {
			break
		}

		log.Printf("[INFO]: Plugin BaseDir is '%s'", bm.loader.GetDir())
	}

	return bm.Error
}

// GetBaseDir implements the interface method
func (bm *PluginManager) GetBaseDir() string {
	return bm.loader.GetDir()
}

// ListPlugins implements the interface method
func (bm *PluginManager) ListPlugins() {
	bm.store.Print()
}

// GetPlugin implements the interface method
func (bm *PluginManager) GetPlugin(p string) (*pluggable.PluginItem, Return.Error) {
	return bm.store.Get(p)
}

// InitPlugin implements the interface method
func (bm *PluginManager) stringToPluginPath(p string) utils.PluginPath {
	path, err := utils.NewFile(p)
	bm.Error = err
	return path
}
