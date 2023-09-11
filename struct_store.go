package GoPlug

import (
	"fmt"
	"sync"

	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoPlug/pluggable"
	"github.com/MickMake/GoUnify/Only"
)


// ---------------------------------------------------------------------------------------------------- //

// Store - Stores loaded plugins
type Store interface {
	IsValid() Return.Error

	// Size - The total count of current items in the store
	Size() uint

	String() string

	Print()

	// Put - Append the plugin item to the store
	// If forced is set to be true, new plugin item will overwrite the existing one
	// Try best to append, ignore any errors
	Put(item *pluggable.PluginItem, forced bool) Return.Error

	// Get - the plugin item by name
	// If existing, return the item and set the bool flag to true
	Get(name string) (*pluggable.PluginItem, Return.Error)

	// Remove - the plugin out of the store and return the removed plugin item
	// If successfully removed, set the bool flag to true
	Remove(name string) (*pluggable.PluginItem, Return.Error)
}


// ---------------------------------------------------------------------------------------------------- //

// BaseStore is the default implementation of Store interface
type BaseStore struct {
	// internal lock
	lock *sync.RWMutex

	// internal list
	hash map[string]*pluggable.PluginItem

	Error Return.Error
}

// NewBaseStore is constructor of BaseStore
func NewBaseStore() *BaseStore {
	var err Return.Error
	err.SetPrefix("BaseStore: ")
	return &BaseStore {
		lock: new(sync.RWMutex),
		hash: make(map[string]*pluggable.PluginItem),
		Error: err,
	}
}

func (bs *BaseStore) IsValid() Return.Error {
	switch {
		case bs == nil:
			bs.Error.SetError("BaseStore is nil")
		case bs.hash == nil:
			bs.Error.SetError("BaseStore map is nil")
	}
	return bs.Error
}

// String - .
func (bs BaseStore) String() string {
	var ret string

	for range Only.Once {
		bs.lock.Lock()
		//goland:noinspection GoDeferInLoop
		defer bs.lock.Unlock()

		ret = fmt.Sprintf("# %d plugins loaded\n", (uint)(len(bs.hash)))
		for _, pl := range bs.hash {
			if pl.Config == nil {
				continue
			}
			ret += fmt.Sprintf("Name: %s\tVersion: %s\tDescription: %s\n\tFile: %s\n",
				pl.Config.Name,
				pl.Config.Version,
				pl.Config.Description,
				pl.File.GetPath(),
			)
		}

		bs.Error = Return.Ok
	}

	return ret
}

// Print - .
func (bs *BaseStore) Print() {
	fmt.Printf("BaseStoreDir: %s\n", bs.String())
}

// Size is the implementation of same method in Store interface
func (bs *BaseStore) Size() uint {
	var size uint

	for range Only.Once {
		bs.lock.Lock()
		//goland:noinspection GoDeferInLoop
		defer bs.lock.Unlock()

		size = (uint)(len(bs.hash))
		bs.Error = Return.Ok
	}

	return size
}

// Put is the implementation of same method in Store interface
func (bs *BaseStore) Put(item *pluggable.PluginItem, forced bool) Return.Error {
	for range Only.Once {
		bs.Error = bs.IsValid()
		if bs.Error.IsError() {
			break
		}

		bs.Error = item.IsValid()
		if bs.Error.IsError() {
			break
		}

		_, ok := bs.hash[item.File.GetPath()]
		if !ok || (ok && forced) {
			bs.lock.Lock()
			//goland:noinspection GoDeferInLoop
			defer bs.lock.Unlock()

			bs.hash[item.File.GetPath()] = item
		}

		bs.Error = Return.Ok
	}

	return bs.Error
}

// Get is the implementation of same method in Store interface
func (bs *BaseStore) Get(name string) (*pluggable.PluginItem, Return.Error) {
	var item *pluggable.PluginItem

	for range Only.Once {
		bs.lock.RLock()
		//goland:noinspection GoDeferInLoop
		defer bs.lock.RUnlock()

		var ok bool
		item, ok = bs.hash[name]
		if ok {
			bs.Error = Return.Ok
			break
		}

		for _, s := range bs.hash {
			if s.Config == nil {
				continue
			}
			if s.Config.Name == name {
				item = s
				break
			}
		}
		if item != nil {
			break
		}

		bs.Error.SetError("plugin %s is not loaded", name)
	}

	return item, bs.Error
}

// Remove is the implementation of same method in Store interface
func (bs *BaseStore) Remove(name string) (*pluggable.PluginItem, Return.Error) {
	var item *pluggable.PluginItem

	for range Only.Once {
		bs.lock.Lock()
		//goland:noinspection GoDeferInLoop
		defer bs.lock.Unlock()

		var ok bool
		item, ok = bs.hash[name]
		if !ok {
			bs.Error.SetError("plugin %s is not loaded", name)
			break
		}

		delete(bs.hash, name)
		bs.Error = Return.Ok
	}

	return item, bs.Error
}
