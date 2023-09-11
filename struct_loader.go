package GoPlug

import (
	"os"
	"path/filepath"
	sysPlugin "plugin"
	"reflect"

	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoPlug/pluggable"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoUnify/Only"
)

// ---------------------------------------------------------------------------------------------------- //

// Loader defines the plugin load flow
type Loader interface {
	// SetDir - sets the plugin base dir
	SetDir(dir string) Return.Error

	// GetDir - get the plugin base dir
	GetDir() string

	// Scan the plugin base dir and get the plugin candidates
	Scan(ext ...string) (utils.PluginPaths, Return.Error)

	// Get the plugin identity config
	Parse(path utils.PluginPath) (*pluggable.PluginIdentity, Return.Error)

	// Load the plugin
	Load(path utils.PluginPath) (*pluggable.PluginItem, Return.Error)

	// Init the plugin
	Init(item *pluggable.PluginItem) Return.Error
}

// ---------------------------------------------------------------------------------------------------- //

// PluginLoader is a default implementation of Loader interface
type PluginLoader struct {
	baseDir utils.PluginPath
	Error   Return.Error
}

// NewPluginLoader is constructor of PluginLoader
func NewPluginLoader() *PluginLoader {
	var err Return.Error
	err.SetPrefix("PluginLoader: ")
	return &PluginLoader{
		baseDir: utils.PluginPath{},
		Error:   err,
	}
}

// SetDir - sets the plugin base dir
func (bl *PluginLoader) SetDir(dir string) Return.Error {
	for range Only.Once {
		if dir == "" {
			var e error
			dir, e = os.Getwd()
			bl.Error.SetError(e)
			if bl.Error.IsError() {
				break
			}
		}

		bl.Error = bl.baseDir.SetDir(dir)
	}

	return bl.Error
}

// GetDir - Gets the plugin base dir
func (bl *PluginLoader) GetDir() string {
	return bl.baseDir.GetPath()
}

// Scan - Implements same method of Loader interface
func (bl *PluginLoader) Scan(ext ...string) (utils.PluginPaths, Return.Error) {
	candidates := utils.NewPluginPaths()

	for range Only.Once {
		bl.Error = bl.baseDir.DirExists()
		if bl.Error.IsError() {
			break
		}

		files, err2 := os.ReadDir(bl.baseDir.GetPath())
		bl.Error.SetError(err2)
		if bl.Error.IsError() {
			break
		}

		for _, f := range files {
			if f.IsDir() {
				dir := filepath.Join(bl.baseDir.GetPath(), f.Name())
				dirs, e2 := os.ReadDir(dir)
				bl.Error.SetError(e2)
				if bl.Error.IsError() {
					continue
				}

				for _, d := range dirs {
					if d.IsDir() {
						continue
					}

					var p utils.PluginPath
					p, bl.Error = bl.baseDir.AppendFile(f.Name(), d.Name())
					if bl.Error.IsError() {
						continue
					}
					if !p.HasExtension(ext...) {
						continue
					}
					p.SetAltPath(bl.baseDir.GetDir(), utils.AltPathString)
					candidates.Add(p)
					// candidates = append(candidates, p)
				}
				continue
			}

			var p utils.PluginPath
			p, bl.Error = bl.baseDir.AppendFile(f.Name())
			if bl.Error.IsError() {
				continue
			}
			if !p.HasExtension(ext...) {
				continue
			}
			p.SetAltPath(bl.baseDir.GetDir(), utils.AltPathString)
			candidates.Add(p)
			// candidates = append(candidates, p)
		}

		bl.Error = Return.Ok
	}

	return candidates, bl.Error
}

// Parse implements same method of Loader interface
func (bl *PluginLoader) Parse(path utils.PluginPath) (*pluggable.PluginIdentity, Return.Error) {
	bl.Error.SetError("not implemented in PluginLoader: %s", path)
	return nil, bl.Error
}

// Load implements same method of Loader interface
func (bl *PluginLoader) Load(path utils.PluginPath) (*pluggable.PluginItem, Return.Error) {
	var plug pluggable.PluginItem

	for range Only.Once {
		blConfig := &pluggable.PluginIdentity{}

		bl.Error = path.FileExists()
		if bl.Error.IsError() {
			break
		}

		p, err2 := sysPlugin.Open(path.GetPath())
		bl.Error.SetError(err2)
		if bl.Error.IsError() {
			break
		}

		var sym sysPlugin.Symbol
		sym, err2 = p.Lookup("GoPlugin")
		bl.Error.SetError(err2)
		if bl.Error.IsError() {
			break
		}

		var ok bool
		blConfig, ok = sym.(*pluggable.PluginIdentity)
		if !ok {
			bl.Error.SetError("plugin structure not defined properly in file '%s' - type is '%s'",
				path, reflect.TypeOf(sym))
			break
		}

		bl.Error = blConfig.IsValid()
		if bl.Error.IsError() {
			break
		}

		plug.File = path
		plug.Config = blConfig
		plug.Plugin = pluggable.NewPlugin(blConfig)
		bl.Error = plug.Plugin.SetConfigDir(path.GetDir())
	}

	return &plug, bl.Error
}

// Init implements same method of Loader interface
func (bl *PluginLoader) Init(item *pluggable.PluginItem) Return.Error {
	for range Only.Once {
		bl.Error = item.Plugin.Init()
		if bl.Error.IsError() {
			break
		}

		bl.Error = item.Initialise()
		if bl.Error.IsError() {
			break
		}
	}

	return bl.Error
}
