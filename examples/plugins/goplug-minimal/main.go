// A minimal, bare-bones native plugin.
package main

import (
	"log"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

// GoPluginIdentity - Set the GoPlugin identity. This is required for a minimal setup.
var GoPluginIdentity = Plugin.Identity{
	Callbacks: Plugin.Callbacks{
		Execute: Minimal,
	},
	Name:        "minimal",
	Version:     "2.0.0",
	Description: "A GoPlug plugin - A minimal plugin",
	Repository:  "https://github.com/MickMake/GoPlug",
	Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
}

// MyPlugin - Define the plugin as a global. Important for native plugins, not required for RPC.
var MyPlugin GoPlugLoader.PluginItem

func init() {
	var err Return.Error
	MyPlugin, err = GoPlugLoader.NewPluginItem(Plugin.NativePluginType, &GoPluginIdentity)
	err.ExitIfError()

	err = MyPlugin.Validate()
	err.ExitIfError()
}

// main - For an RPC plugin, main() will be called. So we can run the RPC server here.
func main() {
	rpcPlugin, err := GoPlugLoader.NewPluginItem(Plugin.RpcPluginType, &GoPluginIdentity)
	err.ExitIfError()

	err = MyPlugin.Validate()
	err.ExitIfError()

	err = rpcPlugin.Serve()
	err.ExitIfError()
}

func Minimal(ctx Plugin.PluginDataInterface, args ...any) Return.Error {
	funcName := utils.GetCaller(0)
	log.Printf("%s() says hi!\n", funcName)
	return Return.Ok
}
