package edge_services

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_services"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/create"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/list"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/resources"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// edgeServicesCmd represents the edgeServices command
	edgeServicesCmd := &cobra.Command{
		Use:   msg.EdgeServiceUsage,
		Short: msg.EdgeServiceShortDescription,
		Long:  msg.EdgeServiceLongDescription,
		Example: heredoc.Doc(`
		$ azioncli edge_services --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	edgeServicesCmd.AddCommand(create.NewCmd(f))
	edgeServicesCmd.AddCommand(update.NewCmd(f))
	edgeServicesCmd.AddCommand(delete.NewCmd(f))
	edgeServicesCmd.AddCommand(list.NewCmd(f))
	edgeServicesCmd.AddCommand(describe.NewCmd(f))
	edgeServicesCmd.AddCommand(resources.NewCmd(f))

	edgeServicesCmd.Flags().BoolP("help", "h", false, msg.EdgeServiceHelpFlag)

	return edgeServicesCmd
}
