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
		Use:           "delete <service_id> <resource_id> [flags]",
		Short:         "Deletes a Resource",
		Long:          `Deletes a Resource based on the service_id and resource_id given`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services resources delete 1234 81234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errmsg.ErrorMissingResourceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0], args[1])
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

			if err := deleteResource(client, f.IOStreams.Out, ids[0], ids[1], verbose); err != nil {
				return err
			}

			return nil
		},
	}
	return deleteCmd
}

func deleteResource(client *sdk.APIClient, out io.Writer, service_id int64, resource_id int64, verbose bool) error {

	c := context.Background()
	api := client.DefaultApi

	httpResp, err := api.DeleteResource(c, service_id, resource_id).Execute()
	if err != nil {
		if httpResp == nil || httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s", errmsg.ErrorDeleteResource, string(body))
	}

	if verbose {
		if httpResp.StatusCode == 204 {
			fmt.Fprintf(out, "Resource %d was successfully deleted\n", resource_id)
		}
	}

	return nil
}
