package presets

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/presets"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/command"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io                    *iostreams.IOStreams
	F                     *cmdutil.Factory
	CommandRunner         func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
	CommandRunInteractive func(f *cmdutil.Factory, comm string) error
	Vulcan                func() *vulcanPkg.VulcanPkg
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		F:  f,
		CommandRunner: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
			return command.CommandRunInteractiveWithOutput(f, comm, envVars)
		},
		CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
			return command.CommandRunInteractive(f, comm)
		},
		Vulcan: vulcanPkg.NewVulcan,
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list presets
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return list.run()
		},
	}

	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	listCmd := NewListCmd(f)
	return NewCobraCmd(listCmd, f)
}

func (list *ListCmd) run() error {
	vulcanVer, err := list.CommandRunner(list.F, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	vul := list.Vulcan()
	if err := vul.CheckVulcanMajor(vulcanVer, list.F, vul); err != nil {
		return err
	}

	logger.FInfo(list.Io.Out, msg.GettingPresets)

	command := vul.Command("--loglevel=error --no-update-notifier", "presets ls", list.F)
	return list.CommandRunInteractive(list.F, command)
}
