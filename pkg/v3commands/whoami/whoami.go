package whoami

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/whoami"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
)

type WhoamiCmd struct {
	Io           *iostreams.IOStreams
	ReadSettings func(string) (token.Settings, error)
	F            *cmdutil.Factory
}

func NewWhoamiCmd(f *cmdutil.Factory) *WhoamiCmd {
	return &WhoamiCmd{
		Io:           f.IOStreams,
		ReadSettings: token.ReadSettings,
		F:            f,
	}
}

func NewCobraCmd(whoami *WhoamiCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azion whoami
		$ azion whoami --help
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return whoami.run()
		},
	}

	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewWhoamiCmd(f), f)
}

func (cmd *WhoamiCmd) run() error {
	activeProfile := cmd.F.GetActiveProfile()
	settings, err := cmd.ReadSettings(activeProfile)
	if err != nil {
		return err
	}

	if settings.Email == "" {
		return msg.ErrorNotLoggedIn
	}

	msg := fmt.Sprintf(" Client ID: %s\n Email: %s\n Active Profile: %s\n", settings.ClientId, settings.Email, activeProfile)
	whoamiOut := output.GeneralOutput{
		Msg:   msg,
		Out:   cmd.Io.Out,
		Flags: cmd.F.Flags,
	}
	return output.Print(&whoamiOut)
}
