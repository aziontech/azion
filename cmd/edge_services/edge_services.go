package edge_services

import (
	"fmt"

	"github.com/aziontech/azion-cli/cmd/edge_services/resources"
	"github.com/spf13/cobra"
)

func NewCmdEdgeServices() *cobra.Command {
	// edgeServicesCmd represents the edgeServices command
	edgeServicesCmd := &cobra.Command{
		Use:   "edge_services",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("edgeServices called")
		},
	}

	edgeServicesCmd.AddCommand(resources.NewCmdResources())

	return edgeServicesCmd
}
