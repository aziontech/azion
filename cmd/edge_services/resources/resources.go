package resources

import (
	"fmt"

	"github.com/aziontech/azion-cli/cmd/edge_services/resources/list"
	"github.com/spf13/cobra"
)

func NewCmdResources() *cobra.Command {
	// resourcesCmd represents the resources command
	resourcesCmd := &cobra.Command{
		Use:   "resources",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("resources called")
		},
	}
	resourcesCmd.AddCommand(list.NewCmdList())
	return resourcesCmd
}
