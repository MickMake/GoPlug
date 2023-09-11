package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/MickMake/GoPlug/Return"

	"github.com/MickMake/GoUnify/Only"
)

const (
	// PluginJSONFileName is the pre-defined filename of plugin metadata json file
	PluginJSONFileName = "plugin.json"

	// PluginSourceModeLocal defines the local mode
	PluginSourceModeLocal = "local_so"

	// PluginSourceModeRemote defines the remote mode
	PluginSourceModeRemote = "remote_git"
)

// FileExists check the existence of the specified file
// If file exists, return true
func FileExists(path string) (time.Time, Return.Error) {
	var mod time.Time
	var err Return.Error

	for range Only.Once {
		if path == "" {
			err.SetError("File '%s' is empty", path)
			break
		}

		fi, e := os.Stat(path)
		if e != nil {
			err.SetError("Error with file '%s': %v", path, e)
			break
		}

		if fi.IsDir() {
			err.SetError("File '%s' is a directory", path)
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
			err.SetError("File '%s' is NOT a directory", dir)
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
