/*
Example plugin3
	Another GoPlug example - this time being built and run from the root of the plugins directory.
*/
package main

import (
	"fmt"

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
			Run:        nil,
			Notify:     nil,
			Execute:    nil,
		},
		Name:        "plugin3",
		Version:     "0.4.2",
		Description: "Another GoLang plugin",
		Repository:  "https://github.com/MickMake/GoPlug",
		Maintainers: []string{"mick@mickmake.com", "mick@boutade.net"},
	}
}

type PluginStruct struct {
	String    string
	Int       int
	Interface interface{}
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
}

// Initialise - the plugin logic here
func (m *PluginStruct) Initialise(ctx pluggable.Plugin, args ...interface{}) Return.Error {
	var err Return.Error

	for range Only.Once {
		fmt.Println("plugin3/Initialise()")
		fmt.Printf("plugin3/Initialise() Map: %v\n", ctx.GetMap())
		fmt.Printf("plugin3/Initialise() args: %v\n", args)

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

		fmt.Println("plugin3/Initialise() DONE!")
	}

	return err
}
