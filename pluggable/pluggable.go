package pluggable

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	// jsonV2 "github.com/go-json-experiment/json"
	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoUnify/Only"
)


// ---------------------------------------------------------------------------------------------------- //

type PluginItem struct {
	File   utils.PluginPath
	Plugin Plugin
	Config *PluginIdentity
}


func (i *PluginItem) IsValid() Return.Error {
	var err Return.Error
	switch {
		case i == nil:
			err.SetError("plugin is nil")
		default:
			err = i.Config.IsValid()
	}
	return err
}


// ---------------------------------------------------------------------------------------------------- //

// Plugin - help to provide related information/parameters to the
// plugin execution entry method.
// Plugin - inherits all from the context.Context
type Plugin interface {
	context.Context
	Values


	// Runtime methods.

	// Init - Create a new instance of this plugin.
	// It should create and fill out the Config structure.
	Init() Return.Error

	// Identify - Return a filled out Config structure.
	Identify() (PluginIdentity, Return.Error)

	// Debug - .
	Debug(ref any)

	GetMap() map[string]interface{}

	// SetConfigDir - Saves a plugin config to disk.
	SetConfigDir(dir string) Return.Error

	// SaveIdentity - Saves a plugin's identity to disk.
	SaveIdentity() Return.Error

	// SaveConfig - Saves a plugin config to disk.
	SaveConfig(fn string, ref interface{}) Return.Error

	// LoadConfig - Loads a plugin config from disk.
	LoadConfig(fn string, ref interface{}) Return.Error

	// SaveJson - Saves json data to disk.
	SaveJson(fn string, data []byte) Return.Error

	// LoadJson - Loads json data from disk.
	LoadJson(fn string) ([]byte, Return.Error)


	// Build methods.

	// Validate - Ensure plugin has a valid Config structure.
	Validate() Return.Error
}

// NewPlugin - build the base plugin context based on the context.
func NewPlugin(config *PluginIdentity) Plugin {
	return &PluggableStruct {
		context: context.Background(),
		Directory: utils.PluginPath{},
		Config:  config,
		values:  make(map[string]interface{}),
	}
}


// ---------------------------------------------------------------------------------------------------- //

// PluggableStruct implemented as default plugin context
//goland:noinspection GoNameStartsWithPackageName
type PluggableStruct struct {
	// For compatible with system context
	context context.Context

	// PluginIdentity - config
	Config *PluginIdentity

	// PluginIdentity - config
	Directory utils.PluginPath

	// For keeping values
	values map[string]interface{}
}


// SetConfigDir - Sets the config directory to be used for saving/loading of files.
func (p *PluggableStruct) SetConfigDir(dir string) Return.Error {
	var err Return.Error
	p.Directory, err = utils.NewDir(dir)
	return err
}

// SaveIdentity - Saves the PluginIdentity struct as a JSON file.
func (p *PluggableStruct) SaveIdentity() Return.Error {
	var err Return.Error

	for range Only.Once {
		err.SetPrefix("SaveIdentity(): ")

		// Expose plugin name to "(c Callbacks) MarshalJSON()"
		p.Config.Callbacks.pluginName = p.Config.Name

		var data []byte
		var e error
		data, e = json.Marshal(p.Config)
		err.SetError(e)
		if err.IsError() {
			break
		}

		err = p.SaveJson(p.Config.Name, data)
		if err.IsError() {
			break
		}
	}

	return err
}

// SaveConfig - Save an arbitrary plugin structure as a JSON file.
func (p *PluggableStruct) SaveConfig(fn string, ref interface{}) Return.Error {
	var err Return.Error

	for range Only.Once {
		err.SetPrefix("SaveConfig(): ")

		var data []byte
		var e error
		data, e = json.Marshal(ref)
		err.SetError(e)
		if err.IsError() {
			break
		}

		err = p.SaveJson(fn, data)
		if err.IsError() {
			break
		}
	}

	return err
}

// LoadConfig - Load a JSON file into an arbitrary plugin structure.
func (p *PluggableStruct) LoadConfig(fn string, ref interface{}) Return.Error {
	var err Return.Error

	for range Only.Once {
		err.SetPrefix("LoadConfig(): ")

		var data []byte
		data, err = p.LoadJson(fn)
		if err.IsError() {
			break
		}

		var e error
		e = json.Unmarshal(data, ref)
		err.SetError(e)
		if err.IsError() {
			break
		}
	}

	return err
}

