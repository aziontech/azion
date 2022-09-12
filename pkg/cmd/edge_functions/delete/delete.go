package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_functions/error_messages"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var function_id int64
	cmd := &cobra.Command{
		Use:           "delete <edge_function_id> [flags]",
		Short:         "Deletes an Edge Function",
		Long:          "Deletes an Edge Function based on the id given",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions delete --function-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("function-id") {
				return errmsg.ErrorMissingFunctionIdArgument
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, function_id)
			if err != nil {
				return fmt.Errorf("%w: %s", errmsg.ErrorFailToDeleteFunction, err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, "Edge Function %d was successfully deleted\n", function_id)

			return nil
		},
	}

	cmd.Flags().Int64VarP(&function_id, "function-id", "f", 0, "Unique identifier of the Edge Function")

	return cmd
}
