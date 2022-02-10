package delete

import (
	"context"
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// deleteCmd represents the delete command
	deleteCmd := &cobra.Command{
		Use:           "delete <service_id> [flags]",
		Short:         "Deletes an Edge Service",
		Long:          `Deletes an Edge Service based on a given service_id`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services delete 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return utils.ErrorMissingServiceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}

			if err := deleteService(client, f.IOStreams.Out, ids[0], verbose); err != nil {
				return err
			}

			return nil
		},
	}
	return deleteCmd
}

func deleteService(client *sdk.APIClient, out io.Writer, service_id int64, verbose bool) error {

	c := context.Background()
	api := client.DefaultApi

	httpResp, err := api.DeleteService(c, service_id).Execute()

	if err != nil {
		if httpResp != nil && httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		return err
	}

	if verbose {
		if httpResp.StatusCode == 204 {
			fmt.Fprintf(out, "Service %d was successfully deleted\n", service_id)
		}
	}

	return nil
}
