package Plugin

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/utils/Return"
)

// ---------------------------------------------------------------------------------------------------- //

const (
	PackageName             = "GoPlugin"
	GoPluginIdentity        = "GoPluginIdentity"
	GoPluginNativeInterface = "GoPluginNativeInterface"
	GoPluginRpcInterface    = "GoPluginRpcInterface"
	HandshakeKey            = "GoPlug"
	HandshakeVersion        = "1.0.0"
	HandshakeProtocol       = 1
)

var (
	AllPluginTypes = Types{
		Rpc:    true,
		Native: true,
	}
	RpcPluginType = Types{
		Rpc:    true,
		Native: false,
	}
	NativePluginType = Types{
		Rpc:    false,
		Native: true,
	}

	OrderedPluginTypes = []Types{NativePluginType, RpcPluginType}
)

//
// Identity is the basic metadata of the plugin
// ---------------------------------------------------------------------------------------------------- //
type Identity struct {
	// Name of the plugin, required.
	Name string `json:"name"`

	// SemVer 2 version, required.
	Version string `json:"version"`

	// The maintainer list with 'Maintainer <email>' format of the plugin
	Maintainers []string `json:"maintainers"`

	// A one line sentence about the function of the plugin - OPTIONAL
	Description string `json:"description"`

	// Icon URL of plugin - OPTIONAL
	Icon string `json:"icon"`

	// IconData build in an icon into the app - OPTIONAL
	IconData []byte `json:"icon_data"`

	// Sets the mechanisms that this plugin supports - OPTIONAL
	PluginTypes Types `json:"types"`

	// The source repository of the plugin code - OPTIONAL
	Repository string `json:"repository"`

	// The source for loading the plugin with specified mode
	Source *Source `json:"source"`

	// The HTTP service should be served by the plugin
	HTTPServices *HTTPServices `json:"HTTPServices,omitempty"`

	// Callbacks - interact with the plugin.
	Callbacks Callbacks `json:"callbacks"`
}

func NewIdentity() *Identity {
	return &Identity{
		Name:         "",
		Version:      "",
		Maintainers:  nil,
		Description:  "",
		Icon:         "",
		IconData:     nil,
		PluginTypes:  NewTypes(),
		Repository:   "",
		Source:       nil,
		HTTPServices: nil,
		Callbacks:    NewCallbacks(),
	}
}

func (i Identity) String() string {
	var ret string
	ret += fmt.Sprintf("Name:\t%s\n", i.Name)
	ret += fmt.Sprintf("\tVersion:\t%s\n", i.Version)
	ret += fmt.Sprintf("\tMaintainers:\t%s\n", strings.Join(i.Maintainers, ", "))
	ret += fmt.Sprintf("\tDescription:\t%s\n", i.Description)
	ret += fmt.Sprintf("\tIcon:\t%s\n", i.Icon)
	ret += fmt.Sprintf("\tRepository:\t%s\n", i.Repository)
	ret += fmt.Sprintf("\tSource:\t%s\n", i.Source)
	ret += fmt.Sprintf("\tCallbacks:\t%v\n", i.Callbacks)
	return ret
}

func (i *Identity) Print() {
	fmt.Println(i.String())
}

func (i *Identity) IsValid() Return.Error {
	var err Return.Error
	for range Only.Once {
		if i == nil {
			err.SetError("PluginIdentity is nil")
			break
		}

		err = i.Callbacks.IsValid()
		if err.IsError() {
			break
		}

		if i.Name == "" {
			err.SetError("plugin config Name is not defined")
		}
		if i.Version == "" {
			err.AddError("plugin config Version is not defined")
		}
		if i.Description == "" {
			err.AddError("plugin config Description is not defined")
		}
		if i.Repository == "" {
			err.AddError("plugin config Repository is not defined")
		}
		if len(i.Maintainers) == 0 {
			err.AddError("plugin config Maintainers is not defined")
		}
	}
	return err
}

