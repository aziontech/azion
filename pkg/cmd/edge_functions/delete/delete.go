package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_functions"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var function_id int64
	cmd := &cobra.Command{
		Use:           msg.EdgeFunctionDeleteUsage,
		Short:         msg.EdgeFunctionDeleteShortDescription,
		Long:          msg.EdgeFunctionDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions delete --function-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("function-id") {
				return msg.ErrorMissingFunctionIdArgument
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, function_id)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorFailToDeleteFunction, err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.EdgeFunctionDeleteOutputSuccess, function_id)

			return nil
		},
	}

	cmd.Flags().Int64VarP(&function_id, "function-id", "f", 0, msg.EdgeFunctionFlagId)
	cmd.Flags().BoolP("help", "h", false, msg.EdgeFunctionDeleteHelpFlag)

	return cmd
}
