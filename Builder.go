package GoPlug

import (
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/MickMake/GoUnify/Only"

	"github.com/MickMake/GoPlug/GoPlugLoader"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoPlug/utils/Return"
)

// ---------------------------------------------------------------------------------------------------- //

// BuildPlugins implements the interface method
func (m *PluginManager) BuildPlugins() Return.Error {
	for range Only.Once {
		loader := m.Loaders.GetLoader(GoPlugLoader.NativeLoaderName)
		if loader == nil {
			break
		}

		// scan plugin base dir
		m.Error = loader.PluginScanByExtension(GoPlugLoader.NativePluginExtensions...)
		if m.Error.IsError() {
			break
		}

		paths := loader.GetFiles()
		log.Printf("[INFO]: Found %d possible plugins", paths.Length())

		if !paths.AnyPaths() {
			// No plugin files found
			break
		}

		for _, pPaths := range paths {
			pPaths.Dir.SetAltPath(loader.GetDir(), "[PluginDir]")

			log.Printf("[INFO]: Plugin(%s): Analyzing plugin files", pPaths.Dir.String())
			fSo := utils.FilePath{}
			fGo := utils.FilePath{}

			for _, pName := range pPaths.Get() {
				if pName.HasExtension("go") {
					fGo = pName
					continue
				}
				if pName.HasExtension("so") {
					fSo = pName
					continue
				}
			}
			fGo.SetAltPath(loader.GetDir(), "[PluginDir]")
			fSo.SetAltPath(loader.GetDir(), "[PluginDir]")

			if fGo.GetMod().IsZero() {
				// We don't have a *.go file, so nothing to build.
				log.Printf("[INFO]: Plugin(%s): No GoLang src code found for '%s'", pPaths.Dir.String(), fSo.GetBase())
				continue
			}

			if fSo.GetMod().IsZero() {
				// If *.go file is more recent than the *.so file.
				log.Printf("[INFO]: Plugin(%s): Plugin file '%s' needs building - *.so file missing", pPaths.Dir.String(), fGo.GetBase())
				m.Error = m.buildPlugin(fGo)
				if m.Error.IsError() {
					continue
				}

				fSo = fGo.ChangeExtension("so")
				m.Error = m.LoadPlugin(fSo)
				if m.Error.IsError() {
					continue
				}

				m.Error = m.UnloadPlugin(fSo)
				if m.Error.IsError() {
					continue
				}

				continue
			}

			if fGo.GetMod().After(fSo.GetMod()) {
				// If *.go file is more recent than the *.so file.
				log.Printf("[INFO]: Plugin(%s): Plugin needs re-building - %s (%s) is newer than %s (%s)",
					pPaths.Dir.String(),
					fGo.GetBase(),
					fGo.GetMod().Format("2006-01-02 15:04:05"),
					fSo.GetBase(),
					fSo.GetMod().Format("2006-01-02 15:04:05"),
				)
				m.Error = m.buildPlugin(fGo)
				if m.Error.IsError() {
					continue
				}

				m.Error = m.LoadPlugin(fSo)
				if m.Error.IsError() {
					continue
				}

				m.Error = m.UnloadPlugin(fSo)
				if m.Error.IsError() {
					continue
				}

				continue
			}

			log.Printf("[INFO]: Plugin(%s): Plugin file '%s' is OK", pPaths.Dir.String(), fGo.GetBase())
		}

		log.Printf("[INFO]: %d plugins loaded", m.Loaders.StoreSize())
	}

	return m.Error
}

// BuildPlugin implements the interface method
func (m *PluginManager) BuildPlugin(pluginPath utils.FilePath) Return.Error {
	for range Only.Once {
		loader := m.Loaders.GetLoader(GoPlugLoader.NativeLoaderName)
		if loader == nil {
			break
		}

		// name := m.stringToPluginPath(p)
		// if m.Error.IsError() {
		// 	break
		// }

		m.Error = m.buildPlugin(pluginPath)
		if m.Error.IsError() {
			break
		}
	}

	return m.Error
}

