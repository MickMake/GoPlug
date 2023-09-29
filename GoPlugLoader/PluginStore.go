package GoPlugLoader

import (
	"fmt"
	"sync"

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
	lock  *sync.RWMutex
	hash  StoreItems
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
			pi := pl.GetItemData()
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

		filename := item.GetFilename()
		_, ok := ps.hash[filename.GetPath()]
		if !ok || (ok && forced) {
			ps.lock.Lock()
			//goland:noinspection GoDeferInLoop
			defer ps.lock.Unlock()

			ps.hash[filename.GetPath()] = item
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
			pi := i.GetItemData()
			if pi == nil {
				continue
			}

			if pi.Dynamic.Identity.Name != name {
				continue
			}

			item = i
			ps.Error = Return.Ok
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
		ps.Error.Clear()
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
		ps.Error.Clear()
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
