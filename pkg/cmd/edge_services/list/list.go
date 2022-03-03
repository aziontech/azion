package list

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/MakeNowJust/heredoc"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
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
		Use:           "list [flags]",
		Short:         "Lists your account's Edge Services",
		Long:          `Lists your account's Edge Services`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services list [--details]
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
		if httpResp != nil && httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s", errmsg.ErrorGetServices, string(body))
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
