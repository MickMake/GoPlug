package GoPlugLoader

import (
	"encoding/gob"
	"net/rpc"
	"os"

	"github.com/MickMake/GoUnify/Only"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
	"github.com/MickMake/GoPlug/utils/store"
)

//
// RpcPluginInterface
// ---------------------------------------------------------------------------------------------------- //
// 4. Plugin setup.
type RpcPluginInterface interface {
	// NewRpcPlugin - Create a new instance of this plugin.
	NewRpcPlugin() Return.Error
	GetRpcPlugin() *RpcPlugin

	// RegisterStructure(ref any) Return.Error
	// SetIdentity(identity *Plugin.Identity) Return.Error
	SetPluginType(types Plugin.Types) Return.Error
	SetInterface(ref any) Return.Error
	SetHandshakeConfig(goplugin.HandshakeConfig) Return.Error

	Hooks() *Plugin.HookStruct
	Values() *store.ValueStruct

	Validate() Return.Error
	Serve() Return.Error

	// IsValid - Validate RpcPluginInterface interface
	IsValid() Return.Error

	Plugin.Interface
}

// NewRpcPluginInterface - Create a new instance of this interface.
func NewRpcPluginInterface() RpcPluginInterface {
	ret := NewRpcPlugin()
	ret.SetPluginType(Plugin.RpcPluginType)
	return ret
}

//
// RpcPlugin
// ---------------------------------------------------------------------------------------------------- //
type RpcPlugin struct {
	Service RpcService

	Plugin.Plugin
}

// NewRpcPlugin - Create a new instance of this structure.
func NewRpcPlugin() *RpcPlugin {
	return &RpcPlugin{
		Service: NewRpcService(),
		Plugin:  *Plugin.NewPlugin(),
	}
}

func (p *RpcPlugin) NewRpcPlugin() Return.Error {
	*p = *NewRpcPlugin()
	return Return.Ok
}

func (p *RpcPlugin) GetRpcPlugin() *RpcPlugin {
	return p
}

func (p *RpcPlugin) Hooks() *Plugin.HookStruct {
	return &p.Dynamic.Hooks
}

func (p *RpcPlugin) Values() *store.ValueStruct {
	return &p.Dynamic.Values
}

func (p *RpcPlugin) Serve() Return.Error {
	for range Only.Once {
		if p.Service.ClientRef == nil {
			p.Error.SetError("ClientRef is nil")
			break
		}

		if p.Service.ClientProtocol == nil {
			p.Error.SetError("ClientProtocol is nil")
			break
		}

		goplugin.Serve(&p.Service.ServerConfig)
	}
	return p.Dynamic.Error
}

func (p *RpcPlugin) Validate() Return.Error {
	for range Only.Once {
		p.Error = Return.Ok

		p.Error = p.Services.SetRpcService(p.Dynamic.Identity.Name, &RpcPlugin{
			Service: p.Service,
			Plugin:  *Plugin.NewPlugin(),
		})
		if p.Error.IsError() {
			break
		}

		if p.Services.CountServices() == 0 {
			p.Dynamic.Error.SetError("No plugin maps defined!")
			break
		}

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

		p.Service.ServerConfig = goplugin.ServeConfig{
			HandshakeConfig: p.Dynamic.HandshakeConfig,
			Plugins:         p.Services.GetAsRpcPluginSet(),
			GRPCServer:      goplugin.DefaultGRPCServer,
			Logger:          p.Common.Logger.Gethclog(),
		}
	}
	return p.Dynamic.Error
}

// IsValid - Validate NativePlugin structure and set p.configured if true
func (p *RpcPlugin) IsValid() Return.Error {
	var err Return.Error

	for range Only.Once {
		if p == nil {
			err.SetError("native plugin structure is nil")
			break
		}

		if p.Service.ClientRef == nil {
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

		raw = p.Common.RawInterface
	}

	return raw, err
}

//
// ---------------------------------------------------------------------------------------------------- //
// Mirror methods of RPC interface structure

func (p *RpcPlugin) Server(_ *goplugin.MuxBroker) (any, error) {
	p.Dynamic.Error = Return.Ok
	impl := Plugin.Plugin{
		Common:   p.Common,
		Services: p.Services,
		Dynamic:  p.Dynamic,
		Error:    p.Error,
	}
	gob.Register(store.ValueStruct{})
	gob.Register(RpcPlugin{})
	return &RpcPluginServer{Impl: &impl}, nil
}
func (p *RpcPlugin) Client(_ *goplugin.MuxBroker, c *rpc.Client) (any, error) {
	p.Dynamic.Error = Return.Ok
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
