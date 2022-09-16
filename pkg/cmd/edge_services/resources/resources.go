package resources

import (
	msg "github.com/aziontech/azion-cli/messages/edge_services"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/resources/create"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/resources/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/resources/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/resources/list"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/resources/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// resourcesCmd represents the resources command
	resourcesCmd := &cobra.Command{
		Use:   msg.EdgeServiceResourceUsage,
		Short: msg.EdgeServiceResourceShortDescription,
		Long:  msg.EdgeServiceResourceLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	resourcesCmd.AddCommand(list.NewCmd(f))
	resourcesCmd.AddCommand(describe.NewCmd(f))
	resourcesCmd.AddCommand(delete.NewCmd(f))
	resourcesCmd.AddCommand(create.NewCmd(f))
	resourcesCmd.AddCommand(update.NewCmd(f))
	return resourcesCmd
}
