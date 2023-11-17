package logout

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/logout"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azion logout --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := token.DeleteToken()
			if err != nil {
				return err
			}

			logger.LogSuccessBad(f.IOStreams.Out, msg.Success)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
