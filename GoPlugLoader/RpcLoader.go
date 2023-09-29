package GoPlugLoader

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

//
// NewRpcLoader - Create a new LoaderInterface interface instance of this structure.
// ---------------------------------------------------------------------------------------------------- //
func NewRpcLoader(dir *utils.FilePath, file *utils.FilePath, cfg *Plugin.Identity, logger *utils.Logger) LoaderInterface {
	var err Return.Error
	err.SetPrefix("RpcLoader: ")

	if dir == nil {
		dir = &utils.FilePath{}
	}
	if file == nil {
		file = &utils.FilePath{}
	}

	return &RpcLoader{
		baseDir: dir,
		Files:   nil,
		logger:  logger,
		store:   NewPluginStore(),
		Error:   err,
	}
}

//
// RpcLoader is a default implementation of LoaderInterface interface
// ---------------------------------------------------------------------------------------------------- //
type RpcLoader ChildLoader

func (l *RpcLoader) SetLogfile(path utils.FilePath) Return.Error {
	l.logfile = &path
	return Return.Ok
}

// SetPluginTypes - Ignored
func (l *RpcLoader) SetPluginTypes(pluginTypes Plugin.Types) Return.Error {
	return Return.Ok
}
func (l *RpcLoader) GetLoader(force string) LoaderInterface {
	if (force == RpcLoaderName) || (force == "") {
		return l
	}
	return nil
}
func (l *RpcLoader) GetLoaderType() string {
	return RpcLoaderName
}
func (l *RpcLoader) IsLoaderType(loaderType string) bool {
	if loaderType == RpcLoaderName {
		return true
	}
	return false
}

func (l *RpcLoader) SetPrefix(prefix string) Return.Error {
	l.prefix = prefix
	return l.Error
}

// SetDir - sets the plugin base dir
func (l *RpcLoader) SetDir(dir string) Return.Error {
	for range Only.Once {
		if dir == "" {
			var e error
			dir, e = os.Getwd()
			l.Error.SetError(e)
			if l.Error.IsError() {
				break
			}
		}

		l.Error = l.baseDir.SetDir(dir)
	}

	return l.Error
}

// GetDir - Gets the plugin base dir
func (l *RpcLoader) GetDir() string {
	return l.baseDir.GetPath()
}

func (l *RpcLoader) GetFiles() utils.FilePaths {
	return l.Files
}

func (l *RpcLoader) NameToPluginPath(id string) (*utils.FilePath, Return.Error) {
	var pluginPath *utils.FilePath

	for range Only.Once {
		var item *PluginItem
		item, l.Error = l.store.StoreGet(id)
		if l.Error.IsError() {
			break
		}

		pluginPath = item.Pluggable.GetPluginPath()
		l.Error = Return.Ok
	}

	return pluginPath, l.Error
}

//
// ---------------------------------------- //
// Plugin methods

func (l *RpcLoader) PluginScan(glob string) Return.Error {
	for range Only.Once {
		l.Files, l.Error = l.baseDir.Scan(glob)
		if l.Error.IsError() {
			break
		}
		l.Files.RemoveExtensions(NativePluginExtensions...)
	}
	return l.Error
}
func (l *RpcLoader) PluginScanByExtension(ext ...string) Return.Error {
	l.Files, l.Error = l.baseDir.ScanForExtension(ext...)
	return l.Error
}

func (l *RpcLoader) PluginRegister() (PluginItems, Return.Error) {
	var items PluginItems
	for range Only.Once {
		for _, pDir := range l.Files {
			log.Printf("[INFO]: %d plugin files found in %s", pDir.Length(), pDir.Dir.String())
			for _, path := range pDir.Get() {
				var item PluginItem
				item, l.Error = l.PluginLoad(path)
				if l.Error.IsError() {
					break
				}
				items = append(items, &item)
			}
		}
	}
	return items, l.Error
}
func (l *RpcLoader) PluginUnregister() Return.Error {
	for range Only.Once {
		for _, pDir := range l.Files {
			for _, path := range pDir.Get() {
				l.Error = l.PluginUnload(path)
				if l.Error.IsError() {
					break
				}
			}
		}
	}
	return l.Error
}

func (l *RpcLoader) PluginLoad(pluginPath utils.FilePath) (PluginItem, Return.Error) {
	var item PluginItem

	for range Only.Once {
		id := strings.TrimPrefix(pluginPath.GetName(), l.prefix)
		item.Pluggable = NewRpcPlugin()

		l.Error = item.Pluggable.PluginLoad(id, pluginPath)
		if l.Error.IsError() {
			break
		}

		l.Error = l.PluginInit(item)
		if l.Error.IsError() {
			break
		}

		l.Error = l.StorePut(&item, true)
	}

	return item, l.Error
}
func (l *RpcLoader) PluginUnload(path utils.FilePath) Return.Error {
	var err Return.Error

	for range Only.Once {
		var plug *PluginItem
		plug, l.Error = l.StoreGet(path.GetPath())
		if l.Error.IsError() {
			break
		}

		item := plug.GetItemData()
		if item == nil {
			l.Error = plug.Error
			break
		}

		_, l.Error = l.store.StoreRemove(path.GetPath())
		if l.Error.IsError() {
			l.Error.SetError("[INFO]: Plugin(%s): Unload FAILED", path.String())
			log.Printf("[INFO]: Plugin(%s): Unload FAILED", path.String())
			break
		}
	}

	return err
}

func (l *RpcLoader) PluginInit(items ...PluginItem) Return.Error {
	for range Only.Once {
		for _, item := range items {
			if !item.IsRpcPlugin() {
				// Silently ignore.
				continue
			}

			itemData := item.GetItemData()
			if itemData == nil {
				l.Error = item.Error
				break
			}

			itemData.SetValue("slave-init-timestamp", time.Now())

			l.Error = item.Initialise()
			if l.Error.IsError() {
				itemData.SetValue("slave-init", l.Error)
				break
			}

			itemData.SetValue("slave-init", "OK")
		}
	}

	return l.Error
}
func (l *RpcLoader) PluginParse(path utils.FilePath) (*Plugin.Identity, Return.Error) {
	l.Error.SetError("Parse() not implemented yet in PluginLoader: %s", path)
	return nil, l.Error
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of PluginStore interface structure

func (l *RpcLoader) StoreIsValid() Return.Error {
	return l.store.StoreIsValid()
}
func (l *RpcLoader) StoreSize() uint {
	return l.store.StoreSize()
}
func (l *RpcLoader) String() string {
	return l.store.String()
}
func (l *RpcLoader) StorePrint() {
	log.Println("# RPC Plugins")
	l.store.StorePrint()
}
func (l *RpcLoader) StorePut(item *PluginItem, forced bool) Return.Error {
	return l.store.StorePut(item, forced)
}
func (l *RpcLoader) StoreGet(name string) (*PluginItem, Return.Error) {
	return l.store.StoreGet(name)
}
func (l *RpcLoader) StoreGetAll() PluginItems {
	return l.store.StoreGetAll()
}
func (l *RpcLoader) StoreRemove(name string) (*PluginItem, Return.Error) {
	return l.store.StoreRemove(name)
}
