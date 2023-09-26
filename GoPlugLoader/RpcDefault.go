package GoPlugLoader

import (
	"net/rpc"

	gplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
	"github.com/MickMake/GoPlug/utils/store"
)

//
// RpcDefaultStruct
// ---------------------------------------------------------------------------------------------------- //
type RpcDefaultStruct struct {
}

func (d RpcDefaultStruct) Server(broker *gplugin.MuxBroker) (interface{}, error) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcDefaultStruct) Client(broker *gplugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}

//
// RpcServerDefaultStruct
// ---------------------------------------------------------------------------------------------------- //
type RpcServerDefaultStruct struct {
}

func (d RpcServerDefaultStruct) New() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) RegisterStructure(ref any) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetIdentity(identity *Plugin.Identity) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetPluginType(types Plugin.Types) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetInterface(ref any) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetHandshakeConfig(config gplugin.HandshakeConfig) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetHookStore(hooks *Plugin.HookStruct) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) Serve() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) Init() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) IsValid() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) Debug(ref any) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) String() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) IsConfigured() bool {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetError() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) IsError() bool {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetRef() *Plugin.Common {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetLogger(logger *utils.Logger) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetLogger() *utils.Logger {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetLogFile(filename string) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetPluginType() Plugin.Types {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetIdentity() *Plugin.Identity {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SaveIdentity() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) Identify() (Plugin.Identity, Return.Error) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetStructName(ref interface{}) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetStructName() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetName() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetVersion() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetIdentityKey(key string) string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetFilename(filename utils.FilePath) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetFilename() utils.FilePath {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetDir() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SetConfigDir(dir string) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) GetConfigDir() utils.FilePath {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SaveJson(filename string, data []byte) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) LoadJson(filename string) ([]byte, Return.Error) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) SaveStruct(filename string, ref any) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) LoadStruct(filename string, ref any) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) Build(args ...string) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) Validate() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcServerDefaultStruct) Values() store.ValueStore {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}

//
// RpcClientDefaultStruct
// ---------------------------------------------------------------------------------------------------- //
type RpcClientDefaultStruct struct {
}

func (d RpcClientDefaultStruct) Init() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) IsValid() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) Debug(ref any) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) String() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) IsConfigured() bool {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetError() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) IsError() bool {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetRef() *Plugin.Common {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SetLogger(logger *utils.Logger) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetLogger() *utils.Logger {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SetLogFile(filename string) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SetPluginType(name Plugin.Types) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetPluginType() Plugin.Types {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SetIdentity(identity *Plugin.Identity) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetIdentity() *Plugin.Identity {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SaveIdentity() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) Identify() (Plugin.Identity, Return.Error) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SetStructName(ref interface{}) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetStructName() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetName() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetVersion() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetIdentityKey(key string) string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SetFilename(filename utils.FilePath) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetFilename() utils.FilePath {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetDir() string {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SetConfigDir(dir string) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) GetConfigDir() utils.FilePath {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SaveJson(filename string, data []byte) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) LoadJson(filename string) ([]byte, Return.Error) {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) SaveStruct(filename string, ref any) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) LoadStruct(filename string, ref any) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) Build(args ...string) Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) Validate() Return.Error {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
func (d RpcClientDefaultStruct) Values() store.ValueStore {
	utils.DEBUG()
	// @TODO - Implement me
	panic("implement me")
}
