package describe

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

	// describeCmd represents the describe command
	describeCmd := &cobra.Command{
		Use:           "describe",
		Short:         "Describes a service based on a given service_id",
		Long:          `FIXME with usage`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return utils.ErrorMissingServiceIdArgument
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

			if err := describeService(client, f.IOStreams.Out, ids[0]); err != nil {
				return err
			}

			return nil

		},
	}
	return describeCmd

}

func describeService(client *sdk.APIClient, out io.Writer, service_id int64) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetService(c, service_id).Execute()
	if err != nil {
		if httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		return err
	}

	fmt.Fprintf(out, "ID: %d\n", resp.Id)
	fmt.Fprintf(out, "Name: %s\n", resp.Name)
	fmt.Fprintf(out, "Last Editor: %s\n", resp.LastEditor)
	fmt.Fprintf(out, "Updated at: %s\n", resp.UpdatedAt)
	fmt.Fprintf(out, "Active: %t\n", resp.Active)
	fmt.Fprintf(out, "Bound Nodes: %d\n", resp.BoundNodes)
	return nil
}
