package GoPlugLoader

import (
	"os"
	sysPlugin "plugin"

	"github.com/MickMake/GoUnify/Only"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
	"github.com/MickMake/GoPlug/utils/store"
)

//
// PluginItemInterface
// ---------------------------------------------------------------------------------------------------- //
type PluginItemInterface interface {
	// NewPlugin - Create a new instance of this plugin.
	NewPlugin() Return.Error
	// IsItemValid - Validate NativePluginInterface interface
	IsItemValid() Return.Error
	GetItemData() *Plugin.PluginData
	GetItemHooks() Plugin.HookStore
	SetItemInterface(ref any) Return.Error
	IsNativePlugin() bool
	IsRpcPlugin() bool
	GetPluginPath() *utils.FilePath

	Validate() Return.Error
	PluginLoad(id string, pluginPath utils.FilePath) Return.Error
	PluginUnload() Return.Error
	Serve() Return.Error

	Execute(args ...any) Return.Error
	Initialise(args ...any) Return.Error
	Run(args ...any) Return.Error
	Notify(args ...any) Return.Error

	SetPluginType(types Plugin.Types) Return.Error
	SetInterface(ref any) Return.Error
	SetHandshakeConfig(goplugin.HandshakeConfig) Return.Error

	Hooks() *Plugin.HookStruct
	Values() *store.ValueStruct

	Plugin.PluginDataInterface
}

// CreatePluginItem -
func CreatePluginItem(types Plugin.Types, identity *Plugin.Identity) PluginItemInterface {
	item, _ := NewPluginItem(types, identity)
	return &item
}

//
// PluginItem
// ---------------------------------------------------------------------------------------------------- //
type PluginItem struct {
	// Stores a pointer to either Native or Rpc plugin loader.
	Pluggable PluginItemInterface
	Error     Return.Error
}

