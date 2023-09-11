package pluggable

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"time"

	"github.com/MickMake/GoPlug/Return"
	"github.com/MickMake/GoUnify/Only"
)


// ---------------------------------------------------------------------------------------------------- //

// PluginCallback - Generic callback function type
type PluginCallback func(ctx Plugin, args ...interface{}) Return.Error

func (c PluginCallback) MarshalJSON() ([]byte, error) {
	// func Greet(name string) {
	//	fmt.Printf("Hello, %s!\n", name)
	// }
	//
	// func main() {
	//	funcValue := reflect.ValueOf(Greet)
	//
	//	parameters := []reflect.Value{
	//		reflect.ValueOf("Alice"),
	//	}
	//
	//	funcValue.Call(parameters) // Output: Hello, Alice!
	// }

	str := `"` + runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name() + `"`
	return []byte(str), nil
}

func (c PluginCallback) GetName() string {
	var name string
	for range Only.Once {
		name = runtime.FuncForPC(reflect.ValueOf(c).Pointer()).Name()

		re := regexp.MustCompile(`^.*?\.(\(\*[A-Za-z0-9_-]+\)\.[A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 1 {
				name = a[1] + "()"
				break
			}
		}

		re = regexp.MustCompile(`^.*?\.([A-Za-z0-9_]+)`)
		if re.MatchString(name) {
			a := re.FindStringSubmatch(name)
			if len(a) >= 1 {
				name = a[1] + "()"
				break
			}
		}
	}
	return name
}

func (c *PluginCallback) IsValid() Return.Error {
	var err Return.Error

	if reflect.ValueOf(*c).IsNil() {
		err.SetWarning("callback is not defined")
	}
	if c == nil {
		err.SetWarning("callback is not defined")
	}
	return err
}


// ---------------------------------------------------------------------------------------------------- //

type Callbacks struct {
	pluginName string

	// Initialise - Called on plugin load.
	Initialise PluginCallback `json:"-"`
	funcInitialise string

	// Run - Execute a function concurrently.
	Run PluginCallback `json:"-"`
	funcRun string

	// Notify - Notify a plugin.
	Notify PluginCallback `json:"-"`
	funcNotify string

	// Execute - Execute a function, should return.
	Execute PluginCallback `json:"-"`
	funcExecute string
}

func (c Callbacks) MarshalJSON() ([]byte, error) {
	str1 := c.pluginName + "." + c.Initialise.GetName()
	str2 := c.pluginName + "." + c.Run.GetName()
	str3 := c.pluginName + "." + c.Notify.GetName()
	str4 := c.pluginName + "." + c.Execute.GetName()

	str := fmt.Sprintf(`{ "Initialise":"%s", "Run":"%s", "Notify":"%s", "Execute":"%s" }`,
		str1, str2, str3, str4,
	)
	return []byte(str), nil
}

func (c *Callbacks) IsValid() Return.Error {
	var err Return.Error
	switch {
		case c == nil:
			err.SetError("callbacks is nil")
		default:
			if c.Initialise == nil {
				err.SetWarning("callback Initialise() is nil")
			}
			if c.Run == nil {
				err.AddWarning("callback Run() is nil")
			}
			if c.Notify == nil {
				err.AddWarning("callback Notify() is nil")
			}
			if c.Execute == nil {
				err.AddWarning("callback Execute() is nil")
			}
	}
	return err
}


// ---------------------------------------------------------------------------------------------------- //

func (i *PluginItem) Execute(args ...interface{}) Return.Error {
	var err Return.Error
	for range Only.Once {
		err.SetPrefix("Execute(): ")
		if i == nil {
			err.SetError("Callbacks.Run() is nil")
			break
		}

		i.Plugin.Set("slave-execute-timestamp", time.Now())

		err = i.IsValid()
		if err.IsError() {
			i.Plugin.Set("slave-execute", err)
			break
		}

		err = i.Config.Callbacks.Execute.IsValid()
		if err.IsError() || err.IsWarning() {
			i.Plugin.Set("slave-execute", err)
			break
		}

		err = i.Config.Callbacks.Execute(i.Plugin, args...)
		if err.IsError() {
			i.Plugin.Set("slave-execute", err)
			break
		}

		i.Plugin.Set("slave-execute", "OK")
	}
	return err
}

func (i *PluginItem) Initialise(args ...interface{}) Return.Error {
	var err Return.Error
	for range Only.Once {
		err.SetPrefix("Initialise(): ")
		if i == nil {
			err.SetError("Callbacks.Run() is nil")
			break
		}

		i.Plugin.Set("slave-init-timestamp", time.Now())

		err = i.IsValid()
		if err.IsError() {
			i.Plugin.Set("slave-init", err)
			break
		}

		err = i.Config.Callbacks.Initialise.IsValid()
		if err.IsError() || err.IsWarning() {
			i.Plugin.Set("slave-init", err)
			break
		}

		err = i.Config.Callbacks.Initialise(i.Plugin, args...)
		if err.IsError() {
			i.Plugin.Set("slave-init", err)
			break
		}

		i.Plugin.Set("slave-init", "OK")
	}
	return err
}

func (i *PluginItem) Run(args ...interface{}) Return.Error {
	var err Return.Error
	for range Only.Once {
		err.SetPrefix("Run(): ")
		if i == nil {
			err.SetError("Callbacks.Run() is nil")
			break
		}

		i.Plugin.Set("slave-run-timestamp", time.Now())

		err = i.IsValid()
		if err.IsError() {
			i.Plugin.Set("slave-run", err)
			break
		}

		err = i.Config.Callbacks.Run.IsValid()
		if err.IsError() || err.IsWarning() {
			i.Plugin.Set("slave-run", err)
			break
		}

		go func() {
			fmt.Println("i.Config.Callbacks.Run")
			err2 := i.Config.Callbacks.Run(i.Plugin, args...)
			i.Plugin.Set("slave-run", err2)
		}()
	}
	return err
}

func (i *PluginItem) Notify(args ...interface{}) Return.Error {
	var err Return.Error
	for range Only.Once {
		err.SetPrefix("Notify(): ")

		if i == nil {
			err.SetError("Callbacks.Run() is nil")
			break
		}

		err = i.IsValid()
		if err.IsError() {
			break
		}

		err = i.Config.Callbacks.Notify.IsValid()
		if err.IsError() || err.IsWarning() {
			err.SetError(err)
			break
		}

		err = i.Config.Callbacks.Notify(i.Plugin, args...)
		if err.IsError() {
			break
		}
	}
	return err
}
