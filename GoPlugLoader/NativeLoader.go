package GoPlugLoader

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

var (
	NativePluginExtensions = []string{".so"}
)

//
// NewNativeLoader - Create a new LoaderInterface interface instance of this structure.
// ---------------------------------------------------------------------------------------------------- //
func NewNativeLoader(dir *utils.FilePath, id *Plugin.Identity, logger *utils.Logger) LoaderInterface {
	var err Return.Error
	err.SetPrefix("NativeLoader: ")
	return &NativeLoader{
		baseDir: dir,
		Files:   nil,
		logger:  logger,
		store:   NewPluginStore(),
		Error:   err,
	}
}

//
// NativeLoader is a default implementation of LoaderInterface interface
// ---------------------------------------------------------------------------------------------------- //
type NativeLoader ChildLoader

func (l *NativeLoader) SetLogfile(path utils.FilePath) Return.Error {
	l.logfile = &path
	return Return.Ok
}

// SetPluginTypes - Ignored
func (l *NativeLoader) SetPluginTypes(pluginTypes Plugin.Types) Return.Error {
	return Return.Ok
}
func (l *NativeLoader) GetLoader(force string) LoaderInterface {
	if (force == NativeLoaderName) || (force == "") {
		return l
	}
	return nil
}
func (l *NativeLoader) GetLoaderType() string {
	return NativeLoaderName
}
func (l *NativeLoader) IsLoaderType(loaderType string) bool {
	if loaderType == NativeLoaderName {
		return true
	}
	return false
}

func (l *NativeLoader) SetPrefix(prefix string) Return.Error {
	l.prefix = prefix
	return l.Error
}