func (i *Identity) GetKey(key string) string {
	var value string
	switch strings.ToLower(key) {
	case "name":
		value = i.Name
	case "version":
		value = i.Version
	case "maintainers":
		value = strings.Join(i.Maintainers, ", ")
	case "description":
		value = i.Description
	case "icon":
		value = i.Icon
	case "repository":
		value = i.Repository
	case "source":
		value = i.Source.String()
	}
	return value
}

func (i *Identity) Callback(callback string, ctx Interface, args ...any) Return.Error {
	callback = strings.ToLower(callback)
	switch callback {
	case CallbackInitialise:
		if i.Callbacks.Initialise == nil {
			return Return.NewWarning("Callback '%s' is not defined", callback)
		}
		return i.Callbacks.Initialise(ctx, args...)
	case CallbackRun:
		if i.Callbacks.Run == nil {
			return Return.NewWarning("Callback '%s' is not defined", callback)
		}
		return i.Callbacks.Run(ctx, args...)
	case CallbackNotify:
		if i.Callbacks.Notify == nil {
			return Return.NewWarning("Callback '%s' is not defined", callback)
		}
		return i.Callbacks.Notify(ctx, args...)
	case CallbackExecute:
		if i.Callbacks.Execute == nil {
			return Return.NewWarning("Callback '%s' is not defined", callback)
		}
		return i.Callbacks.Execute(ctx, args...)
	}
	return Return.NewError("unknown callback name '%s', try '%s', '%s', '%s' or '%s'",
		callback, CallbackInitialise, CallbackRun, CallbackNotify, CallbackExecute)
}

func (i *Identity) SetPluginType(name Types) Return.Error {
	i.PluginTypes = name
	return Return.Ok
}
func (i *Identity) SetPluginTypeNative() Return.Error {
	i.PluginTypes.Native = true
	return Return.Ok
}
func (i *Identity) SetPluginTypeRpc() Return.Error {
	i.PluginTypes.Rpc = true
	return Return.Ok
}

//
// Types
// ---------------------------------------------------------------------------------------------------- //
type Types struct {
	Rpc    bool
	Native bool
}

func NewTypes() Types {
	return Types{
		Rpc:    true,
		Native: true,
	}
}

func (p Types) String() string {
	return fmt.Sprintf("Native: %v, Rpc: %v", p.Native, p.Rpc)
}

//
// HTTPServiceRoute defines the http/rest service endpoint served by the plugin
// ---------------------------------------------------------------------------------------------------- //
type HTTPServiceRoute struct {
	// The service endpoint
	Route string

	// The method of the service
	Method string

	// The label will be set into the plugin context to
	// let the plugin aware what kind of request is incoming for serving
	Label string
}

//
// HTTPServices defines the metadata of http service served by the plugin
// ---------------------------------------------------------------------------------------------------- //
type HTTPServices struct {
	// The http service provider
	// Support 'beego'
	Driver string

	// Routes should be enabled on the driver
	Routes []HTTPServiceRoute
}

//
// Callback - Generic callback function type
// ---------------------------------------------------------------------------------------------------- //
type Callback func(ctx Interface, args ...any) Return.Error

func (c Callback) MarshalJSON() ([]byte, error) {
	str := `"` + runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name() + `"`
	return []byte(str), nil
}

