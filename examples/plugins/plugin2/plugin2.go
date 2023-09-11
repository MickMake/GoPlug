/*
Example plugin2
	Another GoPlug example.
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
			Initialise: Custom.Initialise,
			Run:        Custom.Run,
			Notify:     notifyMe,
			Execute:    Custom.Execute,
		},
		Name:        "plugin2",
		Version:     "0.1.0",
		Description: "A GoLang plugin example2",
		Repository:  "https://github.com/MickMake/GoPlug",
		Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
	}
}

type PluginStruct struct {
	String    string
	Int       int
	Interface interface{}
	Funky     func(ctx pluggable.Plugin, args ...interface{}) Return.Error `json:"-"`
}

var Custom = PluginStruct{
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
	Funky: notifyMe,
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
		fmt.Println("# err = ctx.SaveConfig(\"Custom\", Custom)")
		err = ctx.SaveConfig("Custom", Custom)
		err.Print()
		fmt.Println("################################################")
		fmt.Println("")

		fmt.Println("################################################")
		fmt.Println("# Custom.Int = 4242")
		Custom.Int = 4242
		fmt.Println(Custom)
		fmt.Println("################################################")
		fmt.Println("")

		fmt.Println("################################################")
		fmt.Println("# err = ctx.LoadConfig(\"Custom\", &Custom)")
		err = ctx.LoadConfig("Custom", &Custom)
		err.Print()
		fmt.Println(Custom)
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

func notifyMe(ctx pluggable.Plugin, args ...interface{}) Return.Error {
	fmt.Printf("notifyMe() called with args:\nctx:\t%v\nargs:\t%v\n", ctx, args)
	return Return.Ok
}
