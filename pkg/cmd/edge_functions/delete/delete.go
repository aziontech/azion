package delete

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/edge_functions"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var function_id int64
	cmd := &cobra.Command{
		Use:           edgefunctions.EdgeFunctionDeleteUsage,
		Short:         edgefunctions.EdgeFunctionDeleteShortDescription,
		Long:          edgefunctions.EdgeFunctionDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion edge_functions delete --function-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("function-id") {
				return edgefunctions.ErrorMissingFunctionIdArgumentDelete
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, function_id)
			if err != nil {
				return fmt.Errorf(edgefunctions.ErrorFailToDeleteFunction.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, edgefunctions.EdgeFunctionDeleteOutputSuccess, function_id)

			return nil
		},
	}

	cmd.Flags().Int64VarP(&function_id, "function-id", "f", 0, edgefunctions.EdgeFunctionFlagId)
	cmd.Flags().BoolP("help", "h", false, edgefunctions.EdgeFunctionDeleteHelpFlag)

	return cmd
}
