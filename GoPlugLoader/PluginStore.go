package GoPlugLoader

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

//
// PluginStore - Stores loaded plugins
// ---------------------------------------------------------------------------------------------------- //
type PluginStore interface {
	StoreIsValid() Return.Error

	// StoreSize - The total count of current items in the store
	StoreSize() uint

	String() string

	StorePrint()

	// StorePut - Append the plugin item to the store
	// If forced is set to be true, new plugin item will overwrite the existing one
	// Try best to append, ignore any errors
	StorePut(item *PluginItem, forced bool) Return.Error

	// StoreGet - the plugin item by name
	// If existing, return the item and set the bool flag to true
	StoreGet(name string) (*PluginItem, Return.Error)

	// 	StoreGetAll() - get all stored plugins
	StoreGetAll() PluginItems

	// StoreRemove - the plugin out of the store and return the removed plugin item
	// If successfully removed, set the bool flag to true
	StoreRemove(name string) (*PluginItem, Return.Error)
}

//
// PluginStoreStruct - the default implementation of PluginStore interface
// ---------------------------------------------------------------------------------------------------- //
type PluginStoreStruct struct {
	// internal lock
	lock *sync.RWMutex

	// internal list
	hash StoreItems

	Error Return.Error
}

//
// StoreItems - holds a nice map of PluginItem
// ---------------------------------------------------------------------------------------------------- //
type StoreItems map[string]*PluginItem

// NewPluginStore - Create a new instance of this structure.
func NewPluginStore() *PluginStoreStruct {
	var err Return.Error
	err.SetPrefix("PluginStoreStruct: ")
	return &PluginStoreStruct{
		lock:  new(sync.RWMutex),
		hash:  make(StoreItems),
		Error: err,
	}
}

func (ps *PluginStoreStruct) StoreIsValid() Return.Error {
	switch {
	case ps == nil:
		ps.Error.SetError("PluginStoreStruct is nil")
	case ps.hash == nil:
		ps.Error.SetError("PluginStoreStruct map is nil")
	}
	return ps.Error
}

// String - Basic Stringer method.
func (ps PluginStoreStruct) String() string {
	var ret string

	for range Only.Once {
		ps.lock.Lock()
		//goland:noinspection GoDeferInLoop
		defer ps.lock.Unlock()

		ret = fmt.Sprintf("# %d plugins loaded\n", (uint)(len(ps.hash)))
		for _, pl := range ps.hash {
			pi := pl.GetItemData(nil)
			if pi == nil {
				continue
			}

			ret += fmt.Sprintf("Name: %s\tVersion: %s\tDescription: %s\n\tFile: %s\n",
				pi.Dynamic.Identity.Name,
				pi.Dynamic.Identity.Version,
				pi.Dynamic.Identity.Description,
				pi.Common.Filename.GetPath(),
			)
		}

		ps.Error = Return.Ok
	}

	return ret
}

// StorePrint - Print out the plugin store.
func (ps *PluginStoreStruct) StorePrint() {
	fmt.Print(ps.String())
}

// StoreSize - Return the size of the store.
func (ps *PluginStoreStruct) StoreSize() uint {
	var size uint

	for range Only.Once {
		ps.lock.Lock()
		//goland:noinspection GoDeferInLoop
		defer ps.lock.Unlock()

		size = (uint)(len(ps.hash))
		ps.Error = Return.Ok
	}

	return size
}

// StorePut - Put a PluginItem into the store, (with optional force).
func (ps *PluginStoreStruct) StorePut(item *PluginItem, forced bool) Return.Error {
	for range Only.Once {
		ps.Error = ps.StoreIsValid()
		if ps.Error.IsError() {
			break
		}

		ps.Error = item.IsItemValid()
		if ps.Error.IsError() {
			break
		}

		_, ok := ps.hash[item.Data.Common.Filename.GetPath()]
		if !ok || (ok && forced) {
			ps.lock.Lock()
			//goland:noinspection GoDeferInLoop
			defer ps.lock.Unlock()

			ps.hash[item.Data.Common.Filename.GetPath()] = item
		}

		ps.Error = Return.Ok
	}

	return ps.Error
}

