package edge_functions

import (
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/create"
	del "github.com/aziontech/azion-cli/pkg/cmd/edge_functions/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/list"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	msg "github.com/aziontech/azion-cli/pkg/messages/edge_functions"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	edgeFunctionsCmd := &cobra.Command{
		Use:   msg.EdgeFunctionUsage,
		Short: msg.EdgeFunctionShortDescription,
		Long:  msg.EdgeFunctionLongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	edgeFunctionsCmd.AddCommand(create.NewCmd(f))
	edgeFunctionsCmd.AddCommand(del.NewCmd(f))
	edgeFunctionsCmd.AddCommand(update.NewCmd(f))
	edgeFunctionsCmd.AddCommand(describe.NewCmd(f))
	edgeFunctionsCmd.AddCommand(list.NewCmd(f))
	edgeFunctionsCmd.Flags().BoolP("help", "h", false, msg.EdgeFunctionHelpFlag)

	return edgeFunctionsCmd
}