// NewPluginItem is constructor of PluginManager
func NewPluginItem(types Plugin.Types, identity *Plugin.Identity) (PluginItem, Return.Error) {
	var item PluginItem

	for range Only.Once {
		item.Error = types.IsValid()
		if item.Error.IsError() {
			break
		}

		if identity == nil {
			item.Error.SetError("PluginItem is nil")
			break
		}

		if types.IsNative() {
			item.Pluggable = NewNativePlugin()
			item.Pluggable.SetPluginTypeNative()
		}

		if types.IsRpc() {
			item.Pluggable = NewRpcPlugin()
			item.Pluggable.SetPluginTypeRpc()
		}

		item.Pluggable.SetIdentity(identity)
		item.Pluggable.SetHandshakeConfig(Plugin.HandshakeConfig)

		var fp utils.FilePath
		fp, item.Error = utils.NewFile(os.Args[0])
		if item.Error.IsError() {
			break
		}
		item.Pluggable.SetFilename(fp)

		fp, item.Error = utils.NewDir(fp.GetDir())
		if item.Error.IsError() {
			break
		}
		item.Pluggable.SetDirectory(fp)

		var logname string
		if types.IsNative() {
			logname = identity.Name + "[native]"
		} else if types.IsRpc() {
			logname = identity.Name + "[rpc]"
		} else {
			logname = identity.Name
		}

		var l utils.Logger
		l, item.Error = utils.NewLogger(logname, "")
		if item.Error.IsError() {
			break
		}

		item.Pluggable.SetLogger(&l)
	}

	return item, item.Error
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.Plugin interface structure

func (p *PluginItem) NewPlugin() Return.Error {
	return p.Pluggable.NewPlugin()
}
func (p *PluginItem) SavePluginData() Return.Error {
	return p.Pluggable.SavePluginData()
}
func (p *PluginItem) LoadPluginData() Return.Error {
	return p.Pluggable.LoadPluginData()
}
func (p *PluginItem) SaveObject(filename string, ref any) Return.Error {
	return p.Pluggable.SaveObject(filename, ref)
}
func (p *PluginItem) LoadObject(filename string, ref any) Return.Error {
	return p.Pluggable.LoadObject(filename, ref)
}
func (p *PluginItem) RefPlugin() *Plugin.PluginData {
	return p.Pluggable.RefPlugin()
}
func (p *PluginItem) RefCommon() *Plugin.Common {
	return p.Pluggable.RefCommon()
}
func (p *PluginItem) RefServices() *store.PluginServiceStruct {
	return p.Pluggable.RefServices()
}
func (p *PluginItem) GetData() Plugin.DynamicData {
	return p.Pluggable.GetData()
}
func (p *PluginItem) RefDynamic() *Plugin.DynamicData {
	return p.Pluggable.RefDynamic()
}
func (p *PluginItem) RegisterStructure(ref any) Return.Error {
	return p.Pluggable.RegisterStructure(ref)
}
func (p *PluginItem) Identify() Plugin.Identity {
	return p.Pluggable.Identify()
}
func (p *PluginItem) IdentifyString() string {
	return p.Pluggable.IdentifyString()
}
func (p *PluginItem) SaveIdentity() Return.Error {
	return p.Pluggable.SaveIdentity()
}
func (p *PluginItem) String() string {
	return p.Pluggable.String()
}
func (p *PluginItem) Print() {
	p.Pluggable.Print()
}

func (p *PluginItem) IsItemValid() Return.Error {
	return p.Pluggable.IsItemValid()
}
func (p *PluginItem) GetItemData() *Plugin.PluginData {
	return p.Pluggable.GetItemData()
}
func (p *PluginItem) GetItemHooks() Plugin.HookStore {
	return p.Pluggable.GetItemHooks()
}
func (p *PluginItem) SetItemInterface(ref any) Return.Error {
	return p.Pluggable.SetItemInterface(ref)
}
func (p *PluginItem) IsNativePlugin() bool {
	return p.Pluggable.IsNativePlugin()
}
func (p *PluginItem) IsRpcPlugin() bool {
	return p.Pluggable.IsRpcPlugin()
}
func (p *PluginItem) GetPluginPath() *utils.FilePath {
	return p.Pluggable.GetPluginPath()
}
func (p *PluginItem) Validate() Return.Error {
	return p.Pluggable.Validate()
}
func (p *PluginItem) PluginLoad(id string, pluginPath utils.FilePath) Return.Error {
	return p.Pluggable.PluginLoad(id, pluginPath)
}
func (p *PluginItem) PluginUnload() Return.Error {
	return p.Pluggable.PluginUnload()
}
func (p *PluginItem) Serve() Return.Error {
	return p.Pluggable.Serve()
}
func (p *PluginItem) Execute(args ...any) Return.Error {
	return p.Pluggable.Execute(args...)
}
func (p *PluginItem) Initialise(args ...any) Return.Error {
	return p.Pluggable.Initialise(args...)
}
func (p *PluginItem) Run(args ...any) Return.Error {
	return p.Pluggable.Run(args...)
}
func (p *PluginItem) Notify(args ...any) Return.Error {
	return p.Pluggable.Notify(args...)
}
func (p *PluginItem) Hooks() *Plugin.HookStruct {
	return p.Pluggable.Hooks()
}
func (p *PluginItem) Values() *store.ValueStruct {
	return p.Pluggable.Values()
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.CommonInterface interface structure

func (p *PluginItem) IsCommonValid() Return.Error {
	return p.Pluggable.IsCommonValid()
}
func (p *PluginItem) InitCommon() Return.Error {
	return p.Pluggable.InitCommon()
}
func (p *PluginItem) IsCommonConfigured() bool {
	return p.Pluggable.IsCommonConfigured()
}
func (p *PluginItem) GetCommonError() Return.Error {
	return p.Pluggable.GetCommonError()
}
func (p *PluginItem) IsCommonError() bool {
	return p.Pluggable.IsCommonError()
}
func (p *PluginItem) GetCommonRef() *Plugin.Common {
	return p.Pluggable.GetCommonRef()
}
func (p *PluginItem) SetRawInterface(ref any) {
	p.Pluggable.SetRawInterface(ref)
}
func (p *PluginItem) GetRawInterface() any {
	return p.Pluggable.GetRawInterface()
}
func (p *PluginItem) SetLogger(logger *utils.Logger) {
	p.Pluggable.SetLogger(logger)
}
func (p *PluginItem) GetLogger() *utils.Logger {
	return p.Pluggable.GetLogger()
}
func (p *PluginItem) SetLogFile(filename string) Return.Error {
	return p.Pluggable.SetLogFile(filename)
}
func (p *PluginItem) SetPluginType(name Plugin.Types) Return.Error {
	return p.Pluggable.SetPluginType(name)
}
func (p *PluginItem) SetPluginTypeNative() Return.Error {
	return p.Pluggable.SetPluginTypeNative()
}
func (p *PluginItem) SetPluginTypeRpc() Return.Error {
	return p.Pluggable.SetPluginTypeRpc()
}
func (p *PluginItem) GetPluginType() Plugin.Types {
	return p.Pluggable.GetPluginType()
}
func (p *PluginItem) SetStructName(ref interface{}) {
	p.Pluggable.SetStructName(ref)
}
func (p *PluginItem) GetStructName() string {
	return p.Pluggable.GetStructName()
}
func (p *PluginItem) SetFilename(pluginPath utils.FilePath) Return.Error {
	return p.Pluggable.SetFilename(pluginPath)
}
func (p *PluginItem) GetFilename() utils.FilePath {
	return p.Pluggable.GetFilename()
}
func (p *PluginItem) SetDirectory(pluginPath utils.FilePath) Return.Error {
	return p.Pluggable.SetDirectory(pluginPath)
}
func (p *PluginItem) GetDirectory() utils.FilePath {
	return p.Pluggable.GetDirectory()
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of store.PluginServiceInterface interface structure

func (p *PluginItem) NewPluginService() Return.Error {
	return p.Pluggable.NewPluginService()
}
func (p *PluginItem) GetPluginServiceReference() *store.PluginServiceStruct {
	return p.Pluggable.GetPluginServiceReference()
}
func (p *PluginItem) GetPluginIdentity() string {
	return p.Pluggable.GetPluginIdentity()
}
func (p *PluginItem) SetPluginIdentity(identity string) Return.Error {
	return p.Pluggable.SetPluginIdentity(identity)
}
func (p *PluginItem) ServiceExists(name string) bool {
	return p.Pluggable.ServiceExists(name)
}
func (p *PluginItem) ServiceNotExists(name string) bool {
	return p.Pluggable.ServiceNotExists(name)
}
func (p *PluginItem) SetNativeService(name string, value sysPlugin.Plugin) Return.Error {
	return p.Pluggable.SetNativeService(name, value)
}
func (p *PluginItem) GetNativeService(name string) sysPlugin.Plugin {
	return p.Pluggable.GetNativeService(name)
}
func (p *PluginItem) GetAsNativePluginSet() store.NativeServiceMap {
	return p.Pluggable.GetAsNativePluginSet()
}
func (p *PluginItem) SetRpcService(name string, value goplugin.Plugin) Return.Error {
	return p.Pluggable.SetRpcService(name, value)
}
func (p *PluginItem) GetRpcService(name string) goplugin.Plugin {
	return p.Pluggable.GetRpcService(name)
}
func (p *PluginItem) GetAsRpcPluginSet() goplugin.PluginSet {
	return p.Pluggable.GetAsRpcPluginSet()
}
func (p *PluginItem) ValidateService() Return.Error {
	return p.Pluggable.ValidateService()
}
func (p *PluginItem) CountServices() int {
	return p.Pluggable.CountServices()
}
func (p *PluginItem) ListServices() store.RpcServiceMap {
	return p.Pluggable.ListServices()
}
func (p *PluginItem) PrintServices() {
	p.Pluggable.PrintServices()
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.DynamicDataInterface interface structure

func (p *PluginItem) NewDynamicData(plug Plugin.PluginData) {
	p.Pluggable.NewDynamicData(plug)
}
func (p *PluginItem) RefIdentity() *Plugin.Identity {
	return p.Pluggable.RefIdentity()
}
func (p *PluginItem) RefCallbacks() *Plugin.Callbacks {
	return p.Pluggable.RefCallbacks()
}
func (p *PluginItem) RefHooks() *Plugin.HookStruct {
	return p.Pluggable.RefHooks()
}
func (p *PluginItem) RefValues() *store.ValueStruct {
	return p.Pluggable.RefValues()
}
func (p *PluginItem) GetIdentity() Plugin.Identity {
	return p.Pluggable.GetIdentity()
}
func (p *PluginItem) SetIdentity(identity *Plugin.Identity) Return.Error {
	return p.Pluggable.SetIdentity(identity)
}
func (p *PluginItem) GetName() string {
	return p.Pluggable.GetName()
}
func (p *PluginItem) GetVersion() string {
	return p.Pluggable.GetVersion()
}
func (p *PluginItem) Callback(callback string, ctx Plugin.PluginDataInterface, args ...any) Return.Error {
	return p.Pluggable.Callback(callback, ctx, args...)
}
func (p *PluginItem) SetHookStore(hooks Plugin.HookStore) Return.Error {
	return p.Pluggable.SetHookStore(hooks)
}
func (p *PluginItem) GetHook(name string) *Plugin.Hook {
	return p.Pluggable.GetHook(name)
}
func (p *PluginItem) SetHook(name string, function Plugin.HookFunction, args ...any) Return.Error {
	return p.Pluggable.SetHook(name, function, args...)
}
func (p *PluginItem) CallHook(name string, args ...any) (Plugin.HookResponse, Return.Error) {
	return p.Pluggable.CallHook(name, args...)
}
func (p *PluginItem) ValueExists(key string) bool {
	return p.Pluggable.ValueExists(key)
}
func (p *PluginItem) ValueNotExists(key string) bool {
	return p.Pluggable.ValueNotExists(key)
}
func (p *PluginItem) SetValue(key string, value any) {
	p.Pluggable.SetValue(key, value)
}
func (p *PluginItem) GetValue(key string) any {
	return p.Pluggable.GetValue(key)
}
func (p *PluginItem) SetHandshakeConfig(config goplugin.HandshakeConfig) Return.Error {
	return p.Pluggable.SetHandshakeConfig(config)
}
func (p *PluginItem) SetInterface(ref any) Return.Error {
	return p.Pluggable.SetInterface(ref)
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of Plugin.HookStore interface structure

func (p *PluginItem) NewHookStore() Return.Error {
	return p.Pluggable.NewHookStore()
}
func (p *PluginItem) SetHookPlugin(plugin Plugin.PluginDataInterface) {
	p.Pluggable.SetHookPlugin(plugin)
}
func (p *PluginItem) GetHookReference() *Plugin.HookStruct {
	return p.Pluggable.GetHookReference()
}
func (p *PluginItem) GetHookIdentity() string {
	return p.Pluggable.GetHookIdentity()
}
func (p *PluginItem) SetHookIdentity(identity string) Return.Error {
	return p.Pluggable.SetHookIdentity(identity)
}
func (p *PluginItem) HookExists(hook string) bool {
	return p.Pluggable.HookExists(hook)
}
func (p *PluginItem) HookNotExists(hook string) bool {
	return p.Pluggable.HookNotExists(hook)
}
func (p *PluginItem) GetHookName(name string) (string, Return.Error) {
	return p.Pluggable.GetHookName(name)
}
func (p *PluginItem) GetHookFunction(name string) (Plugin.HookFunction, Return.Error) {
	return p.Pluggable.GetHookFunction(name)
}
func (p *PluginItem) GetHookArgs(name string) (Plugin.HookArgs, Return.Error) {
	return p.Pluggable.GetHookArgs(name)
}
func (p *PluginItem) ValidateHook(args ...any) Return.Error {
	return p.Pluggable.ValidateHook(args...)
}
func (p *PluginItem) CountHooks() int {
	return p.Pluggable.CountHooks()
}
func (p *PluginItem) ListHooks() Plugin.HookMap {
	return p.Pluggable.ListHooks()
}
func (p *PluginItem) PrintHooks() {
	p.Pluggable.PrintHooks()
}

// ---------------------------------------------------------------------------------------------------- //

//
// PluginItems
// ---------------------------------------------------------------------------------------------------- //
type PluginItems []*PluginItem
