package delete

import (
	"context"

	"github.com/aziontech/azion-cli/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	// deleteCmd represents the delete command
	deleteCmd := &cobra.Command{
		Use:           "delete",
		Short:         "Deletes a resource based on a given resource_id",
		Long:          `Deletes a resource when given a service_id and a resource_id.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return utils.ErrorMissingResourceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0], args[1])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			client, err := requests.CreateClient()
			if err != nil {
				return err
			}

			if err := deleteResource(client, ids[0], ids[1]); err != nil {
				return err
			}

			return nil
		},
	}
	return deleteCmd
}

func deleteResource(client *sdk.APIClient, service_id int64, resource_id int64) error {

	c := context.Background()
	api := client.DefaultApi

	httpResp, err := api.DeleteResource(c, service_id, resource_id).Execute()
	if err != nil {
		if httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		return err
	}

	return nil
}
