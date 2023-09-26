// A minimal "hello world" plugin.
package main

import (
	"log"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
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
	Name:        "helloworld",
	Version:     "2.0.0",
	Description: "A GoPlug plugin - Hello World",
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

		err = GoPluginNativeInterface.SetHook("", HelloWorld)
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

// main - For an RPC plugin, main() will be called. So we can run the RPC server here.
func main() {
	var err Return.Error

	for range Only.Once {
		rpc := GoPlugLoader.NewRpcPluginInterface()
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

		err = rpc.SetHook("", HelloWorld)
		if err.IsError() {
			break
		}

		err = rpc.Validate()
		if err.IsError() {
			break
		}

		err = GoPluginRpcInterface.Serve()
	}

	err.Print()
}

// ---------------------------------------------------------------------------------------------------- //

func HelloWorld(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	funcName := utils.GetCaller(0)
	log.Printf("%s() says hi!\n", funcName)
	return Plugin.HookResponseNil, Return.Ok
}
