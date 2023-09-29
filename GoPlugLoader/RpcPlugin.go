package GoPlugLoader

import (
	"encoding/gob"
	"log"
	"net/rpc"
	"os"
	"os/exec"

	"github.com/MickMake/GoUnify/Only"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
	"github.com/MickMake/GoPlug/utils/store"
)

//
// NewRpcPluginInterface - Create a new instance of this interface.
// ---------------------------------------------------------------------------------------------------- //
//goland:noinspection GoUnusedExportedFunction
func NewRpcPluginInterface() PluginItemInterface {
	ret := NewRpcPlugin()
	ret.SetPluginType(Plugin.RpcPluginType)
	return ret
}

//
// RpcPlugin
// ---------------------------------------------------------------------------------------------------- //
type RpcPlugin struct {
	RpcService
	Plugin.PluginData
}

// NewRpcPlugin - Create a new instance of this structure.
func NewRpcPlugin() *RpcPlugin {
	ret := RpcPlugin{
		RpcService: NewRpcService(),
		PluginData: *Plugin.NewPlugin(),
	}
	return &ret
}

func (p *RpcPlugin) NewRpcPlugin() Return.Error {
	*p = *NewRpcPlugin()
	return Return.Ok
}

func (p *RpcPlugin) GetRpcPlugin() *RpcPlugin {
	return p
}

func (p *RpcPlugin) IsItemValid() Return.Error {
	var err Return.Error

	for range Only.Once {
		if p == nil {
			err.SetError("RpcPlugin is nil")
			break
		}

		if p.RpcService.ClientRef == nil {
			err.SetError("RpcService.ClientRef is nil")
			break
		}

		err = p.IsCommonValid()
	}

	return err
}

func (p *RpcPlugin) GetItemData() *Plugin.PluginData {
	return &p.PluginData
}

func (p *RpcPlugin) GetItemHooks() Plugin.HookStore {
	return &p.Dynamic.Hooks
}

func (p *RpcPlugin) SetItemInterface(ref any) Return.Error {
	return p.Dynamic.SetInterface(ref)
}

func (p *RpcPlugin) IsNativePlugin() bool {
	return false
}

func (p *RpcPlugin) IsRpcPlugin() bool {
	return true
}

func (p *RpcPlugin) GetPluginPath() *utils.FilePath {
	return &p.Common.Filename
}

func (p *RpcPlugin) Initialise(args ...any) Return.Error {
	return p.PluginData.Callback(Plugin.CallbackInitialise, &p.PluginData, args...)
}

func (p *RpcPlugin) Execute(args ...any) Return.Error {
	return p.PluginData.Callback(Plugin.CallbackExecute, &p.PluginData, args...)
}

func (p *RpcPlugin) Run(args ...any) Return.Error {
	return p.PluginData.Callback(Plugin.CallbackRun, &p.PluginData, args...)
}

func (p *RpcPlugin) Notify(args ...any) Return.Error {
	return p.PluginData.Callback(Plugin.CallbackNotify, &p.PluginData, args...)
}

// ---------------------------------------------------------------------------------------------------- //

func (p *RpcPlugin) Hooks() *Plugin.HookStruct {
	return &p.Dynamic.Hooks
}

func (p *RpcPlugin) Values() *store.ValueStruct {
	return &p.Dynamic.Values
}

func (p *RpcPlugin) Serve() Return.Error {
	for range Only.Once {
		p.Error = p.Validate()
		if p.Error.IsError() {
			break
		}

		goplugin.Serve(&p.RpcService.ServerConfig)
	}
	return p.Error
}

