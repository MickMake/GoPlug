package GoPlugLoader

import (
	"context"
	"fmt"
	"os"
	sysPlugin "plugin"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/MickMake/GoUnify/Only"
	goplugin "github.com/hashicorp/go-plugin"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/store"

	// jsonV2 "github.com/go-json-experiment/json"
	"github.com/MickMake/GoPlug/utils/Return"
)

//
// NativePluginInterface
// ---------------------------------------------------------------------------------------------------- //
type NativePluginInterface interface {
	// NewNativePlugin - Create a new instance of this plugin.
	NewNativePlugin() Return.Error
	GetNativePlugin() *NativePlugin

	SetPluginType(types Plugin.Types) Return.Error
	SetInterface(ref any) Return.Error
	SetHandshakeConfig(goplugin.HandshakeConfig) Return.Error

	Hooks() *Plugin.HookStruct
	Values() *store.ValueStruct

	Validate() Return.Error
	Serve() Return.Error

	// IsValid - Validate NativePluginInterface interface
	IsValid() Return.Error

	context.Context
	Plugin.Interface
}

// NewNativePluginInterface - Create a new instance of this interface.
func NewNativePluginInterface() NativePluginInterface {
	ret := NewNativePlugin()
	ret.SetPluginType(Plugin.NativePluginType)
	return ret
}

//
// NativePlugin implemented as default plugin context
// ---------------------------------------------------------------------------------------------------- //
type NativePlugin struct {
	context context.Context
	Service NativeService

	Plugin.Plugin
}

// NewNativePlugin - Create a new instance of this structure.
func NewNativePlugin() *NativePlugin {
	return &NativePlugin{
		context: context.Background(),
		Service: NewNativeService(),
		Plugin:  *Plugin.NewPlugin(),
	}
}

func (p *NativePlugin) NewNativePlugin() Return.Error {
	*p = *NewNativePlugin()
	return Return.Ok
}

func (p *NativePlugin) GetNativePlugin() *NativePlugin {
	return p
}

func (p *NativePlugin) Hooks() *Plugin.HookStruct {
	return &p.Dynamic.Hooks
}

func (p *NativePlugin) Values() *store.ValueStruct {
	return &p.Dynamic.Values
}

func (p *NativePlugin) Serve() Return.Error {
	for range Only.Once {
		fmt.Println("goplugin.Serve(&p.Native.ServerConfig)")
	}
	return p.Error
}

