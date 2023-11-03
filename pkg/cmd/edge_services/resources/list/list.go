package list

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/edge_services"
	"io"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/printer"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	var service_id int64

	// listCmd represents the list command
	listCmd := &cobra.Command{
		Use:           edgeservices.EdgeServiceResourceListUsage,
		Short:         edgeservices.EdgeServiceResourceListShortDescription,
		Long:          edgeservices.EdgeServiceResourceListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion edge_services resources list --service-id 1234 [--details]
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("service-id") {
				return edgeservices.ErrorMissingServiceIdArgument
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := listAllResources(client, f.IOStreams.Out, opts, service_id); err != nil {
				return err
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(listCmd, opts)
	listCmd.Flags().Int64VarP(&service_id, "service-id", "s", 0, "Unique identifier of the Edge Service")
	listCmd.Flags().BoolP("help", "h", false, edgeservices.EdgeServiceResourceListFlagHelp)

	return listCmd
}

func listAllResources(client *sdk.APIClient, out io.Writer, opts *contracts.ListOptions, service_id int64) error {
	c := context.Background()
	api := client.DefaultApi

	fields := []string{"Id", "Name"}
	headers := []string{"ID", "NAME"}

	resp, httpResp, err := api.GetResources(c, service_id).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).
		OrderBy(opts.OrderBy).
		Filter(opts.Filter).
		Execute()

	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)

		return fmt.Errorf(edgeservices.ErrorGetResources.Error(), message)
	}

	resources := resp.Resources

	tp := printer.NewTab(out)
	if opts.Details {
		fields = append(fields, "LastEditor", "UpdatedAt", "ContentType", "Trigger")
		headers = append(headers, "LAST EDITOR", "LAST MODIFIED", "CONTENT TYPE", "TRIGGER")
	}

	tp.PrintWithHeaders(resources, fields, headers)

	return nil
}
