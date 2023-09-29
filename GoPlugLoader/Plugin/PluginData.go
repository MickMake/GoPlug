package Plugin

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	sysPlugin "plugin"
	"strings"
	"time"

	"github.com/MickMake/GoUnify/Only"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
	"github.com/MickMake/GoPlug/utils/store"
)

//
// PluginDataInterface
// ---------------------------------------------------------------------------------------------------- //
// 4. PluginData setup.
//goland:noinspection GoNameStartsWithPackageName
type PluginDataInterface interface {
	NewPlugin() Return.Error
	SavePluginData() Return.Error
	LoadPluginData() Return.Error

	SetPluginTypeNative() Return.Error
	SetPluginType(types Types) Return.Error
	SaveObject(filename string, ref any) Return.Error
	LoadObject(filename string, ref any) Return.Error
	// AddService(name string, plugin goplugin.Plugin) Return.Error

	// ---------------------------------------------------------------------------------------------------- //

	RefPlugin() *PluginData
	RefCommon() *Common
	RefServices() *store.PluginServiceStruct
	GetData() DynamicData
	RefDynamic() *DynamicData
	RegisterStructure(ref any) Return.Error

	// ---------------------------------------------------------------------------------------------------- //

	RefIdentity() *Identity
	// Identify - Return the Identity structure.
	Identify() Identity
	IdentifyString() string
	// SaveIdentity - Saves the config.PluginIdentity struct as a JSON file.
	SaveIdentity() Return.Error

	CommonInterface
	store.PluginServiceInterface
	DynamicDataInterface
}

//goland:noinspection GoUnusedExportedFunction
func CreatePlugin() PluginDataInterface {
	return NewPlugin()
}

//
// PluginData
// ---------------------------------------------------------------------------------------------------- //
//goland:noinspection GoNameStartsWithPackageName
type PluginData struct {
	Common   Common                    `json:"common"`
	Services store.PluginServiceStruct `json:"services"`
	Dynamic  DynamicData               `json:"dynamic"`
	Error    Return.Error              `json:"error"`
}

func NewPlugin() *PluginData {
	var ret PluginData
	ret = PluginData{
		Common:   NewPluginCommon(nil),
		Services: store.NewPluginServiceStruct(),
		Dynamic:  *NewDynamicData(PluginData{}),
		Error:    Return.New(),
	}
	ret.Dynamic.SetHookPlugin(&ret)
	return &ret
}

// NewPlugin - Create a new, (or existing), instance of this structure.
func (p *PluginData) NewPlugin() Return.Error {
	*p = *NewPlugin()
	return p.Error
}

func (p *PluginData) RefPlugin() *PluginData {
	return p
}

func (p *PluginData) String() string {
	p.Error = Return.Ok
	var ret string
	ret += p.Common.String()
	ret += p.Services.String()
	ret += p.Dynamic.String()
	return ret
}

func (p *PluginData) Print() {
	fmt.Println(p.String())
}

// SavePluginData - Save an arbitrary plugin structure as a JSON file.
func (p *PluginData) SavePluginData() Return.Error {
	return p.Common.Filename.SaveObject(*p)
}

// LoadPluginData - Load a JSON file into an arbitrary plugin structure.
func (p *PluginData) LoadPluginData() Return.Error {
	return p.Common.Filename.LoadObject(p)
}