func (c *Callback) GetName() string {
	var name string
	for range Only.Once {
		name = runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name()

		re := regexp.MustCompile(`^.*?\.(\(\*[A-Za-z0-9_-]+\)\.[A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 1 {
				name = a[1] + "()"
				break
			}
		}

		re = regexp.MustCompile(`^.*?\.([A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 1 {
				name = a[1] + "()"
				break
			}
		}
	}
	return name
}

func (c *Callback) IsValid() Return.Error {
	var err Return.Error

	// if utils.GetTypeName(c) == "nil" {
	if reflect.ValueOf(*c).IsNil() {
		err.SetWarning("callback is not defined")
	}
	if c == nil {
		err.SetWarning("callback is not defined")
	}
	return err
}

// ---------------------------------------------------------------------------------------------------- //
const (
	CallbackInitialise = "initialise"
	CallbackRun        = "run"
	CallbackNotify     = "notify"
	CallbackExecute    = "execute"
)

//
// Callbacks
// ---------------------------------------------------------------------------------------------------- //
type Callbacks struct {
	PluginName string

	// Initialise - Called on plugin load.
	Initialise     Callback `json:"-"`
	funcInitialise string

	// Run - Execute a function concurrently.
	Run     Callback `json:"-"`
	funcRun string

	// Notify - Notify a plugin.
	Notify     Callback `json:"-"`
	funcNotify string

	// Execute - Execute a function, should return.
	Execute     Callback `json:"-"`
	funcExecute string
}

func NewCallbacks() Callbacks {
	return Callbacks{
		PluginName:     "",
		Initialise:     nil,
		funcInitialise: "",
		Run:            nil,
		funcRun:        "",
		Notify:         nil,
		funcNotify:     "",
		Execute:        nil,
		funcExecute:    "",
	}
}

func (c Callbacks) String() string {
	var ret string
	if c.Initialise != nil {
		ret += fmt.Sprintf(" / Initialise: %s.%s", c.PluginName, c.Initialise.GetName())
	}
	if c.Run != nil {
		ret += fmt.Sprintf(" / Run: %s.%s", c.PluginName, c.Run.GetName())
	}
	if c.Notify != nil {
		ret += fmt.Sprintf(" / Notify: %s.%s", c.PluginName, c.Notify.GetName())
	}
	if c.Execute != nil {
		ret += fmt.Sprintf(" / Execute: %s.%s", c.PluginName, c.Execute.GetName())
	}
	return ret
}

func (c Callbacks) MarshalJSON() ([]byte, error) {
	str1 := c.PluginName + "." + c.Initialise.GetName()
	str2 := c.PluginName + "." + c.Run.GetName()
	str3 := c.PluginName + "." + c.Notify.GetName()
	str4 := c.PluginName + "." + c.Execute.GetName()

	str := fmt.Sprintf(`{ "Initialise":"%s", "Run":"%s", "Notify":"%s", "Execute":"%s" }`,
		str1, str2, str3, str4,
	)
	return []byte(str), nil
}

func (c *Callbacks) IsValid() Return.Error {
	var err Return.Error
	switch {
	case c == nil:
		err.SetError("callbacks is nil")
	default:
		if c.Initialise == nil {
			err.SetWarning("callback Initialise() is nil")
		}
		if c.Run == nil {
			err.AddWarning("callback Run() is nil")
		}
		if c.Notify == nil {
			err.AddWarning("callback Notify() is nil")
		}
		if c.Execute == nil {
			err.AddWarning("callback Execute() is nil")
		}
	}
	return err
}

func (c *Callbacks) SetInitialise(call Callback) Return.Error {
	c.Initialise = call
	return Return.Ok
}

func (c *Callbacks) SetRun(call Callback) Return.Error {
	c.Run = call
	return Return.Ok
}

func (c *Callbacks) SetNotify(call Callback) Return.Error {
	c.Notify = call
	return Return.Ok
}

func (c *Callbacks) SetExecute(call Callback) Return.Error {
	c.Execute = call
	return Return.Ok
}

//
// Source defines the loading mode of the plugin
// ---------------------------------------------------------------------------------------------------- //
type Source struct {
	// The path of the local so file or the URL of the remote git
	Path string

	Version string
}

func (c Source) String() string {
	var ret string
	ret += fmt.Sprintf("Path:\t%s\t", c.Path)
	ret += fmt.Sprintf("\tVersion:\t%s\n", c.Version)
	return ret
}