func (m *PluginManager) buildPlugin(goFile utils.FilePath) Return.Error {
	for range Only.Once {
		loader := m.Loaders.GetLoader(GoPlugLoader.NativeLoaderName)
		if loader == nil {
			break
		}

		m.Error = goFile.FileExists()
		if m.Error.IsError() {
			break
		}

		base := goFile.SetAltPath(loader.GetDir(), "[PluginDir]")
		log.Printf("[INFO]: Plugin(%s): Building", base)

		soFile := goFile.ChangeExtension("so")

		options := ExecOptions{
			Dir:     goFile.GetDir(),
			Command: "go",
			Args: []string{
				"build",
				"-buildmode=plugin",
				"-o", soFile.GetBase(),
				"-gcflags", "all=-N -l",
				goFile.GetBase(),
			},
		}
		e := options.Exec()
		m.Error.SetError(e)
		if m.Error.IsError() {
			log.Printf("[ERROR]: Plugin(%s): Build failed: %s", base, m.Error.String())
			break
		}

		log.Printf("[INFO]: Plugin(%s): Build OK", base)
	}

	return m.Error
}

func (m *PluginManager) CheckPlugin(pluginPath utils.FilePath) (*GoPlugLoader.PluginItem, Return.Error) {
	var plug *GoPlugLoader.PluginItem

	for range Only.Once {
		// name := m.stringToPluginPath(p)
		// if m.Error.IsError() {
		// 	break
		// }

		m.Error = pluginPath.FileExists()
		if m.Error.IsError() {
			break
		}

		plug, m.Error = m.Loaders.StoreGet(pluginPath.GetBase())
		if m.Error.IsError() {
			// bm.Error.SetError("plugin with name '%s' is not existing", name.GetBase())
			break
		}

		log.Printf("[INFO]: Plugin '%s' loaded with version '%s'\n", plug.Data.Dynamic.Identity.Name, plug.Data.Dynamic.Identity.Version)
		// pContext := pluggable.NewPlugin()

		// // plug.Native.Values().KeySet("master", "OK")
		// m.Error = plug.Data.Callback().Init()
		// if m.Error.IsError() {
		// 	break
		// }
		//
		// log.Printf("%v\n", plug.Data.Dynamic.Identity)
		//
		// m.Error = plug.Native.Validate()
		// if m.Error.IsError() {
		// 	break
		// }
	}

	return plug, m.Error
}

type ExecOptions struct {
	Dir     string
	Command string
	Args    []string
}

func (e *ExecOptions) Exec() error {
	var err error

	for range Only.Once {
		if e.Dir != "" {
			err = os.Chdir(e.Dir)
			if err != nil {
				log.Printf("Chdir to %s FAILED\n", e.Dir)
				break
			}
		}

		cmd := exec.Command(e.Command, e.Args...)
		log.Printf("Exec START: %s\n", cmd.String())

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutIn, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()
		err = cmd.Start()
		if err != nil {
			log.Printf("cmd.Start() failed with '%s'\n", err)
			break
		}

		// cmd.Wait() should be called only after we finish reading
		// from stdoutIn and stderrIn.
		// wg ensures that we finish
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			stdout, errStdout = copyAndCapture(os.Stdout, stdoutIn)
			wg.Done()
		}()
		stderr, errStderr = copyAndCapture(os.Stderr, stderrIn)
		wg.Wait()

		err = cmd.Wait()
		if err != nil {
			log.Printf("Exec Error: %s\n", err)
			break
		}

		if errStdout != nil || errStderr != nil {
			log.Printf("failed to capture stdout or stderr\n")
			break
		}

		// outStr, errStr := string(stdout), string(stderr)
		if len(stdout) > 0 {
			log.Printf("stdout:\n%s\n", string(stdout))
		}
		if len(stderr) > 0 {
			log.Printf("stderr:\n%s\n", string(stderr))
		}

		// log.Printf("Exec STOP: %s\n", cmd.String())
	}

	return err
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}
