package GoPlugLoader

import (
	"encoding/gob"
	"net/rpc"

	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
	"github.com/MickMake/GoPlug/utils/store"
)

//
// GoPluginMaster
// ---------------------------------------------------------------------------------------------------- //
// 1. The root - starts it all off.
type GoPluginMaster struct {
	Plugin.Plugin
}

func (g *GoPluginMaster) Server(b *goplugin.MuxBroker) (any, error) {
	utils.DEBUG()
	ret := RpcPlugin{
		Plugin: g.Plugin, // Make local plugin data accessible to client.
	}
	gob.Register(store.ValueStruct{})
	gob.Register(RpcPlugin{})
	return &ret, nil
}

func (g *GoPluginMaster) Client(b *goplugin.MuxBroker, c *rpc.Client) (any, error) {
	utils.DEBUG()
	gob.Register(store.ValueStruct{})
	gob.Register(RpcPlugin{})
	return &RpcPluginClient{Client: c}, nil
}

//
// RpcPluginClient
// ---------------------------------------------------------------------------------------------------- //
// 2. Client sends RPC request.
type RpcPluginClient struct {
	Client *rpc.Client

	Error Return.Error
}

func (g *RpcPluginClient) GetData() Plugin.DynamicData {
	utils.DEBUG()
	g.Error = Return.Ok
	var resp Plugin.DynamicData
	err := g.Client.Call("Plugin.GetData", new(any), &resp)
	if err != nil {
		g.Error.SetError(err)
	}
	return resp
}

func (g *RpcPluginClient) Identify() Plugin.Identity {
	utils.DEBUG()
	g.Error = Return.Ok
	var resp Plugin.Identity
	err := g.Client.Call("Plugin.Identify", new(any), &resp)
	if err != nil {
		g.Error.SetError(err)
	}
	return resp
}

func (g *RpcPluginClient) IdentifyString() string {
	utils.DEBUG()
	g.Error = Return.Ok
	var resp string
	err := g.Client.Call("Plugin.IdentifyString", new(any), &resp)
	if err != nil {
		g.Error.SetError(err)
	}

	return resp
}

func (g *RpcPluginClient) CallHook(name string, args ...any) (Plugin.HookResponse, Return.Error) {
	utils.DEBUG()
	g.Error = Return.Ok
	var resp Plugin.HookResponse
	err := g.Client.Call("Plugin.CallHook", &Plugin.HookCallArgs{Name: name, Args: args}, &resp)
	if err != nil {
		g.Error.SetError(err)
	}
	return resp, g.Error
}

//
// RpcPluginServerInterface
// ---------------------------------------------------------------------------------------------------- //
// 3. Server responds to RPC request.
type RpcPluginServerInterface interface {
	GetData() Plugin.DynamicData
	Identify() Plugin.Identity
	IdentifyString() string
	CallHook(name string, args ...any) (Plugin.HookResponse, Return.Error)
}

//
// RpcPluginServer
// ---------------------------------------------------------------------------------------------------- //
type RpcPluginServer struct {
	Impl RpcPluginServerInterface

	Error Return.Error
}

func (s *RpcPluginServer) GetData(_ any, resp *Plugin.DynamicData) error {
	utils.DEBUG()
	s.Error = Return.Ok
	*resp = s.Impl.GetData()
	return s.Error.GetError()
}

func (s *RpcPluginServer) Identify(_ any, resp *Plugin.Identity) error {
	utils.DEBUG()
	s.Error = Return.Ok
	*resp = s.Impl.Identify()
	return s.Error.GetError()
}

func (s *RpcPluginServer) IdentifyString(_ any, resp *string) error {
	utils.DEBUG()
	s.Error = Return.Ok
	*resp = s.Impl.IdentifyString()
	return s.Error.GetError()
}

func (s *RpcPluginServer) CallHook(args Plugin.HookCallArgs, resp *Plugin.HookResponse) error {
	utils.DEBUG()
	s.Error = Return.Ok
	*resp, s.Error = s.Impl.CallHook(args.Name, args.Args...)
	return s.Error.GetError()
}
