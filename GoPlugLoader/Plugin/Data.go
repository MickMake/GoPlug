package Plugin

// import (
// 	"encoding/gob"
// 	"fmt"
//
// 	"github.com/MickMake/GoUnify/Only"
// 	goplugin "github.com/hashicorp/go-plugin"
//
// 	"github.com/MickMake/GoPlug/utils/Return"
// 	"github.com/MickMake/GoPlug/utils/store"
// )
//
// //
// // DataInterface
// // ---------------------------------------------------------------------------------------------------- //
// // 5. Child plugin structure
// type DataInterface interface {
// 	New()
// 	String() string
// 	SetPluginType(types Types) Return.Error
// 	AddService(name string, plugin goplugin.Plugin) Return.Error
// 	SaveObject() Return.Error
// 	LoadObject(ref any) Return.Error
//
// 	// ---------------------------------------------------------------------------------------------------- //
//
// 	Ref() *Data
// 	RefCommon() *Common
// 	RefServices() *store.PluginServiceStruct
// 	GetData() DynamicData
// 	RefDynamic() *DynamicData
// 	SetInterface(ref any) Return.Error
// 	RegisterStructure(ref any) Return.Error
// 	SetHandshakeConfig(config goplugin.HandshakeConfig) Return.Error
//
// 	// ---------------------------------------------------------------------------------------------------- //
//
// 	RefIdentity() *Identity
// 	// Identify - Return the Identity structure.
// 	Identify() Identity
// 	IdentifyString() string
// 	GetName() string
// 	GetVersion() string
// 	// SaveIdentity - Saves the config.PluginIdentity struct as a JSON file.
// 	SaveIdentity() Return.Error
//
// 	// ---------------------------------------------------------------------------------------------------- //
//
// 	RefHooks() *store.HookStruct
// 	SetHookStore(hooks store.HookStore) Return.Error
// 	SetHook(name string, function store.HookFunction, args ...any) Return.Error
// 	GetHook(name string) *store.Hook
// 	CallHook(name string, args ...any) (store.HookResponse, Return.Error)
//
// 	// ---------------------------------------------------------------------------------------------------- //
//
// 	RefCallbacks() *Callbacks
// 	Callback(callback string, ctx Interface, args ...any) Return.Error
//
// 	// ---------------------------------------------------------------------------------------------------- //
//
// 	RefValues() *store.ValueStruct
// 	ValueExists(key string) bool
// 	ValueNotExists(key string) bool
// 	GetValue(key string) any
// 	SetValue(key string, value any)
//
// 	// ---------------------------------------------------------------------------------------------------- //
//
// 	CommonInterface
// 	store.PluginServiceInterface
// 	DynamicDataInterface
// }
//
// //
// // Data - Holds the transient plugin data.
// // ---------------------------------------------------------------------------------------------------- //
// type Data struct {
// 	Common   Common
// 	Services store.PluginServiceStruct
// 	Dynamic  DynamicData
// 	Error    Return.Error
// }
//
// func NewData() Data {
// 	ret := Data{
// 		Common:   NewPluginCommon(nil),
// 		Services: store.NewPluginServiceStruct(),
// 		Dynamic:  NewDynamicData(),
// 		Error:    Return.New(),
// 	}
// 	return ret
// }
//
// func (d *Data) New() {
// 	*d = NewData()
// }
//
// func (d *Data) String() string {
// 	d.Error = Return.Ok
// 	return d.Dynamic.Values.String()
// }
//
// func (d *Data) Print() {
// 	fmt.Println(d.String())
// }
//
// func (d *Data) SetPluginType(types Types) Return.Error {
// 	d.Error = Return.Ok
// 	d.Common.SetPluginType(types)
// 	d.Dynamic.SetPluginType(types)
// 	return d.Error
// }
//
// func (d *Data) AddService(name string, plugin goplugin.Plugin) Return.Error {
// 	return d.Services.SetRpcService(name, plugin)
// }
//
// // SaveObject - Save an arbitrary plugin structure as a JSON file.
// func (d *Data) SaveObject() Return.Error {
// 	return d.Common.Filename.SaveObject(*d)
// }
//
// // LoadObject - Load a JSON file into an arbitrary plugin structure.
// func (d *Data) LoadObject(ref any) Return.Error {
// 	return d.Common.Filename.LoadObject(*d)
// }
//
// // ---------------------------------------------------------------------------------------------------- //
//
// func (d *Data) Ref() *Data {
// 	return d
// }
//
// func (d *Data) RefCommon() *Common {
// 	return &d.Common
// }
//
// func (d *Data) RefServices() *store.PluginServiceStruct {
// 	return &d.Services
// }
//
// func (d *Data) GetData() DynamicData {
// 	return d.Dynamic
// }
//
// func (d *Data) RefDynamic() *DynamicData {
// 	return &d.Dynamic
// }
//
// func (d *Data) SetInterface(ref any) Return.Error {
// 	return d.Dynamic.SetInterface(ref)
// }
//
// func (d *Data) RegisterStructure(ref any) Return.Error {
// 	d.Error = Return.Ok
// 	gob.Register(ref)
// 	return d.Error
// }
//
// func (d *Data) SetHandshakeConfig(config goplugin.HandshakeConfig) Return.Error {
// 	return d.Dynamic.SetHandshakeConfig(config)
// }
//
// // ---------------------------------------------------------------------------------------------------- //
//
// // RefIdentity - Get the identity of this plugin using the Identity structure.
// func (d *Data) RefIdentity() *Identity {
// 	return &d.Dynamic.Identity
// }
//
// // Identify - Get the identity of this plugin using the Identity structure.
// func (d *Data) Identify() Identity {
// 	d.Error = Return.Ok
// 	return d.Dynamic.Identity
// }
//
// func (d *Data) IdentifyString() string {
// 	d.Error = Return.Ok
// 	return StructToString("Identity", d.Dynamic.Identity)
// }
//
// // SaveIdentity - Saves the PluginIdentity struct as a JSON file.
// func (d *Data) SaveIdentity() Return.Error {
// 	for range Only.Once {
// 		file := d.Common.Filename.ChangeExtension("json")
// 		// Instead, maybe use d.Dynamic.Identity.Name
//
// 		d.Error = file.SaveObject(d.Dynamic.Identity)
// 		if d.Error.IsError() {
// 			break
// 		}
// 	}
// 	return d.Error
// }
//
// // SetIdentity - Set the identity of this plugin using the Identity structure.
// func (d *Data) SetIdentity(identity *Identity) Return.Error {
// 	d.Error = d.Dynamic.SetIdentity(identity)
// 	if d.Error.IsError() {
// 		return d.Error
// 	}
// 	d.Common.Id = identity.Name
// 	d.Common.PluginTypes = identity.PluginTypes
// 	d.Services.Identity = identity.Name
// 	return d.Error
// }
//
// // GetName - Get the name of the plugin.
// func (d *Data) GetName() string {
// 	d.Error = Return.Ok
// 	if d == nil {
// 		return ""
// 	}
// 	return d.Dynamic.GetName()
// }
//
// // GetVersion - Get the version of the plugin.
// func (d *Data) GetVersion() string {
// 	d.Error = Return.Ok
// 	if d == nil {
// 		return ""
// 	}
// 	return d.Dynamic.GetVersion()
// }
//
// // ---------------------------------------------------------------------------------------------------- //
//
// func (d *Data) RefHooks() *store.HookStruct {
// 	return &d.Dynamic.Hooks
// }
//
// func (d *Data) SetHookStore(hooks store.HookStore) Return.Error {
// 	d.Error = Return.Ok
// 	return d.Dynamic.SetHookStore(hooks)
// }
//
// func (d *Data) GetHook(name string) *store.Hook {
// 	d.Error = Return.Ok
// 	return d.Dynamic.GetHook(name)
// }
//
// func (d *Data) SetHook(name string, function store.HookFunction, args ...any) Return.Error {
// 	d.Error = Return.Ok
// 	return d.Dynamic.SetHook(name, function, args...)
// }
//
// func (d *Data) CallHook(name string, args ...any) (store.HookResponse, Return.Error) {
// 	d.Error = Return.Ok
// 	return d.Dynamic.CallHook(name, args...)
// }
//
// // ---------------------------------------------------------------------------------------------------- //
//
// func (d *Data) RefCallbacks() *Callbacks {
// 	return &d.Dynamic.Identity.Callbacks
// }
//
// func (d *Data) Callback(callback string, ctx Interface, args ...any) Return.Error {
// 	return d.Dynamic.Identity.Callback(callback, ctx, args...)
// }
//
// // ---------------------------------------------------------------------------------------------------- //
//
// func (d *Data) RefValues() *store.ValueStruct {
// 	return &d.Dynamic.Values
// }
//
// func (d *Data) ValueExists(key string) bool {
// 	d.Error = Return.Ok
// 	return d.Dynamic.ValueExists(key)
// }
//
// func (d *Data) ValueNotExists(key string) bool {
// 	d.Error = Return.Ok
// 	return d.Dynamic.ValueNotExists(key)
// }
//
// func (d *Data) GetValue(key string) any {
// 	d.Error = Return.Ok
// 	return d.Dynamic.GetValue(key)
// }
//
// func (d *Data) SetValue(key string, value any) {
// 	d.Error = Return.Ok
// 	d.Dynamic.SetValue(key, value)
// }
