package delete

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/edge_functions_instances"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var functionInstID string
	var applicationID string

	cmd := &cobra.Command{
		Use:           edge_functions_instances.EdgeFuncInstanceDeleteUsage,
		Short:         edge_functions_instances.EdgeFuncInstanceDeleteShortDescription,
		Long:          edge_functions_instances.EdgeFuncInstanceDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		  $ azion edge_functions_instances delete --application-id 1673635839 --instance-id 12312
		  $ azion edge_functions_instances delete -a 1673635839 -i 12312
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("instance-id") {
				return edge_functions_instances.ErrorMissingArgumentsDelete
			}
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.DeleteFunctionInstance(ctx, applicationID, functionInstID)
			if err != nil {
				return fmt.Errorf(edge_functions_instances.ErrorFailToDeleteFuncInst.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, edge_functions_instances.EdgeFuncInstanceDeleteOutputSuccess, functionInstID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&applicationID, "application-id", "a", "", edge_functions_instances.ApplicationFlagId)
	cmd.Flags().StringVarP(&functionInstID, "instance-id", "i", "", edge_functions_instances.EdgeFuncInstanceFlagId)
	cmd.Flags().BoolP("help", "h", false, edge_functions_instances.EdgeFuncInstanceDeleteHelpFlag)
	return cmd
}
