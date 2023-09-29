// This program shows how to define a complex setup.
// An alternative structure, ("MyPlugin"), is attached to the Plugin.Interface interface.
// Replacement methods to the Plugin.Interface interface are defined.
// An alternative hook structure is defined.
// Hooks are created and attached to both functions and methods.
// Complex structure is sent back to some requests.
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Cast"
	"github.com/MickMake/GoPlug/utils/Return"
)

// GoPluginIdentity - Set the GoPlugin identity. This is required for a minimal setup.
var GoPluginIdentity = Plugin.Identity{
	Callbacks: Plugin.Callbacks{
		Initialise: InitMe,
		Run:        nil,
		Notify:     nil,
		Execute:    TestExec3,
	},
	Name:        "simple",
	Version:     "2.0.0",
	Description: "A GoPlug plugin - simple",
	Repository:  "https://github.com/MickMake/GoPlug",
	Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
}

// MyNativePlugin - Define the plugin as a global.
var MyNativePlugin GoPlugLoader.PluginItem

// MyRpcPlugin - Define the plugin as a global.
var MyRpcPlugin GoPlugLoader.PluginItem

// ---------------------------------------------------------------------------------------------------- //

// init - For a native plugin, global variables need to be set in init(), because main() is never called.
func init() {
	fmt.Println("INIT()")
	MyRpcPlugin = InitRpc()
	MyNativePlugin = InitNative()
}

// main - For an RPC plugin, main() will be called. So we can run the RPC server here.
func main() {
	err := Return.NewWithPrefix("main")
	for range Only.Once {
		err = MyRpcPlugin.Serve()
		err.Print()
	}
}

// ---------------------------------------------------------------------------------------------------- //

// InitRpc - Set up an RPC based plugin.
func InitRpc() GoPlugLoader.PluginItem {
	var plug GoPlugLoader.PluginItem
	var err Return.Error

	for range Only.Once {
		log.Printf("InitRpc()")

		plug, err = GoPlugLoader.NewPluginItem(Plugin.RpcPluginType, &GoPluginIdentity)
		if err.IsError() {
			break
		}

		plug, err = InitCommon(plug)
		if err.IsError() {
			break
		}

		// Add a hook that is different between RPC & native.
		err = plug.SetHook("TestExec4", TestExecRpc, "", "", "", "", "", 0)

		err = plug.Validate()
		if err.IsError() {
			break
		}

		// p2.KeySet("data", simple)
		// p2.RegisterStructure(simple)
		// p2.KeySet("func", Test)
		// p2.RegisterStructure(Data)
		// p2.RegisterStructure(Test)
	}
	err.Print()

	return plug
}

// InitNative - Set up a native based plugin.
func InitNative() GoPlugLoader.PluginItem {
	var plug GoPlugLoader.PluginItem
	var err Return.Error

	for range Only.Once {
		log.Printf("InitNative()")

		plug, err = GoPlugLoader.NewPluginItem(Plugin.NativePluginType, &GoPluginIdentity)
		if err.IsError() {
			break
		}

		plug, err = InitCommon(plug)
		if err.IsError() {
			break
		}

		// Add a hook that is different between RPC & native.
		err = plug.SetHook("TestExec4", TestExecNative, "", 0)

		err = plug.Validate()
		if err.IsError() {
			break
		}
	}
	err.Print()

	return plug
}

// InitCommon - Set up common plugin parts.
func InitCommon(plug GoPlugLoader.PluginItem) (GoPlugLoader.PluginItem, Return.Error) {
	var err Return.Error

	for range Only.Once {
		log.Printf("InitCommon()")

		err = plug.SetIdentity(&GoPluginIdentity)
		if err.IsError() {
			break
		}

		err = plug.SetPluginType(GoPluginIdentity.PluginTypes)
		if err.IsError() {
			break
		}

		err = plug.SetHandshakeConfig(Plugin.HandshakeConfig)
		if err.IsError() {
			break
		}

		var Data Plugin.HookStore
		Data = new(MyPlugin)
		Data.NewHookStore()
		Data.SetHookPlugin(plug.GetItemData())
		Data.SetHook("TestExec1", TestExec1, "", "", "", "", "", 0)

		err = plug.SetHookStore(Data)
		if err.IsError() {
			break
		}

		err = plug.SetHook("TestExec2", TestExec2, "", "", "", "", "", 0)

		err = plug.Validate()
		if err.IsError() {
			break
		}
	}

	return plug, err
}

