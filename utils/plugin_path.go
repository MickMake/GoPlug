package utils

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/MickMake/GoPlug/Return"

	"github.com/MickMake/GoUnify/Only"
)

type PluginPaths map[string]pluginPathMap
type pluginPathMap struct {
	Dir   PluginPath
	Paths []PluginPath
}

func (p *pluginPathMap) Length() int {
	return len(p.Paths)
}

func (p *pluginPathMap) Get() []PluginPath {
	return p.Paths
}

func NewPluginPaths() PluginPaths {
	return make(PluginPaths)
}

func (p *PluginPaths) Add(path PluginPath) {
	name := path.GetPath()
	ext := filepath.Ext(name)
	name = strings.TrimSuffix(name, ext)

	if _, ok := (*p)[name]; !ok {
		dir := path
		dir.path = path.GetDir()
		dir.base = ""
		(*p)[name] = pluginPathMap{
			Dir: dir,
			Paths: []PluginPath{ path },
		}
		return
	}

	pp := (*p)[name]
	pp.Paths = append(pp.Paths, path)
	(*p)[name] = pp
}

func (p *PluginPaths) Length() int {
	return len(*p)
}

func (p *PluginPaths) AnyPaths() bool {
	if len(*p) > 0 {
		return true
	}
	return false
}

type PluginPath struct {
	path     string
	dir      string
	base     string
	isFile   bool
	isDir    bool
	modified time.Time

	basePath    string
	baseReplace string
	shortenPath bool
}

func NewFile(paths ...string) (PluginPath, Return.Error) {
	var ret PluginPath
	path := filepath.Join(paths...)
	err := ret.SetFile(path)
	return ret, err
}

func NewDir(paths ...string) (PluginPath, Return.Error) {
	var ret PluginPath
	path := filepath.Join(paths...)
	err := ret.SetDir(path)
	return ret, err
}

func (p PluginPath) String() string {
	ret := p.path
	for range Only.Once {
		if p.shortenPath {
			// @TODO - Add the path shortening code in here.
			ret = p.path
			break
		}

		if p.basePath != "" {
			if !p.BeginsWithPath(p.basePath) {
				break
			}

			// re := regexp.MustCompile("^" + path)
			// ret = re.ReplaceAllString(p.base, replace)
			if p.baseReplace != "..." {
				ret = p.baseReplace + strings.TrimPrefix(p.path, p.basePath)
				break
			}

			index := len(p.dir)
			if index < 16 {
				break
			}

			ret = filepath.Join(p.dir[:16] + " ... " + p.dir[index-8:], p.base)
			break
		}
	}
	return ret
}

func (p *PluginPath) GetPath() string {
	return p.path
}

func (p *PluginPath) GetBase() string {
	return p.base
}

func (p *PluginPath) GetDir() string {
	return p.dir
}

func (p *PluginPath) GetMod() time.Time {
	return p.modified
}

func (p *PluginPath) BeginsWithPath(path string) bool {
	if strings.HasPrefix(p.path, path) {
		return true
	}
	return false
}

const AltPathString = "..."

func (p *PluginPath) SetAltPath(path string, replace string) string {
	p.basePath = path
	p.baseReplace = replace
	return p.String()
}

func (p *PluginPath) ShortenPaths() {
	p.shortenPath = true
}

func (p *PluginPath) HasExtension(ext ...string) bool {
	var yes bool

	for range Only.Once {
		if len(ext) == 0 {
			break
		}

		for _, e := range ext {
			e = "." + strings.TrimPrefix(e, ".")
			if strings.HasSuffix(p.base, e) {
				yes = true
				break
			}
		}
	}

	return yes
}

func (p *PluginPath) ChangeExtension(ext string) PluginPath {
	ext = "." + strings.TrimPrefix(ext, ".")
	oldE := filepath.Ext(p.path)
	newP := strings.TrimSuffix(p.path, oldE) + ext
	n, _ := NewFile(newP)
	n.basePath = p.basePath
	n.baseReplace = p.baseReplace
	n.shortenPath = p.shortenPath
	return n
}

