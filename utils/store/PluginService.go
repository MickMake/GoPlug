package store

import (
	"fmt"
	sysPlugin "plugin"
	"strings"

	"github.com/MickMake/GoUnify/Only"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

// ---------------------------------------------------------------------------------------------------- //
// PluginServiceStore interface and methods

//
// PluginServiceInterface - Getter/Setter for string map of interfaces{}
// ---------------------------------------------------------------------------------------------------- //
type PluginServiceInterface interface {
	// NewPluginService - Set up the FuncMap structure.
	NewPluginService() Return.Error

	GetPluginServiceReference() *PluginServiceStruct

	GetPluginIdentity() string
	SetPluginIdentity(identity string) Return.Error

	// ServiceExists - Check if a key exists.
	ServiceExists(name string) bool

	// ServiceNotExists - Inverse of Exists()
	ServiceNotExists(name string) bool

	// GetRpcService - Get a key's value.
	GetRpcService(name string) goplugin.Plugin

	// SetRpcService - Set a key value pair.
	SetRpcService(name string, value goplugin.Plugin) Return.Error

	GetAsRpcPluginSet() goplugin.PluginSet

	// GetNativeService - Get a key's value.
	GetNativeService(name string) sysPlugin.Plugin

	// SetNativeService - Set a key value pair.
	SetNativeService(name string, value sysPlugin.Plugin) Return.Error

	GetAsNativePluginSet() NativeServiceMap

	ValidateService() Return.Error

	// CountServices - Return the number of entries.
	CountServices() int

	// ListServices - Get PluginServiceStruct.
	ListServices() RpcServiceMap

	// PrintServices - Get PluginServiceStruct.
	PrintServices()

	// String - Stringer method.
	String() string
}

// NewPluginServiceStore - Create a PluginServiceInterface interface structure instance.
//goland:noinspection GoUnusedExportedFunction
func NewPluginServiceStore() PluginServiceInterface {
	ret := NewPluginServiceStruct()
	return &ret
}

//
// PluginServiceStruct
// ---------------------------------------------------------------------------------------------------- //
type PluginServiceStruct struct {
	Identity       string `json:"identity,omitempty"`
	rpcServices    RpcServiceMap
	nativeServices NativeServiceMap
	Master         bool         `json:"master,omitempty"`
	Error          Return.Error `json:"error"`
}

// NewPluginServiceStruct - Create a PluginServiceStruct structure instance.
func NewPluginServiceStruct() PluginServiceStruct {
	return PluginServiceStruct{
		Identity:       "",
		rpcServices:    make(RpcServiceMap),
		nativeServices: make(NativeServiceMap),
		Master:         false,
		Error:          Return.New(),
	}
}

// NewPluginService - Create a PluginServiceInterface interface structure instance.
func (p *PluginServiceStruct) NewPluginService() Return.Error {
	*p = NewPluginServiceStruct()
	return Return.Ok
}

func (p *PluginServiceStruct) GetPluginServiceReference() *PluginServiceStruct {
	p.Error = Return.Ok
	return p
}

func (p *PluginServiceStruct) SetPluginIdentity(identity string) Return.Error {
	p.Error = Return.Ok
	p.Identity = identity
	return Return.Ok
}

func (p *PluginServiceStruct) GetPluginIdentity() string {
	p.Error = Return.Ok
	return p.Identity
}

// ServiceExists - Check if a key exists.
func (p *PluginServiceStruct) ServiceExists(name string) bool {
	_, err := p.nativeServices.Get(name)
	if !err.IsError() {
		return true
	}
	_, err = p.rpcServices.Get(name)
	if !err.IsError() {
		return true
	}
	return false
}

// ServiceNotExists - Inverse of Exists()
func (p *PluginServiceStruct) ServiceNotExists(name string) bool {
	return !p.ServiceExists(name)
}

// GetAsRpcPluginSet - Get PluginServiceInterface as goplugin.PluginSet.
func (p *PluginServiceStruct) GetAsRpcPluginSet() goplugin.PluginSet {
	p.Error = Return.Ok
	return goplugin.PluginSet(p.rpcServices)
}

// GetRpcService - Get a key's value.
func (p *PluginServiceStruct) GetRpcService(name string) goplugin.Plugin {
	var service goplugin.Plugin
	service, p.Error = p.rpcServices.Get(name)
	return service
}

// SetRpcService - Set a key value pair.
func (p *PluginServiceStruct) SetRpcService(name string, value goplugin.Plugin) Return.Error {
	p.Error = Return.Ok
	name = strings.TrimSpace(name)
	if p.rpcServices == nil {
		p.rpcServices = make(RpcServiceMap)
	}

	p.rpcServices[name] = value
	return p.Error
}

// GetAsNativePluginSet - Get PluginServiceInterface as goplugin.PluginSet.
func (p *PluginServiceStruct) GetAsNativePluginSet() NativeServiceMap {
	p.Error = Return.Ok
	return p.nativeServices
}

// GetNativeService - Get a key's value.
func (p *PluginServiceStruct) GetNativeService(name string) sysPlugin.Plugin {
	var service sysPlugin.Plugin
	service, p.Error = p.nativeServices.Get(name)
	return service
}

// SetNativeService - Set a key value pair.
func (p *PluginServiceStruct) SetNativeService(name string, value sysPlugin.Plugin) Return.Error {
	p.Error = Return.Ok
	name = strings.TrimSpace(name)
	if p.nativeServices == nil {
		p.nativeServices = make(NativeServiceMap)
	}

	p.nativeServices[name] = value
	return p.Error
}

// CountServices - Return the number of entries.
func (p *PluginServiceStruct) CountServices() int {
	p.Error = Return.Ok
	return len(p.nativeServices) + len(p.rpcServices)
}

func (p *PluginServiceStruct) ListServices() RpcServiceMap {
	p.Error = Return.Ok
	return p.rpcServices
}

func (p *PluginServiceStruct) PrintServices() {
	p.Error = Return.Ok
	fmt.Print(p.String())
}

// String - Stringer interface.
func (p PluginServiceStruct) String() string {
	var ret string
	ret += fmt.Sprintf("# Available plugins from identity '%s'\n", p.Identity)
	for name, plug := range p.rpcServices {
		ret += fmt.Sprintf("\t[%s]: %s\n", name, plug)
	}
	return ret
}

// ValidateService - .
func (p *PluginServiceStruct) ValidateService() Return.Error {
	for range Only.Once {
		p.Error = Return.Ok
		name := utils.GetCallerFunctionName(1)

		p.GetNativeService(name)
		if !p.Error.IsError() {
			break
		}

		p.GetRpcService(name)
		if !p.Error.IsError() {
			break
		}

		p.Error.SetError("plug function mismatch: looking for %s", name)
	}
	return p.Error
}

//
// RpcServiceMap
// ---------------------------------------------------------------------------------------------------- //
type RpcServiceMap map[string]goplugin.Plugin

func (m *RpcServiceMap) Get(name string) (goplugin.Plugin, Return.Error) {
	name = strings.TrimSpace(name)
	if value, ok := (*m)[name]; ok {
		return value, Return.Ok
	}
	return nil, Return.NewError("hook '%s' not found", name)
}

//
// NativeServiceMap
// ---------------------------------------------------------------------------------------------------- //
type NativeServiceMap map[string]sysPlugin.Plugin

func (m *NativeServiceMap) Get(name string) (sysPlugin.Plugin, Return.Error) {
	name = strings.TrimSpace(name)
	if value, ok := (*m)[name]; ok {
		return value, Return.Ok
	}
	return sysPlugin.Plugin{}, Return.NewError("hook '%s' not found", name)
}
