package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var function_id int64
	cmd := &cobra.Command{
		Use:           msg.DeleteUsage,
		Short:         msg.DeleteShortDescription,
		Long:          msg.DeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion edge_functions delete --function-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("function-id") {
				return msg.ErrorMissingFunctionIdArgumentDelete
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, function_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteFunction.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.DeleteOutputSuccess, function_id)

			return nil
		},
	}

	cmd.Flags().Int64VarP(&function_id, "function-id", "f", 0, msg.FlagId)
	cmd.Flags().BoolP("help", "h", false, msg.DeleteHelpFlag)

	return cmd
}
