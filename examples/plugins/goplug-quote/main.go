package main

import (
	"encoding/json"
	"os"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Cast"
	"github.com/MickMake/GoPlug/utils/Return"
)

// GoPluginIdentity - Set the GoPlugin identity. This is required for a minimal setup.
var GoPluginIdentity = Plugin.Identity{
	Name:        "quote",
	Version:     "1.0.0",
	Description: "A GoPlug plugin - Fetch a quote of the day",
	Repository:  "https://github.com/MickMake/GoPlug",
	Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
}

// MyPlugin - Define the plugin as a global. Important for native plugins, not required for RPC.
var MyPlugin GoPlugLoader.PluginItem

// ---------------------------------------------------------------------------------------------------- //

// init - For a native plugin, global variables need to be set in init(), because main() is never called.
func init() {
	var err Return.Error
	MyPlugin, err = CreatePlugin(Plugin.NativePluginType)
	err.Print()
}

// ---------------------------------------------------------------------------------------------------- //

// main - For an RPC plugin, main() will be called. So we can run the RPC server here.
func main() {
	foo, bar := MyPlugin.CallHook("Get", "+Yx9sCPNO2rdepKHAzn23Q==yMrXyexBNoRSxdzP", 1, "men")
	bar.Print()
	foo.Print()

	rpc, err := CreatePlugin(Plugin.RpcPluginType)
	if err.IsError() {
		err.Print()
		os.Exit(1)
	}

	err = rpc.Serve()
	if err.IsError() {
		err.Print()
		os.Exit(1)
	}
}

// ---------------------------------------------------------------------------------------------------- //

func CreatePlugin(types Plugin.Types) (GoPlugLoader.PluginItem, Return.Error) {
	var plug GoPlugLoader.PluginItem
	var err Return.Error

	for range Only.Once {
		plug, err = GoPlugLoader.NewPluginItem(types, &GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = plug.SetIdentity(&GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = plug.SetHandshakeConfig(Plugin.HandshakeConfig)
		if err.IsError() {
			break
		}

		err = plug.SetHook("", Get, "", 0, "")
		if err.IsError() {
			break
		}

		err = plug.Validate()
		if err.IsError() {
			break
		}
	}

	return plug, err
}

type Quote []struct {
	Author   string `json:"author"`
	Category string `json:"category"`
	Quote    string `json:"quote"`
}

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

		var quote Quote
		e := json.Unmarshal(body, &quote)
		if e != nil {
			err.SetError(e)
			break
		}

		if len(quote) == 0 {
			err.SetError("No quotes returned")
			break
		}

		ret, err = Plugin.NewHookResponse(quote[0])
	}

	return ret, err
}
