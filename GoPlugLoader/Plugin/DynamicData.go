package Plugin

import (
	"fmt"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/utils/Return"
	"github.com/MickMake/GoPlug/utils/store"
)

//
// DynamicDataInterface
// ---------------------------------------------------------------------------------------------------- //
type DynamicDataInterface interface {
	NewDynamicData(plug PluginData)
	RefIdentity() *Identity
	SetIdentity(identity *Identity) Return.Error
	GetIdentity() Identity
	// GetName - Get the name of the plugin.
	GetName() string
	// GetVersion - Get the version of the plugin.
	GetVersion() string
	SetPluginType(types Types) Return.Error

	RefCallbacks() *Callbacks
	Callback(callback string, ctx PluginDataInterface, args ...any) Return.Error

	RefHooks() *HookStruct
	SetHookStore(hooks HookStore) Return.Error
	GetHook(name string) *Hook
	// SetHook - Creates a function hook.
	// Leave 'name' blank to autofill with the function/method name.
	// 'args' is a list of arguments that this hook expects.
	// When calling, it will expect the number and type of arguments to match.
	SetHook(name string, function HookFunction, args ...any) Return.Error
	// CallHook - Calls a predefined hook.
	// 'name' has to exist.
	// 'args' also have to match, both into quantity and type.
	CallHook(name string, args ...any) (HookResponse, Return.Error)

	RefValues() *store.ValueStruct
	ValueExists(key string) bool
	ValueNotExists(key string) bool
	SetValue(key string, value any)
	GetValue(key string) any

	SetHandshakeConfig(config goplugin.HandshakeConfig) Return.Error
	SetInterface(ref any) Return.Error
	Print()
	String() string

	HookStore
}

//goland:noinspection GoUnusedExportedFunction
func CreateDynamicData(plug PluginData) DynamicDataInterface {
	return NewDynamicData(plug)
}

//
// DynamicData
// ---------------------------------------------------------------------------------------------------- //
type DynamicData struct {
	Identity        Identity                 `json:"identity"`
	Hooks           HookStruct               `json:"hooks"`
	Values          store.ValueStruct        `json:"values"`
	Interface       any                      `json:"-"`
	HandshakeConfig goplugin.HandshakeConfig `json:"handshakeConfig"`
	Error           Return.Error             `json:"-"`
}

func NewDynamicData(plug PluginData) *DynamicData {
	ret := DynamicData{
		Hooks:           NewHookStruct(),
		Values:          store.NewValueStruct(),
		Interface:       nil,
		Identity:        Identity{},
		HandshakeConfig: HandshakeConfig,
		Error:           Return.New(),
	}
	ret.SetHookPlugin(&plug)
	return &ret
}

func (d *DynamicData) NewDynamicData(plug PluginData) {
	*d = *NewDynamicData(plug)
}

func (d *DynamicData) RefIdentity() *Identity {
	return &d.Identity
}

// SetIdentity - Set the identity of this plugin using the Identity structure.
func (d *DynamicData) SetIdentity(identity *Identity) Return.Error {
	// Expose plugin name to "(c Callbacks) MarshalJSON()"
	if d.Identity.Callbacks.PluginName == "" {
		d.Identity.Callbacks.PluginName = d.Identity.Name
	}
	if d.Hooks.Identity == "" {
		d.Hooks.Identity = identity.Name
	}

	d.Error = identity.IsValid()
	if d.Error.IsError() {
		return d.Error
	}
	d.Identity = *identity
	return d.Error
}

// GetIdentity - Get the identity of this plugin using the Identity structure.
func (d *DynamicData) GetIdentity() Identity {
	if d == nil {
		return Identity{}
	}
	return d.Identity
}

// GetName - Get the name of the plugin.
func (d *DynamicData) GetName() string {
	if d == nil {
		return ""
	}
	return d.Identity.Name
}

// GetVersion - Get the version of the plugin.
func (d *DynamicData) GetVersion() string {
	if d == nil {
		return ""
	}
	return d.Identity.Version
}

