package GoPlugLoader

import (
	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

const (
	MainLoaderName   = "main"
	RpcLoaderName    = "rpc"
	NativeLoaderName = "native"
)

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
// Loader - main implementation of LoaderInterface struct, which pulls in a child LoaderInterface
// ---------------------------------------------------------------------------------------------------- //
type Loader struct {
	Native      LoaderInterface
	Rpc         LoaderInterface
	PluginTypes Plugin.Types
	Error       Return.Error
}

func (l *Loader) SetLogfile(path utils.FilePath) Return.Error {
	for range Only.Once {
		if l.PluginTypes.Native {
			return l.Native.SetLogfile(path)
		}

		if l.PluginTypes.Rpc {
			return l.Rpc.SetLogfile(path)
		}
	}

	return Return.Ok
}

func (l *Loader) SetPluginTypes(pluginTypes Plugin.Types) Return.Error {
	l.PluginTypes = pluginTypes
	return Return.Ok
}

func (l *Loader) GetLoader(force string) LoaderInterface {
	if force == NativeLoaderName {
		return l.Native.GetLoader(NativeLoaderName)
	}

	if force == RpcLoaderName {
		return l.Native.GetLoader(RpcLoaderName)
	}

	ret := l.Native.GetLoader("")
	if ret != nil {
		return ret
	}

	ret = l.Rpc.GetLoader("")
	if ret != nil {
		return ret
	}

	return nil
}

func (l *Loader) GetLoaderType() string {
	return MainLoaderName
}

func (l *Loader) IsLoaderType(loaderType string) bool {
	if loaderType == MainLoaderName {
		return true
	}
	return false
}

func (l *Loader) SetPrefix(prefix string) Return.Error {
	for range Only.Once {
		l.Error = l.Native.SetDir(prefix)
		if l.Error.IsError() {
			break
		}

		l.Error = l.Rpc.SetDir(prefix)
		if l.Error.IsError() {
			break
		}
	}

	return l.Error
}

func (l *Loader) SetDir(dir string) Return.Error {
	for range Only.Once {
		l.Error = l.Native.SetDir(dir)
		if l.Error.IsError() {
			break
		}

		l.Error = l.Rpc.SetDir(dir)
		if l.Error.IsError() {
			break
		}
	}

	return l.Error
}

func (l *Loader) GetDir() string {
	var ret string

	for range Only.Once {
		ret = l.Native.GetDir()
		if ret != "" {
			break
		}

		ret = l.Rpc.GetDir()
	}

	return ret
}

func (l *Loader) GetFiles() utils.FilePaths {
	if l.PluginTypes.Native {
		return l.Native.GetFiles()
	}

	if l.PluginTypes.Rpc {
		return l.Rpc.GetFiles()
	}

	return utils.FilePaths{}
}

func (l *Loader) NativeFiles() utils.FilePaths {
	return l.Native.GetFiles()
}

func (l *Loader) RpcFiles() utils.FilePaths {
	return l.Rpc.GetFiles()
}

func (l *Loader) NameToPluginPath(id string) (*utils.FilePath, Return.Error) {
	var pluginPath *utils.FilePath

	for range Only.Once {
		if l.PluginTypes.Native {
			pluginPath, l.Error = l.Native.NameToPluginPath(id)
			if l.Error.IsError() {
				break
			}
		}

		if l.PluginTypes.Rpc {
			pluginPath, l.Error = l.Rpc.NameToPluginPath(id)
			if l.Error.IsError() {
				break
			}
		}

		l.Error = Return.Ok
	}

	return pluginPath, l.Error
}

//
// ---------------------------------------- //
// Plugin methods

func (l *Loader) PluginScan(glob string) Return.Error {
	for range Only.Once {
		// plugins, l.Error = m.Base.Scan(m.fileGlob)
		// plugins, e := plugin.Discover(m.fileGlob, m.Base.GetPath())
		if l.PluginTypes.Native {
			l.Error = l.Native.PluginScan(glob)
			if l.Error.IsError() {
				break
			}
		}

		if l.PluginTypes.Rpc {
			l.Error = l.Rpc.PluginScan(glob)
			if l.Error.IsError() {
				break
			}
		}

		l.Error = Return.Ok
	}

	return l.Error
}
func (l *Loader) PluginScanByExtension(ext ...string) Return.Error {
	for range Only.Once {
		if l.PluginTypes.Native {
			l.Error = l.Native.PluginScanByExtension(ext...)
			if l.Error.IsError() {
				break
			}
		}

		if l.PluginTypes.Rpc {
			l.Error = l.Rpc.PluginScanByExtension(ext...)
			if l.Error.IsError() {
				break
			}
		}

		l.Error = Return.Ok
	}

	return l.Error
}

func (l *Loader) PluginLoad(path utils.FilePath) (PluginItem, Return.Error) {
	var item PluginItem

	for range Only.Once {
		if l.PluginTypes.Native {
			item, l.Error = l.Native.PluginLoad(path)
			if l.Error.IsError() {
				break
			}
			break
		}

		if l.PluginTypes.Rpc {
			item, l.Error = l.Rpc.PluginLoad(path)
			if l.Error.IsError() {
				break
			}
			break
		}
	}

	return item, l.Error
}
func (l *Loader) PluginUnload(path utils.FilePath) Return.Error {
	for range Only.Once {
		// var pt PluginTypes
		// switch {
		// case path.HasExtension(NativePluginExtensions...):
		// 	pt.Native = true
		// default:
		// 	pt.Rpc = true
		// }

		if l.PluginTypes.Native {
			l.Error = l.Native.PluginUnload(path)
			if l.Error.IsError() {
				break
			}
			break
		}

		if l.PluginTypes.Rpc {
			l.Error = l.Rpc.PluginUnload(path)
			if l.Error.IsError() {
				break
			}
			break
		}
	}

	return l.Error
}

func (l *Loader) PluginRegisterAll() (PluginItems, Return.Error) {
	var items PluginItems

	for range Only.Once {
		if l.PluginTypes.Native {
			var i PluginItems
			i, l.Error = l.Native.PluginRegisterAll()
			if l.Error.IsError() {
				break
			}
			items = append(items, i...)
		}
		if l.PluginTypes.Rpc {
			var i PluginItems
			i, l.Error = l.Rpc.PluginRegisterAll()
			if l.Error.IsError() {
				break
			}
			items = append(items, i...)
		}
	}

	return items, l.Error
}
func (l *Loader) PluginRegister(path utils.FilePath) (PluginItem, Return.Error) {
	var item PluginItem

	for range Only.Once {
		if l.PluginTypes.Native {
			item, l.Error = l.Native.PluginRegister(path)
			if l.Error.IsError() {
				break
			}
		}
		if l.PluginTypes.Rpc {
			item, l.Error = l.Rpc.PluginRegister(path)
			if l.Error.IsError() {
				break
			}
		}
	}

	return item, l.Error
}

func (l *Loader) PluginUnregisterAll() Return.Error {
	for range Only.Once {
		if l.PluginTypes.Native {
			l.Error = l.Native.PluginUnregisterAll()
			if l.Error.IsError() {
				break
			}
		}
		if l.PluginTypes.Rpc {
			l.Error = l.Rpc.PluginUnregisterAll()
			if l.Error.IsError() {
				break
			}
		}
	}

	return l.Error
}
func (l *Loader) PluginUnregister(path utils.FilePath) Return.Error {
	for range Only.Once {
		if l.PluginTypes.Native {
			l.Error = l.Native.PluginUnregister(path)
			if l.Error.IsError() {
				break
			}
		}
		if l.PluginTypes.Rpc {
			l.Error = l.Rpc.PluginUnregister(path)
			if l.Error.IsError() {
				break
			}
		}
	}

	return l.Error
}

func (l *Loader) PluginInit(item ...PluginItem) Return.Error {
	for range Only.Once {
		if l.PluginTypes.Native {
			l.Error = l.Native.PluginInit(item...)
			if l.Error.IsError() {
				break
			}
		}

		if l.PluginTypes.Rpc {
			l.Error = l.Rpc.PluginInit(item...)
			if l.Error.IsError() {
				break
			}
		}
	}

	return l.Error
}
func (l *Loader) PluginParse(path utils.FilePath) (*Plugin.Identity, Return.Error) {
	var identity *Plugin.Identity

	for range Only.Once {
		if l.PluginTypes.Native {
			identity, l.Error = l.Native.PluginParse(path)
			if l.Error.IsError() {
				break
			}
		}

		if l.PluginTypes.Rpc {
			identity, l.Error = l.Rpc.PluginParse(path)
			if l.Error.IsError() {
				break
			}
		}
	}

	return identity, l.Error
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of PluginStore interface structure

func (l *Loader) StoreIsValid() Return.Error {
	l.Error = l.Native.StoreIsValid()
	if l.Error.IsError() {
		return l.Error
	}
	return l.Rpc.StoreIsValid()
}
func (l *Loader) StoreSize() uint {
	count := l.Native.StoreSize()
	return count + l.Rpc.StoreSize()
}
func (l *Loader) String() string {
	ret := l.Native.String()
	return ret + l.Rpc.String()
}
func (l *Loader) StorePrint() {
	if l.PluginTypes.Native {
		l.Native.StorePrint()
	}
	if l.PluginTypes.Rpc {
		l.Rpc.StorePrint()
	}
}
func (l *Loader) StorePut(item *PluginItem, forced bool) Return.Error {
	for range Only.Once {
		if l.PluginTypes.Native && item.IsNativePlugin() {
			l.Error = l.Native.StorePut(item, forced)
			if l.Error.IsError() {
				break
			}
		}
		if l.PluginTypes.Rpc && item.IsRpcPlugin() {
			l.Error = l.Rpc.StorePut(item, forced)
			if l.Error.IsError() {
				break
			}
		}
	}
	return l.Error
}
func (l *Loader) StoreGet(name string) (*PluginItem, Return.Error) {
	var item *PluginItem
	for range Only.Once {
		if l.PluginTypes.Native {
			item, l.Error = l.Native.StoreGet(name)
			if !l.Error.IsError() {
				break
			}
		}
		if l.PluginTypes.Rpc {
			item, l.Error = l.Rpc.StoreGet(name)
			if !l.Error.IsError() {
				break
			}
		}
	}
	return item, l.Error
}
func (l *Loader) StoreGetAll() PluginItems {
	var items PluginItems
	if l.PluginTypes.Native {
		items = append(items, l.Native.StoreGetAll()...)
	}
	if l.PluginTypes.Rpc {
		items = append(items, l.Rpc.StoreGetAll()...)
	}
	return items
}
func (l *Loader) StoreRemove(name string) (*PluginItem, Return.Error) {
	var item *PluginItem
	for range Only.Once {
		if l.PluginTypes.Native {
			item, l.Error = l.Native.StoreRemove(name)
			if !l.Error.IsError() {
				break
			}
		}
		if l.PluginTypes.Rpc {
			item, l.Error = l.Rpc.StoreRemove(name)
			if !l.Error.IsError() {
				break
			}
		}
	}
	return item, l.Error
}
