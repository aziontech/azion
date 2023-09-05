package create

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion create --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