// TestExec1 - A test hook.
func TestExec1(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	funcName := utils.GetCaller(0)
	log.Printf("\nCalled %s(%v)\n", funcName, args)
	return Plugin.NewHookResponse(args)
}

// TestExec2 - Another test hook.
func TestExec2(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	funcName := utils.GetCaller(0)
	log.Printf("\nCalled %s(%v)\n", funcName, args)
	return Plugin.NewHookResponse(args)
}

// TestExec3 - Will be executed when this plugin is loaded and also later by an "Execute" function call.
// See GoPluginIdentity.Callbacks.Execute
func TestExec3(ctx Plugin.PluginDataInterface, args ...any) Return.Error {
	funcName := utils.GetCaller(0)
	log.Printf("\nCalled %s(%v)\n", funcName, args)
	// log.Printf("Called %s(%v)\n # CTX:\n%v\n\n", funcName, args, ctx)
	return Return.Ok
}

// TestExecNative - Hook that is different between RPC & native.
func TestExecNative(ctx Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	funcName := utils.GetCaller(0)
	log.Printf("\nCalled %s(%v)\n", funcName, args)
	log.Printf("CTX:\n%v\n\n", ctx)
	return Plugin.HookResponseNil()
}

// TestExecRpc - Hook that is different between RPC & native.
func TestExecRpc(ctx Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	funcName := utils.GetCaller(0)
	log.Printf("\nCalled %s(%v)\n", funcName, args)
	log.Printf("CTX:\n%v\n\n", ctx)
	return Plugin.HookResponseNil()
}

// InitMe - Will be executed when this plugin is loaded by the "Initialise" Callback.
// See GoPluginIdentity.Callbacks.Initialise
func InitMe(ctx Plugin.PluginDataInterface, args ...any) Return.Error {
	funcName := utils.GetCaller(0)
	log.Printf("\nCalled %s(%v)\n", funcName, args)
	// log.Printf("Called %s(%v)\n # CTX:\n%v\n\n", funcName, args, ctx)
	return Return.Ok
}

// ---------------------------------------------------------------------------------------------------- //

// Test - A test structure, which will be sent to the master.
type Test struct {
	A string
	B string
	C string
}

// Larry - Another test hook, but is a method off MyPlugin.
func (d *MyPlugin) Larry(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	var response Plugin.HookResponse
	var err Return.Error

	for range Only.Once {
		funcName := utils.GetCaller(0)
		log.Printf("\nCalled %s(%v)\n", funcName, args)

		err = d.Data.Hooks.ValidateHook(args...)
		if err.IsError() {
			break
		}

		for index, arg := range args {
			if utils.IsTypeOfName(arg, "string") {
				args[index] = strings.ToUpper(arg.(string))
				continue
			}
			if utils.IsTypeOfName(arg, "int") {
				args[index] = arg.(int) * 100
				continue
			}
		}

		response, err = Plugin.NewHookResponse(args)
		if err.IsError() {
			break
		}
	}

	return response, err
}

