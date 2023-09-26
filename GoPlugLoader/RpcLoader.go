package GoPlugLoader

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/MickMake/GoUnify/Only"
	goplugin "github.com/hashicorp/go-plugin"

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

		pluginPath = &item.Data.Common.Filename
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

func (l *RpcLoader) PluginRegisterAll() (PluginItems, Return.Error) {
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
func (l *RpcLoader) PluginUnregisterAll() Return.Error {
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

func (l *RpcLoader) PluginRegister(path utils.FilePath) (PluginItem, Return.Error) {
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
func (l *RpcLoader) PluginUnregister(path utils.FilePath) Return.Error {
	return l.PluginUnload(path)
}

func (l *RpcLoader) PluginLoad(pluginPath utils.FilePath) (PluginItem, Return.Error) {
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
		var dir utils.FilePath
		dir, l.Error = utils.NewDir(pluginPath.GetDir())
		if l.Error.IsError() {
			break
		}

		var plog utils.Logger
		plog, l.Error = utils.NewLogger(pluginPath.GetName(), "") // @TODO - Config for log file.
		if l.Error.IsError() {
			break
		}
		plog.SetLevel(l.logger.GetLevel())

		id := strings.TrimPrefix(pluginPath.GetName(), l.prefix)
		item.Rpc = NewRpcPlugin()
		item.Rpc.Plugin.Common = Plugin.Common{
			Id:           id,
			PluginTypes:  Plugin.RpcPluginType,
			StructName:   utils.GetStructName(l),
			Directory:    dir,
			Filename:     pluginPath,
			Logger:       &plog,
			Configured:   true,
			RawInterface: nil,
			Error:        Return.New(),
		}

		item.Rpc.Plugin.Dynamic.SetHookPlugin(&item.Rpc.Plugin)
		item.Rpc.Plugin.Common.SetPluginType(Plugin.RpcPluginType)

		item.Rpc.Plugin.Services.SetPluginIdentity(item.Rpc.Plugin.Common.Id)
		item.Rpc.Plugin.Services.SetRpcService(item.Rpc.Plugin.Common.Id, &GoPluginMaster{})
		item.Rpc.Service.ClientConfig = goplugin.ClientConfig{
			HandshakeConfig: Plugin.HandshakeConfig,
			Plugins:         item.Rpc.Plugin.Services.GetAsRpcPluginSet(),
			Cmd:             exec.Command(pluginPath.GetPath()),
		}
		item.Data = &item.Rpc.Plugin

		item.Rpc.Plugin.SetPluginTypeRpc()
		item.Rpc.Plugin.Common.Configured = true
		item.Rpc.Plugin.Common.IsPlugin = true
		item.Rpc.Plugin.Common.Directory = dir
		item.Rpc.Plugin.Common.Filename = pluginPath
		item.Rpc.Plugin.Dynamic.Hooks.SetHookPlugin(&item.Rpc.Plugin)
		if item.Rpc.Plugin.Dynamic.Identity.Callbacks.PluginName == "" {
			item.Rpc.Plugin.Dynamic.Identity.Callbacks.PluginName = item.Rpc.Plugin.Dynamic.Identity.Name
		}
		if item.Rpc.Plugin.Dynamic.Hooks.Identity == "" {
			item.Rpc.Plugin.Dynamic.Hooks.Identity = item.Rpc.Plugin.Dynamic.Identity.Name
		}
		log.Printf("[%s]: Name:%s Path: %s\n",
			item.Rpc.Plugin.Common.Id, item.Rpc.Plugin.Common.Filename.GetName(), item.Rpc.Plugin.Common.Filename.GetPath())

		var e error
		item.Rpc.Service.ClientRef = goplugin.NewClient(&item.Rpc.Service.ClientConfig)
		if item.Rpc.Service.ClientRef == nil {
			l.Error.SetError("[%s]: ERROR: RPC client is nil", item.Rpc.Plugin.Common.Id)
			break
		}

		item.Rpc.Service.ClientProtocol, e = item.Rpc.Service.ClientRef.Client()
		if e != nil {
			l.Error.SetError("[%s]: ERROR: %s", item.Rpc.Common.Id, e.Error())
			break
		}
		defer item.Rpc.Service.ClientProtocol.Close()

		e = item.Rpc.Service.ClientProtocol.Ping()
		if e != nil {
			l.Error.SetError("[%s]: ERROR: %s\n", id, e.Error())
			break
		}

		var raw any
		raw, e = item.Rpc.Service.ClientProtocol.Dispense(item.Rpc.Common.Id)
		if e != nil {
			l.Error.SetError("[%s]: ERROR: %s\n", item.Rpc.Common.Id, e.Error())
			break
		}

		tn := utils.GetTypeName(raw)
		if tn != "*GoPlugLoader.RpcPluginClient" {
			l.Error.SetError("[%s]: ERROR: Invalid type - expecting '*RpcPluginClient', got '%s'", item.Rpc.Common.Id, tn)
			break
		}

		impl := raw.(*RpcPluginClient)
		item.Rpc.Plugin.Dynamic = impl.GetData()
		item.Rpc.Plugin.Dynamic.Identity.Print()

		if item.Rpc == nil {
			item.Rpc = NewRpcPlugin()
		}
		item.Rpc.Plugin.Common.PluginTypes.Rpc = true
		item.Rpc.Plugin.Common.Configured = true
		item.Rpc.Plugin.Common.IsPlugin = true
		item.Rpc.Plugin.Common.Directory = dir
		item.Rpc.Plugin.Common.Filename = pluginPath
		item.Rpc.Plugin.Dynamic.Hooks.SetHookPlugin(&item.Rpc.Plugin)
		if item.Rpc.Plugin.Dynamic.Identity.Callbacks.PluginName == "" {
			item.Rpc.Plugin.Dynamic.Identity.Callbacks.PluginName = item.Rpc.Plugin.Dynamic.Identity.Name
		}
		if item.Rpc.Plugin.Dynamic.Hooks.Identity == "" {
			// item.Rpc.Plugin.Dynamic.Hooks.Identity = identity.Name
		}
		item.Data = &item.Rpc.Plugin

		l.Error = item.IsItemValid()
		if l.Error.IsError() {
			break
		}

		l.Error = item.Rpc.Callback(Plugin.CallbackInitialise, item.Rpc)
		if l.Error.IsError() {
			break
		}

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

		item := plug.GetItemData(&Plugin.Types{
			Rpc:    true,
			Native: false,
		})
		if item == nil {
			l.Error = plug.Error
			break
		}

		_, l.Error = l.store.StoreRemove(item.Common.Filename.GetPath())
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

			itemData := item.GetItemData(&Plugin.Types{
				Rpc:    true,
				Native: false,
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