// StoreGet - Get a PluginItem from the store.
func (ps *PluginStoreStruct) StoreGet(name string) (*PluginItem, Return.Error) {
	var item *PluginItem

	for range Only.Once {
		ps.lock.RLock()
		//goland:noinspection GoDeferInLoop
		defer ps.lock.RUnlock()

		var ok bool
		item, ok = ps.hash[name]
		if ok {
			ps.Error = Return.Ok
			break
		}

		for _, i := range ps.hash {
			pi := i.GetItemData(nil)
			if pi == nil {
				continue
			}

			if pi.Dynamic.Identity.Name != name {
				continue
			}

			item = i
			break
		}

		if item == nil {
			ps.Error.SetError("plugin %s is not loaded", name)
			break
		}
	}

	return item, ps.Error
}

// StoreGetAll - Get all the PluginItem from the store.
func (ps *PluginStoreStruct) StoreGetAll() PluginItems {
	items := make(PluginItems, 0)

	for range Only.Once {
		ps.lock.RLock()
		//goland:noinspection GoDeferInLoop
		defer ps.lock.RUnlock()

		for _, item := range ps.hash {
			items = append(items, item)
		}
	}

	return items
}

// StoreRemove - Remove a PluginItem from the store.
func (ps *PluginStoreStruct) StoreRemove(name string) (*PluginItem, Return.Error) {
	var item *PluginItem

	for range Only.Once {
		ps.lock.Lock()
		//goland:noinspection GoDeferInLoop
		defer ps.lock.Unlock()

		var ok bool
		item, ok = ps.hash[name]
		if !ok {
			ps.Error.SetError("plugin %s is not loaded", name)
			break
		}

		delete(ps.hash, name)
		ps.Error = Return.Ok
	}

	return item, ps.Error
}

//
// PluginInfo
// ---------------------------------------------------------------------------------------------------- //
type PluginInfo struct {
	Types Plugin.Types
	ID    string
	Path  utils.FilePath
	Item  PluginItem
}

//
// PluginInfoMap
// ---------------------------------------------------------------------------------------------------- //
type PluginInfoMap map[string]*PluginInfo

//
// PluginItem
// ---------------------------------------------------------------------------------------------------- //
type PluginItem struct {
	// Stores a pointer to either Native or Rpc plugin loader.
	Data *Plugin.Plugin
	// Native plugin loader.
	Native *NativePlugin
	// Rpc plugin loader.
	Rpc   *RpcPlugin
	Error Return.Error
}

//
// PluginItems
// ---------------------------------------------------------------------------------------------------- //
type PluginItems []*PluginItem

func (i *PluginItem) IsItemValid() Return.Error {
	var err Return.Error
	switch {
	case i == nil:
		err.SetError("plugin is nil")
	case i.Native != nil:
		i.Data = &i.Native.Plugin
		i.Native.Plugin.Common.Configured = true
		break
	case i.Rpc != nil:
		i.Data = &i.Rpc.Plugin
		i.Rpc.Plugin.Common.Configured = true
		break
	default:
		i.Data = nil
		err.SetError("no plugin type defined")
	}
	return err
}

func (i *PluginItem) GetItemData(force *Plugin.Types) *Plugin.Plugin {
	i.Error.Clear()
	if force != nil {
		if force.Native && (i.Native != nil) {
			return &i.Native.Plugin
		}

		if force.Rpc && (i.Rpc != nil) {
			return &i.Rpc.Plugin
		}

		i.Error.SetError("no plugin type defined")
		return nil
	}

	if i.Native != nil {
		return &i.Native.Plugin
	}

	if i.Rpc != nil {
		return &i.Rpc.Plugin
	}

	i.Error.SetError("no plugin type defined")
	return nil
}

func (i *PluginItem) GetItemHooks(force *Plugin.Types) Plugin.HookStore {
	i.Error.Clear()
	if force != nil {
		if force.Native && (i.Native != nil) {
			return &i.Native.Dynamic.Hooks
		}

		if force.Rpc && (i.Rpc != nil) {
			return &i.Rpc.Dynamic.Hooks
		}

		i.Error.SetError("no plugin type defined")
		return nil
	}

	if i.Native != nil {
		return &i.Native.Dynamic.Hooks
	}

	if i.Rpc != nil {
		return &i.Rpc.Dynamic.Hooks
	}

	i.Error.SetError("no plugin type defined")
	return nil
}

