package create

import (
	"context"
	"fmt"

	"github.com/aziontech/azion-cli/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	// listCmd represents the list command
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new edge service",
		Long: `Receives a name as parameter and creates an edge service with the given name
	Usage: azion_cli edge_services create <EDGE_SERVICE_NAME>`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			client, err := requests.CreateClient()
			if err != nil {
				return err
			}

			if err := createNewService(client, name); err != nil {
				return err
			}

			return nil
		},
	}
	createCmd.Flags().StringP("name", "n", "", "<EDGE_SERVICE_NAME>")
	createCmd.MarkFlagRequired("name")

	return createCmd
}

func createNewService(client *sdk.APIClient, name string) error {
	c := context.Background()
	api := client.DefaultApi
	serviceRequest := sdk.CreateServiceRequest{}
	serviceRequest.SetName(name)

	resp, httpResp, err := api.NewService(c).CreateServiceRequest(serviceRequest).Execute()
	if err != nil {
		if httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}

		return err
	}

	fmt.Printf("ID: %d\tName: %s \n", resp.Id, resp.Name)

	return nil
}
