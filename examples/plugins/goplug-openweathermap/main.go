// A simple plugin that fetches the weather for a defined place and returns it.
package main

import (
	"github.com/MickMake/GoUnify/Only"
	owm "github.com/briandowns/openweathermap"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils/Cast"
	"github.com/MickMake/GoPlug/utils/Return"
)

// GoPluginIdentity - Set the GoPlugin identity.
var GoPluginIdentity = Plugin.Identity{
	Callbacks: Plugin.Callbacks{
		Initialise: nil,
		Run:        nil,
		Notify:     nil,
		Execute:    nil,
	},
	Name:        "openweathermap",
	Version:     "1.0.0",
	Description: "A GoPlug plugin - Access openweathermap.org",
	Repository:  "https://github.com/MickMake/GoPlug",
	Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
}

// GoPluginRpcInterface - RPC based plugin.
var GoPluginRpcInterface GoPlugLoader.RpcPluginInterface

// GoPluginNativeInterface - Native based plugin.
var GoPluginNativeInterface GoPlugLoader.NativePluginInterface

// ---------------------------------------------------------------------------------------------------- //

// init - For a native plugin, global variables need to be set in init(), because main() is never called.
func init() {
	var err Return.Error
	var weather Weather

	GoPluginRpcInterface, err = InitRpc(&weather)
	if err.IsError() {
		err.Print()
		return
	}

	GoPluginNativeInterface, err = InitNative(&weather)
	if err.IsError() {
		err.Print()
		return
	}
}

// main - For an RPC plugin, main() will be called. So we can run the RPC server here.
func main() {
	err := GoPluginRpcInterface.Serve()
	err.Print()
}

// InitRpc - Set up an RPC based plugin.
func InitRpc(w *Weather) (GoPlugLoader.RpcPluginInterface, Return.Error) {
	var rpc GoPlugLoader.RpcPluginInterface
	var err Return.Error

	for range Only.Once {
		rpc = GoPlugLoader.NewRpcPluginInterface()
		err = rpc.NewRpcPlugin()
		if err.IsError() {
			break
		}

		err = rpc.SetIdentity(&GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = rpc.SetHandshakeConfig(Plugin.HandshakeConfig)
		if err.IsError() {
			break
		}

		err = rpc.SetInterface(w)
		if err.IsError() {
			break
		}

		err = rpc.SetHook("", w.Get)
		if err.IsError() {
			break
		}

		err = rpc.SetHook("", w.SetLocation, "")
		if err.IsError() {
			break
		}

		err = rpc.SetHook("", w.SetUnit, "")
		if err.IsError() {
			break
		}

		err = rpc.SetHook("", w.SetLanguage, "")
		if err.IsError() {
			break
		}

		err = rpc.SetHook("", w.SetApiKey, "")
		if err.IsError() {
			break
		}

		err = rpc.SetHook("", w.GetApiKey)
		if err.IsError() {
			break
		}

		err = rpc.SetHook("", w.SaveConfig)
		if err.IsError() {
			break
		}

		err = rpc.SetHook("", w.LoadConfig)
		if err.IsError() {
			break
		}

		GoPluginIdentity.Callbacks.SetInitialise(w.Initialise)

		err = rpc.Validate()
		if err.IsError() {
			break
		}
	}

	return rpc, err
}

// InitNative - Set up a native based plugin.
func InitNative(w *Weather) (GoPlugLoader.NativePluginInterface, Return.Error) {
	var native GoPlugLoader.NativePluginInterface
	var err Return.Error

	for range Only.Once {
		native = GoPlugLoader.NewNativePluginInterface()
		err = native.NewNativePlugin()
		if err.IsError() {
			break
		}

		err = native.SetIdentity(&GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = native.SetHandshakeConfig(Plugin.HandshakeConfig)
		if err.IsError() {
			break
		}

		err = native.SetHook("", w.Get)
		if err.IsError() {
			break
		}

		err = native.SetHook("", w.SetLocation, "")
		if err.IsError() {
			break
		}

		err = native.SetHook("", w.SetUnit, "")
		if err.IsError() {
			break
		}

		err = native.SetHook("", w.SetLanguage, "")
		if err.IsError() {
			break
		}

		err = native.SetHook("", w.SetApiKey, "")
		if err.IsError() {
			break
		}

		err = native.SetHook("", w.GetApiKey)
		if err.IsError() {
			break
		}

		err = native.SetHook("", w.SaveConfig)
		if err.IsError() {
			break
		}

		err = native.SetHook("", w.LoadConfig)
		if err.IsError() {
			break
		}

		GoPluginIdentity.Callbacks.SetInitialise(w.Initialise)

		err = native.Validate()
		if err.IsError() {
			break
		}
	}

	return native, err
}

//
// Weather
// ---------------------------------------------------------------------------------------------------- //
type Weather struct {
	ApiKey   string
	Location string
	Unit     string
	Language string
}

func (w *Weather) Get(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	var ret Plugin.HookResponse
	var err Return.Error

	for range Only.Once {
		if w.Location == "" {
			w.Location = "Sydney"
		}

		if w.Unit == "" {
			w.Unit = "C"
		}

		if w.Language == "" {
			w.Language = "en"
		}

		weather, e := owm.NewCurrent(w.Unit, w.Language, w.ApiKey)
		if e != nil {
			err.SetError(e)
			break
		}

		e = weather.CurrentByName(w.Location)
		if e != nil {
			err.SetError(e)
			break
		}

		ret, err = Plugin.NewHookResponse(weather)
	}

	return ret, err
}

func (w *Weather) GetApiKey(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	return Plugin.NewHookResponse(w.ApiKey)
}

func (w *Weather) SaveConfig(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	return Plugin.HookResponseNil, hook.Plugin.SaveObject("openweather", *w)
}

func (w *Weather) LoadConfig(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	err := hook.Plugin.LoadObject("openweather", w)
	if err.IsError() {
		return Plugin.HookResponseNil, err
	}

	resp, _ := Plugin.NewHookResponse(*w)
	return resp, err
}

func (w *Weather) SetApiKey(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	w.ApiKey = Cast.ToString(args[0])
	return Plugin.HookResponseNil, Return.Ok
}

func (w *Weather) SetLocation(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	w.Location = Cast.ToString(args[0])
	return Plugin.HookResponseNil, Return.Ok
}

func (w *Weather) SetUnit(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	w.Unit = Cast.ToString(args[0])
	return Plugin.HookResponseNil, Return.Ok
}

func (w *Weather) SetLanguage(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	w.Language = Cast.ToString(args[0])
	return Plugin.HookResponseNil, Return.Ok
}

func (w *Weather) Initialise(ctx Plugin.Interface, args ...any) Return.Error {
	var err Return.Error

	for range Only.Once {
		err = ctx.SaveIdentity()
		if err.IsError() {
			break
		}
	}

	return err
}
