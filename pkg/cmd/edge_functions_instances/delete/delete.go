package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_functions_instances"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var functionInstID string
	var applicationID string

	cmd := &cobra.Command{
		Use:           msg.EdgeFuncInstanceDeleteUsage,
		Short:         msg.EdgeFuncInstanceDeleteShortDescription,
		Long:          msg.EdgeFuncInstanceDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		  $ azion edge_functions_instances delete --application-id 1673635839 --instance-id 12312
		  $ azion edge_functions_instances delete -a 1673635839 -i 12312
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("instance-id") {
				return msg.ErrorMissingArgumentsDelete
			}
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.DeleteFunctionInstance(ctx, applicationID, functionInstID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteFuncInst.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.EdgeFuncInstanceDeleteOutputSuccess, functionInstID)
			return nil
		},
	}

	cmd.Flags().StringVarP(&applicationID, "application-id", "a", "", msg.ApplicationFlagId)
	cmd.Flags().StringVarP(&functionInstID, "instance-id", "i", "", msg.EdgeFuncInstanceFlagId)
	cmd.Flags().BoolP("help", "h", false, msg.EdgeFuncInstanceDeleteHelpFlag)
	return cmd
}
