package edge_services

import (
	"fmt"

	"github.com/aziontech/azion-cli/cmd/edge_services/create"
	"github.com/aziontech/azion-cli/cmd/edge_services/resources"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

func NewCmdEdgeServices() *cobra.Command {
	// edgeServicesCmd represents the edgeServices command
	edgeServicesCmd := &cobra.Command{
		Use:   "edge_services",
		Short: "Manages all edge-services of a client",
		Long:  `You may create, update, delete, list and describe services of an authenticated azion account.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("edgeServices called")
		},
	}

	edgeServicesCmd.AddCommand(resources.NewCmdResources())
	edgeServicesCmd.AddCommand(create.NewCmdCreate())
	edgeServicesCmd.PersistentFlags().StringP("name", "n", "", "<EDGE_SERVICE_NAME>")
	edgeServicesCmd.MarkPersistentFlagRequired("name")

	return edgeServicesCmd
}
