package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_functions/error_messages"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "delete <edge_function_id> [flags]",
		Short:         "Deletes an Edge Function",
		Long:          "Deletes an Edge Function based on the id given",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions delete 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errmsg.ErrorMissingFunctionIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err = client.Delete(ctx, ids[0])
			if err != nil {
				return fmt.Errorf("%s: %w", errmsg.ErrorFailToDeleteFunction, err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, "Edge Function %s was successfully deleted\n", args[0])

			return nil
		},
	}

	return cmd
}
