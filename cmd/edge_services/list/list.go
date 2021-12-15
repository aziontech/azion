package list

import (
	"context"
	"fmt"
	"io"

	"github.com/aziontech/azion-cli/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// listCmd represents the list command
	listCmd := &cobra.Command{
		Use:           "list",
		Short:         "Lists services of an Azion account",
		Long:          `FIXME with usage`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {

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

			if err := listAllServices(client, f.IOStreams.Out); err != nil {
				return err
			}
			return nil
		},
	}
	return listCmd
}

func listAllServices(client *sdk.APIClient, out io.Writer) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetServices(c).Execute()
	if err != nil {
		if httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		return err
	}

	services := resp.Services

	for _, service := range services {
		fmt.Fprintf(out, "ID: %d     Name: %s \n", service.Id, service.Name)
	}
	return nil
}
