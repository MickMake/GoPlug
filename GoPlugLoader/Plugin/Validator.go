package Plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

// ---------------------------------------------------------------------------------------------------- //

// Validator defines the behaviors of a plugin validator
type Validator interface {
	// Validate - Do validation with the provided params.
	// If meet any issues, an error will be returned.
	// If it succeeds, output the result which depends on the implementations.
	Validate(params ...any) (any, Return.Error)
}

// ---------------------------------------------------------------------------------------------------- //

// BaseValidatorChain build a validation pipeline with 'JSONFileValidator' and 'IdentityValidator'.
type BaseValidatorChain struct {
	// The validator list
	validators []Validator

	Error Return.Error
}

// NewBaseValidatorChain creates a validator chain
func NewBaseValidatorChain(validators ...Validator) Validator {
	var err Return.Error
	err.SetPrefix("Validator: ")

	bvc := &BaseValidatorChain{
		validators: make([]Validator, 0),
		Error:      err,
	}

	if len(validators) > 0 {
		bvc.validators = append(bvc.validators, validators...)
	}

	return bvc
}

// Validate is the implementation of Validator interface
func (bvc *BaseValidatorChain) Validate(params ...any) (any, Return.Error) {
	var ret any
	var err Return.Error

	for range Only.Once {
		ret = nil

		if len(bvc.validators) == 0 {
			err.SetError("no validators")
			break
		}

		if len(params) == 0 {
			err.SetError("missing params")
			break
		}

		for _, vl := range bvc.validators {
			if ret == nil {
				// The first validator
				ret, err = vl.Validate(params...)
			} else {
				extendedParams := []any{ret}
				extendedParams = append(extendedParams, params...)
				ret, err = vl.Validate(extendedParams...)
			}

			if err.IsError() {
				break
			}
		}
	}

	return ret, err
}

// ---------------------------------------------------------------------------------------------------- //

// JSONFileValidator validates the existence of plugin.json.
type JSONFileValidator struct{}

// Validate is the implementation of Validator interface
func (jfv *JSONFileValidator) Validate(params ...any) (any, Return.Error) {
	var ret any
	var err Return.Error

	for range Only.Once {
		ret = nil

		if len(params) == 0 {
			err.SetError("The plugin dir path is required")
			break
		}

		pluginDirPath := fmt.Sprintf("%s", params[0])

		_, err = utils.DirExists(pluginDirPath)
		if err.IsError() {
			break
		}

		if !utils.IsDir(pluginDirPath) {
			err.SetError("File %s is not a dir", pluginDirPath)
			break
		}

		pluginJSONFile := filepath.Join(pluginDirPath, utils.PluginJSONFileName)
		_, err = utils.FileExists(pluginJSONFile)
		if err.IsError() {
			err.AddError("%s is not found under plugin dir %s", utils.PluginJSONFileName, pluginDirPath)
			break
		}

		data, e := os.ReadFile(pluginJSONFile)
		if e != nil {
			err.SetError(e)
			break
		}

		// Load plugin.json
		identity := &Identity{}
		e = json.Unmarshal(data, identity)
		if e != nil {
			err.SetError(e)
			break
		}

		// Plugin dir name should be equal with the name of the plugin
		var fi os.FileInfo
		fi, e = os.Stat(pluginDirPath)
		if e != nil {
			err.SetError(e)
			// Actually, should not come here
			break
		}

		if fi.Name() != identity.Name {
			err.SetError("Name conflicts: expect %s but got %s in the metadata json file", fi.Name(), identity.Name)
			break
		}

		ret = identity
	}

	return ret, err
}

// ---------------------------------------------------------------------------------------------------- //

// IdentityValidator validates the plugin identity.
type IdentityValidator struct{}

// Validate is the implementation of Validator interface
func (sv *IdentityValidator) Validate(params ...any) (any, Return.Error) {
	var ret any
	var err Return.Error

	for range Only.Once {
		ret = nil

		if len(params) == 0 {
			err.SetError("plugin json object is missing")
			break
		}

		identity, ok := params[0].(*Identity)
		if !ok {
			err.SetError("invalid plugin identity object")
			break
		}

		if len(identity.Name) == 0 {
			err.SetError("missing plugin name")
			break
		}

		_, err2 := semver.NewVersion(identity.Version)
		err.SetError(err2)
		if err.IsError() {
			break
		}

		if identity.Source == nil {
			err.SetWarning("plugin source missing")
		} else {
			_, err2 = semver.NewVersion(identity.Source.Version)
			err.SetWarning(err2)
		}

		// if identity.Source.Mode != utils.PluginSourceModeLocal && identity.Source.Mode != utils.PluginSourceModeRemote {
		// 	err.SetError("Only support mode [%s, %s]", utils.PluginSourceModeLocal, utils.PluginSourceModeRemote)
		// 	break
		// }

		ret = identity
	}

	return ret, err
}

//
// LocalSourceValidator validate the local source
// ---------------------------------------------------------------------------------------------------- //
type LocalSourceValidator struct{}

// Validate is the implementation of Validator interface
func (lsv *LocalSourceValidator) Validate(params ...any) (any, Return.Error) {
	var ret any
	var err Return.Error

	for range Only.Once {
		ret = nil
		err.SetPrefix("Validate(): ")

		if len(params) < 2 {
			err.SetError("plugin json object and plugin base dir are required")
			break
		}

		identity, ok := params[0].(*Identity)
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

		// var pluginBaseDir utils.PluginPath
		// pluginBaseDir, err = utils.NewDir(params[1])
		// // plugin so file path
		// var pluginSoFilePath utils.PluginPath
		// if filepath.IsAbs(identity.Source.Path) {
		// 	pluginSoFilePath, err = utils.NewFile(identity.Source.Path)
		// } else {
		// 	pluginSoFilePath, err = utils.NewFile(pluginBaseDir, identity.Source.Path)
		// }
		//
		// _, err = utils.FileExists(pluginSoFilePath)
		// if err.IsError() {
		// 	err.AddError("plugin *.so file")
		// 	break
		// }
		//
		// if filepath.Ext(pluginSoFilePath) != ".so" {
		// 	err.SetError("%s.so file is missing", identity.Name)
		// 	break
		// }
		//
		// // Override the so file path to absolute path
		// identity.Source.Path = pluginSoFilePath
		//
		//
		// // pluginBaseDir := fmt.Sprintf("%s", params[1])
		// // // plugin so file path
		// // var pluginSoFilePath string
		// // if filepath.IsAbs(identity.Source.Path) {
		// // 	pluginSoFilePath = identity.Source.Path
		// // } else {
		// // 	pluginSoFilePath = filepath.Join(pluginBaseDir, identity.Source.Path)
		// // }
		// //
		// // _, err = utils.FileExists(pluginSoFilePath)
		// // if err.IsError() {
		// // 	err.AddError("plugin *.so file")
		// // 	break
		// // }
		// //
		// // if filepath.Ext(pluginSoFilePath) != ".so" {
		// // 	err.SetError("%s.so file is missing", identity.Name)
		// // 	break
		// // }
		// //
		// // // Override the so file path to absolute path
		// // identity.Source.Path = pluginSoFilePath

		ret = identity
	}

	return ret, err
}

//
// RemoteSourceValidator validates the remote source
// ---------------------------------------------------------------------------------------------------- //
type RemoteSourceValidator struct{}

// Validate is the implementation of Validator interface
// TODO:
func (rsv *RemoteSourceValidator) Validate(_ ...any) (any, Return.Error) {
	return nil, Return.Error{}
}
