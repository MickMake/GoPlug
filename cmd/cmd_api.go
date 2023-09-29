package cmd

import (
	"fmt"
	"time"

	"github.com/MickMake/GoUnify/Only"
	"github.com/MickMake/GoUnify/cmdConfig"
	"github.com/MickMake/GoUnify/cmdHelp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagApiUrl        = "host"
	flagApiTimeout    = "timeout"
	flagApiUsername   = "user"
	flagApiPassword   = "password"
	flagApiAppKey     = "appkey"
	flagApiLastLogin  = "token-expiry"
	flagApiOutputType = "out"
	flagApiSaveFile   = "save"
	flagApiDirectory  = "dir"
)

//goland:noinspection GoNameStartsWithPackageName
type CmdApi struct {
	CmdDefault

	// iSolarCloud api
	ApiTimeout   time.Duration
	Url          string
	Username     string
	Password     string
	AppKey       string
	LastLogin    string
	ApiToken     string
	ApiTokenFile string
	OutputType   string
	SaveFile     bool
	Directory    string

	// GoPlug      *iSolarCloud.GoPlug
}

func NewCmdApi() *CmdApi {
	var ret *CmdApi

	for range Only.Once {
		ret = &CmdApi{
			CmdDefault: CmdDefault{
				Error:   nil,
				cmd:     nil,
				SelfCmd: nil,
			},
			// ApiTimeout:   iSolarCloud.DefaultTimeout,
			// Url:          iSolarCloud.DefaultHost,
			Username: "",
			Password: "",
			// AppKey:       iSolarCloud.DefaultApiAppKey,
			LastLogin:    "",
			ApiToken:     "",
			ApiTokenFile: "",
			OutputType:   "",
			// GoPlug:      nil,
		}
	}

	return ret
}

