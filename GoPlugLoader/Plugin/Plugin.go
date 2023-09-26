package Plugin

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"plugin"
	"strings"

	"github.com/MickMake/GoUnify/Only"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
	"github.com/MickMake/GoPlug/utils/store"
)

// Interface
// ---------------------------------------------------------------------------------------------------- //
// 4. Plugin setup.
type Interface interface {
	NewPlugin() Return.Error
	SavePlugin() Return.Error
	LoadPlugin() Return.Error

	SetPluginTypeNative() Return.Error
	SetPluginType(types Types) Return.Error
	AddService(name string, plugin goplugin.Plugin) Return.Error
	SaveObject(filename string, ref any) Return.Error
	LoadObject(filename string, ref any) Return.Error

	// ---------------------------------------------------------------------------------------------------- //

	RefPlugin() *Plugin
	RefCommon() *Common
	RefServices() *store.PluginServiceStruct
	GetData() DynamicData
	RefDynamic() *DynamicData
	// SetInterface(ref any) Return.Error
	RegisterStructure(ref any) Return.Error
	// SetHandshakeConfig(config goplugin.HandshakeConfig) Return.Error
	// ---------------------------------------------------------------------------------------------------- //

	RefIdentity() *Identity
	// Identify - Return the Identity structure.
	Identify() Identity
	IdentifyString() string
	// GetName() string
	// GetVersion() string
	// SaveIdentity - Saves the config.PluginIdentity struct as a JSON file.
	SaveIdentity() Return.Error

	// ---------------------------------------------------------------------------------------------------- //
	// RefHooks() *store.HookStruct
	// SetHookStore(hooks store.HookStore) Return.Error
	// SetHook(name string, function store.HookFunction, args ...any) Return.Error
	// GetHook(name string) *store.Hook
	// CallHook(name string, args ...any) (store.HookResponse, Return.Error)
	// ---------------------------------------------------------------------------------------------------- //
	// RefCallbacks() *Callbacks
	// Callback(callback string, ctx Interface, args ...any) Return.Error
	// ---------------------------------------------------------------------------------------------------- //
	// RefValues() *store.ValueStruct
	// ValueExists(key string) bool
	// ValueNotExists(key string) bool
	// GetValue(key string) any
	// SetValue(key string, value any)
	// ---------------------------------------------------------------------------------------------------- //

	CommonInterface
	store.PluginServiceInterface
	DynamicDataInterface
}

func CreatePlugin() Interface {
	var ret Interface
	// ret = &Plugin{Data: NewData()}
	ret = NewPlugin()

	return ret
}

//
// Plugin
// ---------------------------------------------------------------------------------------------------- //
type Plugin struct {
	Common   Common
	Services store.PluginServiceStruct
	Dynamic  DynamicData
	Error    Return.Error
}

func NewPlugin() *Plugin {
	var ret Plugin
	ret = Plugin{
		Common:   NewPluginCommon(nil),
		Services: store.NewPluginServiceStruct(),
		Dynamic:  NewDynamicData(Plugin{}),
		Error:    Return.New(),
		// Data: NewData(),
	}
	ret.Dynamic.SetHookPlugin(&ret)
	return &ret
}

// NewPlugin - Create a new, (or existing), instance of this structure.
func (p *Plugin) NewPlugin() Return.Error {
	*p = *NewPlugin()
	return p.Error
}

func (p *Plugin) RefPlugin() *Plugin {
	return p
}

func (p *Plugin) String() string {
	p.Error = Return.Ok
	var ret string
	ret += p.Common.String()
	ret += p.Services.String()
	ret += p.Dynamic.String()
	return ret
}

func (p *Plugin) Print() {
	fmt.Println(p.String())
}

// SavePlugin - Save an arbitrary plugin structure as a JSON file.
func (p *Plugin) SavePlugin() Return.Error {
	return p.Common.Filename.SaveObject(*p)
}

// LoadPlugin - Load a JSON file into an arbitrary plugin structure.
func (p *Plugin) LoadPlugin() Return.Error {
	return p.Common.Filename.LoadObject(p)
}

// SaveObject - Save an arbitrary plugin structure as a JSON file.
func (p *Plugin) SaveObject(filename string, ref any) Return.Error {
	var file utils.FilePath
	if strings.Contains(filename, "/") {
		file, _ = utils.NewFile(filename)
	} else {
		file, _ = utils.NewFile(p.Common.Directory.GetPath(), filename)
	}
	file = file.ChangeExtension(".json")
	return file.SaveObject(ref)
}

// LoadObject - Load a JSON file into an arbitrary plugin structure.
func (p *Plugin) LoadObject(filename string, ref any) Return.Error {
	var file utils.FilePath
	if strings.Contains(filename, "/") {
		file, _ = utils.NewFile(filename)
	} else {
		file, _ = utils.NewFile(p.Common.Directory.GetPath(), filename)
	}
	file = file.ChangeExtension(".json")
	return file.LoadObject(ref)
}

// ---------------------------------------------------------------------------------------------------- //

func (p *Plugin) RefCommon() *Common {
	return &p.Common
}

