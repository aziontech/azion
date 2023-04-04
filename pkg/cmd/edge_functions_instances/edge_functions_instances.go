package edge_functions_instances

import (
	msg "github.com/aziontech/azion-cli/messages/edge_functions_instances"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions_instances/list"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	edgeFunctionsCmd := &cobra.Command{
		Use:   msg.EdgeFunctionsInstancesUsage,
		Short: msg.EdgeFunctionsInstancesShortDescription,
		Long:  msg.EdgeFunctionsInstancesLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	edgeFunctionsCmd.AddCommand(list.NewCmd(f))
	edgeFunctionsCmd.Flags().BoolP("help", "h", false, msg.EdgeFunctionsInstancesFlagHelp)

	return edgeFunctionsCmd
}