func (d *DynamicData) SetPluginType(types Types) Return.Error {
	d.Error = Return.Ok
	d.Identity.PluginTypes = types
	return d.Error
}

// ---------------------------------------------------------------------------------------------------- //

func (d *DynamicData) RefCallbacks() *Callbacks {
	return &d.Identity.Callbacks
}

func (d *DynamicData) Callback(callback string, ctx PluginDataInterface, args ...any) Return.Error {
	return d.Identity.Callback(callback, ctx, args...)
}

// ---------------------------------------------------------------------------------------------------- //

func (d *DynamicData) RefHooks() *HookStruct {
	return &d.Hooks
}
func (d *DynamicData) SetHookStore(hooks HookStore) Return.Error {
	d.Error = Return.Ok
	d.Hooks = *hooks.GetHookReference()
	return d.Error
}
func (d *DynamicData) GetHook(name string) *Hook {
	return d.Hooks.GetHook(name)
}
func (d *DynamicData) SetHook(name string, function HookFunction, args ...any) Return.Error {
	return d.Hooks.SetHook(name, function, args...)
}
func (d *DynamicData) CallHook(name string, args ...any) (HookResponse, Return.Error) {
	return d.Hooks.CallHook(name, args...)
}

// ---------------------------------------------------------------------------------------------------- //

func (d *DynamicData) RefValues() *store.ValueStruct {
	return &d.Values
}

func (d *DynamicData) ValueExists(key string) bool {
	return d.Values.ValueExists(key)
}

func (d *DynamicData) ValueNotExists(key string) bool {
	return d.Values.ValueNotExists(key)
}

func (d *DynamicData) SetValue(key string, value any) {
	d.Values.SetValue(key, value)
}

func (d *DynamicData) GetValue(key string) any {
	return d.Values.GetValue(key)
}

// ---------------------------------------------------------------------------------------------------- //

func (d *DynamicData) SetHandshakeConfig(config goplugin.HandshakeConfig) Return.Error {
	d.Error = Return.Ok
	d.HandshakeConfig = config
	return d.Error
}

func (d *DynamicData) SetInterface(ref any) Return.Error {
	d.Error = Return.Ok
	d.Interface = ref
	return d.Error
}

func (d *DynamicData) Print() {
	fmt.Println(d.String())
}

func (d DynamicData) String() string {
	var ret string
	ret += d.Identity.String()
	ret += d.Hooks.String()
	ret += d.Values.String()
	return ret
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.HookStore interface structure

func (d *DynamicData) NewHookStore() Return.Error {
	return d.Hooks.NewHookStore()
}
func (d *DynamicData) SetHookPlugin(plugin PluginDataInterface) {
	d.Hooks.SetHookPlugin(plugin)
}
func (d *DynamicData) GetHookReference() *HookStruct {
	return d.Hooks.GetHookReference()
}
func (d *DynamicData) GetHookIdentity() string {
	return d.Hooks.GetHookIdentity()
}
func (d *DynamicData) SetHookIdentity(identity string) Return.Error {
	return d.Hooks.SetHookIdentity(identity)
}
func (d *DynamicData) HookExists(hook string) bool {
	return d.Hooks.HookExists(hook)
}
func (d *DynamicData) HookNotExists(hook string) bool {
	return d.Hooks.HookNotExists(hook)
}
func (d *DynamicData) GetHookName(name string) (string, Return.Error) {
	return d.Hooks.GetHookName(name)
}
func (d *DynamicData) GetHookFunction(name string) (HookFunction, Return.Error) {
	return d.Hooks.GetHookFunction(name)
}
func (d *DynamicData) GetHookArgs(name string) (HookArgs, Return.Error) {
	return d.Hooks.GetHookArgs(name)
}
func (d *DynamicData) ValidateHook(args ...any) Return.Error {
	return d.Hooks.ValidateHook(args...)
}
func (d *DynamicData) CountHooks() int {
	return d.Hooks.CountHooks()
}
func (d *DynamicData) ListHooks() HookMap {
	return d.Hooks.ListHooks()
}
func (d *DynamicData) PrintHooks() {
	d.Hooks.PrintHooks()
}