// Curly - Another test hook, but is a method off MyPlugin.
func (d *MyPlugin) Curly(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	var response Plugin.HookResponse
	var err Return.Error

	for range Only.Once {
		funcName := utils.GetCaller(0)
		log.Printf("\nCalled %s(%v)\n", funcName, args)

		err = d.Data.Hooks.ValidateHook(args...)
		if err.IsError() {
			break
		}

		if utils.IsTypeOfName(args[0], "int") {
			args[0] = Cast.ToInt(args[0]) * 1000
		}

		if utils.IsTypeOfName(args[0], "string") {
			args[1] = strings.ToLower(Cast.ToString(args[0]))
		}

		if utils.IsTypeOfName(args[0], "Test") {
			args[2] = args[2].(Test)
		}

		response, err = Plugin.NewHookResponse(args)
		if err.IsError() {
			break
		}
	}

	return response, err
}

// Mo - Another test hook, but is a method off MyPlugin.
func (d *MyPlugin) Mo(hook Plugin.HookStruct, args ...any) (Plugin.HookResponse, Return.Error) {
	funcName := utils.GetCaller(0)
	log.Printf("\nCalled %s(%v)\n", funcName, args)

	err := d.Data.Hooks.ValidateHook(args...)
	if err.IsError() {
		return Plugin.HookResponse{}, err
	}
	return Plugin.NewHookResponse(d.Data.Hooks.GetHookIdentity())
}

//
// MyPlugin
// ---------------------------------------------------------------------------------------------------- //
type MyPlugin struct {
	Plugin.HookStore
	Data Plugin.DynamicData
}

// NewHookStore - Alternative method that will install some hooks.
func (d *MyPlugin) NewHookStore() Return.Error {
	var err Return.Error
	for range Only.Once {
		d.Data = *Plugin.NewDynamicData(Plugin.PluginData{})
		err = d.Data.Hooks.SetHook("", d.Larry, "", "", 0)
		if err.IsError() {
			break
		}

		err = d.Data.Hooks.SetHook("", d.Curly, 0, "", Test{})
		if err.IsError() {
			break
		}

		err = d.Data.Hooks.SetHook("", d.Mo)
		if err.IsError() {
			break
		}
	}
	return err
}

func (d *MyPlugin) SetHookPlugin(plugin Plugin.PluginDataInterface) {
	d.Data.Hooks.SetHookPlugin(plugin)
}
func (d *MyPlugin) GetHookReference() *Plugin.HookStruct {
	return d.Data.Hooks.GetHookReference()
}
func (d *MyPlugin) GetHookIdentity() string {
	return d.Data.Hooks.GetHookIdentity()
}
func (d *MyPlugin) SetHookIdentity(identity string) Return.Error {
	return d.Data.Hooks.SetHookIdentity(identity)
}
func (d *MyPlugin) HookExists(hook string) bool {
	return d.Data.Hooks.HookExists(hook)
}
func (d *MyPlugin) HookNotExists(hook string) bool {
	return d.Data.Hooks.HookNotExists(hook)
}
func (d *MyPlugin) GetHook(hook string) *Plugin.Hook {
	return d.Data.Hooks.GetHook(hook)
}
func (d *MyPlugin) GetHookName(name string) (string, Return.Error) {
	return d.Data.Hooks.GetHookName(name)
}
func (d *MyPlugin) GetHookFunction(name string) (Plugin.HookFunction, Return.Error) {
	return d.Data.Hooks.GetHookFunction(name)
}
func (d *MyPlugin) GetHookArgs(name string) (Plugin.HookArgs, Return.Error) {
	return d.Data.Hooks.GetHookArgs(name)
}
func (d *MyPlugin) ValidateHook(args ...any) Return.Error {
	return d.Data.Hooks.ValidateHook(args...)
}
func (d *MyPlugin) SetHook(name string, function Plugin.HookFunction, args ...any) Return.Error {
	return d.Data.Hooks.SetHook(name, function, args...)
}
func (d *MyPlugin) CountHooks() int {
	return d.Data.Hooks.CountHooks()
}
func (d *MyPlugin) ListHooks() Plugin.HookMap {
	return d.Data.Hooks.ListHooks()
}
func (d *MyPlugin) PrintHooks() {
	d.Data.Hooks.PrintHooks()
}
func (d *MyPlugin) String() string {
	return d.Data.String()
}
