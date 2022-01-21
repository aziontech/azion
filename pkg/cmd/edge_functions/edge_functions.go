package edge_functions

import (
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/create"
	del "github.com/aziontech/azion-cli/pkg/cmd/edge_functions/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/list"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	edgeFunctionsCmd := &cobra.Command{
		Use:   "edge_functions",
		Short: "Manages Edge Functions of your Azion account",
		Long:  "You may create, update, delete, list and describe Edge Functions of an Azion account",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		Annotations: map[string]string{
			"IsAPI": "true",
		},
	}

	edgeFunctionsCmd.AddCommand(create.NewCmd(f))
	edgeFunctionsCmd.AddCommand(del.NewCmd(f))
	edgeFunctionsCmd.AddCommand(update.NewCmd(f))
	edgeFunctionsCmd.AddCommand(describe.NewCmd(f))
	edgeFunctionsCmd.AddCommand(list.NewCmd(f))

	return edgeFunctionsCmd
}
