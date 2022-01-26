package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "delete <edge_function_id> [flags]",
		Short:         "Delete an Edge Function",
		Long:          "Delete an Edge Function based on a given id",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions delete 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("missing edge function id argument")
			}

			httpClient, err := f.HttpClient()
			if err != nil {
				return fmt.Errorf("failed to get http client: %w", err)
			}
			client := api.NewClient(httpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err = client.Delete(ctx, args[0])
			if err != nil {
				return err
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, "Edge Function %s was successfully deleted\n", args[0])

			return nil
		},
	}

	return cmd
}