func (p *NativePlugin) Validate() Return.Error {
	for range Only.Once {
		p.Error = Return.Ok

		p.Error = p.Services.SetNativeService(p.Dynamic.Identity.Name, sysPlugin.Plugin{})
		// p.Error = p.Services.SetNativeService(p.Data.Dynamic.Identity.Name, &NativePlugin{
		// 	context: p.context,
		// 	Data:    p.Data,
		// })
		if p.Error.IsError() {
			break
		}

		if p.Services.CountServices() == 0 {
			p.Error.SetError("No plugin maps defined!")
			break
		}

		if p.Common.Logger == nil {
			var l utils.Logger
			l, p.Error = utils.NewLogger(p.GetName()+"[native]", "")
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
	}
	return p.Error
}

// IsValid - Validate NativePlugin structure and set p.configured if true
func (p *NativePlugin) IsValid() Return.Error {
	var err Return.Error

	for range Only.Once {
		if p == nil {
			err.SetError("native plugin structure is nil")
			break
		}

		if p.context == nil {
			err.SetError("native plugin Context is nil")
			break
		}

		err = p.IsCommonValid()
	}

	return err
}

// GetInterface - Get the raw interface.
func (p *NativePlugin) GetInterface() (any, Return.Error) {
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
// Mirror functions of context.Context interface structure

func (p *NativePlugin) Deadline() (deadline time.Time, ok bool) {
	return p.context.Deadline()
}
func (p *NativePlugin) Done() <-chan struct{} {
	return p.context.Done()
}
func (p *NativePlugin) Err() error {
	return p.context.Err()
}
func (p *NativePlugin) Value(key any) any {
	return p.context.Value(key)
}

//
// NativeService
// ---------------------------------------------------------------------------------------------------- //
type NativeService struct {
	pluginPath utils.FilePath
	Object     *sysPlugin.Plugin
	Symbol     any
	Symbols    map[string]string
}

// NewNativeService - Create a new instance of this structure.
func NewNativeService() NativeService {
	return NativeService{
		pluginPath: utils.FilePath{},
		Object:     nil,
		Symbol:     nil,
		Symbols:    make(map[string]string),
	}
}

// Open - Opens a Go plugin file.
func (ns *NativeService) Open(pluginPath utils.FilePath) Return.Error {
	var err Return.Error

	for range Only.Once {
		s, e := sysPlugin.Open(pluginPath.GetPath())
		if e != nil {
			err.SetError(e)
			break
		}
		ns.Object = s
		ns.pluginPath = pluginPath

		err = ns.Scan()
	}

	return err
}

// ListExported - List all exported symbols.
func (ns *NativeService) ListExported() []string {
	var ret []string
	for symbol := range ns.Symbols {
		ret = append(ret, symbol)
	}
	sort.Strings(ret)
	return ret
}

// Scan - Scans all exported symbols and finds type.
func (ns *NativeService) Scan() Return.Error {
	var err Return.Error

	for range Only.Once {
		// Find exported symbols.
		re := regexp.MustCompile("map\\[(.*)\\]")
		str := fmt.Sprintf("%v", ns.Object)
		sa := re.FindStringSubmatch(str)
		if len(sa) < 2 {
			break
		}
		sa2 := strings.Split(sa[1], " ")
		ns.Symbols = make(map[string]string)
		for _, symbol := range sa2 {
			sa3 := strings.Split(symbol, ":")
			if len(sa3) < 2 {
				continue
			}
			ns.Symbols[sa3[0]] = "Unknown"
		}

		// Find types of exported symbols.
		for symbol := range ns.Symbols {
			var sym any
			sym, err = ns.Lookup(symbol)
			if err.IsError() {
				break
			}
			ns.Symbols[symbol] = utils.GetTypeName(sym)
		}
	}

	return err
}

// String - Stringer.
func (ns *NativeService) String() string {
	var ret string
	ret += fmt.Sprintf("Plugin file: %s\n", ns.pluginPath.GetName())
	ret += fmt.Sprintf("- Current symbol type: %s\n", utils.GetTypeName(ns.Symbol))
	ret += fmt.Sprintf("- Exported symbols:\n")

	for _, symbol := range ns.ListExported() {
		ret += fmt.Sprintf("\t- %s\t- %s\n", symbol, ns.Symbols[symbol])
	}
	return ret
}

// Print - .
func (ns *NativeService) Print() {
	fmt.Print(ns.String())
}

// Lookup - Looks up a symbol.
func (ns *NativeService) Lookup(symbol string) (any, Return.Error) {
	var err Return.Error

	for range Only.Once {
		if ns.Object == nil {
			err.SetError("Native plugin file not loaded")
			break
		}

		var e error
		ns.Symbol, e = ns.Object.Lookup(symbol)
		if e != nil {
			err.SetError(e)
			break
		}
	}

	return ns.Symbol, err
}

// LookupType - Looks up a symbol, ensuring it's of a specific type.
func (ns *NativeService) LookupType(symbol string, name string) (any, Return.Error) {
	var err Return.Error

	for range Only.Once {
		if ns.Object == nil {
			err.SetError("Native plugin file not loaded")
			break
		}

		var e error
		ns.Symbol, e = ns.Object.Lookup(symbol)
		if e != nil {
			err.SetError(e)
			break
		}

		n := ns.TypeName()
		if n != name {
			err.SetError("Expected symbol '%s' in Native plugin, but got '%s'", name, ns.TypeName())
			break
		}
	}

	return ns.Symbol, err
}

// TypeName - Looks up a symbol.
func (ns *NativeService) TypeName() string {
	return utils.GetTypeName(ns.Symbol)
}

// IsType - Looks up a symbol.
func (ns *NativeService) IsType(name string) bool {
	return utils.IsTypeOfName(ns.Symbol, name)
}

// GetIdentity - Gets the Identity symbol using several symbol names.
// Will scan exported symbols if none specified.
func (ns *NativeService) GetIdentity(lookups ...string) (*Plugin.Identity, Return.Error) {
	var identity *Plugin.Identity
	var err Return.Error

	for range Only.Once {
		if len(lookups) == 0 {
			lookups = []string{Plugin.GoPluginIdentity}
			lookups = append(lookups, ns.ListExported()...)
		}

		var f *Plugin.Identity
		find := utils.GetTypeName(f)
		for _, lookup := range lookups {
			var sym any
			sym, err = ns.LookupType(lookup, find) // "*Plugin.Identity")
			if err.IsError() {
				err.SetError("%s is not globally defined: %s", Plugin.GoPluginIdentity, err)
				break
			}

			var ok bool
			identity, ok = sym.(*Plugin.Identity)
			if !ok {
				err.SetError("Symbol '%s' had type '%s', was expecting '*NativePlugin.Identity'", lookup, utils.GetTypeName(sym))
				break
			}

			err = identity.IsValid()
			if err.IsError() {
				err.SetError("%s is not globally defined: %s", Plugin.GoPluginIdentity, err)
			}
			break
		}

		if identity == nil {
			err.SetError("%s is not defined as global in plugin", find)
		}
	}

	return identity, err
}

// GetRpcPluginInterface - Gets the RpcPluginInterface symbol using several symbol names.
// Will scan exported symbols if none specified.
func (ns *NativeService) GetRpcPluginInterface(lookups ...string) (*RpcPluginInterface, Return.Error) {
	var rpi *RpcPluginInterface
	var err Return.Error

	for range Only.Once {
		if len(lookups) == 0 {
			lookups = []string{Plugin.GoPluginRpcInterface}
			lookups = append(lookups, ns.ListExported()...)
		}

		// var f RpcPluginInterface
		// f = NewRpcPluginInterface()
		// find := utils.GetTypeName(f)
		// find = strings.ReplaceAll(find, "*", "*GoPlugLoader.")
		for _, lookup := range lookups {
			var sym any
			sym, err = ns.LookupType(lookup, "*GoPlugLoader.RpcPluginInterface")
			if err.IsError() {
				continue
			}

			var ok bool
			rpi, ok = sym.(*RpcPluginInterface)
			if !ok {
				err.SetError("Symbol '%s' had type '%s', was expecting '*GoPlugLoader.%s'",
					lookup, utils.GetTypeName(sym), Plugin.GoPluginRpcInterface)
				continue
			}

			err = (*rpi).IsValid()
			if err.IsError() {
				err.SetError("%s is not valid: %s", Plugin.GoPluginRpcInterface, err)
			}
			break
		}

		if rpi == nil {
			err.SetError("%s is not defined as global in plugin", Plugin.GoPluginRpcInterface)
		}
	}

	return rpi, err
}

// GetNativePluginInterface - Gets the NativePluginInterface symbol using several symbol names.
// Will scan exported symbols if none specified.
func (ns *NativeService) GetNativePluginInterface(lookups ...string) (*NativePluginInterface, Return.Error) {
	var rpi *NativePluginInterface
	var err Return.Error

	for range Only.Once {
		if len(lookups) == 0 {
			lookups = []string{Plugin.GoPluginNativeInterface}
			lookups = append(lookups, ns.ListExported()...)
		}

		// var f NativePluginInterface
		// find := utils.GetTypeName(f)
		// find = strings.ReplaceAll(find, "*", "*GoPlugLoader.")
		for _, lookup := range lookups {
			var sym any
			sym, err = ns.LookupType(lookup, "*GoPlugLoader.NativePluginInterface")
			if err.IsError() {
				continue
			}

			var ok bool
			rpi, ok = sym.(*NativePluginInterface)
			if !ok {
				err.SetError("Symbol '%s' had type '%s', was expecting '*GoPlugLoader.%s'",
					lookup, utils.GetTypeName(sym), Plugin.GoPluginNativeInterface)
				continue
			}

			err = (*rpi).IsValid()
			if err.IsError() {
				err.SetError("%s is not valid: %s", Plugin.GoPluginNativeInterface, err)
			}
			break
		}

		if rpi == nil {
			err.SetError("%s is not defined as global in plugin", Plugin.GoPluginNativeInterface)
		}
	}

	return rpi, err
}
