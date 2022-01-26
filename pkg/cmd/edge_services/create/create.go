package create

import (
	"context"
	"fmt"
	"io"

	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	createCmd := &cobra.Command{
		Use:           "create [flags]",
		Short:         "Creates a new Edge Service",
		Long:          `Creates a new Edge Service with the received name`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}

			if err := createNewService(client, f.IOStreams.Out, name, verbose); err != nil {
				return err
			}

			return nil
		},
	}
	createCmd.Flags().String("name", "", "<EDGE_SERVICE_NAME>")
	_ = createCmd.MarkFlagRequired("name")

	return createCmd
}

func createNewService(client *sdk.APIClient, out io.Writer, name string, verbose bool) error {
	c := context.Background()
	api := client.DefaultApi
	serviceRequest := sdk.CreateServiceRequest{}
	serviceRequest.SetName(name)

	resp, httpResp, err := api.NewService(c).CreateServiceRequest(serviceRequest).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}

		return err
	}
	if verbose {
		fmt.Fprintf(out, "ID: %d\tName: %s \n", resp.Id, resp.Name)
	}

	return nil
}
