package utils

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/utils/Return"
)

//goland:noinspection GoUnusedConst
const (
	// PluginJSONFileName is the pre-defined filename of plugin metadata json file
	PluginJSONFileName = "plugin.json"

	// PluginSourceModeLocal defines the local mode
	PluginSourceModeLocal = "local_so"

	// PluginSourceModeRemote defines the remote mode
	PluginSourceModeRemote = "remote_git"
)

func GetTypeName(ref any) string {
	var name string

	for range Only.Once {
		kind := reflect.ValueOf(ref).Kind()
		if kind == reflect.Invalid {
			name = "nil"
			break
		}

		name = reflect.TypeOf(ref).String()
		name = strings.ReplaceAll(name, "interface {}", "any")
	}

	return name
}

//goland:noinspection GoUnusedExportedFunction
func GetTypeKind(ref any) reflect.Kind {
	return reflect.ValueOf(ref).Kind()
}

func IsTypeOfName(ref any, name string) bool {
	if reflect.TypeOf(ref).String() == name {
		return true
	}
	return false
}

func GetStructName(ref any) string {
	var name string
	for range Only.Once {
		if ref == nil {
			name = "nil"
			break
		}

		name = reflect.TypeOf(ref).String()
		re := regexp.MustCompile(`^.*?\.\(\*([A-Za-z0-9_-]+)\)(\.[A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 2 {
				name = a[1] + a[2]
				break
			}
		}

		re = regexp.MustCompile(`^.*?\.([A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 1 {
				name = a[1]
				break
			}
		}
	}
	return name
}

//goland:noinspection GoUnusedExportedFunction
func GetFunctionName(ref any) string {
	var name string
	for range Only.Once {
		if ref == nil {
			name = "nil"
			break
		}

		name = reflect.TypeOf(ref).String()
		re := regexp.MustCompile(`^.*?\.\(\*([A-Za-z0-9_-]+)\)(\.[A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 2 {
				name = a[1] + a[2]
				break
			}
		}

		re = regexp.MustCompile(`^.*?\.([A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 1 {
				name = a[1]
				break
			}
		}
	}
	return name
}

//goland:noinspection GoUnusedExportedFunction
func GetFunctionNameFromPointer(f interface{}) string {
	p := reflect.ValueOf(f).Pointer()
	rf := runtime.FuncForPC(p)
	name := rf.Name()
	name = filepath.Base(name)
	_, method := SeparatePackageAndFunction(name)
	return method
}

func GetPackageAndFunctionNameFromPointer(f interface{}) (string, string) {
	p := reflect.ValueOf(f).Pointer()
	rf := runtime.FuncForPC(p)
	name := rf.Name()
	name = filepath.Base(name)
	pkg, method := SeparatePackageAndFunction(name)
	return pkg, method
}

func SeparatePackageAndFunction(name string) (string, string) {
	var pkg string
	var method string
	for range Only.Once {
		re := regexp.MustCompile(`^(.*?)\.([A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 2 {
				pkg = a[1]
				method = a[2]
				break
			}
		}
	}
	return pkg, method
}

func GetCallerFunctionName(depth int) string {
	var name string
	for range Only.Once {
		pc, _, _, _ := runtime.Caller(depth + 1)

		name = runtime.FuncForPC(pc).Name()

		re := regexp.MustCompile(`^.*?\.\(\*[A-Za-z0-9_-]+\)\.([A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 1 {
				name = a[1]
				break
			}
		}

		re = regexp.MustCompile(`^.*?\.([A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 1 {
				name = a[1]
				break
			}
		}
	}
	return name
}

func GetCaller(depth int, args ...any) string {
	pc, _, _, _ := runtime.Caller(depth + 1)
	return runtime.FuncForPC(pc).Name() + "()"
}

func GetCallerDebug(depth int, args ...any) string {
	pc, file, line, _ := runtime.Caller(depth + 1)
	ret := fmt.Sprintf("%s:%d <= %s()", file, line, runtime.FuncForPC(pc).Name())
	return ret
}

//goland:noinspection GoUnusedExportedFunction
func MakeServiceCall() string {
	var name string
	for range Only.Once {
		pc, _, _, _ := runtime.Caller(1)

		name = runtime.FuncForPC(pc).Name()

		re := regexp.MustCompile(`^.*?\.\(\*([A-Za-z0-9_-]+)\)(\.[A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 2 {
				name = a[1] + a[2]
				break
			}
		}

		re = regexp.MustCompile(`^.*?\.([A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 1 {
				name = "Plugin." + a[1]
				break
			}
		}
	}
	return name
}

//goland:noinspection GoUnusedExportedFunction
func ToJson(ref any) ([]byte, Return.Error) {
	var err Return.Error
	data, e := json.Marshal(ref)
	if e != nil {
		err.SetError(e)
	}
	return data, err
}

//goland:noinspection GoUnusedExportedFunction
func FromJson(data []byte, ref any) Return.Error {
	var err Return.Error
	e := json.Unmarshal(data, &ref)
	if e != nil {
		err.SetError(e)
	}
	return err
}

func FileToId(file string, glob string) string {
	globAsterix := strings.LastIndex(glob, "*")
	trim := glob[0:globAsterix]
	return strings.TrimPrefix(file, trim)
}
