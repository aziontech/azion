package update

import (
	"context"
	"fmt"

	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	// listCmd represents the list command
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Updates parameters of an edge service",
		Long: `Receives a name as parameter and creates an edge service with the given name
	Usage: azion_cli edge_services create <EDGE_SERVICE_NAME>`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			active, err := cmd.Flags().GetBool("active")
			if err != nil {
				return err
			}

			client, err := utils.CreateClient()
			if err != nil {
				return err
			}

			if err := updateService(client, id[0], name, active); err != nil {
				return err
			}

			return nil
		},
	}
	updateCmd.Flags().StringP("name", "n", "", "<EDGE_SERVICE_NAME>")
	updateCmd.Flags().BoolP("active", "a", false, "<true|false>")

	return updateCmd
}

func updateService(client *sdk.APIClient, id int64, name string, active bool) error {
	c := context.Background()
	api := client.DefaultApi
	serviceRequest := sdk.UpdateServiceRequest{}
	serviceRequest.SetActive(active)
	serviceRequest.SetName(name)
	//serviceRequest.SetVariables(variables)

	resp, httpResp, err := api.PatchService(c, id).UpdateServiceRequest(serviceRequest).Execute()
	if err != nil {
		if httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}

		return err
	}

	fmt.Printf("ID: %d\tName: %s \n", resp.Id, resp.Name)

	return nil
}