func (p *Plugin) RefServices() *store.PluginServiceStruct {
	return &p.Services
}

func (p *Plugin) GetData() DynamicData {
	return p.Dynamic
}

func (p *Plugin) RefDynamic() *DynamicData {
	return &p.Dynamic
}

func (p *Plugin) RegisterStructure(ref any) Return.Error {
	p.Error = Return.Ok
	gob.Register(ref)
	return p.Error
}

// Identify - Get the identity of this plugin using the Identity structure.
func (p *Plugin) Identify() Identity {
	p.Error = Return.Ok
	return p.Dynamic.Identity
}

func (p *Plugin) IdentifyString() string {
	p.Error = Return.Ok
	return StructToString("Identity", p.Dynamic.Identity)
}

// SaveIdentity - Saves the PluginIdentity struct as a JSON file.
func (p *Plugin) SaveIdentity() Return.Error {
	for range Only.Once {
		file := p.Common.Filename.ChangeExtension("json")
		// Instead, maybe use d.Dynamic.Identity.Name

		p.Error = file.SaveObject(p.Dynamic.Identity)
		if p.Error.IsError() {
			break
		}
	}
	return p.Error
}

// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.CommonInterface interface structure

func (p *Plugin) IsCommonValid() Return.Error {
	return p.Common.IsCommonValid()
}
func (p *Plugin) InitCommon() Return.Error {
	return p.Common.InitCommon()
}
func (p *Plugin) IsCommonConfigured() bool {
	return p.Common.IsCommonConfigured()
}
func (p *Plugin) GetCommonError() Return.Error {
	return p.Common.GetCommonError()
}
func (p *Plugin) IsCommonError() bool {
	return p.Common.IsCommonError()
}
func (p *Plugin) GetCommonRef() *Common {
	return p.Common.GetCommonRef()
}
func (p *Plugin) SetRawInterface(ref any) {
	p.Common.SetRawInterface(ref)
}
func (p *Plugin) SetLogger(logger *utils.Logger) {
	p.Common.SetLogger(logger)
}
func (p *Plugin) GetLogger() *utils.Logger {
	return p.Common.GetLogger()
}
func (p *Plugin) SetLogFile(filename string) Return.Error {
	return p.Common.SetLogFile(filename)
}
func (p *Plugin) SetPluginType(name Types) Return.Error {
	p.Dynamic.Identity.SetPluginType(name)
	return p.Common.SetPluginType(name)
}
func (p *Plugin) SetPluginTypeNative() Return.Error {
	p.Dynamic.Identity.SetPluginTypeNative()
	return p.Common.SetPluginTypeNative()
}
func (p *Plugin) SetPluginTypeRpc() Return.Error {
	p.Dynamic.Identity.SetPluginTypeRpc()
	return p.Common.SetPluginTypeRpc()
}
func (p *Plugin) GetPluginType() Types {
	return p.Common.GetPluginType()
}
func (p *Plugin) SetStructName(ref interface{}) {
	p.Common.SetStructName(ref)
}
func (p *Plugin) GetStructName() string {
	return p.Common.GetStructName()
}
func (p *Plugin) SetFilename(pluginPath utils.FilePath) Return.Error {
	return p.Common.SetFilename(pluginPath)
}
func (p *Plugin) GetFilename() utils.FilePath {
	return p.Common.GetFilename()
}
func (p *Plugin) SetDirectory(pluginPath utils.FilePath) Return.Error {
	return p.Common.SetDirectory(pluginPath)
}
func (p *Plugin) GetDirectory() utils.FilePath {
	return p.Common.GetDirectory()
}

// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of store.PluginServiceInterface interface structure

func (p *Plugin) NewPluginService() Return.Error {
	return p.Services.NewPluginService()
}
func (p *Plugin) GetPluginServiceReference() *store.PluginServiceStruct {
	return p.Services.GetPluginServiceReference()
}
func (p *Plugin) GetPluginIdentity() string {
	return p.Services.GetPluginIdentity()
}
func (p *Plugin) SetPluginIdentity(identity string) Return.Error {
	return p.Services.SetPluginIdentity(identity)
}
func (p *Plugin) ServiceExists(name string) bool {
	return p.Services.ServiceExists(name)
}
func (p *Plugin) ServiceNotExists(name string) bool {
	return p.Services.ServiceNotExists(name)
}
func (p *Plugin) SetNativeService(name string, value plugin.Plugin) Return.Error {
	return p.Services.SetNativeService(name, value)
}
func (p *Plugin) GetNativeService(name string) plugin.Plugin {
	return p.Services.GetNativeService(name)
}
func (p *Plugin) GetAsNativePluginSet() store.NativeServiceMap {
	return p.Services.GetAsNativePluginSet()
}
func (p *Plugin) SetRpcService(name string, value goplugin.Plugin) Return.Error {
	return p.Services.SetRpcService(name, value)
}
func (p *Plugin) GetRpcService(name string) goplugin.Plugin {
	return p.Services.GetRpcService(name)
}
func (p *Plugin) GetAsRpcPluginSet() goplugin.PluginSet {
	return p.Services.GetAsRpcPluginSet()
}
func (p *Plugin) ValidateService() Return.Error {
	return p.Services.ValidateService()
}
func (p *Plugin) CountServices() int {
	return p.Services.CountServices()
}
func (p *Plugin) ListServices() store.RpcServiceMap {
	return p.Services.ListServices()
}
func (p *Plugin) PrintServices() {
	p.Services.PrintServices()
}
func (p *Plugin) AddService(name string, plugin goplugin.Plugin) Return.Error {
	return p.Services.SetRpcService(name, plugin)
}

// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.DynamicDataInterface interface structure

func (p *Plugin) NewDynamicData(plug Plugin) {
	p.Dynamic.NewDynamicData(plug)
}
func (p *Plugin) RefIdentity() *Identity {
	return p.Dynamic.RefIdentity()
}
func (p *Plugin) RefCallbacks() *Callbacks {
	return p.Dynamic.RefCallbacks()
}
func (p *Plugin) RefHooks() *HookStruct {
	return p.Dynamic.RefHooks()
}
func (p *Plugin) RefValues() *store.ValueStruct {
	return p.Dynamic.RefValues()
}
func (p *Plugin) GetIdentity() Identity {
	return p.Dynamic.GetIdentity()
}
func (p *Plugin) SetIdentity(identity *Identity) Return.Error {
	return p.Dynamic.SetIdentity(identity)
}
func (p *Plugin) GetName() string {
	return p.Dynamic.GetName()
}
func (p *Plugin) GetVersion() string {
	return p.Dynamic.GetVersion()
}
func (p *Plugin) Callback(callback string, ctx Interface, args ...any) Return.Error {
	return p.Dynamic.Callback(callback, ctx, args...)
}
func (p *Plugin) SetHookStore(hooks HookStore) Return.Error {
	return p.Dynamic.SetHookStore(hooks)
}
func (p *Plugin) GetHook(name string) *Hook {
	return p.Dynamic.GetHook(name)
}
func (p *Plugin) SetHook(name string, function HookFunction, args ...any) Return.Error {
	return p.Dynamic.SetHook(name, function, args...)
}
func (p *Plugin) CallHook(name string, args ...any) (HookResponse, Return.Error) {
	return p.Dynamic.CallHook(name, args...)
}
func (p *Plugin) ValueExists(key string) bool {
	return p.Dynamic.ValueExists(key)
}
func (p *Plugin) ValueNotExists(key string) bool {
	return p.Dynamic.ValueNotExists(key)
}
func (p *Plugin) SetValue(key string, value any) {
	p.Dynamic.SetValue(key, value)
}
func (p *Plugin) GetValue(key string) any {
	return p.Dynamic.GetValue(key)
}
func (p *Plugin) SetHandshakeConfig(config goplugin.HandshakeConfig) Return.Error {
	return p.Dynamic.SetHandshakeConfig(config)
}
func (p *Plugin) SetInterface(ref any) Return.Error {
	return p.Dynamic.SetInterface(ref)
}

// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.HookStore interface structure

func (p *Plugin) NewHookStore() Return.Error {
	return p.Dynamic.Hooks.NewHookStore()
}

func (p *Plugin) SetHookPlugin(plugin Interface) {
	p.Dynamic.Hooks.SetHookPlugin(plugin)
}

func (p *Plugin) GetHookReference() *HookStruct {
	return p.Dynamic.Hooks.GetHookReference()
}

func (p *Plugin) GetHookIdentity() string {
	return p.Dynamic.Hooks.GetHookIdentity()
}

func (p *Plugin) SetHookIdentity(identity string) Return.Error {
	return p.Dynamic.Hooks.SetHookIdentity(identity)
}

func (p *Plugin) HookExists(hook string) bool {
	return p.Dynamic.Hooks.HookExists(hook)
}

func (p *Plugin) HookNotExists(hook string) bool {
	return p.Dynamic.Hooks.HookNotExists(hook)
}

func (p *Plugin) GetHookName(name string) (string, Return.Error) {
	return p.Dynamic.Hooks.GetHookName(name)
}

func (p *Plugin) GetHookFunction(name string) (HookFunction, Return.Error) {
	return p.Dynamic.Hooks.GetHookFunction(name)
}

func (p *Plugin) GetHookArgs(name string) (HookArgs, Return.Error) {
	return p.Dynamic.Hooks.GetHookArgs(name)
}

func (p *Plugin) ValidateHook(args ...any) Return.Error {
	return p.Dynamic.Hooks.ValidateHook(args...)
}

func (p *Plugin) CountHooks() int {
	return p.Dynamic.Hooks.CountHooks()
}

func (p *Plugin) ListHooks() HookMap {
	return p.Dynamic.Hooks.ListHooks()
}

func (p *Plugin) PrintHooks() {
	p.Dynamic.Hooks.PrintHooks()
}

// ---------------------------------------------------------------------------------------------------- //

func StructToString(name string, ref any) string {
	data, err := json.MarshalIndent(ref, "", "\t")
	if err != nil {
		return "ERROR: " + err.Error()
	}
	return fmt.Sprintf("#### JSON[%s] ####\n%s\n", name, data)
}