func (i *PluginItem) SetItemInterface(ref any) Return.Error {
	for range Only.Once {
		// pi := i.GetItemData(&i.Data.Dynamic.Identity.PluginTypes)
		pi := i.GetItemData(&Plugin.NativePluginType)
		if pi == nil {
			break
		}
		pi.Common.SetRawInterface(ref)

		pi = i.GetItemData(&Plugin.RpcPluginType)
		if pi == nil {
			break
		}
		pi.Common.SetRawInterface(ref)
	}
	return i.Error
}

func (i *PluginItem) IsNativePlugin() bool {
	var yes bool
	for range Only.Once {
		err := i.IsItemValid()
		if err.IsError() {
			// Has to be valid.
			break
		}

		if i.Native == nil {
			break
		}

		if !i.Native.Common.PluginTypes.Native {
			// If pluginTypes.Native is false.
			break
		}

		yes = true
	}
	return yes
}

func (i *PluginItem) IsRpcPlugin() bool {
	var yes bool
	for range Only.Once {
		err := i.IsItemValid()
		if err.IsError() {
			// Has to be valid.
			break
		}

		if i.Rpc == nil {
			break
		}

		if !i.Rpc.Common.PluginTypes.Native {
			// If pluginTypes.Native is false.
			break
		}

		yes = true
	}
	return yes
}

// ---------------------------------------------------------------------------------------------------- //

func (i *PluginItem) Execute(args ...any) Return.Error {
	var err Return.Error
	for range Only.Once {
		// pi := i.GetItemData(&i.Data.Dynamic.Identity.PluginTypes)
		for _, pt := range Plugin.OrderedPluginTypes {
			pd := i.GetItemData(&pt)
			if pd == nil {
				break
			}
			i.Error = i.ExecuteWithPlugin(pd, args...)
			if i.Error.IsError() {
				break
			}
		}
	}
	return err
}

func (i *PluginItem) Initialise(args ...any) Return.Error {
	var err Return.Error
	for range Only.Once {
		// pi := i.GetItemData(&i.Data.Dynamic.Identity.PluginTypes)
		for _, pt := range Plugin.OrderedPluginTypes {
			pd := i.GetItemData(&pt)
			if pd == nil {
				break
			}
			i.Error = i.InitialiseWithPlugin(pd, args...)
			if i.Error.IsError() {
				break
			}
		}
	}
	return err
}

func (i *PluginItem) Run(args ...any) Return.Error {
	var err Return.Error
	for range Only.Once {
		// pi := i.GetItemData(&i.Data.Dynamic.Identity.PluginTypes)
		for _, pt := range Plugin.OrderedPluginTypes {
			pd := i.GetItemData(&pt)
			if pd == nil {
				break
			}
			i.Error = i.RunWithPlugin(pd, args...)
			if i.Error.IsError() {
				break
			}
		}
	}
	return err
}

func (i *PluginItem) Notify(args ...any) Return.Error {
	var err Return.Error
	for range Only.Once {
		// pi := i.GetItemData(&i.Data.Dynamic.Identity.PluginTypes)
		for _, pt := range Plugin.OrderedPluginTypes {
			pd := i.GetItemData(&pt)
			if pd == nil {
				break
			}
			i.Error = i.NotifyWithPlugin(pd, args...)
			if i.Error.IsError() {
				break
			}
		}
	}
	return err
}

