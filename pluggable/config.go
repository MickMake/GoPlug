package pluggable

import (
	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoUnify/Only"
)


// ---------------------------------------------------------------------------------------------------- //

// PluginIdentity is the corresponding structure of the 'plugin.json',
// which describe the basic metadata of the plugin
type PluginIdentity struct {
	// Name of the plugin, required.
	Name string `json:"name"`

	// SemVer 2 version, required.
	Version string `json:"version"`

	// The maintainer list with 'Maintainer <email>' format of the plugin
	Maintainers []string `json:"maintainers"`

	// A one line sentence about the function of the plugin - OPTIONAL
	Description string `json:"description"`

	// Icon URL of plugin - OPTIONAL
	Icon string `json:"icon"`

	// The source repository of the plugin code - OPTIONAL
	Repository string `json:"repository"`

	// The source for loading the plugin with specified mode
	Source *Source `json:"source"`

	// The HTTP service should be served by the plugin
	HTTPServices *HTTPServices `json:"HTTPServices,omitempty"`

	// Callbacks - interact with the plugin.
	Callbacks Callbacks `json:"callbacks"`
}


func (c *PluginIdentity) IsValid() Return.Error {
	var err Return.Error
	for range Only.Once {
		if c == nil {
			err.SetError("callbacks is nil")
			break
		}

		err = c.Callbacks.IsValid()
		if err.IsError() {
			break
		}

		if c.Name == "" {
			err.SetError("plugin config Name is not defined")
		}
		if c.Version == "" {
			err.AddError("plugin config Version is not defined")
		}
		if c.Description == "" {
			err.AddError("plugin config Description is not defined")
		}
		if c.Repository == "" {
			err.AddError("plugin config Repository is not defined")
		}
		if len(c.Maintainers) == 0 {
			err.AddError("plugin config Maintainers is not defined")
		}
		// if c.Source == nil {
		// 	err.AddError("plugin config Source is not defined")
		// }
	}
	return err
}


// ---------------------------------------------------------------------------------------------------- //

// HTTPServiceRoute defines the http/rest service endpoint served by the plugin
type HTTPServiceRoute struct {
	// The service endpoint
	Route string

	// The method of the service
	Method string

	// The label will be set into the plugin context to
	// let the plugin aware what kind of request is incoming for serving
	Label string
}


// ---------------------------------------------------------------------------------------------------- //

// HTTPServices defines the metadata of http service served by the plugin
type HTTPServices struct {
	// The http service provider
	// Support 'beego'
	Driver string

	// Routes should be enabled on the driver
	Routes []HTTPServiceRoute
}


// ---------------------------------------------------------------------------------------------------- //

// PluginManagerIdentity is the corresponding structure of the 'plugin.json',
// which describe the basic metadata of the plugin
type PluginManagerIdentity struct {
	// Name of the plugin, required.
	Name string `json:"name"`

	// SemVer 2 version, required.
	Version string `json:"version"`

	// The maintainer list with 'Maintainer <email>' format of the plugin
	Maintainers []string `json:"maintainers"`

	// A one line sentence about the function of the plugin - OPTIONAL
	Description string `json:"description"`

	// Icon URL of plugin - OPTIONAL
	Icon string `json:"icon"`

	// The source repository of the plugin code - OPTIONAL
	Repository string `json:"repository"`

	// The source for loading the plugin with specified mode
	Source *Source `json:"source"`

	// The HTTP service should be served by the plugin
	HTTPServices *HTTPServices `json:"HTTPServices,omitempty"`

	// Callbacks - interact with the plugin.
	Callbacks Callbacks `json:"callbacks"`
}

func (c *PluginManagerIdentity) IsValid() Return.Error {
	var err Return.Error
	for range Only.Once {
		if c == nil {
			err.SetError("plugin manager callbacks is nil")
			break
		}

		err = c.Callbacks.IsValid()
		if err.IsError() {
			break
		}

		if c.Name == "" {
			err.SetError("plugin manager config Name is not defined")
		}
		if c.Version == "" {
			err.AddError("plugin manager config Version is not defined")
		}
		if c.Description == "" {
			err.AddError("plugin manager config Description is not defined")
		}
		if c.Repository == "" {
			err.AddError("plugin manager config Repository is not defined")
		}
		if len(c.Maintainers) == 0 {
			err.AddError("plugin manager config Maintainers is not defined")
		}
		if c.Source == nil {
			err.AddError("plugin manager config Source is not defined")
		}
	}
	return err
}