func (c *CmdApi) AttachCommand(cmd *cobra.Command) *cobra.Command {
	for range Only.Once {
		if cmd == nil {
			break
		}
		c.cmd = cmd

		// ******************************************************************************** //
		var cmdApi = &cobra.Command{
			Use:                   "api",
			Aliases:               []string{},
			Annotations:           map[string]string{"group": "Api"},
			Short:                 fmt.Sprintf("Low-level interface to the GoPlug api."),
			Long:                  fmt.Sprintf("Low-level interface to the GoPlug api."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               nil,
			Run:                   c.CmdApi,
			Args:                  cobra.MinimumNArgs(1),
		}
		cmd.AddCommand(cmdApi)
		cmdApi.Example = cmdHelp.PrintExamples(cmdApi, "get <endpoint>", "put <endpoint>")

		// ******************************************************************************** //
		var cmdApiList = &cobra.Command{
			Use:                   "ls",
			Aliases:               []string{"list"},
			Annotations:           map[string]string{"group": "Api"},
			Short:                 fmt.Sprintf("List GoPlug api endpoints/areas"),
			Long:                  fmt.Sprintf("List GoPlug api endpoints/areas"),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               cmds.GoPlugArgs,
			Run:                   c.CmdApiList,
			Args:                  cobra.RangeArgs(0, 1),
		}
		cmdApi.AddCommand(cmdApiList)
		cmdApiList.Example = cmdHelp.PrintExamples(cmdApiList, "", "areas", "endpoints", "<area name>")

		// ******************************************************************************** //
		var cmdApiLogin = &cobra.Command{
			Use:                   "login",
			Aliases:               []string{},
			Annotations:           map[string]string{"group": "Api"},
			Short:                 fmt.Sprintf("Login to the GoPlug api."),
			Long:                  fmt.Sprintf("Login to the GoPlug api."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               cmds.GoPlugArgs,
			RunE: func(cmd *cobra.Command, args []string) error {
				c.Error = c.ApiLogin(true)
				if c.Error != nil {
					return c.Error
				}
				return nil
			},
			Args: cobra.MinimumNArgs(0),
		}
		cmdApi.AddCommand(cmdApiLogin)
		cmdApiLogin.Example = cmdHelp.PrintExamples(cmdApiLogin, "")

		// ******************************************************************************** //
		var cmdApiStruct = &cobra.Command{
			Use:                   "get",
			Aliases:               []string{},
			Annotations:           map[string]string{"group": "Api"},
			Short:                 fmt.Sprintf("Show response as Go structure (debug)"),
			Long:                  fmt.Sprintf("Show response as Go structure (debug)"),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               cmds.GoPlugArgs,
			RunE:                  c.CmdApiGet,
			Args:                  cobra.MinimumNArgs(1),
		}
		cmdApi.AddCommand(cmdApiStruct)
		cmdApiStruct.Example = cmdHelp.PrintExamples(cmdApiStruct, "[area].<endpoint>")

		// ******************************************************************************** //
		var cmdApiPut = &cobra.Command{
			Use:                   "put",
			Aliases:               []string{"write"},
			Annotations:           map[string]string{"group": "Api"},
			Short:                 fmt.Sprintf("Put details onto the GoPlug api."),
			Long:                  fmt.Sprintf("Put details onto the GoPlug api."),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               cmds.GoPlugArgs,
			Run:                   c.CmdApiPut,
			Args:                  cobra.RangeArgs(0, 1),
		}
		cmdApi.AddCommand(cmdApiPut)
		cmdApiPut.Example = cmdHelp.PrintExamples(cmdApiPut, "[area].<endpoint> <value>")
	}
	return c.SelfCmd
}

func (c *CmdApi) AttachFlags(cmd *cobra.Command, viper *viper.Viper) {
	for range Only.Once {
		cmd.PersistentFlags().StringVarP(&c.Username, flagApiUsername, "u", "", fmt.Sprintf("GoPlug: api username."))
		viper.SetDefault(flagApiUsername, "")
		cmd.PersistentFlags().StringVarP(&c.Password, flagApiPassword, "p", "", fmt.Sprintf("GoPlug: api password."))
		viper.SetDefault(flagApiPassword, "")
		cmd.PersistentFlags().StringVar(&c.LastLogin, flagApiLastLogin, "", "GoPlug: last login.")
		viper.SetDefault(flagApiLastLogin, "")
		// _ = cmd.PersistentFlags().MarkHidden(flagApiLastLogin)

		cmd.PersistentFlags().StringVarP(&c.OutputType, flagApiOutputType, "o", "", fmt.Sprintf("Output type: 'json', 'raw', 'file'"))
		_ = cmd.PersistentFlags().MarkHidden(flagApiOutputType)
		cmd.PersistentFlags().BoolVarP(&c.SaveFile, flagApiSaveFile, "s", false, "Save output as a file.")
		viper.SetDefault(flagApiSaveFile, false)
		cmd.PersistentFlags().StringVarP(&c.Directory, flagApiDirectory, "", "", "Save output base directory.")
		viper.SetDefault(flagApiDirectory, "")
	}
}

func (ca *Cmds) GoPlugArgs(cmd *cobra.Command, args []string) error {
	for range Only.Once {
		ca.Error = cmds.ProcessArgs(cmd, args)
		if ca.Error != nil {
			break
		}
	}

	return ca.Error
}

func (ca *Cmds) SetOutputType(cmd *cobra.Command) error {
	var err error
	for range Only.Once {
		// foo := cmd.Parent()
	}

	return err
}

func (c *CmdApi) CmdApi(cmd *cobra.Command, args []string) {
	for range Only.Once {
		if len(args) == 0 {
			c.Error = cmd.Help()
			break
		}
	}
}

func (c *CmdApi) CmdApiList(cmd *cobra.Command, args []string) {
	for range Only.Once {
		switch {
		case len(args) == 0:
			fmt.Println("Unknown sub-command.")
			_ = cmd.Help()

		case args[0] == "foo":

		case args[0] == "bar":

		default:
		}
	}
}

func (c *CmdApi) CmdApiGet(_ *cobra.Command, args []string) error {
	for range Only.Once {
		args = MinimumArraySize(2, args)
		if args[0] == "all" {
			break
		}

		if c.Error != nil {
			break
		}
	}

	return c.Error
}

func (c *CmdApi) CmdApiPut(_ *cobra.Command, _ []string) {
	for range Only.Once {
		fmt.Println("Not yet implemented.")
		// args = fillArray(1, args)
		// c.Error = GoPlug.Init()
		// if c.Error != nil {
		// 	break
		// }
	}
}

func (c *CmdApi) ApiLogin(force bool) error {
	for range Only.Once {
		if force {
			sf := cmds.Api.SaveFile
			cmds.Api.SaveFile = false // We don't want to lock this in the config.
			c.Error = cmds.Unify.WriteConfig()
			cmds.Api.SaveFile = sf
		}
	}
	return c.Error
}

func MinimumArraySize(count int, args []string) []string {
	var ret []string
	for range Only.Once {
		ret = cmdConfig.FillArray(count, args)
		for i, e := range args {
			if e == "." {
				e = ""
			}
			if e == "-" {
				e = ""
			}
			ret[i] = e
		}
	}
	return ret
}
