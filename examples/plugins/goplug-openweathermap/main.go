// A simple plugin that fetches the weather for a defined place and returns it.
package main

import (
	"os"

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

// MyPlugin - Define the plugin as a global. Important for native plugins, not required for RPC.
var MyPlugin GoPlugLoader.PluginItem

// ---------------------------------------------------------------------------------------------------- //

// init - For a native plugin, global variables need to be set in init(), because main() is never called.
func init() {
	err := CreatePlugin(Plugin.NativePluginType)
	err.Print()
}

// main - For an RPC plugin, main() will be called. So we can run the RPC server here.
func main() {
	err := CreatePlugin(Plugin.RpcPluginType)
	if err.IsError() {
		err.Print()
		os.Exit(1)
	}

	err = MyPlugin.Serve()
	if err.IsError() {
		err.Print()
		os.Exit(1)
	}
}

func CreatePlugin(types Plugin.Types) Return.Error {
	var err Return.Error
	var weather Weather

	for range Only.Once {
		MyPlugin, err = GoPlugLoader.NewPluginItem(types, &GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = MyPlugin.SetIdentity(&GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = MyPlugin.SetHandshakeConfig(Plugin.HandshakeConfig)
		if err.IsError() {
			break
		}

		err = MyPlugin.SetInterface(weather)
		if err.IsError() {
			break
		}

		err = MyPlugin.SetHook("", weather.Get)
		if err.IsError() {
			break
		}

		err = MyPlugin.SetHook("", weather.SetLocation, "")
		if err.IsError() {
			break
		}

		err = MyPlugin.SetHook("", weather.SetUnit, "")
		if err.IsError() {
			break
		}

		err = MyPlugin.SetHook("", weather.SetLanguage, "")
		if err.IsError() {
			break
		}

		err = MyPlugin.SetHook("", weather.SetApiKey, "")
		if err.IsError() {
			break
		}

		err = MyPlugin.SetHook("", weather.GetApiKey)
		if err.IsError() {
			break
		}

		err = MyPlugin.SetHook("", weather.SaveConfig)
		if err.IsError() {
			break
		}

		err = MyPlugin.SetHook("", weather.LoadConfig)
		if err.IsError() {
			break
		}

		GoPluginIdentity.Callbacks.SetInitialise(weather.Initialise)

		err = MyPlugin.Validate()
		if err.IsError() {
			break
		}
	}

	return err
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

		ret, err = Plugin.NewHookResponse(*weather)
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

func (w *Weather) Initialise(ctx Plugin.PluginDataInterface, args ...any) Return.Error {
	var err Return.Error

	for range Only.Once {
		err = ctx.SaveIdentity()
		if err.IsError() {
			break
		}
	}

	return err
}
