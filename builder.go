package GoPlug

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoPlug/pluggable"
	"github.com/MickMake/GoPlug/utils"
	"github.com/MickMake/GoUnify/Only"
)

// BuildPlugins implements the interface method
func (bm *PluginManager) BuildPlugins() Return.Error {
	for range Only.Once {
		// scan plugin base dir
		var paths utils.PluginPaths
		paths, bm.Error = bm.loader.Scan("so", "go")
		if bm.Error.IsError() {
			break
		}

		log.Printf("[INFO]: Found %d possible plugins", paths.Length())

		if !paths.AnyPaths() {
			// No plugin files found
			break
		}

		for _, pPaths := range paths {
			pPaths.Dir.SetAltPath(bm.loader.GetDir(), "[PluginDir]")

			log.Printf("[INFO]: Plugin(%s): Analyzing plugin files", pPaths.Dir.String())
			fSo := utils.PluginPath{}
			fGo := utils.PluginPath{}

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
			fGo.SetAltPath(bm.loader.GetDir(), "[PluginDir]")
			fSo.SetAltPath(bm.loader.GetDir(), "[PluginDir]")

			if fGo.GetMod().IsZero() {
				// We don't have a *.go file, so nothing to build.
				log.Printf("[INFO]: Plugin(%s): No GoLang src code found for '%s'", pPaths.Dir.String(), fSo.GetBase())
				continue
			}

			if fSo.GetMod().IsZero() {
				// If *.go file is more recent than the *.so file.
				log.Printf("[INFO]: Plugin(%s): Plugin file '%s' needs building - *.so file missing", pPaths.Dir.String(), fGo.GetBase())
				bm.Error = bm.buildPlugin(fGo)
				if bm.Error.IsError() {
					continue
				}

				fSo = fGo.ChangeExtension("so")
				bm.Error = bm.loadPlugin(fSo)
				if bm.Error.IsError() {
					continue
				}

				bm.Error = bm.unloadPlugin(fSo)
				if bm.Error.IsError() {
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
				bm.Error = bm.buildPlugin(fGo)
				if bm.Error.IsError() {
					continue
				}

				bm.Error = bm.loadPlugin(fSo)
				// if bm.Error.IsError() {
				// 	continue
				// }

				bm.Error = bm.unloadPlugin(fSo)
				if bm.Error.IsError() {
					continue
				}

				continue
			}

			log.Printf("[INFO]: Plugin(%s): Plugin file '%s' is OK", pPaths.Dir.String(), fGo.GetBase())
		}

		log.Printf("[INFO]: %d plugins loaded", bm.store.Size())
	}

	return bm.Error
}

// BuildPlugin implements the interface method
func (bm *PluginManager) BuildPlugin(p string) Return.Error {
	for range Only.Once {
		name := bm.stringToPluginPath(p)
		if bm.Error.IsError() {
			break
		}

		bm.Error = bm.buildPlugin(name)
		if bm.Error.IsError() {
			break
		}
	}

	return bm.Error
}

func (bm *PluginManager) buildPlugin(goFile utils.PluginPath) Return.Error {
	for range Only.Once {
		bm.Error = goFile.FileExists()
		if bm.Error.IsError() {
			break
		}

		base := goFile.SetAltPath(bm.loader.GetDir(), "[PluginDir]")
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
		bm.Error.SetError(e)
		if bm.Error.IsError() {
			log.Printf("[ERROR]: Plugin(%s): Build failed: %s", base, bm.Error.String())
			break
		}

		log.Printf("[INFO]: Plugin(%s): Build OK", base)
	}

	return bm.Error
}

func (bm *PluginManager) CheckPlugin(p string) (*pluggable.PluginItem, Return.Error) {
	var plug *pluggable.PluginItem

	for range Only.Once {
		name := bm.stringToPluginPath(p)
		if bm.Error.IsError() {
			break
		}

		bm.Error = name.FileExists()
		if bm.Error.IsError() {
			break
		}

		plug, bm.Error = bm.store.Get(name.GetBase())
		if bm.Error.IsError() {
			// bm.Error.SetError("plugin with name '%s' is not existing", name.GetBase())
			break
		}

		log.Printf("[INFO]: Plugin '%s' loaded with version '%s'\n", plug.Config.Name, plug.Config.Version)
		// pContext := pluggable.NewPlugin()

		plug.Plugin.Set("master", "OK")
		bm.Error = plug.Plugin.Init()
		if bm.Error.IsError() {
			break
		}

		fmt.Printf("%v\n", plug.Config)

		bm.Error = plug.Plugin.Validate()
		if bm.Error.IsError() {
			break
		}
	}

	return plug, bm.Error
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
				fmt.Printf("Chdir to %s FAILED\n", e.Dir)
				break
			}
		}

		cmd := exec.Command(e.Command, e.Args...)
		fmt.Printf("Exec START: %s\n", cmd.String())

		var stdout, stderr []byte
		var errStdout, errStderr error
		stdoutIn, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()
		err = cmd.Start()
		if err != nil {
			fmt.Printf("cmd.Start() failed with '%s'\n", err)
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
			fmt.Printf("Exec Error: %s\n", err)
			break
		}

		if errStdout != nil || errStderr != nil {
			fmt.Printf("failed to capture stdout or stderr\n")
			break
		}

		// outStr, errStr := string(stdout), string(stderr)
		if len(stdout) > 0 {
			fmt.Printf("stdout:\n%s\n", string(stdout))
		}
		if len(stderr) > 0 {
			fmt.Printf("stderr:\n%s\n", string(stderr))
		}

		// fmt.Printf("Exec STOP: %s\n", cmd.String())
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
