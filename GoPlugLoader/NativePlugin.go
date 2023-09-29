package GoPlugLoader

import (
	"context"
	"fmt"
	"log"
	"os"
	sysPlugin "plugin"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/store"

	// jsonV2 "github.com/go-json-experiment/json"
	"github.com/MickMake/GoPlug/utils/Return"
)

// NewNativePluginInterface - Create a new instance of this interface.
//goland:noinspection GoUnusedExportedFunction
func NewNativePluginInterface() PluginItemInterface {
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
	Plugin.PluginData
}

// ---------------------------------------------------------------------------------------------------- //

// IsItemValid - Validate NativePlugin structure and set p.configured if true
func (p *NativePlugin) IsItemValid() Return.Error {
	var err Return.Error

	for range Only.Once {
		if p == nil {
			err.SetError("NativePlugin is nil")
			break
		}

		if p.context == nil {
			err.SetError("NativePlugin.Context is nil")
			break
		}

		err = p.IsCommonValid()
	}

	return err
}

func (p *NativePlugin) GetItemData() *Plugin.PluginData {
	return &p.PluginData
}

func (p *NativePlugin) GetItemHooks() Plugin.HookStore {
	return &p.Dynamic.Hooks
}

func (p *NativePlugin) SetItemInterface(ref any) Return.Error {
	return p.Dynamic.SetInterface(ref)
}

func (p *NativePlugin) IsNativePlugin() bool {
	return true
}

func (p *NativePlugin) IsRpcPlugin() bool {
	return false
}

func (p *NativePlugin) GetPluginPath() *utils.FilePath {
	return &p.Common.Filename
}

func (p *NativePlugin) Initialise(args ...any) Return.Error {
	return p.PluginData.Callback(Plugin.CallbackInitialise, &p.PluginData, args...)
}

func (p *NativePlugin) Execute(args ...any) Return.Error {
	return p.PluginData.Callback(Plugin.CallbackExecute, &p.PluginData, args...)
}

func (p *NativePlugin) Run(args ...any) Return.Error {
	return p.PluginData.Callback(Plugin.CallbackRun, &p.PluginData, args...)
}

func (p *NativePlugin) Notify(args ...any) Return.Error {
	return p.PluginData.Callback(Plugin.CallbackNotify, &p.PluginData, args...)
}

// ---------------------------------------------------------------------------------------------------- //

// NewNativePlugin - Create a new instance of this structure.
func NewNativePlugin() *NativePlugin {
	return &NativePlugin{
		context:    context.Background(),
		Service:    NewNativeService(),
		PluginData: *Plugin.NewPlugin(),
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

		p.Error = p.SetPluginTypeNative()
		if p.Error.IsError() {
			break
		}

		p.SetHookPlugin(&p.PluginData)

		// if p.Services.CountServices() == 0 {
		// 	p.Error.SetError("No plugin maps defined!")
		// 	break
		// }

		p.Error = p.Services.SetNativeService(p.Dynamic.Identity.Name, sysPlugin.Plugin{})
		if p.Error.IsError() {
			break
		}

		p.Common.Configured = true
	}
	return p.Error
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

		raw = p.Common.GetRawInterface()
	}

	return raw, err
}

// ---------------------------------------------------------------------------------------------------- //

func (p *NativePlugin) PluginLoad(id string, pluginPath utils.FilePath) Return.Error {
	for range Only.Once {
		p.Error.ReturnClear()
		p.Error.SetPrefix("")
		pluginPath.ShortenPaths()
		p.PluginData.Common.Id = id

		p.Error = pluginPath.FileExists()
		if p.Error.IsError() {
			break
		}

		// ---------------------------------------------------------------------------------------------------- //
		// Load the plugin and pull in configured data.
		p.Error = p.Service.Open(pluginPath)
		if p.Error.IsError() {
			break
		}

		var identity *Plugin.Identity
		identity, p.Error = p.Service.GetIdentity()
		if p.Error.IsError() {
			p.Error.SetError("GoPluginIdentity is not globally defined: %s", p.Error)
			break
		}
		p.SetIdentity(identity) // This will be replaced with a full get of GoPluginInterface

		// ---------------------------------------------------------------------------------------------------- //
		// Reconfigure and check data.
		log.Println("Looking for GoPluginItem")
		var GoPluginInterface *PluginItem
		GoPluginInterface, p.Error = p.Service.GetPluginItem()
		if p.Error.IsError() {
			break
		}
		if GoPluginInterface.Pluggable == nil {
			GoPluginInterface.Pluggable = CreatePluginItem(Plugin.NativePluginType, identity)
		}
		if GoPluginInterface != nil {
			native := (*GoPluginInterface).GetItemData()
			if native == nil {
				p.Error.SetError("GoPluginInterface is defined, but nil!")
				break
			}

			p.PluginData = *native
		}

		p.SetFilename(pluginPath)
		p.SetHookPlugin(&p.PluginData)
		p.SetPluginTypeNative() // Even if the config doesn't set it, do it here.
		p.SetNativeService(p.Common.Id, *p.Service.Object)
		p.SetRawInterface(p.Service.Symbol)
		p.SetStructName(*identity)
		log.Printf("[%s]: Name:%s Path: %s\n",
			p.Common.Id, p.Common.Filename.GetName(), p.Common.Filename.String())

		p.Error = p.Callback(Plugin.CallbackInitialise, p)
		if p.Error.IsError() {
			break
		}
	}

	return p.Error
}

func (p *NativePlugin) PluginUnload() Return.Error {
	for range Only.Once {
		p.Error.ReturnClear()
		p.Error.SetPrefix("")
	}

	return p.Error
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
		re := regexp.MustCompile(`map\[(.*)]`)
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

// GetPluginItem - Gets the PluginItemInterface symbol using several symbol names.
// Will scan exported symbols if none specified.
func (ns *NativeService) GetPluginItem() (*PluginItem, Return.Error) {
	var rpi *PluginItem
	var err Return.Error

	for range Only.Once {
		for name, symType := range ns.Symbols {
			if symType != "*GoPlugLoader.PluginItem" {
				continue
			}

			var sym any
			sym, err = ns.Lookup(name)
			if err.IsError() {
				continue
			}

			var ok bool
			rpi, ok = sym.(*PluginItem)
			if !ok {
				err.SetError("Symbol '%s' had type '%s', was expecting '*GoPlugLoader.%s'",
					name, utils.GetTypeName(sym), Plugin.GoPluginItem)
				continue
			}

			if (*rpi).Error.IsError() {
				err = (*rpi).Error
				break
			}

			if (*rpi).Pluggable == nil {
				// This is actually OK - We will set it up later.
				break
			}

			err = (*rpi).IsItemValid()
			if err.IsError() {
				err.SetError("%s is not valid: %s", Plugin.GoPluginItem, err)
				continue
			}

			if !(*rpi).IsNativePlugin() {
				continue
			}

			break
		}

		if rpi == nil {
			err.SetError("%s is not defined as global in plugin", Plugin.GoPluginItem)
		}
	}

	return rpi, err
}
