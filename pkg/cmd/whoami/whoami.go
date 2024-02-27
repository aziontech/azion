package whoami

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/whoami"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
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

			logger.FInfo(f.IOStreams.Out, settings.Email+"\n")
			return nil
		},
	}

	whoamiCmd.SetIn(f.IOStreams.In)
	whoamiCmd.SetOut(f.IOStreams.Out)
	whoamiCmd.SetErr(f.IOStreams.Err)

	whoamiCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return whoamiCmd
}