func (p *PluginPath) IsValid() Return.Error {
	var err Return.Error

	switch {
	case p == nil:
		err.SetError("PluginPath is nil")
	case p.path == "":
		err.SetError("PluginPath path is empty")
	case p.dir == "":
		err.SetError("PluginPath dir is empty")
	case p.base == "":
		err.SetError("PluginPath base is empty")
	case len(p.path) == 0:
		err.SetError("PluginPath path is empty")
	case len(p.dir) == 0:
		err.SetError("PluginPath dir is empty")
	case len(p.base) == 0:
		err.SetError("PluginPath base is empty")
	}

	return err
}

func (p PluginPath) AppendDir(paths ...string) (PluginPath, Return.Error) {
	var err Return.Error

	for range Only.Once {
		if len(paths) == 0 {
			break
		}

		pa := []string{p.path}
		pa = append(pa, paths...)
		path := filepath.Join(pa...)
		err = p.SetDir(path)
	}

	return p, err
}

func (p PluginPath) PrependDir(paths ...string) (PluginPath, Return.Error) {
	var err Return.Error

	for range Only.Once {
		if len(paths) == 0 {
			break
		}

		paths = append(paths, p.path)
		path := filepath.Join(paths...)
		err = p.SetDir(path)
	}

	return p, err
}

func (p PluginPath) AppendFile(paths ...string) (PluginPath, Return.Error) {
	var err Return.Error

	for range Only.Once {
		if len(paths) == 0 {
			break
		}

		pa := []string{p.path}
		pa = append(pa, paths...)
		path := filepath.Join(pa...)
		err = p.SetFile(path)
	}

	return p, err
}

func (p PluginPath) PrependFile(paths ...string) (PluginPath, Return.Error) {
	var err Return.Error

	for range Only.Once {
		if len(paths) == 0 {
			break
		}

		paths = append(paths, p.path)
		path := filepath.Join(paths...)
		err = p.SetFile(path)
	}

	return p, err
}

func (p *PluginPath) SetDir(dir string) Return.Error {
	var err Return.Error

	for range Only.Once {
		if dir == "" {
			err.SetError("%s is not a valid directory path", dir)
			break
		}

		var e error
		p.path, e = filepath.Abs(dir)
		err.SetError(e)
		if err.IsError() {
			break
		}

		p.dir = filepath.Dir(p.path)
		p.base = filepath.Base(p.path)
		p.basePath = ""
		p.baseReplace = ""

		p.modified, err = DirExists(p.path)
		if err.IsError() {
			break
		}

		p.isDir = true
		p.isFile = false
	}

	return err
}

func (p *PluginPath) SetFile(file string) Return.Error {
	var err Return.Error

	for range Only.Once {
		if file == "" {
			err.SetError("%s is not a valid file path", file)
			break
		}

		var e error
		p.path, e = filepath.Abs(file)
		err.SetError(e)
		if err.IsError() {
			break
		}

		p.dir = filepath.Dir(p.path)
		p.base = filepath.Base(p.path)
		p.basePath = ""
		p.baseReplace = ""

		p.modified, err = FileExists(p.path)
		if err.IsError() {
			break
		}

		p.isDir = false
		p.isFile = true
	}

	return err
}

func (p *PluginPath) DirExists() Return.Error {
	var err Return.Error

	for range Only.Once {
		err = p.IsValid()
		if err.IsError() {
			break
		}

		_, err = DirExists(p.path)
		if err.IsError() {
			break
		}

		p.isDir = true
		p.isFile = false
	}

	return err
}

func (p *PluginPath) FileExists() Return.Error {
	var err Return.Error

	for range Only.Once {
		err = p.IsValid()
		if err.IsError() {
			break
		}

		_, err = FileExists(p.path)
		if err.IsError() {
			break
		}

		p.isDir = false
		p.isFile = true
	}

	return err
}
