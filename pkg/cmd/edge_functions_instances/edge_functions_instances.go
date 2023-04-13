package edge_functions_instances

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_functions_instances"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions_instances/create"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions_instances/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions_instances/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions_instances/list"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions_instances/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	funcInstCmd := &cobra.Command{
		Use:   msg.EdgeFuncInstanceUsage,
		Short: msg.EdgeFuncInstanceShortDescription,
		Long:  msg.EdgeFuncInstanceLongDescription,
		Example: heredoc.Doc(`
		$ azioncli rules_engine --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	funcInstCmd.AddCommand(delete.NewCmd(f))
	funcInstCmd.AddCommand(create.NewCmd(f))
	funcInstCmd.AddCommand(describe.NewCmd(f))
	funcInstCmd.AddCommand(list.NewCmd(f))
	funcInstCmd.AddCommand(update.NewCmd(f))

	funcInstCmd.Flags().BoolP("help", "h", false, msg.EdgeFuncInstanceFlagHelp)
	return funcInstCmd
}