func (p *RpcPlugin) Validate() Return.Error {
	for range Only.Once {
		p.Error = Return.Ok

		if p.Common.Logger == nil {
			var l utils.Logger
			l, p.Error = utils.NewLogger(p.GetName()+"[rpc]", "")
			if p.Error.IsError() {
				break
			}
			p.Common.Logger = &l
		}

		p.Error = p.Common.Filename.IsValid()
		if p.Error.IsError() {
			p.Common.Filename, p.Error = utils.NewFile(os.Args[0])
			if p.Error.IsError() {
				break
			}
		}

		p.Error = p.Common.Directory.IsValid()
		if p.Error.IsError() {
			p.Common.Directory, p.Error = utils.NewDir(p.Common.Filename.GetDir())
			if p.Error.IsError() {
				break
			}
		}

		if p.Common.Id == "" {
			p.Common.Id = p.Dynamic.Identity.Name
		}
		if p.Services.Identity == "" {
			p.Services.Identity = p.Dynamic.Identity.Name
		}
		if p.Dynamic.Hooks.Identity == "" {
			p.Dynamic.Hooks.Identity = p.Dynamic.Identity.Name
		}

		p.Error = p.SetPluginTypeRpc()
		if p.Error.IsError() {
			break
		}

		p.SetHookPlugin(&p.PluginData)

		p.RpcService.ServerConfig = goplugin.ServeConfig{
			HandshakeConfig: p.Dynamic.HandshakeConfig,
			Plugins:         p.Services.GetAsRpcPluginSet(),
			GRPCServer:      goplugin.DefaultGRPCServer,
			Logger:          p.Common.Logger.Gethclog(),
		}

		p.Error = p.Services.SetRpcService(p.Dynamic.Identity.Name, &RpcPlugin{
			RpcService: p.RpcService,
			PluginData: p.PluginData,
		})
		if p.Error.IsError() {
			break
		}

		p.Common.Configured = true
	}
	return p.Error
}

// IsValid - Validate RpcPlugin structure and set p.configured if true
func (p *RpcPlugin) IsValid() Return.Error {
	var err Return.Error

	for range Only.Once {
		if p == nil {
			err.SetError("native plugin structure is nil")
			break
		}

		if p.RpcService.ClientRef == nil {
			err.SetError("RPC plugin Client is nil")
			break
		}

		err = p.IsCommonValid()
	}

	return err
}

// GetInterface - Get the raw interface.
func (p *RpcPlugin) GetInterface() (any, Return.Error) {
	var raw any
	var err Return.Error

	for range Only.Once {
		if !p.Common.IsCommonConfigured() {
			err.SetError("plugin not configured")
			break
		}

		raw = p.Common.GetRawInterface()
	}

	return raw, err
}

