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

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var service_id int64
	// deleteCmd represents the delete command
	deleteCmd := &cobra.Command{
		Use:           msg.EdgeServiceDeleteUsage,
		Short:         msg.EdgeServiceDeleteShortDescription,
		Long:          msg.EdgeServiceDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services delete --service-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("service-id") {
				return msg.ErrorMissingServiceIdArgument
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := deleteService(client, f.IOStreams.Out, service_id); err != nil {
				return err
			}

			return nil
		},
	}

	deleteCmd.Flags().Int64VarP(&service_id, "service-id", "s", 0, msg.EdgeServiceFlagId)
	deleteCmd.Flags().BoolP("help", "h", false, msg.EdgeServiceDeleteFlagHelp)

	return deleteCmd
}

func deleteService(client *sdk.APIClient, out io.Writer, service_id int64) error {

	c := context.Background()
	api := client.DefaultApi

	httpResp, err := api.DeleteService(c, service_id).Execute()

	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)

		return fmt.Errorf(msg.ErrorDeleteService.Error(), message)
	}

	if httpResp.StatusCode == 204 {
		fmt.Fprintf(out, msg.EdgeServiceDeleteOutputSuccess, service_id)
	}

	return nil
}
