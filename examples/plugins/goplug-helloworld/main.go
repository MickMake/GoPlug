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

// GoPluginIdentity - Set the GoPlugin identity. This is required for a minimal setup.
var GoPluginIdentity = Plugin.Identity{
	Name:        "helloworld",
	Version:     "2.0.0",
	Description: "A GoPlug plugin - Hello World",
	Repository:  "https://github.com/MickMake/GoPlug",
	Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
}

// MyNativePlugin - Define the plugin as a global. Important for native plugins, not required for RPC.
var MyNativePlugin GoPlugLoader.PluginItem

// ---------------------------------------------------------------------------------------------------- //

// init - For a native plugin, global variables need to be set in init(), because main() is never called.
func init() {
	var err Return.Error

	for range Only.Once {
		MyNativePlugin, err = GoPlugLoader.NewPluginItem(Plugin.NativePluginType, &GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = MyNativePlugin.SetIdentity(&GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = MyNativePlugin.SetHandshakeConfig(Plugin.HandshakeConfig)
		if err.IsError() {
			break
		}

		err = MyNativePlugin.SetHook("", HelloWorld)
		if err.IsError() {
			break
		}

		err = MyNativePlugin.Validate()
		if err.IsError() {
			break
		}
	}

	err.Print()
}

// main - For an RPC plugin, main() will be called. So we can run the RPC server here.
func main() {
	var MyRpcPlugin GoPlugLoader.PluginItem
	var err Return.Error

	for range Only.Once {
		MyRpcPlugin, err = GoPlugLoader.NewPluginItem(Plugin.RpcPluginType, &GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = MyRpcPlugin.SetIdentity(&GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = MyRpcPlugin.SetHandshakeConfig(Plugin.HandshakeConfig)
		if err.IsError() {
			break
		}

		err = MyRpcPlugin.SetHook("", HelloWorld)
		if err.IsError() {
			break
		}

		err = MyRpcPlugin.Validate()
		if err.IsError() {
			break
		}

		err = MyRpcPlugin.Serve()
	}

	err.Print()
}

// ---------------------------------------------------------------------------------------------------- //

func HelloWorld(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	funcName := utils.GetCaller(0)
	log.Printf("%s() says hi!\n", funcName)
	return Plugin.HookResponseNil()
}
