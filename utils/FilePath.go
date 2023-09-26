package utils

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/MickMake/GoUnify/Only"
	"github.com/h2non/filetype"

	"github.com/MickMake/GoPlug/utils/Return"
)

const (
	ErrorIsNil         = "PluginPath is nil"
	ErrorBaseIsEmpty   = "FilePath[%s]: Error PluginPath base is empty"
	ErrorDirIsEmpty    = "FilePath[%s]: Error PluginPath dir is empty"
	ErrorPathIsEmpty   = "FilePath[%s]: Path is empty"
	ErrorNoGlob        = "FilePath[%s]: Error no file glob specified"
	ErrorIO            = "FilePath[%s]: Error %s"
	ErrorNotDir        = "FilePath[%s]: Is NOT a directory"
	ErrorNotDirButFile = ErrorNotDir + ", but a file"
	ErrorNotFileButDir = "FilePath[%s]: Is NOT a file, but a directory"
)
const AltPathString = "..."

//
// FileInterface
// ---------------------------------------------------------------------------------------------------- //
type FilePathInterface interface {
	NewDynamicData()
	ValueExists(key string) bool
	ValueNotExists(key string) bool
	SetValue(key string, value any)
	GetValue(key string) any
	String() string
}

//
// FilePaths
// ---------------------------------------------------------------------------------------------------- //
type FilePaths map[string]pluginPathMap

func NewFilePaths() FilePaths {
	return make(FilePaths)
}

func (p *FilePaths) AddPaths(paths FilePaths) {
	// @TODO - not tested yet!
	for name, path := range paths {
		if _, ok := (*p)[name]; !ok {
			// Create key/value
			(*p)[name] = path
			continue
		}
		// Else append values.
		tp := (*p)[name]
		tp.Paths = append(tp.Paths, path.Paths...)
	}
}

func (p *FilePaths) Add(paths ...FilePath) {
	for _, path := range paths {
		name := path.GetPath()
		ext := filepath.Ext(name)
		name = strings.TrimSuffix(name, ext)

		if _, ok := (*p)[name]; !ok {
			dir := path
			dir.path = path.GetDir()
			dir.base = ""
			(*p)[name] = pluginPathMap{
				Dir:   dir,
				Paths: []FilePath{path},
			}
			return
		}

		pp := (*p)[name]
		pp.Paths = append(pp.Paths, path)
		(*p)[name] = pp
	}
}

func (p *FilePaths) KeepExtensions(ext ...string) {
	for i, dir := range *p {
		var keep []FilePath
		for _, path := range dir.Paths {
			if path.HasExtension(ext...) {
				keep = append(keep, path)
			}
		}
		if keep == nil {
			delete(*p, i)
		} else {
			dir.Paths = keep
			(*p)[i] = dir
		}
	}
}

func (p *FilePaths) RemoveExtensions(ext ...string) {
	for i, dir := range *p {
		var keep []FilePath
		for _, path := range dir.Paths {
			if path.HasExtension(ext...) {
				continue
			}
			keep = append(keep, path)
		}
		dir.Paths = keep
		(*p)[i] = dir
	}
}

func (p *FilePaths) Length() int {
	return len(*p)
}

func (p *FilePaths) AnyPaths() bool {
	if len(*p) > 0 {
		return true
	}
	return false
}

//
// pluginPathMap
// ---------------------------------------------------------------------------------------------------- //
type pluginPathMap struct {
	Dir   FilePath
	Paths []FilePath
}

func (p *pluginPathMap) Length() int {
	return len(p.Paths)
}

func (p *pluginPathMap) Get() []FilePath {
	return p.Paths
}

//
// FilePath
// ---------------------------------------------------------------------------------------------------- //
type FilePath struct {
	path         string
	dir          string
	base         string
	isFile       bool
	isDir        bool
	isExecutable bool
	fStat        os.FileInfo
	fType        bool
	// modified time.Time

	basePath    string
	baseReplace string
	shortenPath bool

	Error Return.Error
}

func NewFile(path ...string) (FilePath, Return.Error) {
	var ret FilePath
	err := ret.SetFile(path...)
	return ret, err
}

func NewDir(path ...string) (FilePath, Return.Error) {
	var ret FilePath
	err := ret.SetDir(path...)
	return ret, err
}

// NewFile - .
func (p *FilePath) NewFile(filename string) Return.Error {
	var err Return.Error
	*p, err = NewFile(filename)
	return err
}

// NewDir - .
func (p *FilePath) NewDir(filename string) Return.Error {
	var err Return.Error
	*p, err = NewFile(filename)
	return err
}

