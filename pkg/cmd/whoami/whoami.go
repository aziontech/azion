package whoami

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/whoami"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	whoamiCmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azion whoami
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			settings, err := token.ReadSettings()
			if err != nil {
				return err
			}

			if settings.Email == "" {
				return msg.ErrorNotLoggedIn
			}

			whoamiOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(settings.Email + "\n"),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&whoamiOut)
		},
	}

	whoamiCmd.SetIn(f.IOStreams.In)
	whoamiCmd.SetOut(f.IOStreams.Out)
	whoamiCmd.SetErr(f.IOStreams.Err)

	whoamiCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return whoamiCmd
}
