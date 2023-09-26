package main

import (
	"fmt"
	"os"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug"
	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/GoPlugLoader/Plugin"
	"github.com/MickMake/GoPlug/utils/Return"
)

// Manager - Define the manager structure. Can be anything you want,
// but will be based on the ManagerInterface interface.
type Manager struct {
}

// GoPlugin - Define the identity of the manager. Will be used to determine what plugins can run.
var GoPlugin = Plugin.Identity{
	Name:        "master",
	Version:     "1.0.0",
	Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
	Description: "GoPlug master example",
	Icon:        "",
	IconData:    make([]byte, 128),
	PluginTypes: Plugin.Types{
		Rpc:    true,
		Native: false,
	},
	Repository: "https://github.com/MickMake/GoPlug",

	Source: &Plugin.Source{
		Version: "v1.0.0",
		Path:    "https://github.com/MickMake/GoPlug/examples/master/master",
	},
	Callbacks: Plugin.Callbacks{
		Initialise: nil, // pluginNative.Initialise,
		Run:        nil,
		Notify:     nil,
		Execute:    nil, // pluginNative.Execute,
	},
	HTTPServices: &Plugin.HTTPServices{
		Driver: "dummy",
		Routes: nil,
	},
}

// var pluginNative Manager

func main() {
	var err Return.Error

	for range Only.Once {
		var manager GoPlug.Manager
		manager, err = GoPlug.NewPluginManager(&GoPlugin)
		if err.IsError() {
			break
		}

		// Set a log file.
		err = manager.SetLogfile("PluginManager.log")
		if err.IsError() {
			break
		}

		// Define the plugin directory.
		err = manager.SetDir("../plugins")
		if err.IsError() {
			break
		}

		// Define the identity of the manager.
		err = manager.SetIdentity(GoPlugin)
		if err.IsError() {
			break
		}

		// Can be used to define a manager implementor.
		err = manager.SetImplementor(nil)
		if err.IsError() {
			break
		}

		// Define what plugins to load, (native, rpc, or both).
		err = manager.SetPluginTypes(Plugin.NativePluginType) // Plugin.AllPluginTypes
		if err.IsError() {
			break
		}

		// Look for plugins.
		err = manager.Scan()
		if err.IsError() {
			break
		}

		// Register all found plugins.
		err = manager.RegisterPlugins()
		if err.IsError() {
			break
		}

		// Print them all.
		manager.ListPlugins()

		// Run tests on the "quote" plugin.
		err = Quote(manager)
		if err.IsError() {
			break
		}

		// Run tests on the "helloworld" plugin.
		err = HelloWorld(manager)
		if err.IsError() {
			break
		}

		// Run tests on the "openweathermap" plugin.
		err = OpenWeather(manager)
		if err.IsError() {
			break
		}

		// Run tests on the "simple" plugin.
		err = Simple(manager)
		if err.IsError() {
			break
		}
	}

	err.Print()
	if err.IsError() {
		os.Exit(1)
	}

	os.Exit(0)
}

func HelloWorld(manager GoPlug.Manager) Return.Error {
	var err Return.Error

	for range Only.Once {
		var plugin *GoPlugLoader.PluginItem

		plugin, err = manager.GetPluginByName("helloworld")
		if err.IsError() {
			break
		}

		plugin.Data.PrintHooks()

		var response Plugin.HookResponse
		response, err = plugin.Data.CallHook("TestExec2", "can", "you", "see", "these", "args", 42)
		if err.IsError() {
			break
		}
		response.Print()

		err = plugin.Execute()
		if err.IsError() {
			break
		}
	}

	return err
}

func Quote(manager GoPlug.Manager) Return.Error {
	var err Return.Error

	for range Only.Once {
		var plugin *GoPlugLoader.PluginItem

		plugin, err = manager.GetPluginByName("quote")
		if err.IsError() {
			break
		}

		plugin.Data.PrintHooks()

		err = plugin.Execute()
		if err.IsError() {
			break
		}

		fmt.Println("#### Get a quotable quote...")
		var response Plugin.HookResponse

		fmt.Println("Get")
		response, err = plugin.Data.CallHook("Get", "+Yx9sCPNO2rdepKHAzn23Q==yMrXyexBNoRSxdzP", 1, "men")
		if err.IsError() {
			break
		}

		response.Print()
	}

	return err
}

