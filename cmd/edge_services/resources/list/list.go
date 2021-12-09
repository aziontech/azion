package list

import (
	"context"
	"errors"
	"fmt"

	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

func NewCmdList() *cobra.Command {
	// listCmd represents the list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Lists resources in a given service",
		Long: `Lists all resources found in a service by providing a service_id.
	Service_id can be found by listing your services`,
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) < 1 {
				return errors.New(utils.ErrorMissingServiceIdArgument.Error())
			}

			service_id := args[0]

			fmt.Println(service_id)

			client, err := utils.CreateRequest()
			if err != nil {
				return err
			}
			getListAllResources(client, 10)
			return nil
		},
	}
	return listCmd
}

func getListAllResources(client *sdk.APIClient, service_id int64) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpresp, err := api.GetResources(c, service_id).Execute()
	if err != nil {
		return err
	}

	resources := resp.Resources

	for _, resource := range resources {
		fmt.Println()
	}

}
