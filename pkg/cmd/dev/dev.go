package dev

import (
	"io"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/dev"
	"github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DevCmd struct {
	Io                  *iostreams.IOStreams
	CommandRunnerStream func(out io.Writer, cmd string, envvars []string) error
	BuildCmd            func(f *cmdutil.Factory) *build.BuildCmd
	F                   *cmdutil.Factory
}

func NewDevCmd(f *cmdutil.Factory) *DevCmd {
	return &DevCmd{
		F:        f,
		Io:       f.IOStreams,
		BuildCmd: build.NewBuildCmd,
		CommandRunnerStream: func(out io.Writer, cmd string, envs []string) error {
			return utils.RunCommandStreamOutput(f.IOStreams.Out, envs, cmd)
		},
	}
}

func NewCobraCmd(dev *DevCmd) *cobra.Command {
	devCmd := &cobra.Command{
		Use:           msg.DevUsage,
		Short:         msg.DevShortDescription,
		Long:          msg.DevLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`       
        $ azion dev
        $ azion dev --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dev.Run(dev.F)
		},
	}
	devCmd.Flags().BoolP("help", "h", false, msg.DevFlagHelp)
	return devCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDevCmd(f))
}

func (cmd *DevCmd) Run(f *cmdutil.Factory) error {
	logger.Debug("Running dev command")

	// Run build command
	build := cmd.BuildCmd(f)
	err := build.Run()
	if err != nil {
		logger.Debug("Error while running build command called by dev command", zap.Error(err))
		return err
	}

	err = vulcan(cmd)
	if err != nil {
		return err
	}

	return nil
}