func Simple(manager GoPlug.Manager) Return.Error {
	var err Return.Error

	for range Only.Once {
		var plugin *GoPlugLoader.PluginItem

		fmt.Println("\n// ---------------------------------------------------------------------------------------------------- //")
		fmt.Println("Testing Simple plugin")

		plugin, err = manager.GetPluginByName("simple")
		if err.IsError() {
			break
		}

		plugin.Data.PrintHooks()

		var response Plugin.HookResponse

		fmt.Printf("#### Calling TestExec1()\n")
		response, err = plugin.Data.CallHook("TestExec1", "can", "you", "see", "these", "args", 42)
		if err.IsError() {
			break
		}
		response.Print()

		fmt.Printf("#### Calling TestExec2()\n")
		response, err = plugin.Data.CallHook("TestExec2", "can", "you", "see", "these", "args", 42)
		if err.IsError() {
			break
		}
		response.Print()

		fmt.Printf("#### Calling Callbacks.Execute()\n")
		err = plugin.Execute()
		if err.IsError() {
			break
		}

		fmt.Println("Calling Larry(42, \"42\", 42) - invalid args.")
		response, err = plugin.Data.CallHook("Larry", 42, "42", 42)
		response.Print()
		err.Print()

		fmt.Println("Calling Larry(\"hey now\", \"this finally works\", 42) - valid args.")
		response, err = plugin.Data.CallHook("Larry", "hey now", "this finally works", 42)
		response.Print()
		if err.IsError() {
			break
		}

		fmt.Println("Calling Curly(nil, nil) - invalid args.")
		response, err = plugin.Data.CallHook("Curly", nil, nil)
		response.Print()
		err.Print()

		fmt.Println("Calling Curly(nil, nil) - valid args.")
		type test struct {
			A string
			B string
			C string
		}
		response, err = plugin.Data.CallHook("Curly", 42, test{A: "1", B: "2", C: "3"})
		response.Print()
		if err.IsError() {
			break
		}

		fmt.Println("Calling Mo(42) - invalid number of args.")
		response, err = plugin.Data.CallHook("Mo", 42)
		response.Print()
		err.Print()

		fmt.Println("Calling Mo(nil, nil) - valid args.")
		response, err = plugin.Data.CallHook("Mo")
		response.Print()
		if err.IsError() {
			break
		}

		plugin.Data.SetValue(
			"sample", struct {
				string
				int
			}{
				"hello go-plugin",
				42,
			},
		)

		value := plugin.Data.GetValue("sample")
		fmt.Printf("value: %v\n", value)

		hooks := plugin.GetItemHooks(nil)
		hook := hooks.GetHook("Larry")
		fmt.Println("Calling Larry(42, \"42\", 42) - invalid args.")
		err = hook.Validate()
		if err.IsError() {
			break
		}
	}

	return err
}

func OpenWeather(manager GoPlug.Manager) Return.Error {
	var err Return.Error

	for range Only.Once {
		var plugin *GoPlugLoader.PluginItem

		plugin, err = manager.GetPluginByName("openweathermap")
		if err.IsError() {
			break
		}

		plugin.Data.PrintHooks()

		err = plugin.Execute()
		if err.IsError() {
			break
		}

		fmt.Println("#### Check in on the weather...")
		var response Plugin.HookResponse

		fmt.Println("LoadConfig")
		response, err = plugin.Data.CallHook("LoadConfig")
		if err.IsError() {
			fmt.Printf("LoadConfig: %s\n", err)
			fmt.Println("LoadConfig: Setting up variables.")

			fmt.Println("SetApiKey")
			response, err = plugin.Data.CallHook("SetApiKey", "c86a53ab92071be04a20f7580f780fc0")
			if err.IsError() {
				break
			}

			fmt.Println("SetLocation")
			response, err = plugin.Data.CallHook("SetLocation", "Sydney")
			if err.IsError() {
				break
			}

			fmt.Println("SetUnit")
			response, err = plugin.Data.CallHook("SetUnit", "C")
			if err.IsError() {
				break
			}

			fmt.Println("SetLanguage")
			response, err = plugin.Data.CallHook("SetLanguage", "en")
			if err.IsError() {
				break
			}

			fmt.Println("SaveConfig")
			response, err = plugin.Data.CallHook("SaveConfig")
			if err.IsError() {
				break
			}
		}

		fmt.Println("LoadConfig: Checking key")
		response, err = plugin.Data.CallHook("GetApiKey")
		fmt.Printf("GetApiKey: %s\n", response)
		err.Print()
		response.Print()
		foo := response.AsString()
		fmt.Printf("fii:[%s]\n", foo)

		fmt.Println("Get")
		response, err = plugin.Data.CallHook("Get")
		response.Print()
		err.Print()
	}

	return err
}