func (i *PluginItem) ExecuteWithPlugin(item *Plugin.Plugin, args ...any) Return.Error {
	var err Return.Error
	for range Only.Once {
		prefix := "slave-execute"
		err.SetPrefix("Execute(): ")
		if item == nil {
			err.SetError("Plugin.Data is nil")
			break
		}

		if i == nil {
			err.SetError("Callback is nil")
			break
		}

		item.SetValue(prefix+"-timestamp", time.Now())

		err = i.IsItemValid()
		if err.IsError() {
			item.SetValue(prefix, err)
			break
		}

		var ctx Plugin.Interface
		switch {
		case i.IsNativePlugin():
			ctx = i.Native
		case i.IsRpcPlugin():
			ctx = i.Rpc
		default:
			err.SetError("Neither Native nor RPC plugin configured")
		}
		if err.IsError() {
			break
		}

		err = item.Callback(Plugin.CallbackExecute, ctx, args...)
		if err.IsError() || err.IsWarning() {
			item.SetValue(prefix, err)
			break
		}

		item.SetValue(prefix, "OK")
	}
	return err
}

func (i *PluginItem) InitialiseWithPlugin(item *Plugin.Plugin, args ...any) Return.Error {
	var err Return.Error
	for range Only.Once {
		prefix := "slave-initialise"
		err.SetPrefix("Initialise(): ")
		if item == nil {
			err.SetError("Plugin.Data is nil")
			break
		}

		if i == nil {
			err.SetError("Callback is nil")
			break
		}

		item.SetValue(prefix+"-timestamp", time.Now())

		err = i.IsItemValid()
		if err.IsError() {
			item.SetValue(prefix, err)
			break
		}

		var ctx Plugin.Interface
		switch {
		case i.IsNativePlugin():
			ctx = i.Native
		case i.IsRpcPlugin():
			ctx = i.Rpc
		default:
			err.SetError("Neither Native nor RPC plugin configured")
		}
		if err.IsError() {
			break
		}

		err = item.Callback(Plugin.CallbackInitialise, ctx, args...)
		if err.IsError() || err.IsWarning() {
			item.SetValue(prefix, err)
			break
		}

		item.SetValue(prefix, "OK")
	}
	return err
}

func (i *PluginItem) RunWithPlugin(item *Plugin.Plugin, args ...any) Return.Error {
	var err Return.Error
	for range Only.Once {
		prefix := "slave-run"
		err.SetPrefix("Run(): ")
		if item == nil {
			err.SetError("Plugin.Data is nil")
			break
		}

		if i == nil {
			err.SetError("Callback is nil")
			break
		}

		item.SetValue(prefix+"-timestamp", time.Now())

		err = i.IsItemValid()
		if err.IsError() {
			item.SetValue(prefix, err)
			break
		}

		var ctx Plugin.Interface
		switch {
		case i.IsNativePlugin():
			ctx = i.Native
		case i.IsRpcPlugin():
			ctx = i.Rpc
		default:
			err.SetError("Neither Native nor RPC plugin configured")
		}
		if err.IsError() {
			break
		}

		go func() {
			log.Println("i.Common.identity.Callbacks.Run")
			err = item.Callback(Plugin.CallbackRun, ctx, args...)
			if err.IsError() || err.IsWarning() {
				log.Printf("Error: %s\n", err)
				item.SetValue(prefix, err)
				return
			}
			item.SetValue(prefix, "OK")
		}()
	}
	return err
}

func (i *PluginItem) NotifyWithPlugin(item *Plugin.Plugin, args ...any) Return.Error {
	var err Return.Error
	for range Only.Once {
		prefix := "slave-notify"
		err.SetPrefix("Notify(): ")
		if item == nil {
			err.SetError("Plugin.Data is nil")
			break
		}

		if i == nil {
			err.SetError("Callback is nil")
			break
		}

		item.SetValue(prefix+"-timestamp", time.Now())

		err = i.IsItemValid()
		if err.IsError() {
			item.SetValue(prefix, err)
			break
		}

		var ctx Plugin.Interface
		switch {
		case i.IsNativePlugin():
			ctx = i.Native
		case i.IsRpcPlugin():
			ctx = i.Rpc
		default:
			err.SetError("Neither Native nor RPC plugin configured")
		}
		if err.IsError() {
			break
		}

		err = item.Callback(Plugin.CallbackNotify, ctx, args...)
		if err.IsError() || err.IsWarning() {
			item.SetValue(prefix, err)
			break
		}

		item.SetValue(prefix, "OK")
	}
	return err
}
