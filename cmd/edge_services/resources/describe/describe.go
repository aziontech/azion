package describe

import (
	"context"
	"fmt"

	"github.com/aziontech/azion-cli/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {

	// describeCmd represents the describe command
	describeCmd := &cobra.Command{
		Use:           "describe",
		Short:         "Describes a resource based on a given resource_id",
		Long:          `Provides a long desription of a resource based on a given resource_id`,
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
			if err := describeResource(client, ids[0], ids[1]); err != nil {
				return err
			}

			return nil

		},
	}
	return describeCmd

}

func describeResource(client *sdk.APIClient, service_id int64, resource_id int64) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetResource(c, service_id, resource_id).Execute()
	if err != nil {
		if httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		return err
	}

	fmt.Printf("ID: %d\n", resp.Id)
	fmt.Printf("Name: %s\n", resp.Name)
	fmt.Printf("Type: %s\n", resp.Type)
	fmt.Printf("Content type: %s\n", resp.ContentType)
	fmt.Printf("Content: \n")
	fmt.Printf("%s", resp.Content)

	return nil
}