func (p *RpcPlugin) PluginLoad(id string, pluginPath utils.FilePath) Return.Error {

	for range Only.Once {
		p.Error.ReturnClear()
		p.Error.SetPrefix("")

		p.Error = pluginPath.FileExists()
		if p.Error.IsError() {
			break
		}

		// ---------------------------------------------------------------------------------------------------- //
		// Initial setup, before pulling in configured data.
		p.PluginData.Common.Id = id

		var plog utils.Logger
		plog, p.Error = utils.NewLogger(pluginPath.GetName(), "") // @TODO - Config for log file.
		if p.Error.IsError() {
			break
		}
		p.SetLogger(&plog)

		// ---------------------------------------------------------------------------------------------------- //
		// Load the plugin and pull in configured data.
		p.RpcService.ClientConfig = goplugin.ClientConfig{
			HandshakeConfig: Plugin.HandshakeConfig,
			Plugins:         p.PluginData.Services.GetAsRpcPluginSet(),
			Cmd:             exec.Command(pluginPath.GetPath()),
		}
		p.RpcService.ClientConfig.Logger = plog.Gethclog()
		p.SetRpcService(p.Common.Id, &GoPluginMaster{}) // p)

		var e error
		p.RpcService.ClientRef = goplugin.NewClient(&p.RpcService.ClientConfig)
		if p.RpcService.ClientRef == nil {
			p.Error.SetError("[%s]: ERROR: RPC client is nil", p.PluginData.Common.Id)
			break
		}

		p.RpcService.ClientProtocol, e = p.RpcService.ClientRef.Client()
		if e != nil {
			p.Error.SetError("[%s]: ERROR: %s", p.Common.Id, e.Error())
			break
		}
		//goland:noinspection GoDeferInLoop,GoUnhandledErrorResult
		defer p.RpcService.ClientProtocol.Close()

		e = p.RpcService.ClientProtocol.Ping()
		if e != nil {
			p.Error.SetError("[%s]: ERROR: %s\n", id, e.Error())
			break
		}

		var raw any
		raw, e = p.RpcService.ClientProtocol.Dispense(p.Common.Id)
		if e != nil {
			p.Error.SetError("[%s]: ERROR: %s\n", p.Common.Id, e.Error())
			break
		}

		tn := utils.GetTypeName(raw)
		if tn != "*GoPlugLoader.RpcPluginClient" {
			p.Error.SetError("[%s]: ERROR: Invalid type - expecting '*RpcPluginClient', got '%s'", p.Common.Id, tn)
			break
		}

		impl := raw.(*RpcPluginClient)
		p.PluginData.Dynamic = impl.GetData()
		if impl.Error.IsError() {
			p.Error = impl.Error
			break
		}
		p.PluginData.Dynamic.Identity.Print()

		var identity Plugin.Identity
		identity = impl.Identify()
		if p.Error.IsError() {
			p.Error.SetError("GoPluginIdentity is not globally defined: %s", p.Error)
			break
		}
		p.SetIdentity(&identity) // This will be replaced with a full get of GoPluginNativeInterface

		p.SetFilename(pluginPath)
		p.SetHookPlugin(&p.PluginData)
		p.SetPluginTypeRpc() // Even if the config doesn't set it, do it here.
		// l.SetPluginIdentity(l.Common.Id)
		p.SetRpcService(p.Common.Id, &GoPluginMaster{}) // p)
		p.SetRawInterface(p)
		p.SetStructName(identity)
		log.Printf("[%s]: Name:%s Path: %s\n",
			p.Common.Id, p.Common.Filename.GetName(), p.Common.Filename.GetPath())

		p.Error = p.Callback(Plugin.CallbackInitialise, p)
		if p.Error.IsError() {
			break
		}
	}

	return p.Error
}

func (p *RpcPlugin) PluginUnload() Return.Error {
	for range Only.Once {
		p.Error.ReturnClear()
		p.Error.SetPrefix("")
	}

	return p.Error
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of RPC interface structure

func (p *RpcPlugin) Server(_ *goplugin.MuxBroker) (any, error) {
	p.Dynamic.Error = Return.Ok
	impl := Plugin.PluginData{
		Common:   p.Common,
		Services: p.Services,
		Dynamic:  p.Dynamic,
		Error:    p.Error,
	}

	gob.Register(Plugin.PluginData{})
	gob.Register(store.ValueStruct{})
	gob.Register(RpcPlugin{})
	return &RpcPluginServer{Impl: &impl}, nil
}

func (p *RpcPlugin) Client(_ *goplugin.MuxBroker, c *rpc.Client) (any, error) {
	p.Dynamic.Error = Return.Ok

	gob.Register(Plugin.PluginData{})
	gob.Register(store.ValueStruct{})
	gob.Register(RpcPlugin{})
	return &RpcPluginClient{Client: c}, nil
}

//
// RpcService
// ---------------------------------------------------------------------------------------------------- //
type RpcService struct {
	ServerConfig   goplugin.ServeConfig
	ClientConfig   goplugin.ClientConfig
	ClientRef      *goplugin.Client
	ClientProtocol goplugin.ClientProtocol
}

// NewRpcService - Create a new instance of this structure.
func NewRpcService() RpcService {
	return RpcService{
		ServerConfig:   goplugin.ServeConfig{},
		ClientConfig:   goplugin.ClientConfig{},
		ClientRef:      nil,
		ClientProtocol: nil,
	}
}
