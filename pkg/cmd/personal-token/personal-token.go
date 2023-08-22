package personal_token

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/personal-token"
	"github.com/aziontech/azion-cli/pkg/cmd/personal-token/create"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	originsCmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion personal_token --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	originsCmd.AddCommand(create.NewCmd(f))
	originsCmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return originsCmd
}