func (p FilePath) String() string {
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

			ret = filepath.Join(p.dir[:16]+" ... "+p.dir[index-8:], p.base)
			break
		}
	}
	return ret
}

func (p *FilePath) GetName() string {
	return strings.TrimSuffix(p.base, filepath.Ext(p.base))
}

func (p *FilePath) GetPath() string {
	return p.path
}

func (p *FilePath) GetBase() string {
	return p.base
}

func (p *FilePath) GetDir() string {
	return p.dir
}

func (p *FilePath) GetMod() time.Time {
	return p.fStat.ModTime()
}

func (p *FilePath) BeginsWithPath(path string) bool {
	if strings.HasPrefix(p.path, path) {
		return true
	}
	return false
}

func (p *FilePath) SetAltPath(path string, replace string) string {
	p.basePath = path
	p.baseReplace = replace
	return p.String()
}

func (p *FilePath) ShortenPaths() {
	p.shortenPath = true
}

func (p *FilePath) HasExtension(ext ...string) bool {
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

func (p *FilePath) ChangeExtension(ext string) FilePath {
	ext = "." + strings.TrimPrefix(ext, ".")
	oldE := filepath.Ext(p.path)
	newP := strings.TrimSuffix(p.path, oldE) + ext
	n, _ := NewFile(newP)
	n.basePath = p.basePath
	n.baseReplace = p.baseReplace
	n.shortenPath = p.shortenPath
	return n
}

func (p *FilePath) IsValid() Return.Error {
	var err Return.Error

	switch {
	case p == nil:
		err.SetError(ErrorIsNil)
	case p.path == "":
		err.SetError(ErrorPathIsEmpty, p.path)
	case p.dir == "":
		err.SetError(ErrorDirIsEmpty, p.path)
	case p.base == "":
		err.SetError(ErrorBaseIsEmpty, p.path)
	case len(p.path) == 0:
		err.SetError(ErrorPathIsEmpty, p.path)
	case len(p.dir) == 0:
		err.SetError(ErrorDirIsEmpty, p.path)
	case len(p.base) == 0:
		err.SetError(ErrorBaseIsEmpty, p.path)
	}

	return err
}

func (p FilePath) AppendDir(paths ...string) (FilePath, Return.Error) {
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

func (p FilePath) PrependDir(paths ...string) (FilePath, Return.Error) {
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

func (p FilePath) AppendFile(paths ...string) (FilePath, Return.Error) {
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

func (p FilePath) PrependFile(paths ...string) (FilePath, Return.Error) {
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

func (p *FilePath) SetDir(paths ...string) Return.Error {
	var err Return.Error

	for range Only.Once {
		dir := filepath.Join(paths...)

		if dir == "" {
			err.SetError(ErrorPathIsEmpty, dir)
			break
		}

		var e error
		p.path, e = filepath.Abs(dir)
		if e != nil {
			err.SetError(ErrorIO, p.path, e)
			break
		}

		p.dir = filepath.Dir(p.path)
		p.base = filepath.Base(p.path)
		p.basePath = ""
		p.baseReplace = ""

		p.stat(p.path)
		// p.modified, err = DirExists(p.path)
		if err.IsError() {
			break
		}
		if !p.isDir {
			err.SetError(ErrorNotDir, p.path)
			break
		}
		// p.isDir = true
		// p.isFile = false
	}

	return err
}

func (p *FilePath) SetFile(paths ...string) Return.Error {
	var err Return.Error

	for range Only.Once {
		file := filepath.Join(paths...)

		if file == "" {
			err.SetError(ErrorPathIsEmpty, file)
			break
		}

		var e error
		p.path, e = filepath.Abs(file)
		if e != nil {
			err.SetError(ErrorIO, p.path, e)
			break
		}

		p.dir = filepath.Dir(p.path)
		p.base = filepath.Base(p.path)
		p.basePath = ""
		p.baseReplace = ""

		p.stat(p.path)
		// p.modified, err = FileExists(p.path)
		if err.IsError() {
			break
		}
		if p.isDir {
			err.SetError(ErrorNotFileButDir, p.path)
			break
		}
		// p.isDir = false
		// p.isFile = true
	}

	return err
}

func (p *FilePath) DirExists() Return.Error {
	var err Return.Error

	for range Only.Once {
		err = p.IsValid()
		if err.IsError() {
			break
		}

		p.stat(p.path)
		// p.modified, err = DirExists(p.path)
		if err.IsError() {
			break
		}
		if !p.isDir {
			err.SetError(ErrorNotDirButFile, p.path)
			break
		}
		// p.isDir = true
		// p.isFile = false
	}

	return err
}

func (p *FilePath) FileExists() Return.Error {
	var err Return.Error

	for range Only.Once {
		err = p.IsValid()
		if err.IsError() {
			break
		}

		p.stat(p.path)
		// p.modified, err = FileExists(p.path)
		if err.IsError() {
			break
		}
		if p.isDir {
			err.SetError(ErrorNotFileButDir, p.path)
			break
		}
		// p.isDir = false
		// p.isFile = true
	}

	return err
}

func (p *FilePath) Scan(glob string) (FilePaths, Return.Error) {
	candidates := NewFilePaths()
	var err Return.Error

	for range Only.Once {
		// if glob == "" {
		// 	err.SetError(ErrorNoGlob, p.path)
		// 	break
		// }

		err = p.DirExists()
		if err.IsError() {
			break
		}

		files, e := os.ReadDir(p.GetPath())
		if e != nil {
			err.SetError(ErrorIO, p.path, e)
			break
		}

		re := regexp.MustCompile(glob)

		for _, f := range files {
			if f.IsDir() {
				dir := filepath.Join(p.GetPath(), f.Name())
				dirs, e2 := os.ReadDir(dir)
				if e2 != nil {
					err.SetError(ErrorIO, p.path, e2)
					continue
				}

				for _, d := range dirs {
					if d.IsDir() {
						continue
					}

					var pp FilePath
					pp, err = p.AppendFile(f.Name(), d.Name())
					if err.IsError() {
						continue
					}
					if glob != "" {
						if !re.MatchString(pp.base) {
							continue
						}
					}
					// if !pp.HasExtension(ext...) {
					// 	continue
					// }
					pp.SetAltPath(p.GetDir(), AltPathString)
					candidates.Add(pp)
				}
				continue
			}

			var pp FilePath
			pp, err = p.AppendFile(f.Name())
			if err.IsError() {
				continue
			}
			if glob != "" {
				if !re.MatchString(pp.base) {
					continue
				}
			}
			// if !pp.HasExtension(ext...) {
			// 	continue
			// }
			pp.SetAltPath(p.GetDir(), AltPathString)
			candidates.Add(pp)
		}

		err = Return.Ok
	}

	return candidates, err
}

func (p *FilePath) ScanForExtension(ext ...string) (FilePaths, Return.Error) {
	candidates := NewFilePaths()
	var err Return.Error

	for range Only.Once {
		if len(ext) == 0 {
			err.SetError(ErrorNoGlob, p.path)
			break
		}

		glob := "(" + strings.Join(ext, "|") + ")$"
		candidates, err = p.Scan(glob)
	}

	return candidates, err
}

func (p *FilePath) ScanForExecutable() (FilePaths, Return.Error) {
	candidates := NewFilePaths()
	var err Return.Error

	for range Only.Once {
		candidates, err = p.Scan("")
		if err.IsError() {
			break
		}

		re := regexp.MustCompile(`\.\w+$`)

		for name, candidate := range candidates {
			var pp []FilePath
			for _, c := range candidate.Paths {
				if re.MatchString(c.base) {
					if !c.HasExtension(".exe") {
						// Accommodate Windows.
						continue
					}
				}
				if c.fStat.IsDir() {
					// Just to be sure.
					continue
				}
				if !strings.Contains(c.fStat.Mode().String(), "x") {
					// Execute bit not set - @TODO - check Windows.
					continue
				}
				pp = append(pp, c)
			}
			if len(pp) == 0 {
				delete(candidates, name)
				continue
			}
			candidate.Paths = pp
		}
	}

	return candidates, err
}

// SaveObject - Save an arbitrary plugin structure as a JSON file.
func (p *FilePath) SaveObject(ref any) Return.Error {
	var err Return.Error

	for range Only.Once {
		data, e := json.Marshal(ref)
		if e != nil {
			err.SetError(e)
			break
		}

		file := p.ChangeExtension("json")
		err = WriteFile(file.path, data)
		if err.IsError() {
			break
		}
	}

	return err
}

// LoadObject - Load a JSON file into an arbitrary plugin structure.
func (p *FilePath) LoadObject(ref any) Return.Error {
	var err Return.Error

	for range Only.Once {
		var data []byte
		file := p.ChangeExtension("json")
		data, err = ReadFile(file.path)
		if err.IsError() {
			break
		}

		e := json.Unmarshal(data, ref)
		if e != nil {
			err.SetError(e)
			break
		}
	}

	return err
}

// stat check the existence of the specified file
// If file exists, return true
func (p *FilePath) stat(path string) Return.Error {
	var err Return.Error

	for range Only.Once {
		if path == "" {
			err.SetError(ErrorPathIsEmpty, path)
			break
		}

		var e error
		p.fStat, e = os.Stat(path)
		if e != nil {
			err.SetError(ErrorIO, p.path, e)
			break
		}

		p.isDir = false
		p.isFile = false
		p.isExecutable = false
		if p.fStat.IsDir() {
			p.isDir = true
			err = Return.Ok
			break
		}

		// Must be a file then.
		p.isFile = true

		if !strings.Contains(p.fStat.Mode().String(), "x") {
			// Execute bit not set, ignoring it as an executable.
			// @TODO - check Windows.
			break
		}

		// Determine if it's an executable based on mimetype of file.
		// Limit the number of bytes read - we only need at most 32 bytes.
		var f *os.File
		f, e = os.Open(p.path)
		if e != nil {
			err.SetError(ErrorIO, p.path, e)
			break
		}
		defer f.Close()

		var buf [256]byte
		_, e = io.ReadFull(f, buf[:])
		if e != nil {
			if errors.Is(e, io.ErrUnexpectedEOF) {
			} else if errors.Is(e, io.EOF) {
			} else {
				err.SetError(ErrorIO, p.path, e)
				break
			}
		}

		kind, _ := filetype.Match(buf[:])

		switch {
		case kind.MIME.Type == "application":
			fallthrough
		case strings.HasSuffix(path, ".exe"):
			fallthrough
		case strings.HasSuffix(path, ".sh"):
			fallthrough
		case strings.HasSuffix(path, ".bash"):
			p.isExecutable = true
		}
		// fmt.Printf("FILE[%s]: '%v' %v\n", p.path, p.isExecutable, kind)
		err = Return.Ok
	}

	return err
}

// ---------------------------------------------------------------------------------------------------- //

// FileExists check the existence of the specified file
// If file exists, return true
func FileExists(path string) (time.Time, Return.Error) {
	var mod time.Time
	var err Return.Error

	for range Only.Once {
		if path == "" {
			err.SetError("FilePath '%s' is empty", path)
			break
		}

		fi, e := os.Stat(path)
		if e != nil {
			err.SetError("Error with file '%s': %v", path, e)
			break
		}

		if fi.IsDir() {
			err.SetError("FilePath '%s' is a directory", path)
			break
		}

		mod = fi.ModTime()
		err = Return.Ok
	}

	return mod, err
}

// DirExists check the existence of the specified dir
// If dir exists, return true
func DirExists(dir string) (time.Time, Return.Error) {
	var mod time.Time
	var err Return.Error

	for range Only.Once {
		if dir == "" {
			err.SetError("Directory '%s' is empty", dir)
			break
		}

		fi, e := os.Stat(dir)
		if e != nil {
			err.SetError("Error with directory '%s': %v", dir, e)
			break
		}

		if !fi.IsDir() {
			err.SetError("FilePath '%s' is NOT a directory", dir)
			break
		}

		mod = fi.ModTime()
		err = Return.Ok
	}

	return mod, err
}

// IsDir checks if the file is a dir
func IsDir(filePath string) bool {
	fi, err := os.Stat(filePath)

	return err == nil && fi.Mode().IsDir()
}

// ReadFile - read file
func ReadFile(filePath string) ([]byte, Return.Error) {
	var data []byte
	var err Return.Error

	for range Only.Once {
		fi, e := os.Stat(filePath)
		err.SetError(e)
		if err.IsError() {
			break
		}

		if fi.Mode().IsDir() {
			err.SetError("file '%s' is a directory", filePath)
			break
		}

		data, e = os.ReadFile(filePath)
		err.SetError(e)
		if err.IsError() {
			break
		}
	}

	return data, err
}

// WriteFile - read file
func WriteFile(filePath string, data []byte) Return.Error {
	var err Return.Error

	for range Only.Once {
		dir := filepath.Dir(filePath)
		fi, e := os.Stat(dir)
		err.SetError(e)
		if err.IsError() {
			break
		}

		if !fi.Mode().IsDir() {
			err.SetError("path '%s' is not a directory", dir)
			break
		}

		e = os.WriteFile(filePath, data, 0644)
		err.SetError(e)
		if err.IsError() {
			break
		}
	}

	return err
}

// GetFiles - Gets files that are in a given directory.
// The directory doesn't need to be absolute. For example, "." will work fine.
func GetFiles(glob, dir string) ([]string, error) {
	var err error

	// Make the directory absolute if it isn't already
	if !filepath.IsAbs(dir) {
		dir, err = filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
	}

	return filepath.Glob(filepath.Join(dir, glob))
}
