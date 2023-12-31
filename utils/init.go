package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func init() {
	// Exec functions and cache results on load.
	IsPlugin()
	GetArgs()
}

// ---------------------------------------------------------------------------------------------------- //

var cacheIsPlugin *bool
var True = true
var False = false

func IsPlugin() bool {
	if cacheIsPlugin != nil {
		return *cacheIsPlugin
	}

	cacheIsPlugin = &False
	callers := GetCallers(0)
	for _, caller := range callers {
		if strings.HasPrefix(caller, "plugin") {
			cacheIsPlugin = &True
			break
		}
		if strings.HasPrefix(caller, "main") {
			cacheIsPlugin = &False
			break
		}
	}

	return *cacheIsPlugin
}

// ---------------------------------------------------------------------------------------------------- //

var cacheGetArgs []string

func GetArgs() []string {
	if len(cacheGetArgs) > 0 {
		return cacheGetArgs
	}

	cacheGetArgs = os.Args
	dir, err := filepath.Abs(cacheGetArgs[0])
	if err == nil {
		cacheGetArgs[0] = dir
	}
	return cacheGetArgs
}

// ---------------------------------------------------------------------------------------------------- //

func GetCallers(depth int) []string {
	var callers []string
	if depth == 0 {
		depth = 32
	}
	for index := 0; index < depth; index++ {
		pc, _, _, _ := runtime.Caller(index + 1)
		name := runtime.FuncForPC(pc).Name()
		if name == "" {
			break
		}
		callers = append(callers, name)
	}
	return callers
}

// ---------------------------------------------------------------------------------------------------- //
