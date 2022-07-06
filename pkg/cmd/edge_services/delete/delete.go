package delete

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/MakeNowJust/heredoc"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
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
		Long:          `Deletes an Edge Service based on the id given`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services delete 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errmsg.ErrorMissingServiceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := deleteService(client, f.IOStreams.Out, ids[0]); err != nil {
				return err
			}

			return nil
		},
	}
	return deleteCmd
}

func deleteService(client *sdk.APIClient, out io.Writer, service_id int64) error {

	c := context.Background()
	api := client.DefaultApi

	httpResp, err := api.DeleteService(c, service_id).Execute()

	if err != nil {
		if httpResp == nil || httpResp.StatusCode >= 500 {
			err := utils.CheckStatusCode500Error(err)
			return err
		}
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s", errmsg.ErrorDeleteService, string(body))
	}

	if httpResp.StatusCode == 204 {
		fmt.Fprintf(out, "Service %d was successfully deleted\n", service_id)
	}

	return nil
}
