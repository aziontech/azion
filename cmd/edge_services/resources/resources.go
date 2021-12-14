package resources

import (
	"fmt"

	"github.com/aziontech/azion-cli/cmd/edge_services/resources/create"
	"github.com/aziontech/azion-cli/cmd/edge_services/resources/delete"
	"github.com/aziontech/azion-cli/cmd/edge_services/resources/describe"
	"github.com/aziontech/azion-cli/cmd/edge_services/resources/list"
	"github.com/spf13/cobra"
)

func NewCmdResources() *cobra.Command {
	// resourcesCmd represents the resources command
	resourcesCmd := &cobra.Command{
		Use:   "resources",
		Short: "Manages resources in a given edge-service",
		Long:  `You may create, update, delete, list and describe resources in a given edge-service.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("resources called")
		},
	}
	resourcesCmd.AddCommand(list.NewCmd())
	resourcesCmd.AddCommand(describe.NewCmd())
	resourcesCmd.AddCommand(delete.NewCmd())
	resourcesCmd.AddCommand(create.NewCmd())
	return resourcesCmd
}
