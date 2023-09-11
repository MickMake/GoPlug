package main

import (
	"fmt"
	"runtime"

	"github.com/MickMake/GoPlug"
	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoPlug/pluggable"
	"github.com/MickMake/GoUnify/Only"
)

var GoPlugin = pluggable.PluginManagerIdentity{
	Callbacks: pluggable.Callbacks{
		Initialise: master.Initialise,
		Run:        nil,
		Notify:     nil,
		Execute:    master.Execute,
	},
	Name:        "master",
	Version:     "1.0.0",
	Description: "GoPlug master example",
	Repository:  "https://github.com/MickMake/GoPlug",
	Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
	Source: &pluggable.Source{
		Version: "v1.0.0",
		Path:    "https://github.com/MickMake/GoPlug/examples/master/master",
	},
}

var master ManagerStruct

type ManagerStruct struct{}

// Initialise - the plugin logic here
func (m *ManagerStruct) Initialise(ctx pluggable.Plugin, args ...interface{}) Return.Error {
	var err Return.Error

	for range Only.Once {
		fmt.Printf("Initialise(%v, %v)\n", ctx, args)
		ctx.Init()

		runtime.Breakpoint()

		label := ctx.Get("master")
		if label == nil {
			fmt.Println("label is not set in the plugin context")
			err.SetError("label is not set in the plugin context")
			break
		}

		fmt.Printf("label == '%s'\n", label)
		labelV := fmt.Sprintf("%s", label)
		switch labelV {
		case "case1":
			fmt.Println("return call_subMethod1()")
		case "case2":
			fmt.Println("return call_subMethod12()")
		default:
			fmt.Println("not supported")
			err.SetError("not supported")
		}
	}

	return err
}

// Execute - the plugin logic here
func (m *ManagerStruct) Execute(ctx pluggable.Plugin, args ...interface{}) Return.Error {
	var err Return.Error

	for range Only.Once {
		fmt.Printf("Execute(%v, %v)\n", ctx, args)
		ctx.Init()

		label := ctx.Get("master")
		if label == nil {
			fmt.Println("label is not set in the plugin context")
			err.SetError("label is not set in the plugin context")
			break
		}

		fmt.Printf("label == '%s'\n", label)
		labelV := fmt.Sprintf("%s", label)
		switch labelV {
		case "case1":
			fmt.Println("return call_subMethod1()")
		case "case2":
			fmt.Println("return call_subMethod12()")
		default:
			fmt.Println("not supported")
			err.SetError("not supported")
		}
	}

	return err
}

func main() {
	for range Only.Once {
		pluginManager := GoPlug.NewPluginManager()
		err := pluginManager.SetConfig(GoPlugin)
		err.Print()
		if err.IsError() {
			break
		}

		err = pluginManager.SetBaseDir("../plugins")
		err.Print()
		if err.IsError() {
			break
		}

		// err = pluginManager.BuildPlugins()
		// err.Print()

		pluginManager.ListPlugins()

		err = pluginManager.LoadPlugins()
		err.Print()

		var plugin *pluggable.PluginItem
		plugin, err = pluginManager.GetPlugin("plugin1")
		err.Print()
		if err.IsError() {
			break
		}

		pluginManager.ListPlugins()

		err = plugin.Run("hey", "there")
		err.Print()
		if err.IsError() {
			break
		}

		err = plugin.Execute()
		err.Print()
		if err.IsError() {
			break
		}

		plugin.Plugin.Set("sample", struct {
			string
			int
		}{
			"hello go-plugin",
			42,
		})
		err = plugin.Plugin.Init()
		if err.IsError() {
			break
		}

		value := plugin.Plugin.Get("sample")
		fmt.Printf("value: %v\n", value)
	}
}
