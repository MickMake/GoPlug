package GoPlugLoader

import (
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

//
// LoaderInterface - main implementation of LoaderInterface interface, which pulls in a child LoaderInterface
// ---------------------------------------------------------------------------------------------------- //
type LoaderInterface interface {
	// SetPrefix - sets the prefix of files to include as a plugin
	SetPrefix(prefix string) Return.Error

	// SetDir - sets the plugin base dir
	SetDir(dir string) Return.Error

	// GetDir - get the plugin base dir
	GetDir() string

	SetLogfile(path utils.FilePath) Return.Error

	// PluginInit the plugin
	PluginInit(item ...PluginItem) Return.Error

	// PluginScan the plugin base dir and get the plugin candidates
	PluginScan(glob string) Return.Error

	PluginScanByExtension(ext ...string) Return.Error

	GetFiles() utils.FilePaths
	NameToPluginPath(id string) (*utils.FilePath, Return.Error)

	PluginRegister() (PluginItems, Return.Error)
	PluginUnregister() Return.Error

	PluginLoad(path utils.FilePath) (PluginItem, Return.Error)
	PluginUnload(path utils.FilePath) Return.Error

	// PluginParse the plugin identity config
	PluginParse(path utils.FilePath) (*Plugin.Identity, Return.Error)

	SetPluginTypes(pluginTypes Plugin.Types) Return.Error
	GetLoader(force string) LoaderInterface
	GetLoaderType() string
	IsLoaderType(loaderType string) bool

	PluginStore
}

// NewLoaders - Create a new instance of this structure.
func NewLoaders(dir *utils.FilePath, file *utils.FilePath, cfg *Plugin.Identity, logger *utils.Logger) LoaderInterface {
	var err Return.Error
	err.SetPrefix("Loader: ")

	return &Loader{
		Native:      NewNativeLoader(dir, cfg, logger),
		Rpc:         NewRpcLoader(dir, file, cfg, logger),
		PluginTypes: Plugin.AllPluginTypes,
		Error:       err,
	}
}

//
// ChildLoader - child implementation of LoaderInterface struct, used by child loaders
// ---------------------------------------------------------------------------------------------------- //
type ChildLoader struct {
	baseDir *utils.FilePath
	glob    string
	prefix  string
	Files   utils.FilePaths `json:"files"`
	logger  *utils.Logger
	logfile *utils.FilePath
	store   PluginStore
	Error   Return.Error
}
