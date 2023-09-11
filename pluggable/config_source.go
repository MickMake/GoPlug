package pluggable

import (
	"fmt"
	"path/filepath"

	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoPlug/utils"
	"github.com/Masterminds/semver"
	"github.com/MickMake/GoUnify/Only"
)


// ---------------------------------------------------------------------------------------------------- //

// Source defines the loading mode of the plugin
type Source struct {
	// The path of the local so file or the URL of the remote git
	Path string

	Version string
}


// ---------------------------------------------------------------------------------------------------- //

// LocalSourceValidator validate the local source
type LocalSourceValidator struct{}

// Validate is the implementation of Validator interface
func (lsv *LocalSourceValidator) Validate(params ...interface{}) (interface{}, Return.Error) {
	var ret interface{}
	var err Return.Error

	for range Only.Once {
		ret = nil
		err.SetPrefix("Validate(): ")

		if len(params) < 2 {
			err.SetError("plugin json object and plugin base dir are required")
			break
		}

		identity, ok := params[0].(*PluginIdentity)
		if !ok {
			err.SetError("invalid plugin identity")
			break
		}

		if identity.Source == nil {
			err.SetWarning("plugin source missing")
		} else {
			_, err2 := semver.NewVersion(identity.Source.Version)
			if err2 != nil {
				err.SetWarning("plugin source version missing: %s", err2)
			}
		}

		pluginBaseDir := fmt.Sprintf("%s", params[1])
		// plugin so file path
		var pluginSoFilePath string
		if filepath.IsAbs(identity.Source.Path) {
			pluginSoFilePath = identity.Source.Path
		} else {
			pluginSoFilePath = filepath.Join(pluginBaseDir, identity.Source.Path)
		}

		_, err = utils.FileExists(pluginSoFilePath)
		if err.IsError() {
			err.AddError("plugin *.so file")
			break
		}

		if filepath.Ext(pluginSoFilePath) != ".so" {
			err.SetError("%s.so file is missing", identity.Name)
			break
		}

		// Override the so file path to absolute path
		identity.Source.Path = pluginSoFilePath

		ret = identity
	}

	return ret, err
}


// ---------------------------------------------------------------------------------------------------- //

// RemoteSourceValidator validates the remote source
type RemoteSourceValidator struct{}

// Validate is the implementation of Validator interface
// TODO:
func (rsv *RemoteSourceValidator) Validate(_ ...interface{}) (interface{}, Return.Error) {
	return nil, Return.Error{}
}