//
// -------------------------------------------------------------------------------- //
//

// // Initialise - the plugin logic here
// func (m *Manager) Initialise(ctx Plugin.Interface, args ...any) Return.Error {
// 	var err Return.Error
//
// 	// GoPlugConfig.PluginCallbackContext
// 	for range Only.Once {
// 		fmt.Printf("Initialise(%v, %v)\n", ctx, args)
// 		ctx.Callback(Plugin.CallbackInitialise, ctx, args...)
//
// 		label := "case1"
// 		// label := ctx.Values().KeyGet("master")
// 		// if label == nil {
// 		// 	fmt.Println("label is not set in the plugin context")
// 		// 	err.SetError("label is not set in the plugin context")
// 		// 	break
// 		// }
//
// 		fmt.Printf("label == '%s'\n", label)
// 		labelV := fmt.Sprintf("%s", label)
// 		switch labelV {
// 		case "case1":
// 			fmt.Println("return call_subMethod1()")
// 		case "case2":
// 			fmt.Println("return call_subMethod12()")
// 		default:
// 			fmt.Println("not supported")
// 			err.SetError("not supported")
// 		}
// 	}
//
// 	return err
// }
//
// // Execute - the plugin logic here
// func (m *Manager) Execute(ctx Plugin.Interface, args ...any) Return.Error {
// 	var err Return.Error
//
// 	for range Only.Once {
// 		fmt.Printf("Execute(%v, %v)\n", ctx, args)
// 		ctx.Callback(Plugin.CallbackExecute, ctx, args...)
//
// 		label := "case1"
// 		// label := ctx.Values().KeyGet("master")
// 		// if label == nil {
// 		// 	fmt.Println("label is not set in the plugin context")
// 		// 	err.SetError("label is not set in the plugin context")
// 		// 	break
// 		// }
//
// 		fmt.Printf("label == '%s'\n", label)
// 		labelV := fmt.Sprintf("%s", label)
// 		switch labelV {
// 		case "case1":
// 			fmt.Println("return call_subMethod1()")
// 		case "case2":
// 			fmt.Println("return call_subMethod12()")
// 		default:
// 			fmt.Println("not supported")
// 			err.SetError("not supported")
// 		}
// 	}
//
// 	return err
// }
//
// func (info *PluginInfo) Init(pluginPath utils.FilePath, prefix string) {
// 	for range Only.Once {
// 		info.Error.ReturnClear()
// 		info.Error.SetPrefix("")
//
// 		id := strings.TrimPrefix(pluginPath.GetName(), prefix)
// 		*info = PluginInfo{
// 			Id:   id,
// 			Path: pluginPath,
// 			Client: goplugin.NewClient(
// 				&goplugin.ClientConfig{
// 					HandshakeConfig: Plugin.HandshakeConfig,
// 					Plugins: map[string]goplugin.Plugin{
// 						id: &GoPlugLoader.GoPluginMaster{},
// 					},
// 					Cmd: exec.Command(pluginPath.GetPath()),
// 				},
// 			),
// 		}
// 	}
// }
//
// func (info *PluginInfo) GetIdentity() (Plugin.Identity, Return.Error) {
// 	var ret Plugin.Identity
//
// 	for range Only.Once {
// 		info.Error.ReturnClear()
// 		info.Error.SetPrefix("")
//
// 		log.Printf("[%s]: Name:%s Path: %s\n", info.Id, info.Path.GetName(), info.Path)
// 		fmt.Printf("[%s]: Client - ID:%s Exited:%v Protocol:%s Version:%d\n",
// 			info.Id, info.Client.ID(), info.Client.Exited(), info.Client.Protocol(),
// 			info.Client.NegotiatedVersion(),
// 		)
//
// 		var rpcClient goplugin.ClientProtocol
// 		var e error
// 		rpcClient, e = info.Client.Client()
// 		if e != nil {
// 			info.Error.SetError("[%s]: ERROR: %s\n", info.Id, e.Error())
// 			break
// 		}
// 		defer rpcClient.Close()
//
// 		e = rpcClient.Ping()
// 		if e != nil {
// 			info.Error.SetError("[%s]: ERROR: %s\n", info.Id, e.Error())
// 			break
// 		}
//
// 		var raw any
// 		raw, e = rpcClient.Dispense(info.Id)
// 		if e != nil {
// 			info.Error.SetError("[%s]: ERROR: %s\n", info.Id, e.Error())
// 			break
// 		}
//
// 		impl := raw.(*GoPlugLoader.RpcPluginClient)
// 		ret = impl.Identify()
// 		ret.Print()
// 		data := impl.GetData()
// 		fmt.Println("data:")
// 		data.Print()
// 		fmt.Println("Identity:")
// 		data.Identity.Print()
// 		fmt.Println("Identity:")
// 		data.Hooks.PrintHooks()
// 		if data.Hooks.HookExists("Larry") {
// 			name, err := data.Hooks.GetHookName("Larry")
// 			fmt.Printf("hookName: %s (%s)\n", name, err)
//
// 			args, err := data.Hooks.GetHookArgs("Larry")
// 			fmt.Printf("hookArgs: %s (%s)\n", args, err)
//
// 			fnpReturn, err := impl.CallHook("Larry", "does this work", "42", 42)
// 			fmt.Printf("Output: %s (%s)\n", fnpReturn, err)
// 		}
//
// 		fmt.Printf("")
// 	}
//
// 	return ret, info.Error
// }
//
// func (info *PluginInfo) Close() Return.Error {
// 	for range Only.Once {
// 		info.Error.ReturnClear()
// 		info.Error.SetPrefix("")
//
// 		log.Printf("[%s]: Name:%s Path: %s\n", info.Id, info.Path.GetName(), info.Path)
// 		fmt.Printf("[%s]: Client - ID:%s Exited:%v Protocol:%s Version:%d\n",
// 			info.Id, info.Client.ID(), info.Client.Exited(), info.Client.Protocol(),
// 			info.Client.NegotiatedVersion(),
// 		)
//
// 		var rpcClient goplugin.ClientProtocol
// 		var e error
// 		rpcClient, e = info.Client.Client()
// 		if e != nil {
// 			log.Printf("[%s]: ERROR: %s\n", info.Id, e.Error())
// 			break
// 		}
//
// 		e = rpcClient.Ping()
// 		if e != nil {
// 			log.Printf("[%s]: ERROR: %s\n", info.Id, e.Error())
// 			break
// 		}
//
// 		var raw any
// 		raw, e = rpcClient.Dispense(info.Id)
// 		if e != nil {
// 			log.Printf("[%s]: ERROR: %s\n", info.Id, e.Error())
// 			break
// 		}
//
// 		name := utils.GetStructName(raw)
// 		log.Printf("[%s]: Structure: %s\n", info.Id, name)
//
// 		impl := raw.(*GoPlugLoader.RpcPluginClient)
// 		result := impl.GetData()
// 		log.Print(Plugin.StructToString(fmt.Sprintf("[%s]:impl.Greet()", info.Id), result))
//
// 		result2 := impl.Identify()
// 		log.Print(Plugin.StructToString(fmt.Sprintf("[%s]:impl.Identify()", info.Id), result2))
//
// 		result3 := impl.IdentifyString()
// 		log.Print(Plugin.StructToString(fmt.Sprintf("[%s]:impl.IdentifyString()", info.Id), result3))
//
// 		if !info.Client.Exited() {
// 			log.Printf("[%s]: Close()\n", info.Id)
// 			e = rpcClient.Close()
// 		}
// 		if e != nil {
// 			log.Printf("[%s]: ERROR: %s\n", info.Id, e.Error())
// 			break
// 		}
//
// 		if info.Client.Exited() {
// 			continue
// 		}
// 		info.Client.ReattachConfig()
// 		log.Printf("[%s]: Kill()\n", info.Id)
// 		info.Client.Kill()
// 	}
//
// 	return info.Error
// }
//
// type PluginInfo struct {
// 	Id             string
// 	Path           utils.FilePath
// 	Client         *goplugin.Client
// 	clientProtocol goplugin.ClientProtocol
// 	Error          Return.Error
// }
//
// type Manager struct {
// 	Type    string
// 	Plugins map[string]*PluginInfo
// }
//
// func (m *Manager) Dispose() {
// 	var wg sync.WaitGroup
// 	for _, pinfo := range m.Plugins {
// 		wg.Add(1)
// 		go func(client *goplugin.Client) {
// 			client.Kill()
// 			wg.Done()
// 		}(pinfo.Client)
// 	}
// 	wg.Wait()
// }
