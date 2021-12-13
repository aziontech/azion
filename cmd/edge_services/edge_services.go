package edge_services

import (
	"fmt"

	"github.com/aziontech/azion-cli/cmd/edge_services/create"
	"github.com/aziontech/azion-cli/cmd/edge_services/resources"
	"github.com/aziontech/azion-cli/cmd/edge_services/update"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

func NewCmdEdgeServices() *cobra.Command {
	// edgeServicesCmd represents the edgeServices command
	edgeServicesCmd := &cobra.Command{
		Use:   "edge_services",
		Short: "Manages edge services of an Azion account",
		Long:  `You may create, update, delete, list and describe services of an Azion account.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("edgeServices called")
		},
	}

	edgeServicesCmd.AddCommand(create.NewCmd())
	edgeServicesCmd.AddCommand(update.NewCmd())
	edgeServicesCmd.AddCommand(resources.NewCmdResources())

	return edgeServicesCmd
}
