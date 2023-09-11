/*
Example plugin1
	This code can be run as both a GoPlug plugin and a standalone executable.
*/
package main

import (
	"fmt"
	"time"

	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoPlug/pluggable"
	"github.com/MickMake/GoUnify/Only"
)

//goland:noinspection GoUnusedGlobalVariable
var GoPlugin pluggable.PluginIdentity

func init() {
	GoPlugin = pluggable.PluginIdentity{
		Callbacks: pluggable.Callbacks{
			Initialise: CustomConfig.Initialise,
			Run:        CustomConfig.Run,
			Notify:     notify,
			Execute:    CustomConfig.Execute,
		},
		Name:        "plugin1",
		Version:     "0.1.0",
		Description: "A GoLang plugin example1",
		Repository:  "https://github.com/MickMake/GoPlug",
		Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
		Source: &pluggable.Source{
			Path:    "https://github.com/MickMake/GoPlug/examples/plugins/plugin1/plugin1.so",
			Version: "v0.1.0",
		},
		HTTPServices: &pluggable.HTTPServices{
			Driver: "ShutUp",
			Routes: []pluggable.HTTPServiceRoute{},
		},
	}
}

type PluginStruct struct {
	String    string
	Int       int
	Interface interface{}
	Funky     func(ctx pluggable.Plugin, args ...interface{}) Return.Error `json:"-"`
}

var CustomConfig = PluginStruct{
	String: "Oh yeah!",
	Int:    42,
	Interface: struct {
		One   string
		Two   string
		Three string
		Four  string
	}{
		One:   "anonymous",
		Two:   "structures",
		Three: "are",
		Four:  "cool",
	},
	Funky: notify,
}

// Initialise - the plugin logic here
func (m *PluginStruct) Initialise(ctx pluggable.Plugin, args ...interface{}) Return.Error {
	var err Return.Error

	for range Only.Once {
		fmt.Println("plugin1/Initialise()")
		fmt.Printf("plugin1/Initialise() Map: %v\n", ctx.GetMap())
		fmt.Printf("plugin1/Initialise() args: %v\n", args)

		fmt.Println("################################################")
		fmt.Println("# err = ctx.SaveIdentity()")
		err = ctx.SaveIdentity()
		err.Print()
		fmt.Println("################################################")
		fmt.Println("")

		fmt.Println("################################################")
		fmt.Println("# err = ctx.SaveConfig(\"CustomConfig\", CustomConfig)")
		err = ctx.SaveConfig("CustomConfig", CustomConfig)
		err.Print()
		fmt.Println("################################################")
		fmt.Println("")

		fmt.Println("################################################")
		fmt.Println("# CustomConfig.Int = 4242")
		CustomConfig.Int = 4242
		fmt.Println(CustomConfig)
		fmt.Println("################################################")
		fmt.Println("")

		fmt.Println("################################################")
		fmt.Println("# err = ctx.LoadConfig(\"CustomConfig\", &CustomConfig)")
		err = ctx.LoadConfig("CustomConfig", &CustomConfig)
		err.Print()
		fmt.Println(CustomConfig)
		fmt.Println("################################################")

		label := ctx.Get("master")
		if label == nil {
			fmt.Println("plugin1/Initialise() label == nil")
			err.SetWarning("label is not set in the plugin context")
			break
		}

		fmt.Printf("master == '%s'\n", label)
		labelV := fmt.Sprintf("%s", label)
		switch labelV {
		case "OK":
			fmt.Println("return call_subMethod1()")
		case "case2":
			fmt.Println("return call_subMethod12()")
		default:
			fmt.Println("not supported")
			err.SetWarning("not supported")
		}
	}

	return err
}

// Execute - the plugin logic here
func (m *PluginStruct) Execute(ctx pluggable.Plugin, args ...interface{}) Return.Error {
	var err Return.Error

	for range Only.Once {
		fmt.Println("plugin1/Execute()")
		fmt.Printf("plugin1/Execute() Map: %v\n", ctx.GetMap())
		fmt.Printf("plugin1/Execute() args: %v\n", args)
	}

	return err
}

func (m *PluginStruct) Run(ctx pluggable.Plugin, args ...interface{}) Return.Error {
	var err Return.Error

	for range Only.Once {
		fmt.Println("plugin1/Run()")
		fmt.Printf("plugin1/Run() Map: %v\n", ctx.GetMap())
		fmt.Printf("plugin1/Run() args: %v\n", args)

		time.Sleep(time.Second * 5)
		fmt.Printf("plugin1/Run() Now: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(time.Second * 5)
		fmt.Printf("plugin1/Run() Now: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(time.Second * 5)
		fmt.Printf("plugin1/Run() Now: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		time.Sleep(time.Second * 5)
		fmt.Printf("plugin1/Run() Now: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	return err
}

func notify(ctx pluggable.Plugin, args ...interface{}) Return.Error {
	fmt.Printf("plugin1/Notify() Map: %v\n", ctx.GetMap())
	return Return.Ok
}

func main() {
	foo := pluggable.NewPlugin(&GoPlugin)
	foo2 := foo.GetMap()
	fmt.Print(foo2)
}
