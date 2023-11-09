package edge_functions

import (
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/create"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	edgeFunctionsCmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	edgeFunctionsCmd.AddCommand(create.NewCmd(f))
	edgeFunctionsCmd.AddCommand(update.NewCmd(f))
	edgeFunctionsCmd.AddCommand(describe.NewCmd(f))
	edgeFunctionsCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	edgeFunctionsCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return edgeFunctionsCmd
}
