package whoami

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/whoami"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	versionCmd := &cobra.Command{
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

	versionCmd.SetIn(f.IOStreams.In)
	versionCmd.SetOut(f.IOStreams.Out)
	versionCmd.SetErr(f.IOStreams.Err)

	versionCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return versionCmd
}
