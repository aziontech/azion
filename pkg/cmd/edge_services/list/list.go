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

	// listCmd represents the list command
	listCmd := &cobra.Command{
		Use:           edgeservices.EdgeServiceListUsage,
		Short:         edgeservices.EdgeServiceListShortDescription,
		Long:          edgeservices.EdgeServiceListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion edge_services list [--details]
		$ azion edge_services list --order_by "id"
		$ azion edge_services list --page 1  
		$ azion edge_services list --page_size 5
		$ azion edge_services list --sort "asc" 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := listAllServices(client, f.IOStreams.Out, opts); err != nil {
				return err
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(listCmd, opts)
	listCmd.Flags().BoolP("help", "h", false, edgeservices.EdgeServiceListFlagHelp)

	return listCmd
}

func listAllServices(client *sdk.APIClient, out io.Writer, opts *contracts.ListOptions) error {
	c := context.Background()
	api := client.DefaultApi

	fields := []string{"Id", "Name"}
	headers := []string{"ID", "NAME"}

	resp, httpResp, err := api.GetServices(c).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).
		OrderBy(opts.OrderBy).
		Filter(opts.Filter).
		Execute()

	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)

		return fmt.Errorf(edgeservices.ErrorGetServices.Error(), message)
	}

	services := resp.Services

	tp := printer.NewTab(out)
	if opts.Details {
		fields = append(fields, "LastEditor", "UpdatedAt", "Active", "BoundNodes")
		headers = append(headers, "LAST EDITOR", "LAST MODIFIED", "ACTIVE", "BOUND NODES")
	}

	tp.PrintWithHeaders(services, fields, headers)

	return nil
}
