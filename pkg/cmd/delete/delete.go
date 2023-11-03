package delete

import (
	"github.com/MakeNowJust/heredoc"
	domain "github.com/aziontech/azion-cli/pkg/cmd/delete/domain"
	edgeApplication "github.com/aziontech/azion-cli/pkg/cmd/delete/edge_application"
	origin "github.com/aziontech/azion-cli/pkg/cmd/delete/origin"
	token "github.com/aziontech/azion-cli/pkg/cmd/delete/personal_token"
	rulesEngine "github.com/aziontech/azion-cli/pkg/cmd/delete/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	msg "github.com/aziontech/azion-cli/pkg/messages/delete"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion delete --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplication.NewCmd(f))
	cmd.AddCommand(rulesEngine.NewCmd(f))
	cmd.AddCommand(domain.NewCmd(f))
	cmd.AddCommand(token.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
