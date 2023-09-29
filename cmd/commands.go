package cmd

import (
	"github.com/MickMake/GoUnify/Only"
	"github.com/spf13/cobra"

	// 	"github.com/MickMake/GoUnify/Unify"

	"github.com/MickMake/GoUnify/Unify"

	"github.com/MickMake/GoPlug/defaults"
)

type Cmds struct {
	Unify *Unify.Unify
	Api   *CmdApi

	ConfigDir   string
	CacheDir    string
	ConfigFile  string
	WriteConfig bool
	Quiet       bool
	Debug       bool

	Args []string

	Error error
}

//goland:noinspection GoNameStartsWithPackageName
type CmdDefault struct {
	Error   error
	cmd     *cobra.Command
	SelfCmd *cobra.Command
}

var cmds Cmds

func init() {
	for range Only.Once {
		cmds.Unify = Unify.New(
			Unify.Options{
				Description:   defaults.Description,
				BinaryName:    defaults.BinaryName,
				BinaryVersion: defaults.BinaryVersion,
				SourceRepo:    defaults.SourceRepo,
				BinaryRepo:    defaults.BinaryRepo,
				EnvPrefix:     defaults.EnvPrefix,
				HelpSummary:   defaults.HelpSummary,
				ReadMe:        defaults.Readme,
				Examples:      defaults.Examples,
			},
			Unify.Flags{},
		)

		cmdRoot := cmds.Unify.GetCmd()

		cmds.Api = NewCmdApi()
		cmds.Api.AttachCommand(cmdRoot)
		cmds.Api.AttachFlags(cmdRoot, cmds.Unify.GetViper())

		// cmds.Info = NewCmdInfo()
		// cmds.Info.AttachCommand(cmdRoot)
	}
}

func Execute() error {
	var err error

	for range Only.Once {
		// Execute adds all child commands to the root command and sets flags appropriately.
		// This is called by main.main(). It only needs to happen once to the rootCmd.
		err = cmds.Unify.Execute()
		if err != nil {
			break
		}
	}

	return err
}

func (ca *Cmds) ProcessArgs(_ *cobra.Command, args []string) error {
	for range Only.Once {
		ca.Args = args

		ca.ConfigDir = cmds.Unify.GetConfigDir()
		ca.ConfigFile = cmds.Unify.GetConfigFile()
		ca.CacheDir = cmds.Unify.GetCacheDir()
		ca.Debug = cmds.Unify.Flags.Debug
		ca.Quiet = cmds.Unify.Flags.Quiet
	}

	return ca.Error
}