// SetDir - sets the plugin base dir
func (l *NativeLoader) SetDir(dir string) Return.Error {
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
func (l *NativeLoader) GetDir() string {
	return l.baseDir.GetPath()
}

func (l *NativeLoader) GetFiles() utils.FilePaths {
	return l.Files
}

func (l *NativeLoader) NameToPluginPath(id string) (*utils.FilePath, Return.Error) {
	var pluginPath *utils.FilePath

	for range Only.Once {
		var item *PluginItem
		item, l.Error = l.store.StoreGet(id)
		if l.Error.IsError() {
			break
		}

		pluginPath = &item.Data.Common.Filename
		l.Error = Return.Ok
	}

	return pluginPath, l.Error
}

//
// ---------------------------------------- //
// Plugin methods

func (l *NativeLoader) PluginScan(glob string) Return.Error {
	for range Only.Once {
		l.Files, l.Error = l.baseDir.Scan(glob)
		if l.Error.IsError() {
			break
		}
		l.Files.KeepExtensions(NativePluginExtensions...)
	}
	return l.Error
}
func (l *NativeLoader) PluginScanByExtension(ext ...string) Return.Error {
	l.Files, l.Error = l.baseDir.ScanForExtension(ext...)
	return l.Error
}

func (l *NativeLoader) PluginRegisterAll() (PluginItems, Return.Error) {
	var items PluginItems
	for range Only.Once {
		for _, pDir := range l.Files {
			log.Printf("[INFO]: %d plugin files found in %s", pDir.Length(), pDir.Dir.String())
			for _, path := range pDir.Get() {
				var item PluginItem
				item, l.Error = l.PluginRegister(path)
				if l.Error.IsError() {
					break
				}
				items = append(items, &item)
			}
		}
	}
	return items, l.Error
}
func (l *NativeLoader) PluginUnregisterAll() Return.Error {
	for range Only.Once {
		for _, pDir := range l.Files {
			for _, path := range pDir.Get() {
				l.Error = l.PluginUnregister(path)
				if l.Error.IsError() {
					break
				}
			}
		}
	}
	return l.Error
}

func (l *NativeLoader) PluginRegister(path utils.FilePath) (PluginItem, Return.Error) {
	var item PluginItem
	for range Only.Once {
		item, l.Error = l.PluginLoad(path)
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
func (l *NativeLoader) PluginUnregister(path utils.FilePath) Return.Error {
	return l.PluginUnload(path)
}

func (l *NativeLoader) PluginLoad(pluginPath utils.FilePath) (PluginItem, Return.Error) {
	var item PluginItem

	for range Only.Once {
		l.Error.ReturnClear()
		l.Error.SetPrefix("")

		l.Error = pluginPath.FileExists()
		if l.Error.IsError() {
			break
		}

		// ---------------------------------------------------------------------------------------------------- //
		// Initial setup, before pulling in configured data.
		// var dir utils.FilePath
		// dir, l.Error = utils.NewDir(pluginPath.GetDir())
		// if l.Error.IsError() {
		// 	break
		// }

		id := strings.TrimPrefix(pluginPath.GetName(), l.prefix)
		item.Native = NewNativePlugin()
		item.Native.Plugin.Common = Plugin.Common{
			Id: id,
			// PluginTypes:  Plugin.NativePluginType,
			// StructName:   "unknown",
			// Directory:    dir,
			// Filename:     pluginPath,
			// Logger:       &plog,
			// Configured:   true,
			// RawInterface: nil,
			// Error:        Return.New(),
		}

		// ---------------------------------------------------------------------------------------------------- //
		// Load the plugin and pull in configured data.
		l.Error = item.Native.Service.Open(pluginPath)
		if l.Error.IsError() {
			break
		}

		var identity *Plugin.Identity
		identity, l.Error = item.Native.Service.GetIdentity()
		if l.Error.IsError() {
			l.Error.SetError("GoPluginIdentity is not globally defined: %s", l.Error)
			break
		}
		item.Native.SetIdentity(identity) // This will be replaced with a full get of GoPluginNativeInterface

		// ---------------------------------------------------------------------------------------------------- //
		// Reconfigure and check data.
		log.Println("Looking for GoPluginNativeInterface")
		var GoPluginNativeInterface *NativePluginInterface
		GoPluginNativeInterface, l.Error = item.Native.Service.GetNativePluginInterface()
		if l.Error.IsError() {
			l.Error.SetWarning("GoPluginNativeInterface is not globally defined: %s", l.Error)
			// We can continue with a minimally defined plugin, (with just GoPluginIdentity defined).
			// break
		}
		if GoPluginNativeInterface != nil {
			native := (*GoPluginNativeInterface).RefPlugin()
			if native == nil {
				l.Error.SetError("GoPluginNativeInterface is defined, but nil!")
				break
			}

			item.Native.Plugin = *native
			// (*GoPluginNativeInterface).SetNativeService(identity.Name, *item.Native.Service.Object)
			// log.Println((*GoPluginNativeInterface).String())
		}

		// var plog utils.Logger
		// plog, l.Error = utils.NewLogger(pluginPath.GetName(), "") // @TODO - Config for log file.
		// if l.Error.IsError() {
		// 	break
		// }
		// plog.SetLevel(l.logger.GetLevel())
		// item.Native.SetLogger(&plog)

		item.Native.SetFilename(pluginPath)
		item.Native.SetHookPlugin(&item.Native.Plugin)
		item.Native.SetPluginTypeNative() // Even if the config doesn't set it, do it here.
		// item.Native.SetPluginIdentity(item.Native.Common.Id)
		item.Native.SetNativeService(item.Native.Common.Id, *item.Native.Service.Object)
		item.Native.SetRawInterface(item.Native.Service.Symbol)
		item.Native.SetStructName(*identity)
		item.Data = &item.Native.Plugin
		log.Printf("[%s]: Name:%s Path: %s\n",
			item.Native.Common.Id, item.Native.Common.Filename.GetName(), item.Native.Common.Filename.GetPath())

		// ---------------------------------------------------------------------------------------------------- //
		// Reconfigure and check data.
		// item.Native.Common.Id = item.Native.Dynamic.Identity.Name

		// // ---------------------------------------------------------------------------------------------------- //
		// // Reconfigure and check data.
		// log.Println("Looking for GoPluginRpcInterface")
		// var GoPluginRpcInterface *RpcPluginInterface
		// GoPluginRpcInterface, l.Error = item.Native.Service.GetRpcPluginInterface()
		// if !l.Error.IsError() {
		// 	// (*GoPluginRpcInterface).SetRpcService(identity.Name, *item.Rpc.Service.ClientProtocol)
		// 	// log.Println((*GoPluginRpcInterface).String())
		// 	rpc := (*GoPluginRpcInterface).RefPlugin()
		// 	if rpc != nil {
		// 		if item.Rpc == nil {
		// 			item.Rpc = NewRpcPlugin()
		// 		}
		// 		item.Rpc.Plugin = *rpc
		// 		item.Rpc.Plugin.Common.PluginTypes.Rpc = true
		// 		item.Rpc.Plugin.Common.Configured = true
		// 		item.Rpc.Plugin.Common.IsPlugin = true
		// 		item.Rpc.Plugin.Common.Directory = dir
		// 		item.Rpc.Plugin.Common.Filename = pluginPath
		// 		item.Rpc.Plugin.Dynamic.Hooks.SetHookPlugin(&item.Rpc.Plugin)
		// 		if item.Rpc.Plugin.Dynamic.Identity.Callbacks.PluginName == "" {
		// 			item.Rpc.Plugin.Dynamic.Identity.Callbacks.PluginName = item.Rpc.Plugin.Dynamic.Identity.Name
		// 		}
		// 		if item.Rpc.Plugin.Dynamic.Hooks.Identity == "" {
		// 			item.Rpc.Plugin.Dynamic.Hooks.Identity = identity.Name
		// 		}
		// 		item.Data = &item.Rpc.Plugin
		// 	}
		// }

		l.Error = item.IsItemValid()
		if l.Error.IsError() {
			break
		}

		// item.Native.Dynamic.Identity.Print()
		l.Error = item.Native.Callback(Plugin.CallbackInitialise, item.Native)
		if l.Error.IsError() {
			break
		}

		// // ---------------------------------------------------------------------------------------------------- //
		// fmt.Println("Looking for GoPluginLoader")
		// sym, l.Error = item.Native.Service.Lookup("GoPluginLoader")
		// if l.Error.IsError() {
		// 	break
		// }
		// fmt.Printf("sym: '%T'\n", sym)
		// init, ok := sym.(func(Plugin.Interface, ...interface{}) Return.Error)
		// if !ok {
		// 	l.Error.SetError("plugin structure not defined properly in file '%s' - type is '%s'",
		// 		pluginPath, utils.GetTypeName(sym))
		// 	break
		// }
		// l.Error = init(item.Native)
	}

	return item, l.Error
}
func (l *NativeLoader) PluginUnload(path utils.FilePath) Return.Error {
	var err Return.Error

	for range Only.Once {
		var plug *PluginItem
		plug, l.Error = l.StoreGet(path.GetPath())
		if l.Error.IsError() {
			break
		}

		_, l.Error = l.store.StoreRemove(plug.Data.Common.Filename.GetPath())
		if l.Error.IsError() {
			l.Error.SetError("[INFO]: Plugin(%s): Unload FAILED", path.String())
			// log.Printf("[INFO]: Plugin(%s): Unload FAILED", path.String())
			break
		}
	}

	return err
}

func (l *NativeLoader) PluginInit(items ...PluginItem) Return.Error {
	for range Only.Once {
		for _, item := range items {
			if !item.IsNativePlugin() {
				// Silently ignore.
				continue
			}

			itemData := item.GetItemData(&Plugin.Types{
				Rpc:    false,
				Native: true,
			})
			if itemData == nil {
				l.Error = item.Error
				break
			}

			itemData.SetValue("slave-init-timestamp", time.Now())

			l.Error = item.InitialiseWithPlugin(itemData)
			if l.Error.IsError() {
				itemData.SetValue("slave-init", l.Error)
				break
			}

			itemData.SetValue("slave-init", "OK")
		}
	}

	return l.Error
}

func (l *NativeLoader) PluginParse(path utils.FilePath) (*Plugin.Identity, Return.Error) {
	l.Error.SetError("Parse() not implemented yet in PluginLoader: %s", path)
	return nil, l.Error
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of PluginStore interface structure

func (l *NativeLoader) StoreIsValid() Return.Error {
	return l.store.StoreIsValid()
}
func (l *NativeLoader) StoreSize() uint {
	return l.store.StoreSize()
}
func (l *NativeLoader) String() string {
	return l.store.String()
}
func (l *NativeLoader) StorePrint() {
	fmt.Println("# Native Plugins")
	l.store.StorePrint()
}
func (l *NativeLoader) StorePut(item *PluginItem, forced bool) Return.Error {
	return l.store.StorePut(item, forced)
}
func (l *NativeLoader) StoreGet(name string) (*PluginItem, Return.Error) {
	return l.store.StoreGet(name)
}
func (l *NativeLoader) StoreGetAll() PluginItems {
	return l.store.StoreGetAll()
}
func (l *NativeLoader) StoreRemove(name string) (*PluginItem, Return.Error) {
	return l.store.StoreRemove(name)
}