// SaveJson - .
func (p *PluggableStruct) SaveJson(fn string, data []byte) Return.Error {
	var err Return.Error
	for range Only.Once {
		var fp utils.PluginPath
		fp, _ = p.Directory.AppendFile(fn)
		fnp := fp.ChangeExtension("json")
		fn = fnp.GetPath()

		err = utils.WriteFile(fn, data)
		if err.IsError() {
			break
		}
	}
	return err
}

// LoadJson - .
func (p *PluggableStruct) LoadJson(fn string) ([]byte, Return.Error) {
	var data []byte
	var err Return.Error
	for range Only.Once {
		var fp utils.PluginPath
		fp, _ = p.Directory.AppendFile(fn)
		fnp := fp.ChangeExtension("json")
		fn = fnp.GetPath()

		data, err = utils.ReadFile(fn)
		if err.IsError() {
			break
		}
	}
	return data, err
}

// IsValid - .
func (p *PluggableStruct) IsValid() Return.Error {
	var err Return.Error
	switch {
		case p == nil:
			err.SetError("callbacks is nil")
		default:
			if p.Config.Callbacks.Initialise == nil {
				err.SetWarning("callback Initialise() is nil")
			}
			if p.Config.Callbacks.Run == nil {
				err.AddWarning("callback Run() is nil")
			}
			if p.Config.Callbacks.Notify == nil {
				err.AddWarning("callback Notify() is nil")
			}
			if p.Config.Callbacks.Execute == nil {
				err.AddWarning("callback Execute() is nil")
			}
	}
	return err
}

// GetMap - .
func (p *PluggableStruct) GetMap() map[string]interface{} {
	return p.values
}

// Debug - .
func (p *PluggableStruct) Debug(ref any) {
	fmt.Printf("\tp.Debug: %v\n", ref)
}

// Init - .
func (p *PluggableStruct) Init() Return.Error {
	p.Set("master-init", "OK")
	p.Set("master-init-timestamp", time.Now())

	// TODO implement me
	fmt.Printf("Init() - implement me\n")
	fmt.Printf("\tp.Config: %v\n", p.Config)
	fmt.Printf("\tp.valueStore: %v\n", p.values)
	fmt.Printf("\tp.pluginContext: %v\n", p.context)
	return Return.Ok
}

// Identify - .
func (p *PluggableStruct) Identify() (PluginIdentity, Return.Error) {
	// TODO implement me
	fmt.Printf("Identify() - implement me\n")
	fmt.Printf("\tp.Config: %v\n", p.Config)
	fmt.Printf("\tp.valueStore: %v\n", p.values)
	fmt.Printf("\tp.pluginContext: %v\n", p.context)
	return PluginIdentity{}, Return.Ok
}

// Validate - .
func (p *PluggableStruct) Validate() Return.Error {
	p.Set("master-validate", "OK")
	p.Set("master-validate-timestamp", time.Now())

	// TODO implement me
	fmt.Printf("Validate() - implement me\n")
	fmt.Printf("\tp.Config: %v\n", p.Config)
	fmt.Printf("\tp.valueStore: %v\n", p.values)
	fmt.Printf("\tp.pluginContext: %v\n", p.context)
	return Return.Ok
}


// ---------------------------------------------------------------------------------------------------- //
// Mirror functions of context.Context

// Deadline - Mirrors context.Context
func (p *PluggableStruct) Deadline() (deadline time.Time, ok bool) {
	return p.context.Deadline()
}

// Done - Mirrors context.Context
func (p *PluggableStruct) Done() <-chan struct{} {
	return p.context.Done()
}

// Err - Mirrors context.Context
func (p *PluggableStruct) Err() error {
	return p.context.Err()
}

// Value - implements 'Value' in context.Context
func (p *PluggableStruct) Value(key interface{}) interface{} {
	return p.context.Value(key)
}


// ---------------------------------------------------------------------------------------------------- //

// Values - Getter/Setter for string map of interfaces{}
type Values interface {
	// Get - Get value by the key
	Get(key string) interface{}

	// Set - Set value to context
	Set(key string, value interface{})
}

// Get - implements 'GetValue' in Values interface
func (p *PluggableStruct) Get(key string) interface{} {
	return p.values[key]
}

// Set - implements 'SetValue' in Values interface
func (p *PluggableStruct) Set(key string, value interface{}) {
	if len(strings.TrimSpace(key)) > 0 {
		// nil value is allowed
		p.values[key] = value
	}
}
