package delete

import (
	"context"
	"fmt"
	"io"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_services"
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
		Use:           msg.EdgeServiceResourceDeleteUsage,
		Short:         msg.EdgeServiceResourceDeleteShortDescription,
		Long:          msg.EdgeServiceResourceDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services resources delete --service-id 1234 --resource-id 81234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("service-id") || !cmd.Flags().Changed("resource-id") {
				return msg.ErrorMissingResourceIdArgument
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

	deleteCmd.Flags().Int64VarP(&fields.ServiceId, "service-id", "s", 0, msg.EdgeServiceFlagId)
	deleteCmd.Flags().Int64VarP(&fields.ResourceId, "resource-id", "r", 0, msg.EdgeServiceResourceFlagId)
	deleteCmd.Flags().BoolP("help", "h", false, msg.EdgeServiceResourceDeleteFlagHelp)

	return deleteCmd
}

func deleteResource(client *sdk.APIClient, out io.Writer, service_id int64, resource_id int64) error {

	c := context.Background()
	api := client.DefaultApi

	httpResp, err := api.DeleteResource(c, service_id, resource_id).Execute()
	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)

		return fmt.Errorf("%w: %s", msg.ErrorDeleteResource, message)
	}

	if httpResp.StatusCode == 204 {
		fmt.Fprintf(out, msg.EdgeServiceResourceDeleteOutputSuccess, resource_id)
	}

	return nil
}
