package Plugin

import (
	"fmt"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

const (
	ErrorIsNil = "PluginCommon is nil"
)

//
// CommonInterface
// ---------------------------------------------------------------------------------------------------- //
type CommonInterface interface {
	// ---------------------------------------- //
	// Structure validity methods

	// Init - Initialise this plugin.
	InitCommon() Return.Error
	SetRawInterface(ref any)
	GetRawInterface() any
	// IsCommonValid - Validate Common structure and set p.configured if true
	IsCommonValid() Return.Error
	// IsCommonConfigured - Quick check to determine if the plugin has been configured properly.
	IsCommonConfigured() bool
	// IsCommonError - Is the interface Errored?
	IsCommonError() bool
	// GetCommonError - Get the interface Error structure.
	GetCommonError() Return.Error
	// String - Standard stringer method.
	String() string
	// GetCommonRef - Get the underlying Common structure.
	GetCommonRef() *Common

	// ---------------------------------------- //
	// Logger methods

	// SetLogger - Configure the logger.
	SetLogger(logger *utils.Logger)
	// GetLogger - Get the logger.
	GetLogger() *utils.Logger
	// SetLogFile - Set an alternative file location for this plugin's logfiles.
	SetLogFile(filename string) Return.Error

	// ---------------------------------------- //
	// Identity based methods

	// SetPluginType - Set the plugin type.
	SetPluginType(name Types) Return.Error
	// SetPluginTypeNative - Set the plugin type to Native.
	SetPluginTypeNative() Return.Error
	// SetPluginTypeRpc - Set the plugin type to RPC.
	SetPluginTypeRpc() Return.Error
	// GetPluginType - Get the plugin type.
	GetPluginType() Types
	// SetStructName - Set the implementor structure name.
	SetStructName(ref interface{})
	// GetStructName - Get the implementor structure name, (used in RPC.Call).
	GetStructName() string

	// ---------------------------------------- //
	// File based methods

	// SetFilename - Set the file location of this plugin.
	SetFilename(pluginPath utils.FilePath) Return.Error
	// GetFilename - Get the file location of this plugin.
	GetFilename() utils.FilePath
	// SetDirectory - Set the directory that this plugin resides - to be used for saving/loading of files.
	SetDirectory(pluginPath utils.FilePath) Return.Error
	// GetDirectory - Get the directory that this plugin resides - to be used for saving/loading of files.
	GetDirectory() utils.FilePath
}

// NewPluginCommonInterface - Create a Common interface structure instance.
func NewPluginCommonInterface(common *Common) CommonInterface {
	ret := NewPluginCommon(common)
	return &ret
}

//
// Common
// ---------------------------------------------------------------------------------------------------- //
type Common struct {
	Id           string         `json:"id,omitempty"`
	PluginTypes  Types          `json:"plugin_types"`
	StructName   string         `json:"struct_name,omitempty"`
	Directory    utils.FilePath `json:"directory"`
	Filename     utils.FilePath `json:"filename"`
	Logger       *utils.Logger  `json:"-"`
	Configured   bool           `json:"configured,omitempty"`
	rawInterface any
	IsPlugin     bool           `json:"is_plugin,omitempty"`
	OsArgs       []string       `json:"os_args"`
	OsExecPath   utils.FilePath `json:"os_execpath"`
	Error        Return.Error   `json:"error"`
}

// NewPluginCommon - Create a store.ValueStore interface structure instance.
func NewPluginCommon(common *Common) Common {
	var err Return.Error

	osargs := utils.GetArgs()
	var execpath utils.FilePath
	if len(osargs) > 0 {
		execpath, err = utils.NewFile(osargs[0])
	}

	if common != nil {
		common.OsArgs = osargs
		common.OsExecPath = execpath
		// if common.Logger == nil {
		// 	// common.Logger = &l
		// }
		return *common
	}

	return Common{
		PluginTypes: Types{
			Rpc:    true,
			Native: true,
		},
		StructName:   GoPluginIdentity,
		Directory:    utils.FilePath{},
		Filename:     utils.FilePath{},
		Logger:       nil,
		Configured:   false,
		rawInterface: nil,
		IsPlugin:     utils.IsPlugin(),
		OsArgs:       osargs,
		OsExecPath:   execpath,
		Error:        err,
	}
}

// InitCommon - Initialise this plugin.
func (p *Common) InitCommon() Return.Error {
	if p == nil {
		return Return.NewError(ErrorIsNil)
	}

	// @TODO - What else?
	return Return.Ok
}

// SetRawInterface - .
func (p *Common) SetRawInterface(ref any) {
	if p == nil {
		return
	}
	p.rawInterface = ref
}

// GetRawInterface - .
func (p *Common) GetRawInterface() any {
	return p.rawInterface
}

// IsCommonValid - Validate Common structure and set p.configured if true
func (p *Common) IsCommonValid() Return.Error {
	if p == nil {
		return Return.NewError(ErrorIsNil)
	}

	for range Only.Once {
		p.Error.ReturnClear()
		if p.Configured {
			break
		}

		if p.StructName == "" {
			p.Error.SetError("plugin implementor structure not specified")
			break
		}

		p.Error = p.Directory.IsValid()
		if p.Error.IsError() {
			break
		}

		p.Error = p.Filename.IsValid()
		if p.Error.IsError() {
			break
		}

		p.Error = p.Logger.IsValid()
		if p.Error.IsError() {
			break
		}

		p.Configured = true
	}

	return p.Error
}

// IsCommonConfigured - Quick check to determine if the plugin has been configured properly.
func (p *Common) IsCommonConfigured() bool {
	if p == nil {
		return false
	}
	return p.Configured
}

// IsCommonError - Is the interface Errored?
func (p *Common) IsCommonError() bool {
	if p == nil {
		return true
	}
	return p.Error.IsError()
}

// GetCommonError - Get the interface Error structure.
func (p *Common) GetCommonError() Return.Error {
	if p == nil {
		return Return.NewError(ErrorIsNil)
	}
	return p.Error
}

// String - Stringer method.
func (p Common) String() string {
	var ret string

	ret += fmt.Sprintf("#### Plugin Type:\n%v\n", p.PluginTypes)
	ret += fmt.Sprintf("#### Is Configured?: %v\n", p.Configured)
	ret += fmt.Sprintf("#### Implementor Structure: %s\n", p.StructName)
	ret += fmt.Sprintf("#### Plugin Filename: %s\n", p.Filename.String())
	ret += fmt.Sprintf("#### Plugin Directory: %s\n", p.Directory.GetDir())
	ret += fmt.Sprintf("#### Logger Name:")
	if p.Logger == nil {
		ret += fmt.Sprintf(" Unknown\n")
	} else {
		ret += fmt.Sprintf(" %s\n", p.Logger.Gethclog().Name())
	}

	return ret
}

func (p *Common) Print() {
	fmt.Println(p.String())
}

// GetCommonRef - Get the underlying Common structure.
func (p *Common) GetCommonRef() *Common {
	return p
}

// ---------------------------------------------------------------------------------------------------- //
// Logger methods
//

// SetLogger - Configure the logger.
func (p *Common) SetLogger(logger *utils.Logger) {
	if p == nil {
		return
	}
	p.Logger = logger
}

// GetLogger - Get the logger.
func (p *Common) GetLogger() *utils.Logger {
	if p == nil {
		return nil
	}
	return p.Logger
}

// SetLogFile - Set an alternative file location for this plugin's logfiles.
func (p *Common) SetLogFile(filename string) Return.Error {
	if p == nil {
		p.Error.SetError("Common is nil")
		return p.Error
	}
	p.Error = p.Logger.SetLogFile(filename)
	p.Logger.Info("g.Logger.Info(): Logfile set to '%s'", filename)
	return p.Error
}

// ---------------------------------------------------------------------------------------------------- //
// Identity based methods
//

// SetPluginType - Set the plugin type.
func (p *Common) SetPluginType(types Types) Return.Error {
	p.Error = Return.Ok
	p.PluginTypes = types
	p.IsPlugin = true
	return p.Error
}

// SetPluginTypeNative - Set the plugin type to Native.
func (p *Common) SetPluginTypeNative() Return.Error {
	p.Error = Return.Ok
	p.PluginTypes.Native = true
	p.IsPlugin = true
	return p.Error
}

// SetPluginTypeRpc - Set the plugin type to RPC.
func (p *Common) SetPluginTypeRpc() Return.Error {
	p.Error = Return.Ok
	p.PluginTypes.Rpc = true
	p.IsPlugin = true
	return p.Error
}

// GetPluginType - Get the plugin type.
func (p *Common) GetPluginType() Types {
	if p == nil {
		return Types{}
	}
	return p.PluginTypes
}

// SetStructName - Set the implementor structure name.
func (p *Common) SetStructName(ref interface{}) {
	if p == nil {
		return
	}
	p.StructName = utils.GetStructName(ref)
}

// GetStructName - Get the implementor structure name, (used in RPC.Call).
func (p *Common) GetStructName() string {
	if p == nil {
		return ""
	}
	return p.StructName
}

// ---------------------------------------------------------------------------------------------------- //
// File based methods
//

// SetFilename - Set the file location of this plugin.
func (p *Common) SetFilename(pluginPath utils.FilePath) Return.Error {
	for range Only.Once {
		if p == nil {
			p.Error.SetError("Common struct is nil")
			break
		}

		p.Error = pluginPath.FileExists()
		if p.Error.IsError() {
			break
		}

		p.Filename = pluginPath

		var dir utils.FilePath
		dir, p.Error = utils.NewDir(pluginPath.GetDir())
		if p.Error.IsError() {
			break
		}

		p.Error = p.SetDirectory(dir)
	}

	return p.Error
}

// GetFilename - Get the file location of this plugin.
func (p *Common) GetFilename() utils.FilePath {
	if p == nil {
		return utils.FilePath{}
	}
	return p.Filename
}

// SetDirectory - Get the directory that this plugin resides - to be used for saving/loading of files.
func (p *Common) SetDirectory(pluginPath utils.FilePath) Return.Error {
	var err Return.Error
	for range Only.Once {
		if p == nil {
			p.Error.SetError("Common struct is nil")
			break
		}

		p.Error = pluginPath.DirExists()
		if p.Error.IsError() {
			break
		}

		p.Directory = pluginPath
	}
	return err
}

// GetDirectory - Get the directory that this plugin resides - to be used for saving/loading of files.
func (p *Common) GetDirectory() utils.FilePath {
	if p == nil {
		return utils.FilePath{}
	}
	return p.Directory
}

// ---------------------------------------------------------------------------------------------------- //
// Other methods
//

// ---------------------------------------------------------------------------------------------------- //

type GlobArgs []interface{}

type PutArgs struct {
	Key   string
	Value any
}

type GetArgs struct {
	Key string
}
