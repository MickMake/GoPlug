package main

import (
	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
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
	Name:        "quote",
	Version:     "1.0.0",
	Description: "A GoPlug plugin - Fetch a quote of the day",
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

	for range Only.Once {
		// Native plugin
		GoPluginNativeInterface = GoPlugLoader.NewNativePluginInterface()
		err = GoPluginNativeInterface.NewNativePlugin()
		if err.IsError() {
			break
		}

		err = GoPluginNativeInterface.SetIdentity(&GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = GoPluginNativeInterface.SetHandshakeConfig(Plugin.HandshakeConfig)
		if err.IsError() {
			break
		}

		err = GoPluginNativeInterface.SetHook("", Get, "", 0, "")
		if err.IsError() {
			break
		}

		err = GoPluginNativeInterface.Validate()
		if err.IsError() {
			break
		}
	}

	err.Print()
}

// ---------------------------------------------------------------------------------------------------- //

// main - For an RPC plugin, main() will be called. So we can run the RPC server here.
func main() {
	var err Return.Error

	for range Only.Once {
		// RPC plugin
		GoPluginRpcInterface = GoPlugLoader.NewRpcPluginInterface()
		err = GoPluginRpcInterface.NewRpcPlugin()
		if err.IsError() {
			break
		}

		err = GoPluginRpcInterface.SetIdentity(&GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = GoPluginRpcInterface.SetHandshakeConfig(Plugin.HandshakeConfig)
		if err.IsError() {
			break
		}

		err = GoPluginRpcInterface.SetHook("", Get)
		if err.IsError() {
			break
		}

		err = GoPluginRpcInterface.Validate()
		if err.IsError() {
			break
		}

		err = GoPluginRpcInterface.Serve()
	}

	err.Print()
}

// ---------------------------------------------------------------------------------------------------- //

// Get - This program has only one call - which fetches a web page.
func Get(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	var ret Plugin.HookResponse
	var err Return.Error

	for range Only.Once {
		// Expecting args in this order: ApiKey, Limit, Category
		ApiKey := Cast.ToString(args[0])
		if ApiKey == "" {
			err.SetError("Need an API key!")
			break
		}

		Limit := Cast.ToInt(args[1])
		if Limit == 0 {
			Limit = 1
		}

		Category := Cast.ToString(args[2])
		if Category == "" {
			err.SetError("Need a category!")
			break
		}

		fetch := utils.NewHttp()
		err = fetch.SetUrl("https://api.api-ninjas.com/v1/quotes?limit=%d&category=%s", Limit, Category)
		if err.IsError() {
			break
		}

		err = fetch.SetHeader("X-Api-Key", ApiKey)
		if err.IsError() {
			break
		}

		var body []byte
		body, err = fetch.Get()
		if err.IsError() {
			break
		}

		ret, err = Plugin.NewHookResponse(string(body))
	}

	return ret, err
}
