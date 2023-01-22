package domains

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/domains"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/list"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	edge_applicationsCmd := &cobra.Command{
		Use:   msg.DomainsUsage,
		Short: msg.DomainsShortDescription,
		Long:  msg.DomainsLongDescription,
		Example: heredoc.Doc(`
		$ azioncli domains --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	edge_applicationsCmd.AddCommand(list.NewCmd(f))
	edge_applicationsCmd.Flags().BoolP("help", "h", false, msg.DomainsFlagHelp)

	return edge_applicationsCmd
}
