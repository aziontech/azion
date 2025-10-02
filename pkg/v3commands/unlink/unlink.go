package unlink

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/unlink"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	app "github.com/aziontech/azion-cli/pkg/v3commands/delete/edge_application"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type UnlinkCmd struct {
	ShouldClean func(f *cmdutil.Factory) bool
	Clean       func(f *cmdutil.Factory, cmd *UnlinkCmd) error
	IsDirEmpty  func(dirpath string) (bool, error)
	CleanDir    func(dirpath string) error
	F           *cmdutil.Factory
	DeleteCmd   func(f *cmdutil.Factory) *app.DeleteCmd
}

func NewUnlinkCmd(f *cmdutil.Factory) *UnlinkCmd {
	return &UnlinkCmd{
		F:           f,
		ShouldClean: shouldClean,
		Clean:       clean,
		IsDirEmpty:  utils.IsDirEmpty,
		CleanDir:    utils.CleanDirectory,
		DeleteCmd:   app.NewDeleteCmd,
	}
}

func NewCobraCmd(unlink *UnlinkCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion unlink
		$ azion unlink --help
		$ azion unlink --yes
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return unlink.run()
		},
	}

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewUnlinkCmd(f), f)
}

func (cmd *UnlinkCmd) run() error {
	if shouldClean(cmd.F) {
		err := cmd.Clean(cmd.F, cmd)
		if err != nil {
			return err
		}
		fmt.Fprint(cmd.F.IOStreams.Out, msg.UnlinkSuccess)
	}
	return nil
}
