package delete

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/edge_services"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

type Fields struct {
	ServiceId  int64
	ResourceId int64
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}
	// deleteCmd represents the delete command
	deleteCmd := &cobra.Command{
		Use:           edgeservices.EdgeServiceResourceDeleteUsage,
		Short:         edgeservices.EdgeServiceResourceDeleteShortDescription,
		Long:          edgeservices.EdgeServiceResourceDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion edge_services resources delete --service-id 1234 --resource-id 81234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("service-id") || !cmd.Flags().Changed("resource-id") {
				return edgeservices.ErrorMissingResourceIdArgument
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := deleteResource(client, f.IOStreams.Out, fields.ServiceId, fields.ResourceId); err != nil {
				return err
			}

			return nil
		},
	}

	deleteCmd.Flags().Int64VarP(&fields.ServiceId, "service-id", "s", 0, edgeservices.EdgeServiceFlagId)
	deleteCmd.Flags().Int64VarP(&fields.ResourceId, "resource-id", "r", 0, edgeservices.EdgeServiceResourceFlagId)
	deleteCmd.Flags().BoolP("help", "h", false, edgeservices.EdgeServiceResourceDeleteFlagHelp)

	return deleteCmd
}

func deleteResource(client *sdk.APIClient, out io.Writer, service_id int64, resource_id int64) error {

	c := context.Background()
	api := client.DefaultApi

	httpResp, err := api.DeleteResource(c, service_id, resource_id).Execute()
	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)
		return fmt.Errorf(edgeservices.ErrorDeleteResource.Error(), message)
	}

	if httpResp.StatusCode == 204 {
		fmt.Fprintf(out, edgeservices.EdgeServiceResourceDeleteOutputSuccess, resource_id)
	}

	return nil
}
