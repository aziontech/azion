package delete

import (
	"context"
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var service_id int64
	// deleteCmd represents the delete command
	deleteCmd := &cobra.Command{
		Use:           "delete --service-id <service_id> [flags]",
		Short:         "Deletes an Edge Service",
		Long:          `Deletes an Edge Service based on the id given`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services delete --service-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("service-id") {
				return errmsg.ErrorMissingServiceIdArgument
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := deleteService(client, f.IOStreams.Out, service_id); err != nil {
				return err
			}

			return nil
		},
	}

	deleteCmd.Flags().Int64VarP(&service_id, "service-id", "s", 0, "Unique identifier of the Edge Service")

	return deleteCmd
}

func deleteService(client *sdk.APIClient, out io.Writer, service_id int64) error {

	c := context.Background()
	api := client.DefaultApi

	httpResp, err := api.DeleteService(c, service_id).Execute()

	if err != nil {
		errMsg := utils.ErrorPerStatusCode(httpResp, err)

		return fmt.Errorf("%w: %s", errmsg.ErrorDeleteService, errMsg)
	}

	if httpResp.StatusCode == 204 {
		fmt.Fprintf(out, "Service %d was successfully deleted\n", service_id)
	}

	return nil
}
