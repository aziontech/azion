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
		Use:           "list <service_id> [flags]",
		Short:         "Lists the Resources in a given Edge Service",
		Long:          `Lists the Resources in a given Edge Service`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services resources list 1234 [--details]
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errmsg.ErrorMissingServiceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := listAllResources(client, f.IOStreams.Out, opts, ids[0]); err != nil {
				return err
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(listCmd, opts)

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
		if httpResp == nil || httpResp.StatusCode >= 500 {
			err := utils.CheckStatusCode500Error(err)
			return err
		}
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s", errmsg.ErrorGetResources, string(body))
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
