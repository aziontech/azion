package delete

import (
	"context"

	"github.com/aziontech/azion-cli/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// deleteCmd represents the delete command
	deleteCmd := &cobra.Command{
		Use:           "delete",
		Short:         "Deletes a service based on a given service_id",
		Long:          `FIXME with USAGE`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return utils.ErrorMissingResourceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			tok, err := cmd.Flags().GetString("token")
			if err != nil {
				return err
			}

			httpClient, err := f.HttpClient()
			if err != nil {
				return err
			}

			client, err := requests.CreateClient(httpClient, tok)
			if err != nil {
				return err
			}

			if err := deleteService(client, ids[0]); err != nil {
				return err
			}

			return nil
		},
	}
	return deleteCmd
}

func deleteService(client *sdk.APIClient, service_id int64) error {

	c := context.Background()
	api := client.DefaultApi

	httpResp, err := api.DeleteService(c, service_id).Execute()
	if err != nil {
		if httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		return err
	}

	return nil
}