// SaveObject - Save an arbitrary plugin structure as a JSON file.
func (p *PluginData) SaveObject(filename string, ref any) Return.Error {
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
func (p *PluginData) LoadObject(filename string, ref any) Return.Error {
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

func (p *PluginData) RefCommon() *Common {
	return &p.Common
}

func (p *PluginData) RefServices() *store.PluginServiceStruct {
	return &p.Services
}

func (p *PluginData) GetData() DynamicData {
	return p.Dynamic
}

func (p *PluginData) RefDynamic() *DynamicData {
	return &p.Dynamic
}

func (p *PluginData) RegisterStructure(ref any) Return.Error {
	p.Error = Return.Ok
	gob.Register(ref)
	return p.Error
}

// Identify - Get the identity of this plugin using the Identity structure.
func (p *PluginData) Identify() Identity {
	p.Error = Return.Ok
	return p.Dynamic.Identity
}

func (p *PluginData) IdentifyString() string {
	p.Error = Return.Ok
	return StructToString("Identity", p.Dynamic.Identity)
}

// SaveIdentity - Saves the PluginIdentity struct as a JSON file.
func (p *PluginData) SaveIdentity() Return.Error {
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

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.CommonInterface interface structure

func (p *PluginData) IsCommonValid() Return.Error {
	return p.Common.IsCommonValid()
}
func (p *PluginData) InitCommon() Return.Error {
	return p.Common.InitCommon()
}
func (p *PluginData) IsCommonConfigured() bool {
	return p.Common.IsCommonConfigured()
}
func (p *PluginData) GetCommonError() Return.Error {
	return p.Common.GetCommonError()
}
func (p *PluginData) IsCommonError() bool {
	return p.Common.IsCommonError()
}
func (p *PluginData) GetCommonRef() *Common {
	return p.Common.GetCommonRef()
}
func (p *PluginData) SetRawInterface(ref any) {
	p.Common.SetRawInterface(ref)
}
func (p *PluginData) GetRawInterface() any {
	return p.Common.GetRawInterface()
}
func (p *PluginData) SetLogger(logger *utils.Logger) {
	p.Common.SetLogger(logger)
}
func (p *PluginData) GetLogger() *utils.Logger {
	return p.Common.GetLogger()
}
func (p *PluginData) SetLogFile(filename string) Return.Error {
	return p.Common.SetLogFile(filename)
}
func (p *PluginData) SetPluginType(name Types) Return.Error {
	p.Dynamic.Identity.SetPluginType(name)
	return p.Common.SetPluginType(name)
}
func (p *PluginData) SetPluginTypeNative() Return.Error {
	p.Dynamic.Identity.SetPluginTypeNative()
	return p.Common.SetPluginTypeNative()
}
func (p *PluginData) SetPluginTypeRpc() Return.Error {
	p.Dynamic.Identity.SetPluginTypeRpc()
	return p.Common.SetPluginTypeRpc()
}
func (p *PluginData) GetPluginType() Types {
	return p.Common.GetPluginType()
}
func (p *PluginData) SetStructName(ref interface{}) {
	p.Common.SetStructName(ref)
}
func (p *PluginData) GetStructName() string {
	return p.Common.GetStructName()
}
func (p *PluginData) SetFilename(pluginPath utils.FilePath) Return.Error {
	return p.Common.SetFilename(pluginPath)
}
func (p *PluginData) GetFilename() utils.FilePath {
	return p.Common.GetFilename()
}
func (p *PluginData) SetDirectory(pluginPath utils.FilePath) Return.Error {
	return p.Common.SetDirectory(pluginPath)
}
func (p *PluginData) GetDirectory() utils.FilePath {
	return p.Common.GetDirectory()
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of store.PluginServiceInterface interface structure

func (p *PluginData) NewPluginService() Return.Error {
	return p.Services.NewPluginService()
}
func (p *PluginData) GetPluginServiceReference() *store.PluginServiceStruct {
	return p.Services.GetPluginServiceReference()
}
func (p *PluginData) GetPluginIdentity() string {
	return p.Services.GetPluginIdentity()
}
func (p *PluginData) SetPluginIdentity(identity string) Return.Error {
	return p.Services.SetPluginIdentity(identity)
}
func (p *PluginData) ServiceExists(name string) bool {
	return p.Services.ServiceExists(name)
}
func (p *PluginData) ServiceNotExists(name string) bool {
	return p.Services.ServiceNotExists(name)
}
func (p *PluginData) SetNativeService(name string, value sysPlugin.Plugin) Return.Error {
	return p.Services.SetNativeService(name, value)
}
func (p *PluginData) GetNativeService(name string) sysPlugin.Plugin {
	return p.Services.GetNativeService(name)
}
func (p *PluginData) GetAsNativePluginSet() store.NativeServiceMap {
	return p.Services.GetAsNativePluginSet()
}
func (p *PluginData) SetRpcService(name string, value goplugin.Plugin) Return.Error {
	return p.Services.SetRpcService(name, value)
}
func (p *PluginData) GetRpcService(name string) goplugin.Plugin {
	return p.Services.GetRpcService(name)
}
func (p *PluginData) GetAsRpcPluginSet() goplugin.PluginSet {
	return p.Services.GetAsRpcPluginSet()
}
func (p *PluginData) ValidateService() Return.Error {
	return p.Services.ValidateService()
}
func (p *PluginData) CountServices() int {
	return p.Services.CountServices()
}
func (p *PluginData) ListServices() store.RpcServiceMap {
	return p.Services.ListServices()
}
func (p *PluginData) PrintServices() {
	p.Services.PrintServices()
}
func (p *PluginData) AddService(name string, plugin goplugin.Plugin) Return.Error {
	return p.Services.SetRpcService(name, plugin)
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.DynamicDataInterface interface structure

func (p *PluginData) NewDynamicData(plug PluginData) {
	p.Dynamic.NewDynamicData(plug)
}
func (p *PluginData) RefIdentity() *Identity {
	return p.Dynamic.RefIdentity()
}
func (p *PluginData) RefCallbacks() *Callbacks {
	return p.Dynamic.RefCallbacks()
}
func (p *PluginData) RefHooks() *HookStruct {
	return p.Dynamic.RefHooks()
}
func (p *PluginData) RefValues() *store.ValueStruct {
	return p.Dynamic.RefValues()
}
func (p *PluginData) GetIdentity() Identity {
	return p.Dynamic.GetIdentity()
}
func (p *PluginData) SetIdentity(identity *Identity) Return.Error {
	return p.Dynamic.SetIdentity(identity)
}
func (p *PluginData) GetName() string {
	return p.Dynamic.GetName()
}
func (p *PluginData) GetVersion() string {
	return p.Dynamic.GetVersion()
}
func (p *PluginData) Callback(callback string, ctx PluginDataInterface, args ...any) Return.Error {
	for range Only.Once {
		p.Error.SetPrefix("Callback(%s): ", callback)
		prefix := "callback-" + callback
		p.SetValue(prefix+"-timestamp", time.Now())

		p.Error = p.IsCommonValid()
		if p.Error.IsError() {
			p.SetValue(prefix, p.Error)
			break
		}

		p.Error = p.Dynamic.Callback(callback, ctx, args...)
		if p.Error.IsError() || p.Error.IsWarning() {
			p.SetValue(prefix, p.Error)
			break
		}

		p.SetValue(prefix, "OK")
	}
	return p.Error
}
func (p *PluginData) SetHookStore(hooks HookStore) Return.Error {
	return p.Dynamic.SetHookStore(hooks)
}
func (p *PluginData) GetHook(name string) *Hook {
	return p.Dynamic.GetHook(name)
}
func (p *PluginData) SetHook(name string, function HookFunction, args ...any) Return.Error {
	return p.Dynamic.SetHook(name, function, args...)
}
func (p *PluginData) CallHook(name string, args ...any) (HookResponse, Return.Error) {
	return p.Dynamic.CallHook(name, args...)
}
func (p *PluginData) ValueExists(key string) bool {
	return p.Dynamic.ValueExists(key)
}
func (p *PluginData) ValueNotExists(key string) bool {
	return p.Dynamic.ValueNotExists(key)
}
func (p *PluginData) SetValue(key string, value any) {
	p.Dynamic.SetValue(key, value)
}
func (p *PluginData) GetValue(key string) any {
	return p.Dynamic.GetValue(key)
}
func (p *PluginData) SetHandshakeConfig(config goplugin.HandshakeConfig) Return.Error {
	return p.Dynamic.SetHandshakeConfig(config)
}
func (p *PluginData) SetInterface(ref any) Return.Error {
	return p.Dynamic.SetInterface(ref)
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.HookStore interface structure

func (p *PluginData) NewHookStore() Return.Error {
	return p.Dynamic.Hooks.NewHookStore()
}
func (p *PluginData) SetHookPlugin(plugin PluginDataInterface) {
	p.Dynamic.Hooks.SetHookPlugin(plugin)
}
func (p *PluginData) GetHookReference() *HookStruct {
	return p.Dynamic.Hooks.GetHookReference()
}
func (p *PluginData) GetHookIdentity() string {
	return p.Dynamic.Hooks.GetHookIdentity()
}
func (p *PluginData) SetHookIdentity(identity string) Return.Error {
	return p.Dynamic.Hooks.SetHookIdentity(identity)
}
func (p *PluginData) HookExists(hook string) bool {
	return p.Dynamic.Hooks.HookExists(hook)
}
func (p *PluginData) HookNotExists(hook string) bool {
	return p.Dynamic.Hooks.HookNotExists(hook)
}
func (p *PluginData) GetHookName(name string) (string, Return.Error) {
	return p.Dynamic.Hooks.GetHookName(name)
}
func (p *PluginData) GetHookFunction(name string) (HookFunction, Return.Error) {
	return p.Dynamic.Hooks.GetHookFunction(name)
}
func (p *PluginData) GetHookArgs(name string) (HookArgs, Return.Error) {
	return p.Dynamic.Hooks.GetHookArgs(name)
}
func (p *PluginData) ValidateHook(args ...any) Return.Error {
	return p.Dynamic.Hooks.ValidateHook(args...)
}
func (p *PluginData) CountHooks() int {
	return p.Dynamic.Hooks.CountHooks()
}
func (p *PluginData) ListHooks() HookMap {
	return p.Dynamic.Hooks.ListHooks()
}
func (p *PluginData) PrintHooks() {
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
