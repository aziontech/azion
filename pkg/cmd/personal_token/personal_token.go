package personal_token

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/personal-token"
	"github.com/aziontech/azion-cli/pkg/cmd/personal_token/create"
	"github.com/aziontech/azion-cli/pkg/cmd/personal_token/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/personal_token/list"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion personal_token --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(create.NewCmd(f))
	cmd.AddCommand(list.NewCmd(f))
	cmd.AddCommand(delete.NewCmd(f))
	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
