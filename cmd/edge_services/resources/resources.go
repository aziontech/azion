package resources

import (
	"fmt"

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
	resourcesCmd.AddCommand(list.NewCmdList())
	resourcesCmd.AddCommand(describe.NewCmdDescribe())
	// resourcesCmd.AddCommand(delete.NewCmdDelete())
	return resourcesCmd
}
